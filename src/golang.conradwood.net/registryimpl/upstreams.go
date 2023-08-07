package registryimpl

import (
	"context"
	"fmt"
	//	"golang.conradwood.net/apis/common"
	"flag"
	reg "golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/cache"
	"golang.conradwood.net/go-easyops/client"
	"strings"
	"sync"
	"time"
)

var (
	upstream_cache  = cache.New("upstream_cache", time.Duration(15)*time.Second, 1000)
	NO_UPSTREAM     = []string{"registry.Registry"}
	no_upstream     = flag.String("no_upstream", "", "comma delimited list of services never to resolve via upstream regi1stry")
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

type upstream_cache_entry struct {
	err error
	res []*reg.DownstreamTargetResponse
}

// error if _any_ upstream registry is upset. partial results are returned
func (s *V2Registry) getFromUpstream(ctx context.Context, req *reg.UpstreamTargetRequest) ([]*reg.DownstreamTargetResponse, error) {

	key := strings.Join(req.Request.ServiceName, "_") + "_" + req.Request.Partition
	key = key + fmt.Sprintf("_%v", req.Request.ApiType)
	uceo := upstream_cache.Get(key)
	if uceo != nil {
		uce := uceo.(*upstream_cache_entry)
		return uce.res, uce.err
	}
	req.Request.ServiceName = filter_for_upstream(req.Request.ServiceName)
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
	upstream_cache.Put(key, &upstream_cache_entry{err: terr, res: res})
	return res, terr
}

func filter_for_upstream(servicenames []string) []string {
	var nup []string
	for _, np := range strings.Split(*no_upstream, ",") {
		np = strings.Trim(np, " ")
		nup = append(nup, np)
	}
	filter := append(NO_UPSTREAM, nup...)
	var fsvc []string
	for _, svc := range servicenames {
		add := true
		for _, no := range filter {
			if svc == no {
				add = false
				break
			}
		}
		if add {
			fsvc = append(fsvc, svc)
		}
	}
	return fsvc
}
