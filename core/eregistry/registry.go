package eregistry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"google.golang.org/grpc/resolver"

	"github.com/gotomicro/ego/server"
)

// Registry register/unregister service
// registry impl should control rpc timeout
type Registry interface {
	RegisterService(context.Context, *server.ServiceInfo) error
	UnregisterService(context.Context, *server.ServiceInfo) error
	ListServices(context.Context, Target) ([]*server.ServiceInfo, error)
	WatchServices(context.Context, Target) (chan Endpoints, error)
	SyncServices(context.Context, SyncServicesOptions) error
	io.Closer
}

const (
	ProtocolGRPC = "grpc"
	ProtocolHTTP = "http"
)

type Target struct {
	Protocol  string // "http"|"grpc"
	Scheme    string // "etcd"|"k8s"|"dns"
	Endpoint  string // "<SVC-NAME>:<PORT>"
	Authority string
}

type SyncServicesOptions struct {
	GrpcResolverNowOptions resolver.ResolveNowOptions
}

// GetServiceKey ..
func GetServiceKey(prefix string, s *server.ServiceInfo) string {
	return fmt.Sprintf("/%s/%s/%s/%s://%s", prefix, s.Name, s.Kind.String(), s.Scheme, s.Address)
}

// GetServiceValue ..
func GetServiceValue(s *server.ServiceInfo) string {
	val, _ := json.Marshal(s)
	return string(val)
}

// GetService ..
func GetService(s string) *server.ServiceInfo {
	var si server.ServiceInfo
	_ = json.Unmarshal([]byte(s), &si)
	return &si
}

// Nop registry, used for local development/debugging
type Nop struct{}

// ListServices ...
func (n Nop) ListServices(ctx context.Context, target Target) ([]*server.ServiceInfo, error) {
	panic("implement me")
}

// WatchServices ...
func (n Nop) WatchServices(ctx context.Context, target Target) (chan Endpoints, error) {
	panic("implement me")
}

// RegisterService ...
func (n Nop) RegisterService(context.Context, *server.ServiceInfo) error { return nil }

// UnregisterService ...
func (n Nop) UnregisterService(context.Context, *server.ServiceInfo) error { return nil }

func (n Nop) SyncServices(context.Context, SyncServicesOptions) error { return nil }

// Close ...
func (n Nop) Close() error { return nil }

// Configuration ...
type Configuration struct {
	Routes []Route           `json:"routes"` // 配置客户端路由策略
	Labels map[string]string `json:"labels"` // 配置服务端标签: 分组
}

// Route represents route configuration
type Route struct {
	// 路由方法名
	Method string `json:"method" toml:"method"`
	// 路由权重组, 按比率在各个权重组中分配流量
	WeightGroups []WeightGroup `json:"weightGroups" toml:"weightGroups"`
	// 路由部署组, 将流量导入部署组
	Deployment string `json:"deployment" toml:"deployment"`
}

// WeightGroup ...
type WeightGroup struct {
	Group  string `json:"group" toml:"group"`
	Weight int    `json:"weight" toml:"weight"`
}
