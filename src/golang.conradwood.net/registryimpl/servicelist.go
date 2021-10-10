package registryimpl

import (
	"context"
	"flag"
	"fmt"
	"time"
	//	"golang.conradwood.net/apis/common"
	reg "golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/utils"
	"sync"
)

var (
	targetable_timeout = flag.Int("refresh_timeout", 60, "timeout in `seconds` after which a registration becomes stale and will not be offered as target any more")
	instancesequence   = uint64(1)
	instanceModifyLock sync.Mutex // locking when we modify instances
	infoLock           sync.Mutex
)

type ServiceList struct {
	instances []*serviceInstance
	services  map[string]*serviceInfo
}
type serviceInfo struct {
	list      *ServiceList
	created   time.Time
	ip        IP
	createdAs *reg.CreateServiceRequest
	lastUsed  time.Time
}
type serviceInstance struct {
	list                 *ServiceList // backpointer
	sequencenumber       uint64
	IP                   IP
	pid                  uint64
	createdAs            *reg.CreateServiceRequest
	registeredAs         *reg.RegisterServiceRequest
	refreshed            time.Time
	isLocal              bool
	isFallback           bool
	isRemote             bool // someone else added it, e.g. "grpc.conradwood.net:..". remotes are also fallback-only
	deregistered         bool
	didQueryAutodeployer bool
	serviceCheckFailures int
	lastServiceCheck     time.Time
}

func (si *serviceInstance) IncludesApiType(apitype reg.Apitype) bool {
	ras := si.registeredAs
	if ras == nil {
		return false
	}
	for _, a := range ras.ApiType {
		if a == apitype {
			return true
		}
	}
	return false

}
func (si *serviceInstance) BuildID() uint64 {
	csr := si.createdAs
	if csr != nil && csr.DeployInfo != nil {
		return csr.DeployInfo.BuildID
	}
	return 0
}
func (si *serviceInstance) Port() uint32 {
	ras := si.registeredAs
	if ras != nil {
		return ras.Port
	}
	return 0
}
func (si *serviceInstance) ServiceName() string {
	ras := si.registeredAs
	if ras != nil {
		return ras.ServiceName
	}
	return ""
}

// may this instance be used to connect to?
func (si *serviceInstance) Targetable() bool {
	if si.isLocal {
		return true
	}
	if si.isFallback || si.isRemote {
		// serve only if we have zero non-fall back ones available
		sameName := si.list.FilterByNames([]string{si.ServiceName()})
		ct := 0
		for _, sc := range sameName.Instances() {
			if sc.isFallback || si.isRemote {
				continue
			}
			ct++
		}
		if ct == 0 {
			return true
		} else {
			return false
		}
	}
	if time.Since(si.refreshed) > time.Duration(*targetable_timeout)*time.Second {
		return false
	}
	return true
}

func (si *serviceInstance) RemoveMe() bool {
	if si.isLocal || si.isFallback || si.isRemote {
		return false
	}
	if si.serviceCheckFailures > 10 {
		return true
	}
	if time.Since(si.refreshed) < time.Duration(*targetable_timeout*5)*time.Second {
		return false
	}
	return true
}
func (si *serviceInstance) Running() bool {
	if si.isLocal || si.isFallback || si.isRemote {
		return true
	}
	if time.Since(si.refreshed) > time.Duration(*targetable_timeout*3)*time.Second {
		return false
	}
	return true
}
func (si *serviceInstance) Registration() *reg.Registration {
	r := &reg.Registration{
		Target:        si.Target(),
		Targetable:    si.Targetable(),
		Pid:           si.pid,
		LastRefreshed: uint32(si.refreshed.Unix()),
		Running:       si.Running(),
	}
	csr := si.createdAs
	if si.isLocal {
		r.LastRefreshed = uint32(time.Now().Unix())
	}
	ras := si.registeredAs
	if ras != nil {
		r.ProcessID = ras.ProcessID
	}

	if csr != nil {
		r.ProcessID = csr.ProcessID
		r.DeployInfo = csr.DeployInfo
		r.Partition = csr.Partition
	}
	return r
}

func (si *serviceInstance) Target() *reg.Target {
	t := &reg.Target{ServiceName: si.ServiceName()}

	ras := si.registeredAs
	if ras != nil {
		t.IP = si.IP.ExposeAs()
		t.Port = ras.Port
		t.ApiType = ras.ApiType
		t.RoutingInfo = ras.RoutingInfo
	}
	return t
}
func sequence() uint64 {
	instanceModifyLock.Lock()
	instancesequence++
	res := instancesequence
	instanceModifyLock.Unlock()
	return res
}

