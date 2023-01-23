package registryimpl

import (
	"flag"
	"fmt"
	prom "golang.conradwood.net/apis/promconfig"
	reg "golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/authremote"
	"time"
)

var (
	promClient        prom.PromConfigServiceClient
	update_prometheus = flag.Bool("prometheus_update", false, "if true, automatically update prometheus")
)

type promUpdateRequest struct {
}

// called on each update
func (s *V2Registry) promUpdate() {
	if !*update_prometheus {
		return
	}
	s.promCtr++
	pur := &promUpdateRequest{}
	if len(s.promChan) > 50 {
		return
	}
	s.promChan <- pur
}

// runs async
func (s *V2Registry) promUpdater() {
	if promClient == nil {
		promClient = prom.GetPromConfigServiceClient()
	}
	for {
		b := false
		select {
		case _ = <-s.promChan:
			b = true
		case <-time.After(5 * time.Second):
			//
		}
		//		fmt.Printf("Running promupdate (%s)\n", s.promName)
		err := s.rebuildPrometheus(b)
		if err != nil {
			fmt.Printf("Rebuild of prometheus failed: %s\n", err)
		}
	}
}

func (s *V2Registry) rebuildPrometheus(must bool) error {
	if !must && s.promCtr <= s.promDoneCtr {
		return nil
	}
	do := s.promCtr
	tl := &prom.TargetList{
		Reporter: &prom.Reporter{Reporter: s.promName},
	}
	sis := s.serviceList.FilterByApiType(reg.Apitype_status).Instances()
	addrs := make(map[string][]string)
	for _, si := range sis {
		if !si.Targetable() {
			continue
		}
		addrs[si.ServiceName()] = append(addrs[si.ServiceName()], fmt.Sprintf("%s:%d", si.IP.ExposeAs(), si.Port()))
	}
	for k, v := range addrs {
		t := &prom.Target{Name: k, Addresses: v}
		tl.Targets = append(tl.Targets, t)
	}
	ctx := authremote.Context()
	_, err := promClient.NewTargets(ctx, tl)
	if err != nil {
		return err
	}
	s.promDoneCtr = do
	time.Sleep(1 * time.Second) // next update in 1 second earliest..
	return nil
}
