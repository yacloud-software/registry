package registryimpl

import (
	"context"
	"golang.conradwood.net/apis/common"
	reg "golang.conradwood.net/apis/registry"
	"time"
)

var (
	missedLookups []*missedLookup
)

type missedLookup struct {
	serviceName string
	last        time.Time
	missed      uint32
	found       uint32
}

func reportMissedLookup(serviceName string) {
	for _, m := range missedLookups {
		if m.serviceName == serviceName {
			m.last = time.Now()
			m.missed++
			return
		}
	}
	m := &missedLookup{serviceName: serviceName, last: time.Now()}
	missedLookups = append(missedLookups, m)
}
func reportFoundUpstream(serviceName string) {
	for _, m := range missedLookups {
		if m.serviceName == serviceName {
			m.last = time.Now()
			m.found++
			return
		}
	}
}

func (s *V2Registry) GetMissedLookups(ctx context.Context, req *common.Void) (*reg.MissedLookupList, error) {
	res := &reg.MissedLookupList{}
	for _, m := range missedLookups {
		if time.Since(m.last) > time.Duration(20)*time.Minute {
			continue
		}
		res.Lookups = append(res.Lookups, &reg.MissedLookup{
			ServiceName: m.serviceName,
			Last:        uint32(m.last.Unix()),
			Missed:      m.missed,
			Found:       m.found,
		})
	}
	return res, nil
}
