package registryimpl

/*
this checks for the healthz status of registered services and removes them if they repeatedly fail
to answer altogether
*/
import (
	"fmt"
	reg "golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/http"
	"golang.conradwood.net/go-easyops/prometheus"
	"time"
)

const (
	WORKERS        = 10
	CHECK_INTERVAL = 30
)

var (
	healthzChecks = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "registry_total_healthzChecks",
			Help: "V=1 UNIT=ops DESC=number of healthzchecks",
		},
		[]string{"servicename"},
	)
	failedHealthzChecks = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "registry_failed_healthzChecks",
			Help: "V=1 UNIT=ops DESC=number of failed healthzchecks",
		},
		[]string{"servicename"},
	)
	verify_chan = make(chan *verifyWork, WORKERS)
)

func init() {
	prometheus.MustRegister(healthzChecks, failedHealthzChecks)
}
func (rv *V2Registry) VerifyStatusLoop() {
	for i := 0; i < WORKERS; i++ {
		go rv.verifyStatusWorker()
	}
	for {
		time.Sleep(CHECK_INTERVAL * time.Second)
		sl := rv.serviceList
		if sl == nil {
			continue
		}
		instances := sl.Instances()
		for _, i := range instances {
			// does this expose a 'status'?
			if !i.IncludesApiType(reg.Apitype_status) {
				// if not, ignore ig
				continue
			}
			// check synchronously AS WELL as asynchronously in the worker
			if time.Since(i.lastServiceCheck) < time.Duration(CHECK_INTERVAL)*time.Second {
				continue
			}
			verify_chan <- &verifyWork{instance: i}
		}

	}
}

type verifyWork struct {
	instance *serviceInstance
}

func (rv *V2Registry) verifyStatusWorker() {
	for {
		w := <-verify_chan
		si := w.instance
		if si == nil {
			continue
		}
		reg := si.registeredAs
		if reg == nil {
			continue
		}
		if time.Since(si.lastServiceCheck) < time.Duration(CHECK_INTERVAL)*time.Second {
			continue
		}
		si.lastServiceCheck = time.Now()
		l := prometheus.Labels{"servicename": reg.ServiceName}
		healthzChecks.With(l).Inc()
		//		fmt.Printf("Checking %s:%d\n", si.IP.ExposeAs(), reg.Port)
		h := &http.HTTP{}
		url := fmt.Sprintf("https://%s:%d/internal/healthz", si.IP.ExposeAs(), reg.Port)
		hr := h.Get(url)
		if hr.Error() != nil {
			if hr.HTTPCode() == 400 || hr.HTTPCode() == 404 {
				url = fmt.Sprintf("https://%s:%d/internal/service-info/name", si.IP.ExposeAs(), reg.Port)
				hr = h.Get(url)
			}
			if hr.Error() != nil {
				failedHealthzChecks.With(l).Inc()
				si.serviceCheckFailures++
				fmt.Printf("%d Failure (#%d) %s: %s\n", hr.HTTPCode(), si.serviceCheckFailures, hr.FinalURL(), hr.Error())
			}
			si.lastServiceCheck = time.Now()
			continue
		}
		b := string(hr.Body())
		//	fmt.Printf("Body: \"%s\"\n", b)
		if b == "OK" {
			si.serviceCheckFailures = 0
		} else {
			si.serviceCheckFailures++
		}
		si.lastServiceCheck = time.Now()
	}
}
