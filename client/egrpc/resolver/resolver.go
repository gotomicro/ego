package resolver

import (
	"context"
	"reflect"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"

	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/eregistry"
	"github.com/gotomicro/ego/core/util/xgo"
	"github.com/gotomicro/ego/server"
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

	br := &baseResolver{
		target: target,
		cc:     cc,
		stop:   make(chan struct{}),
		reg:    b.reg,
		cancel: cancel,
		attrs:  make(map[string]*attributes.Attributes),
	}
	br.run(endpoints)
	return br, nil
}

// Scheme ...
func (b baseBuilder) Scheme() string {
	return b.name
}

type baseResolver struct {
	target resolver.Target
	cc     resolver.ClientConn
	stop   chan struct{}
	reg    eregistry.Registry
	cancel context.CancelFunc
	attrs  map[string]*attributes.Attributes
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

func (b *baseResolver) run(endpoints chan eregistry.Endpoints) {
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
				b.tryUpdateAttrs(endpoint.Nodes)
				for key, node := range endpoint.Nodes {
					var address resolver.Address
					address.Addr = node.Address
					address.ServerName = b.target.Endpoint
					address.Attributes = b.attrs[key]
					state.Addresses = append(state.Addresses, address)
				}
				b.cc.UpdateState(state)
			case <-b.stop:
				return
			}
		}
	})
}

func attrEqual(oldAttr *attributes.Attributes, node server.ServiceInfo) bool {
	oldNode := oldAttr.Value(constant.KeyServiceInfo)
	// NOTICE:目前暂时未用Services和Metadata，所以可以使用reflect.DeepEqual
	return reflect.DeepEqual(oldNode, node)
}

func (b *baseResolver) tryUpdateAttrs(nodes map[string]server.ServiceInfo) {
	for addr, node := range nodes {
		oldAttr, ok := b.attrs[addr]
		if !ok || !attrEqual(oldAttr, node) {
			attr := attributes.New(constant.KeyServiceInfo, node)
			b.attrs[addr] = attr
		}
	}
	for addr := range b.attrs {
		if _, ok := nodes[addr]; !ok {
			delete(b.attrs, addr)
		}
	}
}
