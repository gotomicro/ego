package ejaeger

import (
	"os"

	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/econf"
	jaegerv2 "go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

// Config ...
type Config struct {
	ServiceName  string
	Addr         string
	Fraction     float64
	PanicOnError bool
	options      []tracesdk.TracerProviderOption
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
	agentAddr := "127.0.0.1:6831"
	if addr := os.Getenv("JAEGER_AGENT_ADDR"); addr != "" {
		agentAddr = addr
	}
	return &Config{
		ServiceName:  eapp.Name(),
		Addr:         agentAddr,
		PanicOnError: true,
	}
}

// WithTracerProviderOption ...
func (config *Config) WithTracerProviderOption(options ...tracesdk.TracerProviderOption) *Config {
	config.options = append(config.options, options...)
	return config
}

// Build ...
func (config *Config) Build(ops ...Option) trace.TracerProvider {
	// Create the Jaeger exporter
	exp, err := jaegerv2.New(jaegerv2.WithCollectorEndpoint(jaegerv2.WithEndpoint(config.Addr)))
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
