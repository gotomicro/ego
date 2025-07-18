package egin

import (
	"fmt"

	healthcheck "github.com/RaMin0/gin-health-check"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	rpcpb "google.golang.org/genproto/googleapis/rpc/context/attribute_context"

	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/etrace"
	"github.com/gotomicro/ego/core/util/xnet"
)

// Container defines a component instance.
type Container struct {
	config *Config
	name   string
	logger *elog.Component
}

// DefaultContainer returns an default container.
func DefaultContainer() *Container {
	return &Container{
		config: DefaultConfig(),
		logger: elog.EgoLogger.With(elog.FieldComponent(PackageName)),
	}
}

// Load parses container configuration from configuration provider, such as a toml file,
// then use the configuration to construct a component container.
func Load(key string) *Container {
	c := DefaultContainer()
	c.logger = c.logger.With(elog.FieldComponentName(key))
	if err := econf.UnmarshalKey(key, &c.config); err != nil {
		c.logger.Panic("parse config error", elog.FieldErr(err), elog.FieldKey(key))
		return c
	}
	var host string
	var err error
	if err := c.setAiReqResCelPrg(); err != nil {
		c.logger.Warn("init AccessInterceptorReqResFilter fail", elog.FieldErr(err), elog.String("AccessInterceptorReqResFilter", c.config.AccessInterceptorReqResFilter))
	}
	// 获取网卡ip
	if c.config.EnableLocalMainIP {
		host, _, err = xnet.GetLocalMainIP()
		if err != nil {
			elog.Error("get local main ip error", elog.FieldErr(err))
		} else {
			c.config.Host = host
		}
	}
	c.name = key
	return c
}

var aiReqResCelEnv *cel.Env

func init() {
	var err error
	aiReqResCelEnv, err = cel.NewEnv(
		cel.Types(&rpcpb.AttributeContext_Request{}),
		cel.Types(&rpcpb.AttributeContext_Response{}),
		cel.Declarations(
			decls.NewVar("request",
				decls.NewObjectType("google.rpc.context.AttributeContext.Request"),
			),
			decls.NewVar("response",
				decls.NewObjectType("google.rpc.context.AttributeContext.Response"),
			),
		),
	)
	if err != nil {
		elog.Warn("invalid aiReqResCelEnv", elog.FieldErr(err))
	}
}

func (c *Container) setAiReqResCelPrg() error {
	if c.config.AccessInterceptorReqResFilter != "" {
		c.logger.Info("load new AccessInterceptorReqResFilter", elog.String("filter", c.config.AccessInterceptorReqResFilter))
		ast, iss := aiReqResCelEnv.Compile(c.config.AccessInterceptorReqResFilter)
		if iss.Err() != nil {
			return fmt.Errorf("invalid AccessInterceptorReqResFilter, %w", iss.Err())
		}
		prg, err := aiReqResCelEnv.Program(ast)
		if err != nil {
			return fmt.Errorf("build cel program fail , %w", err)
		}
		c.config.aiReqResCelPrg = prg
		return nil
	}
	c.config.aiReqResCelPrg = nil
	return nil
}

// Build constructs a specific component from container.
func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}

	server := newComponent(c.name, c.config, c.logger)
	server.Use(healthcheck.Default())
	server.Use(c.defaultServerInterceptor())
	server.Use(NewXResCostTimer(c.name, c.config.EnableResHeaderApp))
	if c.config.ContextTimeout > 0 {
		server.Use(timeoutMiddleware(c.config.ContextTimeout))
	}

	// if c.config.EnableMetricInterceptor {
	//	server.Use(metricServerInterceptor())
	// }
	if len(c.config.customGinMiddleware) > 0 {
		server.Use(c.config.customGinMiddleware...)
	}

	if c.config.EnableTraceInterceptor && etrace.IsGlobalTracerRegistered() {
		server.Use(traceServerInterceptor(c.config.compatibleTrace))
	}

	if c.config.EnableSentinel {
		server.Use(c.sentinelMiddleware())
	}

	econf.OnChange(func(newConf *econf.Configuration) {
		c.config.mu.Lock()
		cf := newConf.Sub(c.name)
		c.config.EnableAccessInterceptorReq = cf.GetBool("enableAccessInterceptorReq")
		c.config.EnableAccessInterceptorRes = cf.GetBool("enableAccessInterceptorRes")
		if c.config.AccessInterceptorReqResFilter != cf.GetString("accessInterceptorReqResFilter") {
			c.config.AccessInterceptorReqResFilter = cf.GetString("accessInterceptorReqResFilter")
			if err := c.setAiReqResCelPrg(); err != nil {
				c.logger.Warn("init AccessInterceptorReqResFilter fail", elog.FieldErr(err), elog.String("AccessInterceptorReqResFilter", c.config.AccessInterceptorReqResFilter))
			}
		}
		c.config.mu.Unlock()
	})

	return server
}
