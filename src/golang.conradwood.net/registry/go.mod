module golang.conradwood.net/registry

go 1.21.1

toolchain go1.22.2

replace golang.conradwood.net/registryimpl => ../registryimpl

require (
	golang.conradwood.net/apis/common v1.1.3077
	golang.conradwood.net/apis/registry v1.1.3071
	golang.conradwood.net/go-easyops v0.1.29674
	golang.conradwood.net/registryimpl v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.66.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/grafana/pyroscope-go v1.1.2 // indirect
	github.com/grafana/pyroscope-go/godeltaprof v0.1.8 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_golang v1.20.3 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.59.1 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	golang.conradwood.net/apis/auth v1.1.3077 // indirect
	golang.conradwood.net/apis/autodeployer v1.1.3071 // indirect
	golang.conradwood.net/apis/deploymonkey v1.1.3071 // indirect
	golang.conradwood.net/apis/echoservice v1.1.3071 // indirect
	golang.conradwood.net/apis/errorlogger v1.1.3071 // indirect
	golang.conradwood.net/apis/framework v1.1.3071 // indirect
	golang.conradwood.net/apis/goeasyops v1.1.3077 // indirect
	golang.conradwood.net/apis/grafanadata v1.1.3071 // indirect
	golang.conradwood.net/apis/h2gproxy v1.1.3071 // indirect
	golang.conradwood.net/apis/htmlserver v1.1.3071 // indirect
	golang.conradwood.net/apis/objectstore v1.1.3071 // indirect
	golang.conradwood.net/apis/promconfig v1.1.3071 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	golang.yacloud.eu/apis/autodeployer2 v1.1.3071 // indirect
	golang.yacloud.eu/apis/fscache v1.1.3071 // indirect
	golang.yacloud.eu/apis/session v1.1.3077 // indirect
	golang.yacloud.eu/apis/unixipc v1.1.3071 // indirect
	golang.yacloud.eu/apis/urlcacher v1.1.3071 // indirect
	golang.yacloud.eu/unixipc v0.1.26852 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240604185151-ef581f913117 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace golang.conradwood.net/apis/registry => ../apis/registry
