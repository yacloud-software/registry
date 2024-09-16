package registryimpl

import (
	"context"
	"flag"
	"fmt"
	"time"

	//	"golang.conradwood.net/apis/common"
	"sync"

	"golang.conradwood.net/apis/common"
	reg "golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/utils"
)

var (
	start_ready         = flag.Bool("assume_service_is_ready", false, "if true assume service is ready as soon as it registers,otherwise do an http call for /health to it first")
	do_dump_servicelist = flag.Bool("dump_service_list", false, "if true, dump the entire service list sometimes. lots of debug")
	targetable_timeout  = flag.Int("refresh_timeout", 60, "timeout in `seconds` after which a registration becomes stale and will not be offered as target any more")
	instancesequence    = uint64(1)
	instanceModifyLock  sync.Mutex // locking when we modify instances
	infoLock            sync.Mutex
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
	//	serviceReady         bool // set by "verify status", based on service health exposed http value
	health           common.Health
	lastServiceCheck time.Time
}

func (si *serviceInstance) String() string {
	x := "(noname)"
	if si != nil {
		if si.registeredAs != nil {
			x = si.registeredAs.ServiceName
		}
	}
	return fmt.Sprintf("[%s:seq=%d,ip=%s,pid=%d]", x, si.sequencenumber, si.IP.String(), si.pid)
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
	if !si.serviceReady() && si.IncludesApiType(reg.Apitype_status) {
		return false
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
		r.UserID = ras.UserID
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
	s.dump("After Addfallbackwithhost:")
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
	s.dump("After Addfallback:")
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
	s.dump("After Addremotewithhost:")
}
func (s *ServiceList) AddLocal(Name string, port uint32) {
	s.debugf(debugnd(Name), "adding local instance %s at port %d\n", Name, port)
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
	s.dump("After addlocal:")
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
	s.debugf(debugnd("none"), "Create service processid=%s (got deployinfo: %v)\n", req.ProcessID, (req.DeployInfo != nil))
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
	s.dump("After create:")
	return nil
}

// find a given service instance and update it or create it. (true if it is a new service, false if refresh)
// TODO: we should verify ports here if they match create request
func (s *ServiceList) Registration(ctx context.Context, req *reg.RegisterServiceRequest) (*serviceInstance, bool, error) {
	if req.ProcessID == "" {
		return nil, false, errors.InvalidArgs(ctx, "missing processid", "missing processid for \"%s\" (pid=%d)", req.ServiceName, req.Pid)
	}
	s.debugf(req, "Refreshing %s with processid \"%s\"\n", req.ServiceName, req.ProcessID)
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
		s.debugf(req, "Refreshed #%d (%s)\n", si.sequencenumber, si.ServiceName())
		si.registeredAs = req
		s.refreshAllWithProcessID(req.ProcessID)
		return si, false, nil
	}
	s.debugf(req, "Refreshed processid (%s) - new service instance\n", req.ProcessID)
	si = &serviceInstance{list: s, pid: req.Pid, sequencenumber: sequence(), refreshed: time.Now()}
	if *start_ready {
		si.health = common.Health_READY
	} else {
		si.health = common.Health_STARTING
	}
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
	s.debugf(req, "Registered #%d (%s)\n", si.sequencenumber, si.ServiceName())
	s.dump("After registration:")
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
	procid := req.ProcessID
	for _, si := range s.instances {
		ras := si.registeredAs
		if ras != nil && ras.ProcessID != req.ProcessID {
			s.debugf(req, "processid %s does not match %s (processid mismatch \"%s\" != \"%s\")\n", procid, si.String(), req.ProcessID, ras.ProcessID)
			continue
		}
		if si.ServiceName() != servicename {
			s.debugf(req, "processid %s does not match %s (servicename mismatch)\n", procid, si.String())
			continue
		}
		if !si.IP.Equals(ip) {
			s.debugf(req, "processid %s does not match %s (ip mismatch)\n", procid, si.String())
			continue
		}
		if si.Port() != port {
			s.debugf(req, "processid %s does not match %s (port mismatch (%d != %d))\n", procid, si.String(), si.Port(), port)
			continue
		}
		s.debugf(req, "found: %s\n", si.String())
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
		s.debugf(req, "found: %s\n", si.String())
		return si
	}
	s.debugf(req, "no instance found\n")
	return nil
}

