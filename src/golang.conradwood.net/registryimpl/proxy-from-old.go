package registryimpl

import (
	"context"
	"fmt"
	"golang.conradwood.net/apis/common"
	reg "golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/errors"
	//	"golang.conradwood.net/go-easyops/utils"
)

func (s *V2Registry) DeregisterService(ctx context.Context, req *reg.DeregisterRequest) (*reg.EmptyResponse, error) {
	s.Printf("v1-deregister serviceid %s\n", req.ServiceID)
	r := &reg.DeregisterServiceRequest{ProcessID: req.ServiceID}
	_, err := s.V2DeregisterService(ctx, r)
	if err != nil {
		return nil, err
	}
	return &reg.EmptyResponse{}, nil
}
func (s *V2Registry) RegisterService(ctx context.Context, req *reg.ServiceLocation) (*reg.GetResponse, error) {
	myip, err := IPFromContext(ctx) // use this for pid
	if err != nil {
		return nil, err
	}
	var pid string
	// old style has no processid
	for _, svc := range req.Address {
		s.Printf("v1-Registration for %s: %#v\n", req.Service.Name, svc)
		r := &reg.RegisterServiceRequest{
			Port:        uint32(svc.Port),
			ApiType:     svc.ApiType,
			ServiceName: req.Service.Name,
			RoutingInfo: svc.RoutingInfo,
		}
		pid = fmt.Sprintf("V1-registration_of_%s@%s:%d", r.ServiceName, myip.ExposeAs(), r.Port)
		r.ProcessID = pid
		_, err := s.V2RegisterService(ctx, r)
		if err != nil {
			return nil, err
		}
	}
	return &reg.GetResponse{ServiceID: pid}, nil
}
func (s *V2Registry) ListServices(ctx context.Context, req *reg.ListRequest) (*reg.ListResponse, error) {
	v, err := s.ListRegistrations(ctx, &reg.V2ListRequest{})
	if err != nil {
		return nil, err
	}
	l := &reg.ListResponse{}
	for _, rs := range v.Registrations {
		gr := &reg.GetResponse{
			Service:   &reg.ServiceDescription{Name: req.Name},
			ServiceID: rs.ProcessID,
			YourIP:    "[noip-in-v1]",
		}
		gr.Location = &reg.ServiceLocation{
			Service: gr.Service,
		}
		t := rs.Target
		gr.Service.Name = t.ServiceName
		sl := &reg.ServiceAddress{
			Host:        t.IP,
			Port:        int32(t.Port),
			ApiType:     t.ApiType,
			Filtered:    false, //not used
			RoutingInfo: t.RoutingInfo,
		}
		gr.Location.Address = append(gr.Location.Address, sl)

		l.Service = append(l.Service, gr)
	}
	return l, nil
}
func (s *V2Registry) ShutdownService(ctx context.Context, req *reg.ShutdownRequest) (*reg.EmptyResponse, error) {
	return nil, errors.NotImplemented(ctx, "ShutdownService")
}

func (s *V2Registry) GetUpstreamTarget(ctx context.Context, req *reg.GetUpstreamTargetRequest) (*reg.ListResponse, error) {
	buh, err := s.getTarget(ctx, req.TargetRequest, false)
	if err != nil {
		return nil, err
	}
	return buh, nil
}

func (s *V2Registry) GetTarget(ctx context.Context, req *reg.GetTargetRequest) (*reg.ListResponse, error) {
	return s.getTarget(ctx, req, true)
}
func (s *V2Registry) getTarget(ctx context.Context, req *reg.GetTargetRequest, use_upstream bool) (*reg.ListResponse, error) {
	//s.Printf("v1 gettarget for \"%s\"\n", req.Name)
	v, err := s.V2getTarget(ctx, &reg.V2GetTargetRequest{ApiType: req.ApiType, ServiceName: []string{req.Name}}, use_upstream, nil)
	if err != nil {
		return nil, err
	}
	gr := &reg.GetResponse{
		Service:   &reg.ServiceDescription{Name: req.Name},
		ServiceID: "[v1]",
		YourIP:    "[v1]",
	}
	gr.Location = &reg.ServiceLocation{
		Service: gr.Service,
	}
	for _, t := range v.Targets {
		gr.Service.Name = t.ServiceName
		sl := &reg.ServiceAddress{
			Host:        t.IP,
			Port:        int32(t.Port),
			ApiType:     t.ApiType,
			Filtered:    false, //not used
			RoutingInfo: t.RoutingInfo,
		}
		gr.Location.Address = append(gr.Location.Address, sl)

	}
	l := &reg.ListResponse{}
	if len(gr.Location.Address) > 0 {
		l.Service = append(l.Service, gr)
	}
	return l, nil
}

func (s *V2Registry) HideService(ctx context.Context, req *reg.HideServiceRequest) (*reg.ServiceLocation, error) {
	return nil, errors.NotImplemented(ctx, "HideService")
}

func (s *V2Registry) HostHiddenStatus(ctx context.Context, req *reg.HideUpdateRequest) (*common.Void, error) {
	return nil, errors.NotImplemented(ctx, "HostHiddenStatus")
}

func (s *V2Registry) ByIPPort(ctx context.Context, req *reg.ByIPPortRequest) (*reg.ByIPPortResponse, error) {
	return nil, errors.NotImplemented(ctx, "ByIPPort")
}

func (s *V2Registry) InformProcessShutdown(ctx context.Context, req *reg.ProcessShutdownRequest) (*reg.EmptyResponse, error) {
	ip, err := IPFromContext(ctx)
	if err != nil {
		return nil, err
	}
	for _, p := range req.Port {
		s.serviceList.removeIPPort(ip, uint32(p))
	}
	return &reg.EmptyResponse{}, nil
}
