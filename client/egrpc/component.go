package egrpc

import (
	"context"
	"time"

	"github.com/gotomicro/ego/internal/egrpclog"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"github.com/gotomicro/ego/core/elog"
)

// PackageName 设置包名
const PackageName = "client.egrpc"

// Component 组件
type Component struct {
	name   string
	config *Config
	logger *elog.Component
	*grpc.ClientConn
	err error
}

func newComponent(name string, config *Config, logger *elog.Component) *Component {
	if config.EnableOfficialGrpcLog {
		// grpc框架日志，因为官方grpc日志是单例，所以这里要处理下
		grpclog.SetLoggerV2(zapgrpc.NewLogger(egrpclog.Build().ZapLogger()))
	}
	var ctx = context.Background()
	var dialOptions = config.dialOptions
	// 默认配置使用block
	if config.EnableBlock {
		if config.DialTimeout > time.Duration(0) {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, config.DialTimeout)
			defer cancel()
		}

		dialOptions = append(dialOptions, grpc.WithBlock())
	}

	if config.EnableWithInsecure {
		dialOptions = append(dialOptions, grpc.WithInsecure())
	}

	if config.keepAlive != nil {
		dialOptions = append(dialOptions, grpc.WithKeepaliveParams(*config.keepAlive))
	}

	// 因为默认是开启这个配置
	// 并且开启后，在grpc 1.40以上会导致dns多一次解析txt内容（目测是为了做grpc的load balance策略，但我们实际上不会用到）
	// 因为这个service config dns域名通常是没有设置dns解析，所以会跳过k8s的dns，穿透到上一级的dns，而如果dns配置有问题或者不存在，那么会查询非常长的时间（通常在20s或者更长）
	// 那么为false的时候，禁用他，可以加快我们的启动时间或者提升我们的性能
	if !config.EnableServiceConfig {
		dialOptions = append(dialOptions, grpc.WithDisableServiceConfig())
	}

	dialOptions = append(dialOptions, grpc.WithBalancerName(config.BalancerName)) //nolint
	dialOptions = append(dialOptions, grpc.FailOnNonTempDialError(config.EnableFailOnNonTempDialError))

	startTime := time.Now()
	cc, err := grpc.DialContext(ctx, config.Addr, dialOptions...)

	component := &Component{
		name:       name,
		config:     config,
		logger:     logger,
		ClientConn: cc,
	}

	if err != nil {
		component.err = err
		if config.OnFail == "panic" {
			logger.Panic("dial grpc server", elog.FieldErrKind("request err"), elog.FieldErr(err), elog.FieldKey(name), elog.FieldAddr(config.Addr), elog.FieldCost(time.Since(startTime)))
			return component
		}
		logger.Error("dial grpc server", elog.FieldErrKind("request err"), elog.FieldErr(err), elog.FieldKey(name), elog.FieldAddr(config.Addr), elog.FieldCost(time.Since(startTime)))
		return component
	}
	logger.Info("start grpc client", elog.FieldName(name), elog.FieldCost(time.Since(startTime)))
	return component
}

// Error 错误信息
func (c *Component) Error() error {
	return c.err
}
