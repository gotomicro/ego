package egrpc

import (
	"context"
	"net"

	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/internal/egrpclog"
	"github.com/gotomicro/ego/server"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
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
	invokers   []func() error // 用户初始化函数
}

func newComponent(name string, config *Config, logger *elog.Component) *Component {
	if config.EnableOfficialGrpcLog {
		// grpc框架日志，因为官方grpc日志是单例，所以这里要处理下
		grpclog.SetLoggerV2(zapgrpc.NewLogger(egrpclog.Build().ZapLogger()))
	}
	newServer := grpc.NewServer(config.serverOptions...)
	reflection.Register(newServer)
	healthSvc := health.NewServer()
	// server should register all the services manually
	// use empty service name for all etcd services' health status,
	// see https://github.com/grpc/grpc/blob/master/doc/health-checking.md for more
	healthSvc.SetServingStatus(eapp.Name(), healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(newServer, healthSvc)
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

// Prepare 预准备一些数据
func (c *Component) Prepare() error {
	for _, fn := range c.invokers {
		err := fn()
		if err != nil {
			return err
		}
	}
	return nil
}

// Init 初始化一些信息
func (c *Component) Init() error {
	info := server.ApplyOptions(
		server.WithScheme("grpc"),
		server.WithAddress(c.config.Address()),
		server.WithKind(constant.ServiceProvider),
	)
	c.serverInfo = &info
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
	tcpInfo, flag := listener.Addr().(*net.TCPAddr)
	if flag {
		c.config.Port = tcpInfo.Port
	}
	c.listener = listener
	return nil
}

// Start implements server.Component interface.
func (c *Component) Start() error {
	return c.Server.Serve(c.listener)
}

// Health implements server.Component interface.
// Experimental
func (c *Component) Health() bool {
	addr := c.config.Address()
	if c.config.Network == "unix" {
		addr = "unix:" + addr
	}
	cc, err := grpc.DialContext(context.Background(), addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.logger.Error("health connection err", elog.FieldErr(err))
		return false
	}
	healthClient := healthpb.NewHealthClient(cc)
	resp, err := healthClient.Check(context.Background(), &healthpb.HealthCheckRequest{
		Service: eapp.Name(),
	})
	if err != nil {
		c.logger.Error("health rpc err", elog.FieldErr(err))
		return false
	}
	c.logger.Info("grpc health connection OK")
	return resp.Status == healthpb.HealthCheckResponse_SERVING
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

// Invoker returns server info, used by governor and consumer balancer
func (c *Component) Invoker(fns ...func() error) {
	c.invokers = append(c.invokers, fns...)
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
