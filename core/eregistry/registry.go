package eregistry

import (
	"context"
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
	// ProtocolGRPC ...
	ProtocolGRPC = "grpc"
	// ProtocolHTTP ...
	ProtocolHTTP = "http"
)

// Target ...
type Target struct {
	Protocol  string // "http"|"grpc"
	Scheme    string // "etcd"|"k8s"|"dns"
	Endpoint  string // "<SVC-NAME>:<PORT>"
	Authority string
}

// SyncServicesOptions ...
type SyncServicesOptions struct {
	GrpcResolverNowOptions resolver.ResolveNowOptions
}

// Deprecated: Use *server.ServiceInfo.GetServiceKey()
// GetServiceKey ETCD注册需要使用
func GetServiceKey(prefix string, s *server.ServiceInfo) string {
	return s.GetServiceKey(prefix)
}

// Deprecated: Use *server.ServiceInfo.GetServiceValue()
// GetServiceValue ETCD注册需要使用
func GetServiceValue(s *server.ServiceInfo) string {
	return s.GetServiceValue()
}

//func GetService(s string) *server.ServiceInfo {
//	var si server.ServiceInfo
//	_ = json.Unmarshal([]byte(s), &si)
//	return &si
//}

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
