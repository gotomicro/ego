package otel

import (
	"context"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cast"
	"go.opentelemetry.io/otel/attribute"
	//lint:ignore SA1019
	jaegerv2 "go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/internal/ienv"
)

type attr struct {
	Key   string      // attribute 的key
	Type  string      // attribute 的类型, 可选： bool |  int64 | float64 | string | boolSlice | int64Slice | float64Slice | stringSlice
	Value interface{} // attribute 的值，需要与类型向匹配，支持通过$操作符读取环境变量
}

// Config ...
type Config struct {
	ServiceName       string                          // 服务名
	OtelType          string                          // type: otlp ,jaeger
	Fraction          float64                         // 采样率： 默认0不会采集
	PanicOnError      bool                            // 异常处理模式
	DefaultAttributes []attr                          // 默认attributes
	options           []tracesdk.TracerProviderOption // 额外的traceProviderOptions
	Jaeger            jaegerConfig                    // otel jaeger 配置
	Otlp              otlpConfig                      // otel otlp 配置
}

// otlpConfig otlp上报协议配置
type otlpConfig struct {
	Endpoint       string                 // oltp endpoint
	Headers        map[string]string      // 默认提供一个 请求头的参数配置
	EnableInsecure bool                   // 是否启用不安全连接 default: true
	options        []otlptracegrpc.Option // 预留自定义配置：   例如 grpc WithGRPCConn
	resOptions     []resource.Option      // res 预留自定以配置
}

// jaegerConfig jaeger上报协议配置
type jaegerConfig struct {
	EndpointType      string // type: agent,collector
	AgentHost         string // agent host
	AgentPort         string // agent port
	CollectorEndpoint string // collector endpoint
	CollectorUser     string // collector user
	CollectorPassword string // collector password
}

