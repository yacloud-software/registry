package main

/*********************************************************************************






This is obsolete.

replaced by ../registryimpl








*********************************************************************************/

import (
	"container/list"
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/prometheus"
	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/registryimpl"
	"google.golang.org/grpc"
)

// static variables for flag parser
var (
	myid                 string
	allow_localhost      = flag.Bool("allow_localhost", false, "if true, registrations at localhost will be allowed (and exposed)")
	callerid             uint64
	debug                = flag.Bool("debug_old", false, "Enable debugging")
	debug_match          = flag.Bool("debug_match", false, "Enable debugging for the matching algorithm (LOTS OF LOGS)")
	port                 = flag.Int("port", 5000, "The server port (non-tls)")
	tlsport              = flag.Int("tlsport", 5001, "The server port (tls)")
	max_failures         = flag.Int("max_failures", 10, "max failures after which service will be deregistered")
	upstream             = flag.String("upstream", "", "the `ip:port` of an upstream registry. \"port\" defaults to 5000. The upstream registry is consulted by this instance of the registry whenever an instance of a service is requested, but not registered on this instance. CAUTION: this can lead to confused developers because instances registered on the upstream registry won't ever contact this instance and thus won't use services registered here.")
	services             *list.List
	idCtr                = 1
	missedLookupsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "registry_missed_lookups_requests",
			Help: "number of lookup requests received and not answered",
		},
		[]string{"servicename"},
	)
	non_tls_grpc_server_requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "registry_non_tls_requests_received",
			Help: "requests to log stuff received",
		},
		[]string{"servicename", "method"},
	)
	//forceip    = flag.String("forceip", "", "all services will be exposed at this ip address. (testing only)")
	deployPath = "registry/myself"
	v2server   pb.RegistryServer
)

type serviceEntry struct {
	desc      *pb.ServiceDescription
	instances []*serviceInstance
}
type serviceInstance struct {
	serviceID       int
	failures        int
	disabled        bool
	firstRegistered time.Time
	lastSuccess     time.Time
	lastRefresh     time.Time
	address         pb.ServiceAddress
	apitype         []pb.Apitype
}

func (si *serviceInstance) toString() string {
	s := fmt.Sprintf("%d %s:%d", si.serviceID, si.address.Host, si.address.Port)
	return s
}
func main() {
	flag.Parse() // parse stuff. see "var" section above
   server.SetHealth(common.Health_STARTING)
	myid = "V1-" + utils.RandomString(96)

	v2server = registryimpl.New(*upstream)
	prometheus.MustRegister(non_tls_grpc_server_requests)
	go RegisterNonTLS()

	sd := server.NewServerDef()
	sd.SetNoAuth()
	sd.SetPort(*tlsport)
sd.SetOnStartupCallback(startup)
	sd.DontRegister()
	//	sd.DeployPath = deployPath
	sd.SetRegister(server.Register(
		func(server *grpc.Server) error {
			pb.RegisterRegistryServer(server, v2server) // created by proto
			return nil
		},
	))
	fmt.Printf("Starting Registry Server...\n")
	err := server.ServerStartup(sd)
	fmt.Printf("Failed to start server: %s\n", err)
	os.Exit(10)

}
func startup() {
	server.SetHealth(common.Health_READY)
}


func RegisterNonTLS() {
	listenAddr := fmt.Sprintf(":%d", *port)
	fmt.Println("Starting Registry Service on ", listenAddr)
	lis, err := net.Listen("tcp4", fmt.Sprintf(":%d", *port))
	utils.Bail("Failed to listen", err)

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(NONTLSUnaryAuthInterceptor))
	pb.RegisterRegistryServer(grpcServer, v2server) // created by proto

	fmt.Printf("Serving...\n")
	grpcServer.Serve(lis)

}

func NONTLSUnaryAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	method := server.MethodNameFromUnaryInfo(info)
	service := server.ServiceNameFromUnaryInfo(info)
	if *debug {
		fmt.Printf("Non-TLS Method: \"%s\"\n", method)
	}
	non_tls_grpc_server_requests.With(prometheus.Labels{
		"method":      method,
		"servicename": service,
	}).Inc()

	/*
		fmt.Printf("NONTLS-req: %v\n", req)
		fmt.Printf("NONTLS-inf: %v\n", info)
	*/
	return handler(ctx, req)
}

/**********************************
* check registered servers regularly
***********************************/
func (regService *RegistryService) CheckRegistry() {
	return
}

// true if some were removed
func (regService *RegistryService) removeInvalidInstances() bool {
	return false
}

func (si *serviceEntry) toString() string {
	name := fmt.Sprintf("%s:%s", si.desc.Name, si.desc.Path)
	return name
}

func (regService *RegistryService) isValid(si *serviceInstance) bool {
	return false
}
func (regService *RegistryService) CheckService(ctx context.Context, desc *serviceEntry, addr *serviceInstance) error {

	return nil
}

/**********************************
* helpers
***********************************/
func (regService *RegistryService) FindInstanceById(id int) *serviceInstance {
	for e := services.Front(); e != nil; e = e.Next() {
		sl := e.Value.(*serviceEntry)
		for _, si := range sl.instances {
			if si.serviceID == id {
				return si
			}
		}
	}
	return nil
}

