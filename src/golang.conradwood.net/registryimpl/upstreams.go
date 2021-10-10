package registryimpl

import (
	"context"
	"fmt"
	//	"golang.conradwood.net/apis/common"
	reg "golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/client"
	"sync"
)

var (
	NO_UPSTREAM     = []string{"registry.Registry"}
	upstreamClients = make(map[string]reg.RegistryClient)
	regLock         sync.Mutex
	connectingip    string
)

func (s *V2Registry) getUpstreamClients() []reg.RegistryClient {
	var res []reg.RegistryClient
	for _, ip := range s.upstreamips {
		r, found := upstreamClients[ip]
		if found {
			res = append(res, r)
			continue
		}
		if connectingip == ip {
			fmt.Printf("Upstream registry at \"%s\" blocked whilst connecting...\n", ip)
			continue
		}
		regLock.Lock()
		connectingip = ip
		r, found = upstreamClients[ip]
		if found {
			res = append(res, r)
			connectingip = ""
			regLock.Unlock()
			continue
		}
		fmt.Printf("Connecting registry client to %s\n", ip)
		conn, err := client.ConnectWithIP(ip)
		if err != nil {
			fmt.Printf("Failed to connect to registry %s: %s\n", ip, err)
			connectingip = ""
			regLock.Unlock()
			continue
		}
		fmt.Printf("Connected registry client to %s\n", ip)
		r = reg.NewRegistryClient(conn)
		upstreamClients[ip] = r
		res = append(res, r)
		connectingip = ""
		regLock.Unlock()

	}
	return res
}

// error if _any_ upstream registry is upset. partial results are returned
func (s *V2Registry) getFromUpstream(ctx context.Context, req *reg.UpstreamTargetRequest) ([]*reg.DownstreamTargetResponse, error) {
	var fsvc []string
	for _, svc := range req.Request.ServiceName {
		add := true
		for _, no := range NO_UPSTREAM {
			if svc == no {
				add = false
				break
			}
		}
		if add {
			fsvc = append(fsvc, svc)
		}
	}
	nreq := *req.Request
	wg := &sync.WaitGroup{}
	var terr error
	var res []*reg.DownstreamTargetResponse
	cls := s.getUpstreamClients()
	//	s.Printf("Requesting data from %d upstream registries...\n", len(cls))
	for _, upclient := range cls {
		wg.Add(1)
		go func(upregistry reg.RegistryClient) {
			defer wg.Done()
			/*
				fmt.Printf("Querying upstream...\n")
				defer fmt.Printf("Querying upstream done...\n")
			*/
			utr := &reg.UpstreamTargetRequest{
				TTL:         req.TTL,
				Request:     &nreq,
				RegistryIDs: append(req.RegistryIDs, s.myid),
			}
			r, err := upregistry.RequestForTarget(ctx, utr)
			if err != nil {
				fmt.Printf("upstream registry error: %s\n", err)
				terr = err
				return
			}
			res = append(res, r)

		}(upclient)
	}
	wg.Wait()
	//	s.Printf("Requested data from %d upstream registries.\n", len(cls))
	return res, terr
}
