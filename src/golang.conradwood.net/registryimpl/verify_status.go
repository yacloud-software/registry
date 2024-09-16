package registryimpl

/*
this checks for the healthz status of registered services and removes them if they repeatedly fail
to answer altogether
*/
import (
	"flag"
	"fmt"
	"time"

	"golang.conradwood.net/apis/common"
	reg "golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/http"
	"golang.conradwood.net/go-easyops/prometheus"
)

const (
	WORKERS = 10

// CHECK_INTERVAL = 30
)

var (
	keepAlive     = flag.Duration("keepalive", time.Duration(30)*time.Second, "keep alive interval in seconds to check each registered service")
	healthzChecks = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "registry_total_healthzChecks",
			Help: "V=1 UNIT=ops DESC=number of healthzchecks",
		},
		[]string{"servicename"},
	)
	healthzChecksCur = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "registry_cur_failed_healthzCheck",
			Help: "V=1 UNIT=ops DESC=current number of failed healthchecks",
		},
		[]string{"servicename", "serviceinstance"},
	)
	healthzChecksQ = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "registry_healthzChecks_queued",
			Help: "V=1 UNIT=ops DESC=number of healthzchecks queued for processing",
		},
	)
	failedHealthzChecks = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "registry_failed_healthzChecks",
			Help: "V=1 UNIT=ops DESC=number of failed healthzchecks",
		},
		[]string{"servicename"},
	)
	verify_chan = make(chan *verifyWork, WORKERS*10)
)

func init() {
	prometheus.MustRegister(healthzChecks, failedHealthzChecks, healthzChecksQ, healthzChecksCur)
}
func updateqcounter() {
	healthzChecksQ.Set(float64(len(verify_chan)))
}
func (rv *V2Registry) VerifyStatusLoop() {
	for i := 0; i < WORKERS; i++ {
		go rv.verifyStatusWorker()
	}
	for {
		updateqcounter()
		time.Sleep(*keepAlive)
		updateqcounter()
		sl := rv.serviceList
		if sl == nil {
			continue
		}
		instances := sl.Instances()
		for _, i := range instances {
			updateqcounter()
			// does this expose a 'status'?
			if !i.IncludesApiType(reg.Apitype_status) {
				// if not, ignore ig
				continue
			}
			// check synchronously AS WELL as asynchronously in the worker
			if time.Since(i.lastServiceCheck) < *keepAlive {
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
		updateqcounter()
		w := <-verify_chan
		si := w.instance
		if si == nil {
			continue
		}
		regas := si.registeredAs
		if regas == nil {
			continue
		}
		if !si.IncludesApiType(reg.Apitype_status) {
			// if not, ignore ig
			continue
		}
		if time.Since(si.lastServiceCheck) < *keepAlive {
			continue
		}
		si.lastServiceCheck = time.Now()
		l := prometheus.Labels{"servicename": regas.ServiceName}
		healthzChecks.With(l).Inc()
		//		fmt.Printf("Checking %s:%d\n", si.IP.ExposeAs(), reg.Port)
		h := &http.HTTP{}
		url := fmt.Sprintf("https://%s:%d/internal/health", si.IP.ExposeAs(), regas.Port)
		hr := h.Get(url)
		if hr.Error() != nil && hr.HTTPCode() == 400 || hr.HTTPCode() == 404 {
			xh := &http.HTTP{}
			xurl := fmt.Sprintf("https://%s:%d/internal/healthz", si.IP.ExposeAs(), regas.Port)
			hr = xh.Get(xurl)
		}
		if hr.Error() != nil {
			if hr.HTTPCode() == 400 || hr.HTTPCode() == 404 {
				url = fmt.Sprintf("https://%s:%d/internal/service-info/name", si.IP.ExposeAs(), regas.Port)
				hr = h.Get(url)
			}
			if hr.Error() != nil {
				failedHealthzChecks.With(l).Inc()
				setServiceCheckFailure(si, si.serviceCheckFailures+1)
				fmt.Printf("%d Failure (#%d) %s: %s\n", hr.HTTPCode(), si.serviceCheckFailures, hr.FinalURL(), hr.Error())
			}
			si.lastServiceCheck = time.Now()
			continue
		}
		b := string(hr.Body())
		if b == "OK" { // very old client
			fmt.Printf("WARNING: Service %s reported very old style \"OK\"\n", si.String())
			b = "READY"
		}
		//	fmt.Printf("Body: \"%s\"\n", b)
		health, found := common.Health_value[b]
		if !found {
			fmt.Printf("WARNING: Service %s reported unsupported health \"%s\"\n", si.String(), b)
		} else {
			si.setHealth(common.Health(health))
		}
		si.lastServiceCheck = time.Now()
	}
}
func (si *serviceInstance) setHealth(h common.Health) {
	if h != si.cur_health {
		fmt.Printf("Service %s changed health from %v to %v\n", si.String(), si.cur_health, h)
	}
	si.cur_health = h
}
func (si *serviceInstance) getHealth() common.Health {
	return si.cur_health
}
func (si *serviceInstance) serviceReady() bool {
	if si.getHealth() == common.Health_READY {
		return true
	}
	return false
}
func setServiceCheckFailure(si *serviceInstance, newvalue int) {
	si.serviceCheckFailures = newvalue
	reg := si.registeredAs
	l := prometheus.Labels{
		"servicename":     reg.ServiceName,
		"serviceinstance": fmt.Sprintf("%s:%d", si.IP.ExposeAs(), reg.Port),
	}
	healthzChecksCur.With(l).Set(float64(newvalue))
}
func TriggerVerifyStatus(si *serviceInstance) {
	select {
	case verify_chan <- &verifyWork{instance: si}:
		//
	default:
		//
	}
}
