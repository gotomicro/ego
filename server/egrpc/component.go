package egrpc

import (
	"context"
	"net"

	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/test/bufconn"
)

// PackageName 包名
const PackageName = "server.egrpc"

// Component ...
type Component struct {
	name   string
	config *Config
	logger *elog.Component
	*grpc.Server
	listener   net.Listener
	serverInfo *server.ServiceInfo
	quit       chan error
}

func newComponent(name string, config *Config, logger *elog.Component) *Component {
	newServer := grpc.NewServer(config.serverOptions...)
	reflection.Register(newServer)
	healthpb.RegisterHealthServer(newServer, health.NewServer())

	return &Component{
		name:       name,
		config:     config,
		logger:     logger,
		Server:     newServer,
		listener:   nil,
		serverInfo: nil,
		quit:       make(chan error),
	}
}

// Name 配置名称
func (c *Component) Name() string {
	return c.name
}

// PackageName 包名
func (c *Component) PackageName() string {
	return PackageName
}

// Init 初始化
func (c *Component) Init() error {
	var (
		listener net.Listener
		err      error
	)
	// gRPC测试listener
	if c.config.Network == "bufnet" {
		listener = bufconn.Listen(1024 * 1024)
		c.listener = listener
		return nil
	}
	// 正式listener
	listener, err = net.Listen(c.config.Network, c.config.Address())
	if err != nil {
		c.logger.Panic("new grpc server err", elog.FieldErrKind("listen err"), elog.FieldErr(err))
	}
	c.config.Port = listener.Addr().(*net.TCPAddr).Port

	info := server.ApplyOptions(
		server.WithScheme("grpc"),
		server.WithAddress(listener.Addr().String()),
		server.WithKind(constant.ServiceProvider),
	)
	c.listener = listener
	c.serverInfo = &info
	return nil
}

// Start implements server.Component interface.
func (c *Component) Start() error {
	err := c.Server.Serve(c.listener)
	return err
}

// Stop implements server.Component interface
// it will terminate echo server immediately
func (c *Component) Stop() error {
	c.Server.Stop()
	return nil
}

// GracefulStop implements server.Component interface
// it will stop echo server gracefully
func (c *Component) GracefulStop(ctx context.Context) error {
	go func() {
		c.Server.GracefulStop()
		close(c.quit)
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c.quit:
			return nil
		}
	}
}

// Info returns server info, used by governor and consumer balancer
func (c *Component) Info() *server.ServiceInfo {
	return c.serverInfo
}

// Address 服务地址
func (c *Component) Address() string {
	return c.config.Address()
}

// Listener listener信息
func (c *Component) Listener() net.Listener {
	return c.listener
}
