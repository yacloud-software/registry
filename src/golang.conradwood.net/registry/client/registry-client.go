package main

import (
	"context"
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/client"
	"golang.conradwood.net/go-easyops/tokens"
	"golang.conradwood.net/go-easyops/utils"
	"os"
	"time"
)

// static variables for flag parser
var (
	iponly  = flag.Bool("iponly", false, "only list ip addresses")
	target  = flag.String("target", "", "attempt to connect to a service (with go-easyops.client.Connect())")
	filter  string
	rclient pb.RegistryClient
	long    = flag.Bool("long", false, "long output")
	missed  = flag.Bool("missed", false, "list missed lookups")
)

type Apitypes []pb.Apitype

func main() {
	flag.Parse()
	fs := flag.Args()
	if len(fs) > 0 {
		filter = fs[0]
	}
	rclient = client.GetRegistryClient()
	if *missed {
		showMissed()
		os.Exit(0)
	}
	if *target != "" {
		targetConnect()
		os.Exit(0)
	}
	if !*iponly {
		fmt.Printf("ListRegistrations()...\n\n")
	}
	lr := &pb.V2ListRequest{}
	if len(flag.Args()) != 0 {
		lr.NameMatch = flag.Args()[0]
	}
	resp, err := rclient.ListRegistrations(context.Background(), lr)
	utils.Bail("failed to list services: %v", err)
	printList(resp)
}
func showMissed() {
	ctx := tokens.ContextWithToken()
	ml, err := rclient.GetMissedLookups(ctx, &common.Void{})
	utils.Bail("Failed to get lookups", err)
	for _, l := range ml.Lookups {
		fmt.Printf("%30s %s (%d/%d)\n", l.ServiceName, utils.TimestampString(l.Last), l.Missed, l.Found)
	}
}

func (a Apitypes) String() string {
	deli := ""
	res := ""
	for _, apitype := range a {
		res = fmt.Sprintf("%s%s%s", res, deli, apitype)
		deli = ", "
	}
	return res
}

func lookup(name string) {
	gt := &pb.V2GetTargetRequest{
		ServiceName: []string{name},
		ApiType:     pb.Apitype_grpc,
	}
	fmt.Printf("GetTarget()...\n")
	lr, err := rclient.V2GetTarget(context.Background(), gt)
	utils.Bail("failed to lookup target", err)
	printTargets(lr.Targets)
}

func ApiToString(pa []pb.Apitype) string {
	deli := ""
	res := ""
	for _, apitype := range pa {
		res = fmt.Sprintf("%s%s%s", res, deli, apitype)
		deli = ", "
	}
	return res
}

func targetConnect() {
	fmt.Printf("Connecting to %s\n", *target)
	lookup(*target)
	client.Connect(*target)
	fmt.Printf("Connect successful\n")
}

func printTargets(targets []*pb.Target) {
	for _, t := range targets {
		for _, a := range t.ApiType {
			fmt.Printf("[%6v] %s %s:%d\n", a, t.ServiceName, t.IP, t.Port)
		}
	}
}
func printList(list *pb.RegistrationList) {
	longestName := 0
	for _, r := range list.Registrations {
		if len(r.Target.ServiceName) > longestName {
			longestName = len(r.Target.ServiceName)
		}
	}
	lg := fmt.Sprintf("%%%ds", longestName)
	for _, r := range list.Registrations {
		t := r.Target
		ipport := fmt.Sprintf("%s:%d", t.IP, t.Port)
		apitypes := Apitypes(t.ApiType).String()
		ris := ""
		if t.RoutingInfo != nil {
			u := t.RoutingInfo.RunningAs
			if u != nil {
				gws := ""
				if t.RoutingInfo.GatewayID != "" {
					gws = " " + t.RoutingInfo.GatewayID + " "
				}
				ris = fmt.Sprintf(" [user: %s (#%s)%s]", auth.Description(u), u.ID, gws)
			}
			if len(t.RoutingInfo.Tags) != 0 {
				ris = ris + fmt.Sprintf(" Tags: %v", t.RoutingInfo.Tags)
			}
		}
		flags := flagsFromRegistration(r)
		details := getDetailString(r)
		di := getDeployString(r)
		if *iponly {
			fmt.Printf("%s\n", ipport)
		} else {
			fmt.Printf("%s "+lg+" %20s %15s %6d %s%s%s\n", flags, t.ServiceName, ipport, apitypes, r.Pid, details, ris, di)
		}
	}
	fmt.Println()
	if !*iponly {
		fmt.Println(flagsHelp())
	}
}
func getDeployString(r *pb.Registration) string {
	if r.DeployInfo == nil || r.DeployInfo.AppReference == nil {
		return ""
	}
	d := r.DeployInfo
	ar := d.AppReference
	ad := ar.AppDef
	if ad == nil {
		return ""
	}
	return fmt.Sprintf(" (Did:%s V:%d)", d.DeploymentID, d.BuildID)
}
func getDetailString(r *pb.Registration) string {
	if !*long {
		return ""
	}
	ts := time.Since(time.Unix(int64(r.LastRefreshed), 0)).Seconds()
	s := fmt.Sprintf(" %3.0fs ago ", ts)
	if ts > 999 {
		s = ">999s ago"
	}
	prid := r.ProcessID
	if len(prid) > 20 {
		prid = prid[:20]
	}
	for len(prid) < 20 {
		prid = prid + " "
	}
	s = s + " " + prid
	return s
}
func flagsHelp() string {
	return "Flags: T=targetable, R=running, D=deployinfo available, (lowercase==false)"
}
func flagsFromRegistration(t *pb.Registration) string {
	res := "t"
	if t.Targetable {
		res = "T"
	}
	if t.Running {
		res = res + "R"
	} else {
		res = res + "r"
	}
	if t.DeployInfo == nil {
		res = res + "d"
	} else {
		res = res + "D"
	}
	return res
}