// a fallback is a service which is _only_ served as long as there are no
// targetable registrations for this service available
func (s *ServiceList) AddFallbackWithHost(Name string, host string, port uint32) {
	si := &serviceInstance{
		list:           s,
		isFallback:     true,
		sequencenumber: sequence(),
		registeredAs: &reg.RegisterServiceRequest{
			ProcessID:   "fallback",
			Port:        port,
			ApiType:     []reg.Apitype{reg.Apitype_grpc, reg.Apitype_status},
			ServiceName: Name,
			RoutingInfo: nil,
		},
	}
	si.IP = IPFromLiteralString(host)
	s.instances = append(s.instances, si)
}

// a fallback is a service which is _only_ served as long as there are no
// targetable registrations for this service available
func (s *ServiceList) AddFallback(Name string, port uint32) {
	si := &serviceInstance{
		list:           s,
		isFallback:     true,
		sequencenumber: sequence(),
		registeredAs: &reg.RegisterServiceRequest{
			ProcessID:   "fallback",
			Port:        port,
			ApiType:     []reg.Apitype{reg.Apitype_grpc, reg.Apitype_status},
			ServiceName: Name,
			RoutingInfo: nil,
		},
	}
	si.IP = IPLocal()
	s.instances = append(s.instances, si)
}
func (s *ServiceList) AddRemoteWithHost(Name string, host string, port uint32) {
	si := &serviceInstance{
		list:           s,
		isRemote:       true,
		sequencenumber: sequence(),
		registeredAs: &reg.RegisterServiceRequest{
			ProcessID:   fmt.Sprintf("R-%s-%s:%d", Name, host, port),
			Port:        port,
			ApiType:     []reg.Apitype{reg.Apitype_grpc, reg.Apitype_status},
			ServiceName: Name,
			RoutingInfo: nil,
		},
	}
	si.IP = IPFromLiteralString(host)
	s.instances = append(s.instances, si)
}
func (s *ServiceList) AddLocal(Name string, port uint32) {
	si := &serviceInstance{
		list:           s,
		isLocal:        true,
		sequencenumber: sequence(),
		registeredAs: &reg.RegisterServiceRequest{
			ProcessID:   "local",
			Port:        port,
			ApiType:     []reg.Apitype{reg.Apitype_grpc, reg.Apitype_status},
			ServiceName: Name,
			RoutingInfo: nil,
		},
	}
	si.IP = IPLocal()
	s.instances = append(s.instances, si)
}

func (s *ServiceList) Instances() []*serviceInstance {
	if s == nil {
		return make([]*serviceInstance, 0)
	}
	return s.instances
}

// only return instances that are active (e.g. refreshed recently),
// not being shutdown or otherwise stupid
func (s *ServiceList) TargetableInstances() []*serviceInstance {
	if s == nil {
		return make([]*serviceInstance, 0)
	}
	var res []*serviceInstance
	for _, si := range s.instances {
		if !si.Targetable() {
			continue
		}
		res = append(res, si)
	}
	return res
}

// autodeployer send us a thing
func (s *ServiceList) Create(ctx context.Context, req *reg.CreateServiceRequest) error {
	s.Printf("Create service processid=%s (got deployinfo: %v)\n", req.ProcessID, (req.DeployInfo != nil))
	ip, err := IPFromContext(ctx)
	if err != nil {
		return err
	}
	infoLock.Lock()
	defer infoLock.Unlock()
	sin := s.services[req.ProcessID]
	if sin != nil {
		return errors.InvalidArgs(ctx, "exists already", "process %s exists already", req.ProcessID)
	}
	si := &serviceInfo{
		ip:        ip,
		list:      s,
		createdAs: req,
		created:   time.Now(),
		lastUsed:  time.Now(),
	}
	s.services[req.ProcessID] = si
	return nil
}

// find a given service instance and update it or create it. (true if it is a new service, false if refresh)
// TODO: we should verify ports here if they match create request
func (s *ServiceList) Registration(ctx context.Context, req *reg.RegisterServiceRequest) (*serviceInstance, bool, error) {
	if req.ProcessID == "" {
		return nil, false, errors.InvalidArgs(ctx, "missing processid", "missing processid for \"%s\" (pid=%d)", req.ServiceName, req.Pid)
	}
	s.Printf("Refreshing %s with processid \"%s\"\n", req.ServiceName, req.ProcessID)
	if req.ServiceName == "" {
		return nil, false, errors.InvalidArgs(ctx, "missing servicename", "missing servicename")
	}
	ip, err := IPFromContext(ctx)
	if err != nil {
		return nil, false, err
	}
	sin := s.services[req.ProcessID]
	if sin != nil {
		sin.lastUsed = time.Now()
	}
	si := s.findExisting(req, ip)
	if si != nil {
		s.Printf("Refreshed #%d (%s)\n", si.sequencenumber, si.ServiceName())
		si.registeredAs = req
		s.refreshAllWithProcessID(req.ProcessID)
		return si, false, nil
	}
	si = &serviceInstance{list: s, pid: req.Pid, sequencenumber: sequence(), refreshed: time.Now()}
	si.registeredAs = req
	si.IP = ip
	s.instances = append(s.instances, si)
	if sin != nil {
		si.createdAs = sin.createdAs
		if !sin.ip.Equals(ip) {
			return nil, false, errors.AccessDenied(ctx, "service %s was created from host %s, registration came from %s", si.ServiceName(), sin.ip.ExposeAs(), ip.ExposeAs())
		}
	}
	s.refreshAllWithProcessID(req.ProcessID)
	s.Printf("Registered #%d (%s)\n", si.sequencenumber, si.ServiceName())
	return si, true, nil
}
func (s *ServiceList) refreshAllWithProcessID(procid string) {
	for _, si := range s.instances {
		ras := si.registeredAs
		if ras == nil || ras.ProcessID != procid {
			continue
		}
		si.refreshed = time.Now()
	}
}

