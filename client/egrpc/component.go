package egrpc

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"

	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/internal/egrpclog"
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
			ctx, cancel = context.WithTimeoutCause(ctx, config.DialTimeout, fmt.Errorf("grpc client conn dial timeout"))
			defer cancel()
		}

		dialOptions = append(dialOptions, grpc.WithBlock())
	}

	if config.EnableWithInsecure {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if config.keepAlive != nil {
		dialOptions = append(dialOptions, grpc.WithKeepaliveParams(*config.keepAlive))
	}

	// service config 默认开启，且 grpc 1.46 及以上版本废弃了 WithBalancer 方法，改用 service config 配置 lb，但是开启后，在
	// grpc 1.40 以上会导致 dns 多一次解析 txt 内容（目测是为了做 grpc 的 load balance 策略，但我们实际上不会用到）
	// 因为这个 service config dns域名通常是没有设置dns解析，所以会跳过k8s的dns，穿透到上一级的dns，而如果dns配置有问题或者不存在，那么会查询非常长的时间（通常在20s或者更长）
	// 在上面场景下，配置为 false 禁用 service config，可以加快我们的启动时间或者提升我们的性能，**但请注意，禁用后，我们配置的 LB 将无法生效，默认为 pick_first 策略**
	if !config.EnableServiceConfig {
		dialOptions = append(dialOptions, grpc.WithDisableServiceConfig())
		if config.BalancerName != "pick_first" {
			elog.Warn(fmt.Sprintf("The LB policy `%s` will be ignored and use `pick_first` as default since you disabled service config", config.BalancerName))
		}
	} else {
		dialOptions = append(dialOptions, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, config.BalancerName)))
	}

	dialOptions = append(dialOptions, grpc.FailOnNonTempDialError(config.EnableFailOnNonTempDialError))

	if config.MaxCallRecvMsgSize != DefaultMaxCallRecvMsgSize {
		dialOptions = append(dialOptions, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(config.MaxCallRecvMsgSize)))
	}

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
