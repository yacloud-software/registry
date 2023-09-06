package registryimpl

import (
	"context"
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	reg "golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/utils"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	nond = []string{ // do not print debug info for these services
		"quickdev.QuickDevService",
		"autodeployer.AutoDeployer",
	}
	promNameCtrLock sync.Mutex
	promNameCtr     = 0
	debug           = flag.Bool("debug_registry", false, "debug v2 code")
	print_lookups   = flag.Bool("print_lookups", false, "if true, print all lookup requests")
)

func New(upstream string) *V2Registry {
	res := &V2Registry{serviceList: &ServiceList{}}
	res.serviceList.services = make(map[string]*serviceInfo)
	res.myid = "V2-" + utils.RandomString(96)
	res.upstreamips = strings.Split(upstream, ",")
	for i, s := range res.upstreamips {
		if strings.Contains(s, ":") {
			continue
		}
		res.upstreamips[i] = fmt.Sprintf("%s:5001", res.upstreamips[i])
	}
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Unable to get hostname: %s\n", err)
		hostname = "NOHOST"
	}
	promNameCtrLock.Lock()
	res.promName = hostname
	if promNameCtr != 0 {
		res.promName = fmt.Sprintf("%s_%d", hostname, promNameCtr)
	}
	promNameCtr++
	promNameCtrLock.Unlock()
	// we always include ourselves:
	res.serviceList.AddLocal("registry.Registry", 5001)
	res.promChan = make(chan *promUpdateRequest, 100)
	go res.promUpdater()
	go res.Cleaner()
	go res.VerifyStatusLoop()
	return res

}

type V2Registry struct {
	fallbackhost          string
	fallbackport          uint32
	promDoneCtr           int
	promCtr               int
	promName              string
	promChan              chan *promUpdateRequest
	myid                  string
	serviceList           *ServiceList
	upstreamips           []string
	NewServiceListener    func(ctx context.Context, req *reg.Registration)
	RemoveServiceListener func(ctx context.Context, req *reg.Registration)
}

func (s *V2Registry) ListRegistrations(ctx context.Context, req *reg.V2ListRequest) (*reg.RegistrationList, error) {
	s.Printf("listing all registrations...\n")
	res := &reg.RegistrationList{}
	for _, si := range s.serviceList.Instances() {
		if req.NameMatch != "" {
			if !strings.Contains(si.ServiceName(), req.NameMatch) {
				continue
			}
		}
		r := si.Registration()
		res.Registrations = append(res.Registrations, r)
	}
	return res, nil
}
func (s *V2Registry) V2GetTarget(ctx context.Context, req *reg.V2GetTargetRequest) (*reg.V2GetTargetResponse, error) {
	return s.V2getTarget(ctx, req, true, nil)
}
func (s *V2Registry) V2getTarget(ctx context.Context, req *reg.V2GetTargetRequest, use_upstream bool, utr *reg.UpstreamTargetRequest) (*reg.V2GetTargetResponse, error) {
	if *print_lookups {
		fmt.Printf("gettarget(%s)\n", req.ServiceName)
	}
	sl := s.serviceList.FilterByNames(req.ServiceName)
	sl = sl.FilterByApiType(req.ApiType)
	res := &reg.V2GetTargetResponse{}
	// also consider buildid - only get highest buildid _AND_ all with user routinginfo or no routinginfo
	// on defaultroutes, we also prefer the ones with a buildid of "0", presumed to be a test version
	highestBuild := uint64(0)
	hasnull := false
	for _, si := range sl.TargetableInstances() {
		if !isDefaultRoute(si.Target().RoutingInfo) {
			continue
		}
		if si.BuildID() == 0 {
			hasnull = true
		}

		if si.BuildID() > highestBuild {
			highestBuild = si.BuildID()
		}
	}
	for _, si := range sl.TargetableInstances() {
		t := si.Target()
		if !isDefaultRoute(t.RoutingInfo) {
			res.Targets = append(res.Targets, si.Target())
			continue
		}
		if hasnull {
			if si.BuildID() == 0 {
				res.Targets = append(res.Targets, si.Target())
			}
			continue
		}
		if si.BuildID() == highestBuild {
			res.Targets = append(res.Targets, si.Target())
		}
	}
	if len(res.Targets) != 0 || use_upstream == false {
		return res, nil
	}
	for _, s := range req.ServiceName {
		reportMissedLookup(s)
	}
	//s.Printf("No targets registered for %s - asking upstream\n", req.ServiceName)
	// if we have no targets ourselves, ask upstream
	upstreamResponses, err := s.get_upstream(ctx, req, utr)
	if err != nil {
		s.Printf("Upstream error. (%d responses): %s\n", len(upstreamResponses), err)
	}

	// take note if we want any of them upstream
	for _, u := range upstreamResponses {
		for _, t := range u.Response.Targets {
			reportFoundUpstream(t.ServiceName)
		}
	}
	// merge upstream into our response
	new_targets := 0
	for _, u := range upstreamResponses {
		new_targets = new_targets + len(u.Response.Targets)
		res.Targets = append(res.Targets, u.Response.Targets...)
	}
	//	s.Printf("Checked upstream and got %d new targets\n", new_targets)

	if len(res.Targets) == 0 && s.fallbackhost != "" {
		// fallback target?
		for _, sv := range req.ServiceName {
			t := &reg.Target{ServiceName: sv, IP: s.fallbackhost, Port: s.fallbackport}
			res.Targets = append(res.Targets, t)
		}
	}
	return res, nil
}
func (s *V2Registry) get_upstream(ctx context.Context, req *reg.V2GetTargetRequest, utr *reg.UpstreamTargetRequest) ([]*reg.DownstreamTargetResponse, error) {
	usr := &reg.UpstreamTargetRequest{Request: req, TTL: 10}
	if utr != nil {
		usr = utr
	}
	upstreamResponses, err := s.getFromUpstream(ctx, usr)
	return upstreamResponses, err

}
func (s *V2Registry) V2DeregisterService(ctx context.Context, req *reg.DeregisterServiceRequest) (*common.Void, error) {
	sis := s.serviceList.Deregister(ctx, req.ProcessID)
	s.promUpdate()
	s.callRemoveServiceListener(ctx, sis)
	return &common.Void{}, nil
}
func (s *V2Registry) callRemoveServiceListener(ctx context.Context, sis []*serviceInstance) {
	l := s.RemoveServiceListener
	if l == nil {
		return
	}
	for _, si := range sis {
		l(ctx, si.Registration())
	}
}

