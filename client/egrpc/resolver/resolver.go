package resolver

import (
	"context"
	"reflect"
	"strings"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"

	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/eregistry"
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

	// grpc新版本1.40以后，会采用url parse解析，获取endpoint，但这个方法官方说了会有些问题。
	// 而 Ego 支持 unix socket，所以需要做一些兼容处理，详情请看 grpc.ClientConn.parseTarget 方法
	// For targets of the form "[scheme]://[authority]/endpoint, the endpoint
	// value returned from url.Parse() contains a leading "/". Although this is
	// in accordance with RFC 3986, we do not want to break existing resolver
	// implementations which expect the endpoint without the leading "/". So, we
	// end up stripping the leading "/" here. But this will result in an
	// incorrect parsing for something like "unix:///path/to/socket". Since we
	// own the "unix" resolver, we can workaround in the unix resolver by using
	// the `URL` field instead of the `Endpoint` field.

	endpoint := target.URL.Path
	if endpoint == "" {
		endpoint = target.URL.Opaque
	}
	endpoint = strings.TrimPrefix(endpoint, "/")
	endpoints, err := b.reg.WatchServices(ctx, eregistry.Target{
		Protocol:  eregistry.ProtocolGRPC,
		Scheme:    target.URL.Scheme,
		Endpoint:  endpoint,
		Authority: target.URL.Host,
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
	go func() {
		for {
			select {
			case endpoint := <-endpoints:
				var state = resolver.State{
					Addresses: make([]resolver.Address, 0),
					Attributes: attributes.New(constant.KeyRouteConfig, endpoint.RouteConfigs). // 路由配置
															WithValue(constant.KeyProviderConfig, endpoint.ProviderConfigs). // 服务提供方元信息
															WithValue(constant.KeyConsumerConfig, endpoint.ConsumerConfigs), // 服务消费方配置信息
				}
				b.tryUpdateAttrs(endpoint.Nodes)
				for key, node := range endpoint.Nodes {
					var address resolver.Address
					address.Addr = node.Address
					address.ServerName = b.target.URL.Path
					address.Attributes = b.attrs[key]
					state.Addresses = append(state.Addresses, address)
				}
				_ = b.cc.UpdateState(state)
			case <-b.stop:
				return
			}
		}
	}()
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