// Load parses container configuration from configuration provider, such as a toml file,
// then use the configuration to construct a component container.
func Load(key string) *Config {
	var config = DefaultConfig()
	if err := econf.UnmarshalKey(key, config); err != nil {
		panic(err)
	}
	return config
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		ServiceName: ienv.EnvOrStr("OTEL_SERVICE_NAME", eapp.Name()),
		Jaeger: jaegerConfig{
			AgentHost:         ienv.EnvOrStr("OTEL_EXPORTER_JAEGER_AGENT_HOST", "localhost"),
			AgentPort:         ienv.EnvOrStr("OTEL_EXPORTER_JAEGER_AGENT_PORT", "6831"),
			CollectorEndpoint: ienv.EnvOrStr("OTEL_EXPORTER_JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
			CollectorUser:     ienv.EnvOrStr("OTEL_EXPORTER_JAEGER_USER", ""),
			CollectorPassword: ienv.EnvOrStr("OTEL_EXPORTER_JAEGER_PASSWORD", ""),
			EndpointType:      "collector",
		},
		Otlp: otlpConfig{
			Endpoint:       ienv.EnvOrStr("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
			EnableInsecure: true,
		},
		OtelType:     "otlp",
		PanicOnError: true,
	}
}

// WithTracerProviderOption ...
func (config *Config) WithTracerProviderOption(options ...tracesdk.TracerProviderOption) *Config {
	config.options = append(config.options, options...)
	return config
}

// WithOtlpTraceGrpcOption 自定义otlp Option
func (config *Config) WithOtlpTraceGrpcOption(options ...otlptracegrpc.Option) *Config {
	config.Otlp.options = append(config.Otlp.options, options...)
	return config
}

// WithOtlpResourceOption 自定义otlp resource Option
func (config *Config) WithOtlpResourceOption(options ...resource.Option) *Config {
	config.Otlp.resOptions = append(config.Otlp.resOptions, options...)
	return config
}

// Build ...
func (config *Config) Build(options ...Option) trace.TracerProvider {
	for _, option := range options {
		option(config)
	}
	switch config.OtelType {
	case "otlp":
		return config.buildOtlpTP()
	case "jaeger":
		return config.buildJaegerTP()
	default:
		elog.Panic("otel type error", elog.FieldName(config.OtelType))
	}
	return nil
}

func (config *Config) buildJaegerTP() trace.TracerProvider {
	var endpoint jaegerv2.EndpointOption
	switch config.Jaeger.EndpointType {
	case "agent":
		// Create the Jaeger exporter
		endpoint = jaegerv2.WithAgentEndpoint(
			jaegerv2.WithAgentHost(config.Jaeger.AgentHost),
			jaegerv2.WithAgentPort(config.Jaeger.AgentPort),
		)
	case "collector":
		endpoint = jaegerv2.WithCollectorEndpoint(
			jaegerv2.WithEndpoint(config.Jaeger.CollectorEndpoint),
			jaegerv2.WithUsername(config.Jaeger.CollectorUser),
			jaegerv2.WithPassword(config.Jaeger.CollectorPassword),
		)
	default:
		elog.Panic("jaeger type error", elog.FieldName(config.Jaeger.EndpointType))
	}

	exp, err := jaegerv2.New(endpoint)
	if err != nil {
		return nil
	}

	extraDefaultAttributes := generateDefaultAttributes(config.ServiceName, config.DefaultAttributes...)
	options := []tracesdk.TracerProviderOption{
		// Set the sampling rate based on the parent span to 100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(config.Fraction))),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewSchemaless(extraDefaultAttributes...)),
	}
	options = append(options, config.options...)
	tp := tracesdk.NewTracerProvider(options...)
	return tp
}

// envRgx 是判断某个字符串是否是环境变量的正则表达式
var envRgx = regexp.MustCompile(`^\$([a-zA-Z_][a-zA-Z0-9_]*)$`)

// resolveVal 检查字符串是否为 $ENV 格式，若是则返回环境变量值，否则返回原值
func resolveVal(val string) string {
	// 快速判断，减少正则开销（可选）
	if !strings.HasPrefix(val, "$") {
		return val
	}

	match := envRgx.FindStringSubmatch(val)
	if len(match) <= 1 {
		return val
	}

	envVal := os.Getenv(match[1])
	// 注意：如果环境变量为空，则返回原val
	if envVal == "" {
		return val
	}

	return envVal
}

// generateDefaultAttributes 生成默认的attributes
func generateDefaultAttributes(serviceName string, attrs ...attr) []attribute.KeyValue {
	res := make([]attribute.KeyValue, 0, len(attrs)+1)
	res = append(res, semconv.ServiceNameKey.String(serviceName))
	// 目前仅支持对bool\int64\float64\string这四种类型的val尝试从环境变量取值
	for _, a := range attrs {
		switch a.Type {
		case "bool":
			val := resolveVal(cast.ToString(a.Value))
			res = append(res, attribute.KeyValue{
				Key:   attribute.Key(a.Key),
				Value: attribute.BoolValue(cast.ToBool(val)),
			})
		case "int64":
			val := resolveVal(cast.ToString(a.Value))
			res = append(res, attribute.KeyValue{
				Key:   attribute.Key(a.Key),
				Value: attribute.Int64Value(cast.ToInt64(val)),
			})
		case "float64":
			val := resolveVal(cast.ToString(a.Value))
			res = append(res, attribute.KeyValue{
				Key:   attribute.Key(a.Key),
				Value: attribute.Float64Value(cast.ToFloat64(val)),
			})
		case "string":
			val := resolveVal(cast.ToString(a.Value))
			res = append(res, attribute.KeyValue{
				Key:   attribute.Key(a.Key),
				Value: attribute.StringValue(val),
			})
		case "boolSlice":
			res = append(res, attribute.KeyValue{
				Key:   attribute.Key(a.Key),
				Value: attribute.BoolSliceValue(a.Value.([]bool)),
			})
		case "int64Slice":
			res = append(res, attribute.KeyValue{
				Key:   attribute.Key(a.Key),
				Value: attribute.Int64SliceValue(a.Value.([]int64)),
			})
		case "float64Slice":
			res = append(res, attribute.KeyValue{
				Key:   attribute.Key(a.Key),
				Value: attribute.Float64SliceValue(a.Value.([]float64)),
			})
		case "stringSlice":
			res = append(res, attribute.KeyValue{
				Key:   attribute.Key(a.Key),
				Value: attribute.StringSliceValue(a.Value.([]string)),
			})
		default:
			elog.Panic("otel attribute type error", elog.String("key", a.Key), elog.String("type", a.Type), elog.Any("value", a.Value))
		}
	}
	return res
}

func (config *Config) buildOtlpTP() trace.TracerProvider {
	var tpOptions []tracesdk.TracerProviderOption
	// 当开启了采集，才会设置采样率
	// traceExp 为 nil，不采集，但是目前来看在ego 1.2 版本，怀疑otel升级后，traceExp确实不为空，导致有采集行为，影响性能，所以需要关闭这个
	// 所以通过这个方式，屏蔽采集
	if config.Fraction == 0 {
		tp := tracesdk.NewTracerProvider(tpOptions...)
		return tp
	}

	// otlp exporter
	options := []otlptracegrpc.Option{
		otlptracegrpc.WithHeaders(config.Otlp.Headers),   // WithHeaders will send the provided headers with each gRPC requests.
		otlptracegrpc.WithEndpoint(config.Otlp.Endpoint), // WithEndpoint sets the target endpoint the exporter will connect to. If unset, localhost:4317 will be used as a default.
		// otlptracegrpc.WithDialOption(grpc.WithBlock()), //默认不设置 同步状态，会产生阻塞等待 Ready
	}
	if config.Otlp.EnableInsecure {
		// WithInsecure disables client transport security for the exporter's gRPC
		options = append(options, otlptracegrpc.WithInsecure())
	}

	options = append(options, config.Otlp.options...)
	traceClient := otlptracegrpc.NewClient(options...)
	ctx := context.Background()
	traceExp, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		elog.Error("otlp exporter error", elog.FieldErr(err))
		return nil
	}

	// res
	extraDefaultAttributes := generateDefaultAttributes(config.ServiceName, config.DefaultAttributes...)
	resOptions := []resource.Option{
		resource.WithTelemetrySDK(), // WithTelemetrySDK adds TelemetrySDK version info to the configured resource.
		resource.WithHost(),         // WithHost adds attributes from the host to the configured resource.
		resource.WithAttributes(extraDefaultAttributes...),
	}
	resOptions = append(resOptions, config.Otlp.resOptions...)
	res, err := resource.New(ctx, resOptions...)
	if err != nil {
		elog.Error("otlp resource New error", elog.FieldErr(err))
		return nil
	}

	// tp
	tpOptions = []tracesdk.TracerProviderOption{
		// Set the sampling rate based on the parent span to 100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(config.Fraction))),
		// WithSpanProcessor registers the SpanProcessor with a TracerProvider.
		// traceExp 为 nil，不采集，但是目前来看在ego 1.2 版本，怀疑otel升级后，traceExp确实不为空，导致有采集行为，影响性能，所以需要关闭这个
		tracesdk.WithSpanProcessor(tracesdk.NewBatchSpanProcessor(traceExp)),
		// Record information about this application in a Resource.
		tracesdk.WithResource(res),
	}

	tpOptions = append(tpOptions, config.options...)
	tp := tracesdk.NewTracerProvider(tpOptions...)
	return tp
}

// Stop 停止
func (config *Config) Stop() error {
	return nil
}
