package resolver

import (
	"context"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"

	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/eregistry"
	"github.com/gotomicro/ego/core/util/xgo"
)

// Register ...
func Register(name string, reg eregistry.Registry) {
	resolver.Register(&baseBuilder{
		name: name,
		reg:  reg,
	})
}

type baseBuilder struct {
	name string
	reg  eregistry.Registry
}

// Build ...
func (b *baseBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	ctx, cancel := context.WithCancel(context.Background())
	endpoints, err := b.reg.WatchServices(ctx, eregistry.Target{
		Protocol:  eregistry.ProtocolGRPC,
		Scheme:    target.Scheme,
		Endpoint:  target.Endpoint,
		Authority: target.Authority,
	})
	if err != nil {
		cancel()
		return nil, err
	}

	var stop = make(chan struct{})
	xgo.Go(func() {
		for {
			select {
			case endpoint := <-endpoints:
				var state = resolver.State{
					Addresses: make([]resolver.Address, 0),
					Attributes: attributes.New(
						constant.KeyRouteConfig, endpoint.RouteConfigs, // 路由配置
						constant.KeyProviderConfig, endpoint.ProviderConfigs, // 服务提供方元信息
						constant.KeyConsumerConfig, endpoint.ConsumerConfigs, // 服务消费方配置信息
					),
				}
				for _, node := range endpoint.Nodes {
					var address resolver.Address
					address.Addr = node.Address
					address.ServerName = target.Endpoint
					address.Attributes = attributes.New(constant.KeyServiceInfo, node)
					state.Addresses = append(state.Addresses, address)
				}
				cc.UpdateState(state)
			case <-stop:
				return
			}
		}
	})

	return &baseResolver{
		stop:   stop,
		reg:    b.reg,
		cancel: cancel,
	}, nil
}

// Scheme ...
func (b baseBuilder) Scheme() string {
	return b.name
}

type baseResolver struct {
	stop   chan struct{}
	reg    eregistry.Registry
	cancel context.CancelFunc
}

// ResolveNow ...
func (b *baseResolver) ResolveNow(options resolver.ResolveNowOptions) {
	if err := b.reg.SyncServices(context.Background(), eregistry.SyncServicesOptions{GrpcResolverNowOptions: options}); err != nil {
		elog.Error("ResolveNow fail", elog.FieldErr(err))
	}
}

// Close ...
func (b *baseResolver) Close() {
	b.stop <- struct{}{}
	b.cancel()
}