func (s *V2Registry) V2HeartBeat(ctx context.Context, req *reg.HeartBeatRequest) (*common.Void, error) {
	return nil, errors.NotImplemented(ctx, "V2HeartBeat")
}
func (s *V2Registry) V2CreateService(ctx context.Context, req *reg.CreateServiceRequest) (*common.Void, error) {
	debugf(debugnd(req.DeployInfo.Binary), "V2CreateService request received , processid \"%s\"\n", req.ProcessID)
	err := s.serviceList.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &common.Void{}, nil
}
func (s *V2Registry) V2RegisterService(ctx context.Context, req *reg.RegisterServiceRequest) (*reg.RegisterServiceResponse, error) {
	debugf(req, "V2RegisterService request received %s, processid \"%s\"\n", req.ServiceName, req.ProcessID)
	si, isnew, err := s.serviceList.Registration(ctx, req)
	if err != nil {
		fmt.Printf("Failed to register: %s\n", utils.ErrorString(err))
		return nil, err
	}
	debugf(req, "RegisterService processid %s, isnew: %v\n", req.ProcessID, isnew)
	if si.createdAs == nil && !si.didQueryAutodeployer {
		// maybe we just restarted. so we get lots of registrations without create-service.
		// we try to contact the autodeployer at the host for more information
		err = s.queryAutodeployer(ctx, si)
		if err != nil {
			s.Printf("could not query autodeployer: %s\n", err)
		}
	}

	if isnew {
		s.promUpdate()
		s.callNewServiceListener(ctx, si)
		go func(sx *serviceInstance) {
			time.Sleep(time.Duration(2) * time.Second)
			TriggerVerifyStatus(sx)
		}(si)
	}
	res := &reg.RegisterServiceResponse{}
	return res, nil
}
func (s *V2Registry) callNewServiceListener(ctx context.Context, si *serviceInstance) {
	l := s.NewServiceListener
	if l == nil {
		return
	}
	l(ctx, si.Registration())
}
func (s *V2Registry) RequestForTarget(ctx context.Context, req *reg.UpstreamTargetRequest) (*reg.DownstreamTargetResponse, error) {
	ip, err := IPFromContext(ctx)
	if err != nil {
		fmt.Printf("Invalid ip address: %s\n", err)
		return nil, err
	}
	s.Printf("Requested target %s from downstream registry @%s (TTL=%d)\n", req.Request.ServiceName, ip.ExposeAs(), req.TTL)
	for _, i := range req.RegistryIDs {
		if i == s.myid {
			return &reg.DownstreamTargetResponse{Response: &reg.V2GetTargetResponse{}}, nil
		}
	}
	req.RegistryIDs = append(req.RegistryIDs, s.myid)
	resp, err := s.V2getTarget(ctx, req.Request, true, req)
	if err != nil {
		return nil, err
	}
	res := &reg.DownstreamTargetResponse{}
	res.Response = resp
	res.RegistryIDs = append(req.RegistryIDs, s.myid)
	return res, nil
}

// if no target is available then one might specify a fallback host
// this host should handle _all_ services that we might come across.
// for example a proxy server is a good example
// if there are upstream targets, these will be served first
func (s *V2Registry) SetFallbackTarget(host string, port uint32) {
	s.fallbackhost = host
	s.fallbackport = port
}

// add remote targets (which are served as normal targets AS LONG AS no local ones are available)
func (s *V2Registry) AddRemote(name string, host string, port uint32) error {
	s.serviceList.AddRemoteWithHost(name, host, port)
	return nil
}

// add remote targets (which are served as normal targets AS LONG AS no local ones are available)
func (s *V2Registry) AddFallbackWithHost(name string, host string, port uint32) error {
	s.serviceList.AddFallbackWithHost(name, host, port)
	return nil
}

func (s *V2Registry) Printf(format string, args ...interface{}) {
	if !*debug {
		return
	}
	txt := "[v2] "
	fmt.Printf(txt+format, args...)
}

func isDefaultRoute(t *reg.RoutingInfo) bool {
	if t == nil {
		return true
	}
	if t.RunningAs == nil && t.GatewayID == "" {
		return true
	}
	return false
}

type debugif interface {
	GetServiceName() string
}
type debugifs struct {
	name string
}

func (d *debugifs) GetServiceName() string {
	return d.name
}
func debugnd(name string) debugif {
	return &debugifs{name: name}
}

func debugf(d debugif, format string, args ...interface{}) {
	if !*debug {
		return
	}
	for _, n := range nond {
		if d.GetServiceName() == n {
			return
		}
	}
	txt := fmt.Sprintf("[v2 registry %s] ", d.GetServiceName())
	ftxt := fmt.Sprintf(format, args...)
	fmt.Print(txt + ftxt)
}
