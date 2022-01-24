package ejaeger

import (
	"os"

	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	jaegerv2 "go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

// Config ...
type Config struct {
	ServiceName        string
	AgentHost          string // agent host
	AgentPort          string // agent port
	JaegerEndpointType string // type: agent,collector
	CollectorEndpoint  string // collector endpoint
	CollectorUser      string // collector user
	CollectorPassword  string // collector password
	Fraction           float64
	PanicOnError       bool
	options            []tracesdk.TracerProviderOption
}

// Load 加载配置key
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
		ServiceName:        eapp.Name(),
		AgentHost:          envOr("OTEL_EXPORTER_JAEGER_AGENT_HOST", "localhost"),
		AgentPort:          envOr("OTEL_EXPORTER_JAEGER_AGENT_PORT", "6831"),
		CollectorEndpoint:  envOr("OTEL_EXPORTER_JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
		CollectorUser:      envOr("OTEL_EXPORTER_JAEGER_USER", ""),
		CollectorPassword:  envOr("OTEL_EXPORTER_JAEGER_PASSWORD", ""),
		JaegerEndpointType: "agent",
		PanicOnError:       true,
	}
}

// envOr returns an env variable's value if it is exists or the default if not
func envOr(key, defaultValue string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return defaultValue
}

// WithTracerProviderOption ...
func (config *Config) WithTracerProviderOption(options ...tracesdk.TracerProviderOption) *Config {
	config.options = append(config.options, options...)
	return config
}

// Build ...
func (config *Config) Build(ops ...Option) trace.TracerProvider {
	var endpoint jaegerv2.EndpointOption
	switch config.JaegerEndpointType {
	case "agent":
		// Create the Jaeger exporter
		endpoint = jaegerv2.WithAgentEndpoint(
			jaegerv2.WithAgentHost(config.AgentHost),
			jaegerv2.WithAgentPort(config.AgentPort),
		)
	case "collector":
		endpoint = jaegerv2.WithCollectorEndpoint(
			jaegerv2.WithEndpoint(config.CollectorEndpoint),
			jaegerv2.WithUsername(config.CollectorUser),
			jaegerv2.WithPassword(config.CollectorPassword),
		)
	default:
		elog.Panic("jaeger type error", elog.FieldName(config.JaegerEndpointType))
	}

	exp, err := jaegerv2.New(endpoint)
	if err != nil {
		return nil
	}
	options := []tracesdk.TracerProviderOption{
		// Set the sampling rate based on the parent span to 100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(config.Fraction))),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(config.ServiceName),
		)),
	}
	options = append(options, config.options...)
	tp := tracesdk.NewTracerProvider(options...)
	return tp
}

// Stop 停止
func (config *Config) Stop() error {
	return nil
}
