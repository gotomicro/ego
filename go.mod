module github.com/gotomicro/ego

go 1.16

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/RaMin0/gin-health-check v0.0.0-20180807004848-a677317b3f01
	github.com/StackExchange/wmi v0.0.0-20210224194228-fe8f1750fd46 // indirect
	github.com/alibaba/sentinel-golang v1.0.3
	github.com/codegangsta/inject v0.0.0-20150114235600-33e0aa1cb7c0
	github.com/dave/dst v0.26.2
	github.com/davecgh/go-spew v1.1.1
	github.com/felixge/fgprof v0.9.1
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gin-gonic/gin v1.7.7
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/go-playground/validator/v10 v10.9.0 // indirect
	github.com/go-resty/resty/v2 v2.6.0
	github.com/golang/protobuf v1.5.2
	github.com/google/cel-go v0.11.2
	github.com/gorilla/websocket v1.4.2
	github.com/gotomicro/logrotate v0.0.0-20211108024517-45d1f9a03ff5
	github.com/iancoleman/strcase v0.2.0
	github.com/json-iterator/go v1.1.12
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mitchellh/mapstructure v1.3.2
	github.com/modern-go/reflect2 v1.0.2
	github.com/prometheus/client_golang v1.12.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/shirou/gopsutil v3.21.3+incompatible
	github.com/spf13/cast v1.3.1
	github.com/stretchr/testify v1.7.1
	github.com/ugorji/go v1.2.6 // indirect
	github.com/wk8/go-ordered-map v0.2.0
	go.opencensus.io v0.22.4
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/exporters/jaeger v1.4.1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.7.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.7.0
	go.opentelemetry.io/otel/sdk v1.7.0
	go.opentelemetry.io/otel/trace v1.7.0
	go.uber.org/automaxprocs v1.3.0
	go.uber.org/zap v1.17.0
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/tools v0.1.6
	google.golang.org/genproto v0.0.0-20220310185008-1973136f34c6
	google.golang.org/grpc v1.46.0
	google.golang.org/protobuf v1.28.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace (
    //todo undefined: grpc.WithBalancerName
	google.golang.org/grpc v1.46.0 => google.golang.org/grpc v1.44.0
	google.golang.org/protobuf v1.28.0 => google.golang.org/protobuf v1.27.1
)
