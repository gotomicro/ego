package egrpc

import (
	"context"
	"time"

	"google.golang.org/grpc"

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

	dialOptions = append(dialOptions, grpc.WithBalancerName(config.BalancerName)) //nolint
	dialOptions = append(dialOptions, grpc.FailOnNonTempDialError(config.EnableFailOnNonTempDialError))

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
			logger.Panic("dial grpc server", elog.FieldErrKind("request err"), elog.FieldErr(err), elog.FieldKey(name), elog.FieldAddr(config.Addr))
			return component
		}
		logger.Error("dial grpc server", elog.FieldErrKind("request err"), elog.FieldErr(err), elog.FieldKey(name), elog.FieldAddr(config.Addr))
		return component
	}
	logger.Info("start grpc client", elog.FieldName(name))
	return component
}

// Error 错误信息
func (c *Component) Error() error {
	return c.err
}
