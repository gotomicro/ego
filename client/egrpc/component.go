package egrpc

import (
	"context"
	"github.com/gotomicro/ego/core/ecode"
	"github.com/gotomicro/ego/core/elog"
	"time"

	"google.golang.org/grpc"
)

type Component struct {
	name   string
	config *Config
	logger *elog.Component
	*grpc.ClientConn
}

func newComponent(name string, config *Config, logger *elog.Component) *Component {
	var ctx = context.Background()
	var dialOptions = config.dialOptions
	// 默认配置使用block
	if config.Block {
		if config.DialTimeout > time.Duration(0) {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, config.DialTimeout)
			defer cancel()
		}

		dialOptions = append(dialOptions, grpc.WithBlock())
	}

	if config.KeepAlive != nil {
		dialOptions = append(dialOptions, grpc.WithKeepaliveParams(*config.KeepAlive))
	}

	dialOptions = append(dialOptions, grpc.WithBalancerName(config.BalancerName))

	cc, err := grpc.DialContext(ctx, config.Address, dialOptions...)

	if err != nil {
		if config.OnDialError == "panic" {
			logger.Panic("dial grpc server", elog.FieldErrKind(ecode.ErrKindRequestErr), elog.FieldErr(err))
		} else {
			logger.Error("dial grpc server", elog.FieldErrKind(ecode.ErrKindRequestErr), elog.FieldErr(err))
		}
	}
	logger.Info("start grpc client", elog.FieldName(name))
	return &Component{
		name:       name,
		config:     config,
		logger:     logger,
		ClientConn: cc,
	}
}
