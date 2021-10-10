package registryimpl

import (
	"context"
	"flag"
	"fmt"
	ad "golang.conradwood.net/apis/autodeployer"
	reg "golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/client"
	"golang.conradwood.net/go-easyops/tokens"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc"
	"sync"
	"time"
)

var (
	connectLock sync.Mutex
	ad_clients  = make(map[string]*adcinfo)
	qaud        = flag.Bool("query_autodeployer", true, "query autodeployer if we have registrations without createservice")
)

type adcinfo struct {
	host   string
	adc    ad.AutoDeployerClient
	conn   *grpc.ClientConn
	failed time.Time
}

func (s *V2Registry) queryAutodeployer(ctx context.Context, si *serviceInstance) error {
	if !*qaud {
		return nil
	}
	host := si.IP.ExposeAs()
	adc, err := getad(host)
	if err != nil {
		fmt.Printf("Error getting autodeploy client: %s\n", utils.ErrorString(err))
		return err
	}
	if time.Since(adc.failed) < time.Duration(5)*time.Minute {
		return nil
	}
	ctx = tokens.ContextWithToken()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(4)*time.Second)
	defer cancel()
	des, err := adc.adc.GetDeployments(ctx, &ad.InfoRequest{})
	if err != nil {
		adc.failed = time.Now()
		//fmt.Printf("get deployments failed...(%s)\n", utils.ErrorString(err))
		return err
	}
	infoLock.Lock()
	defer infoLock.Unlock()
	ip := si.IP
	for _, dapp := range des.Apps {
		di := dapp.Deployment // deployinfo
		processid := dapp.ID
		sin := s.serviceList.services[processid]
		if sin == nil {
			si := &serviceInfo{
				ip:   ip,
				list: s.serviceList,
				createdAs: &reg.CreateServiceRequest{
					ProcessID:  processid,
					Partition:  "",
					Pid:        0,
					DeployInfo: di,
				},
				created:  time.Now(),
				lastUsed: time.Now(),
			}
			s.serviceList.services[processid] = si
			continue
		}
	}

	// now go through the instances and upate them
	instanceModifyLock.Lock()
	defer instanceModifyLock.Unlock()
	for _, sis := range s.serviceList.instances {
		if !sis.Targetable() {
			continue
		}
		if sis.createdAs != nil {
			continue
		}
		ras := sis.registeredAs
		if ras == nil {
			continue
		}
		procid := ras.ProcessID
		sin := s.serviceList.services[procid]
		if sin == nil {
			continue
		}
		sis.createdAs = sin.createdAs
	}
	si.didQueryAutodeployer = true
	return nil
}

func getad(host string) (*adcinfo, error) {
	res := ad_clients[host]
	if res != nil {
		return res, nil
	}
	connectLock.Lock()
	defer connectLock.Unlock()
	fmt.Printf("Connecting to autodeployer on host %s\n", host)
	conn, err := client.ConnectWithIPNoBlock(fmt.Sprintf("%s:4000", host))
	if err != nil {
		fmt.Printf("Failed to connect to autodeployer on host %s: %s\n", host, utils.ErrorString(err))
		return nil, err
	}
	fmt.Printf("Connected to autodeployer on host %s\n", host)
	ai := &adcinfo{adc: ad.NewAutoDeployerClient(conn)}
	ad_clients[host] = ai
	return ai, nil
}