func (regService *RegistryService) FindServices(sd *pb.ServiceDescription) []*serviceEntry {
	var res []*serviceEntry
	for e := services.Front(); e != nil; e = e.Next() {
		sl := e.Value.(*serviceEntry)
		if (sl.desc.Path != "") && (sd.Path != "") {
			if sl.desc.Path != sd.Path {
				continue
			}
		}
		if sl.desc.Name == sd.Name {
			res = append(res, sl)
		}
	}
	return res
}

// this is not a good thing - it finds the FIRST entry by name
func (regService *RegistryService) FindService(sd *pb.ServiceDescription) *serviceEntry {

	return nil
}

func (regService *RegistryService) AddService(sd *pb.ServiceDescription, hostname string, port int32, apitype []pb.Apitype, ri *pb.RoutingInfo) *serviceInstance {
	return nil

}

/**********************************
* implementing the functions here:
***********************************/
type RegistryService struct {
	TLS bool
}

func (regService *RegistryService) DEPRECATED_GetServiceAddress(ctx context.Context, gr *pb.GetRequest) (*pb.GetResponse, error) {
	return nil, errors.NotImplemented(ctx, "GetServiceAddress")
}

func (regService *RegistryService) GetServiceAddress(ctx context.Context, gr *pb.GetRequest) (*pb.GetResponse, error) {
	return nil, errors.NotImplemented(ctx, "registryv1")
}
func (regService *RegistryService) DeregisterService(ctx context.Context, pr *pb.DeregisterRequest) (*pb.EmptyResponse, error) {
	return nil, errors.NotImplemented(ctx, "registryv1")
}
func (regService *RegistryService) RegisterService(ctx context.Context, pr *pb.ServiceLocation) (*pb.GetResponse, error) {
	return nil, errors.NotImplemented(ctx, "registryv1")
}

func (regService *RegistryService) ListServices(ctx context.Context, pr *pb.ListRequest) (*pb.ListResponse, error) {
	return nil, errors.NotImplemented(ctx, "registryv1")
}
func (regService *RegistryService) ShutdownService(ctx context.Context, pr *pb.ShutdownRequest) (*pb.EmptyResponse, error) {
	return nil, errors.NotImplemented(ctx, "registryv1")
}
func (regService *RegistryService) GetUpstreamTarget(ctx context.Context, pr *pb.GetUpstreamTargetRequest) (*pb.ListResponse, error) {
	return nil, errors.NotImplemented(ctx, "registryv1")
}

func (regService *RegistryService) GetTarget(ctx context.Context, pr *pb.GetTargetRequest) (*pb.ListResponse, error) {
	return nil, errors.NotImplemented(ctx, "registryv1")
}

func (regService *RegistryService) HideService(ctx context.Context, pr *pb.HideServiceRequest) (*pb.ServiceLocation, error) {
	return nil, errors.NotImplemented(ctx, "registryv1")
}

func (regService *RegistryService) HostHiddenStatus(ctx context.Context, pr *pb.HideUpdateRequest) (*common.Void, error) {
	return nil, errors.NotImplemented(ctx, "registryv1")
}
func (regService *RegistryService) ByIPPort(ctx context.Context, pr *pb.ByIPPortRequest) (*pb.ByIPPortResponse, error) {
	return nil, errors.NotImplemented(ctx, "registryv1")
}
func (regService *RegistryService) InformProcessShutdown(ctx context.Context, pr *pb.ProcessShutdownRequest) (*pb.EmptyResponse, error) {
	return nil, errors.NotImplemented(ctx, "registryv1")
}

func (r *RegistryService) ListRegistrations(ctx context.Context, req *pb.V2ListRequest) (*pb.RegistrationList, error) {
	return nil, errors.NotImplemented(ctx, "ListRegistrations")
}
func (r *RegistryService) V2GetTarget(ctx context.Context, req *pb.V2GetTargetRequest) (*pb.V2GetTargetResponse, error) {
	return nil, errors.NotImplemented(ctx, "V2GetTarget")
}
func (r *RegistryService) V2DeregisterService(ctx context.Context, req *pb.DeregisterServiceRequest) (*common.Void, error) {
	return nil, errors.NotImplemented(ctx, "V2DeregisterService")
}
func (r *RegistryService) V2HeartBeat(ctx context.Context, req *pb.HeartBeatRequest) (*common.Void, error) {
	return nil, errors.NotImplemented(ctx, "V2HeartBeat")
}
func (r *RegistryService) V2CreateService(ctx context.Context, req *pb.CreateServiceRequest) (*common.Void, error) {
	return nil, errors.NotImplemented(ctx, "V2CreateService")
}
func (r *RegistryService) V2RegisterService(ctx context.Context, req *pb.RegisterServiceRequest) (*pb.RegisterServiceResponse, error) {
	return nil, errors.NotImplemented(ctx, "V2RegisterService")
}

func (r RegistryService) Printf(format string, args ...interface{}) {
	if !*debug {
		return
	}
	txt := "[v1] "
	fmt.Printf(txt+format, args...)
}
func (r *RegistryService) RequestForTarget(ctx context.Context, req *pb.UpstreamTargetRequest) (*pb.DownstreamTargetResponse, error) {
	return nil, errors.NotImplemented(ctx, "registryv1")
}