// returns all service instances which were deregistered.
func (s *ServiceList) Deregister(ctx context.Context, processid string) []*serviceInstance {
	s.debugf(debugnd("dereg"), "Deregistering processpid \"%s\"\n", processid)
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
	nd := debugnd("removeipport")
	s.debugf(nd, "Request to deregister service @%s:%d\n", ip.ExposeAs(), port)
	instanceModifyLock.Lock()
	defer instanceModifyLock.Unlock()
	var newi []*serviceInstance
	for _, si := range s.instances {
		if si.Port() == port && si.IP.Equals(ip) {
			s.debugf(nd, "IP/Port based removal of service %s@%s:%d\n", si.ServiceName(), ip.ExposeAs(), port)
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

func (s *ServiceList) debugf(d debugif, format string, args ...interface{}) {
	if !*debug {
		return
	}
	for _, n := range nond {
		if d.GetServiceName() == n {
			return
		}
	}
	txt := fmt.Sprintf("[v2 servicelist %s] ", d.GetServiceName())
	ftxt := fmt.Sprintf(format, args...)
	fmt.Print(txt + ftxt)
}

func (s *ServiceList) dump(txt string) {
	if !*do_dump_servicelist {
		return
	}
	fmt.Printf("BEGIN SERVICELIST ------- %s ----------- \n", txt)
	fmt.Printf("Instances:\n")
	for _, i := range s.instances {
		fmt.Printf("   %s", i.LongString())
		fmt.Printf("   %#v\n", i.Registration())
	}
	fmt.Printf("Services:\n")
	for k, v := range s.services {
		fmt.Printf("   %s: %s", k, v.LongString())
	}
	fmt.Printf("END SERVICELIST ------- %s ----------- \n", txt)
}
func (si *serviceInstance) LongString() string {
	ca := si.createdAs
	cr := "none"
	if ca != nil {
		di := ca.DeployInfo
		cr = fmt.Sprintf("{processid: %s, Binary: %s, deploymentid: %s}", ca.ProcessID, di.Binary, di.DeploymentID)
	}

	ra := si.registeredAs
	rr := "none"
	if ra != nil {
		rr = fmt.Sprintf("{processid: %s, Port: %d, ServiceName: %s, Pid:%d, user: %s", ra.ProcessID, ra.Port, ra.ServiceName, ra.Pid, ra.UserID)
	}

	n := "none"
	if si.registeredAs != nil {
		n = si.registeredAs.ServiceName
	}
	s := fmt.Sprintf("Name=%s, IP=%v,local=%v,fallback=%v,remote=%v,queried=%v,failures=%d\n", n, si.IP, si.isLocal, si.isFallback, si.isRemote, si.didQueryAutodeployer, si.serviceCheckFailures)
	s = s + fmt.Sprintf("      createdas=%s\n", cr)
	s = s + fmt.Sprintf("      registeredas=%s\n", rr)
	return s
}
func (si *serviceInfo) LongString() string {
	cr := "none"
	ca := si.createdAs
	if ca != nil {
		di := ca.DeployInfo
		cr = fmt.Sprintf("{processid: %s, Binary: %s, deploymentid: %s}", ca.ProcessID, di.Binary, di.DeploymentID)
	}
	s := fmt.Sprintf("Created=%s, ip=%v, lastused=%s\n", utils.TimeString(si.created), si.ip, utils.TimeString(si.lastUsed))
	s = s + fmt.Sprintf("     createdAs=%s\n", cr)
	return s
}