// prefer targetable ones :)
func (s *ServiceList) findExisting(req *reg.RegisterServiceRequest, ip IP) *serviceInstance {
	// find by processid
	servicename := req.ServiceName
	port := req.Port
	for _, si := range s.instances {
		ras := si.registeredAs
		if ras != nil && ras.ProcessID != req.ProcessID {
			continue
		}
		if si.ServiceName() != servicename {
			continue
		}
		if !si.IP.Equals(ip) {
			continue
		}
		if si.Port() != port {
			continue
		}
		return si
	}

	// find by port name ip only
	for _, si := range s.instances {
		if si.ServiceName() != servicename {
			continue
		}
		if !si.IP.Equals(ip) {
			continue
		}
		if si.Port() != port {
			continue
		}
		return si
	}
	return nil
}

// returns all service instances which were deregistered.
func (s *ServiceList) Deregister(ctx context.Context, processid string) []*serviceInstance {
	s.Printf("Deregistering processpid \"%s\"\n", processid)
	var res []*serviceInstance
	instanceModifyLock.Lock()
	defer instanceModifyLock.Unlock()
	var newi []*serviceInstance
	for _, si := range s.instances {
		ras := si.registeredAs
		if ras != nil && ras.ProcessID == processid {
			if si.deregistered == false {
				res = append(res, si)
			}
			si.deregistered = true
			continue
		}
		newi = append(newi, si)
	}
	s.instances = newi
	return res
}

// old v1 style compatibility - remove ip / port (instead of processpid)
func (s *ServiceList) removeIPPort(ip IP, port uint32) {
	s.Printf("Request to deregister service @%s:%d\n", ip.ExposeAs(), port)
	instanceModifyLock.Lock()
	defer instanceModifyLock.Unlock()
	var newi []*serviceInstance
	for _, si := range s.instances {
		if si.Port() == port && si.IP.Equals(ip) {
			s.Printf("IP/Port based removal of service %s@%s:%d\n", si.ServiceName(), ip.ExposeAs(), port)
			continue
		}
		newi = append(newi, si)
	}
	s.instances = newi
}

/********************************************************************************
* filters
********************************************************************************/
func (s *ServiceList) Clone() *ServiceList {
	res := &ServiceList{instances: s.instances}
	return res
}
func (s *ServiceList) FilterByNames(servicenames []string) *ServiceList {
	var newi []*serviceInstance
	for _, si := range s.instances {
		keep := false
		for _, sn := range servicenames {
			if sn == si.ServiceName() {
				keep = true
				break
			}
		}
		if keep {
			newi = append(newi, si)
		}
	}
	res := s.Clone()
	res.instances = newi
	return res
}
func (s *ServiceList) FilterByApiType(ApiType reg.Apitype) *ServiceList {
	var newi []*serviceInstance
	for _, si := range s.instances {
		if si.IncludesApiType(ApiType) {
			newi = append(newi, si)
		}
	}
	res := s.Clone()
	res.instances = newi
	return res
}

// regularly clean out (deregister) stuff that hasn't been registered with us for a while
func (s *V2Registry) Cleaner() {
	for {
		utils.RandomStall(1)
		s.clean()
	}
}
func (rv *V2Registry) clean() {
	s := rv.serviceList
	var newi []*serviceInstance
	instanceModifyLock.Lock()
	b := false
	for _, si := range s.instances {
		if si.RemoveMe() {
			b = true
			continue
		}
		newi = append(newi, si)
	}
	s.instances = newi
	instanceModifyLock.Unlock()
	if b {
		rv.promUpdate()
	}
}

func (s *ServiceList) Printf(format string, args ...interface{}) {
	if !*debug {
		return
	}
	txt := "[v2 servicelist] "
	fmt.Printf(txt+format, args...)
}
