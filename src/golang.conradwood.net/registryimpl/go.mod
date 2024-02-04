module golang.conradwood.net/registryimpl

go 1.18

require (
	golang.conradwood.net/apis/autodeployer v1.1.2643
	golang.conradwood.net/apis/common v1.1.2778
	golang.conradwood.net/apis/promconfig v1.1.2643
	golang.conradwood.net/apis/registry v1.1.2503
	golang.conradwood.net/go-easyops v0.1.24785
	google.golang.org/grpc v1.60.1
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/prometheus/client_golang v1.17.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	golang.conradwood.net/apis/auth v1.1.2778 // indirect
	golang.conradwood.net/apis/deploymonkey v1.1.2643 // indirect
	golang.conradwood.net/apis/framework v1.1.2643 // indirect
	golang.conradwood.net/apis/goeasyops v1.1.2778 // indirect
	golang.conradwood.net/apis/h2gproxy v1.1.2643 // indirect
	golang.conradwood.net/apis/htmlserver v1.1.2643 // indirect
	golang.conradwood.net/apis/objectstore v1.1.2643 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.yacloud.eu/apis/fscache v1.1.2643 // indirect
	golang.yacloud.eu/apis/session v1.1.2778 // indirect
	golang.yacloud.eu/apis/urlcacher v1.1.2643 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231002182017-d307bd883b97 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace golang.conradwood.net/apis/registry => ../apis/registry
