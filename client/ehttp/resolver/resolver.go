package resolver

import (
	"context"
	"net/url"
	"strings"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"

	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/eregistry"
	"github.com/gotomicro/ego/server"
)

var (
	// m is a map from scheme to resolver builder.
	m = make(map[string]Builder)
	// defaultScheme is the default scheme to use.
	// defaultScheme = "http"
)

// Builder creates a resolver that will be used to watch name resolution updates.
type Builder interface {
	// Build creates a new resolver for the given target.
	// gRPC dial calls Build synchronously, and fails if the returned error is
	// not nil.
	Build(addr string) (Resolver, error)
	// Scheme returns the scheme supported by this resolver.
	Scheme() string
}

type Resolver interface {
	GetAddr() string
}

// Register ...
func Register(name string, reg eregistry.Registry) {
	b := &baseBuilder{
		name: name,
		reg:  reg,
	}
	m[b.Scheme()] = b
}

func Get(scheme string) Builder {
	if b, ok := m[scheme]; ok {
		return b
	}
	return nil
}

type baseBuilder struct {
	name string
	reg  eregistry.Registry
}

// Build ...
func (b *baseBuilder) Build(addr string) (Resolver, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	target, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	endpoint := target.Path
	if endpoint == "" {
		endpoint = target.Opaque
	}
	endpoint = strings.TrimPrefix(endpoint, "/")

	egoTarget := eregistry.Target{
		Protocol:  eregistry.ProtocolHTTP,
		Scheme:    target.Scheme,
		Endpoint:  endpoint,
		Authority: target.Host,
	}

	endpoints, err := b.reg.WatchServices(ctx, egoTarget)
	if err != nil {
		cancel()
		return nil, err
	}

	br := &baseResolver{
		target:   egoTarget,
		stop:     make(chan struct{}),
		reg:      b.reg,
		cancel:   cancel,
		nodeInfo: make(map[string]*attributes.Attributes),
	}
	br.run(endpoints)
	return br, nil
}

// Scheme ...
func (b *baseBuilder) Scheme() string {
	return b.name
}

type baseResolver struct {
	target eregistry.Target // 使用ego的target，因为官方的target后续会不兼容
	stop   chan struct{}
	reg    eregistry.Registry
	cancel context.CancelFunc
	// addrSlices []string
	nodeInfo map[string]*attributes.Attributes // node节点的属性
}

func (b *baseResolver) GetAddr() string {
	for key := range b.nodeInfo {
		return "http://" + key
	}
	return ""
}

// Close ...
func (b *baseResolver) Close() {
	b.stop <- struct{}{}
	b.cancel()
}

// run 更新节点信息
// State
//
//	     Addresses   []Address{  IP列表
//	                 	Addr: IP 地址,
//							ServerName: 应用名称, 如：svc-user
//							Attributes: 节点基本信息： server.ServiceInfo
//	                 }
//	     Attributes： {  用于负载均衡的配置，目前需要通过后台来设置
//							constant.KeyRouteConfig    路由配置
//							constant.KeyProviderConfig 服务提供方元信息
//							constant.KeyConsumerConfig 服务消费方配置信息
//	                  }
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
				// 如果node信息有变更，那么就添加，更新或者删除
				b.tryUpdateAttrs(endpoint.Nodes)
				for key, node := range endpoint.Nodes {
					var address resolver.Address
					address.Addr = node.Address
					address.ServerName = b.target.Endpoint
					address.Attributes = b.nodeInfo[key]
					state.Addresses = append(state.Addresses, address)
				}
			case <-b.stop:
				return
			}
		}
	}()
}

// tryUpdateAttrs 更新节点数据
func (b *baseResolver) tryUpdateAttrs(nodes map[string]server.ServiceInfo) {
	for addr, node := range nodes {
		oldAttr, ok := b.nodeInfo[addr]
		newAttr := attributes.New(constant.KeyServiceInfo, node)
		if !ok || !oldAttr.Equal(newAttr) {
			b.nodeInfo[addr] = newAttr
		}
	}
	for addr := range b.nodeInfo {
		if _, ok := nodes[addr]; !ok {
			delete(b.nodeInfo, addr)
		}
	}
}
