package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/eflag"
	"github.com/gotomicro/ego/core/standard"
)

// Option overrides a Container's default configuration.
type Option func(c *ServiceInfo)

// ConfigInfo represents service configurator
type ConfigInfo struct {
	Routes []Route
}

// ServiceInfo represents service info
type ServiceInfo struct {
	Name       string               `json:"name"`
	Scheme     string               `json:"scheme"`
	Address    string               `json:"address"`
	Weight     float64              `json:"weight"`
	Enable     bool                 `json:"enable"`
	Healthy    bool                 `json:"healthy"`
	Metadata   map[string]string    `json:"metadata"`
	Region     string               `json:"region"`
	Zone       string               `json:"zone"`
	Kind       constant.ServiceKind `json:"kind"`
	Deployment string               `json:"deployment"` // Deployment 部署组: 不同组的流量隔离,  比如某些服务给内部调用和第三方调用，可以配置不同的deployment,进行流量隔离
	Group      string               `json:"group"`      // Group 流量组: 流量在Group之间进行负载均衡
	Services   map[string]*Service  `json:"services" toml:"services"`
}

// Service ...
type Service struct {
	Namespace string            `json:"namespace" toml:"namespace"`
	Name      string            `json:"name" toml:"name"`
	Labels    map[string]string `json:"labels" toml:"labels"`
	Methods   []string          `json:"methods" toml:"methods"`
}

// Label ...
func (si ServiceInfo) Label() string {
	return fmt.Sprintf("%s://%s", si.Scheme, si.Address)
}

// Equal 一定要实现这个方法，在gRPC的attributes里会使用该方法断言，判断是否相等
func (si ServiceInfo) Equal(o interface{}) bool {
	return reflect.DeepEqual(si, o)
}

// GetServiceValue ETCD注册需要使用的服务信息
func (si *ServiceInfo) GetServiceValue() string {
	val, _ := json.Marshal(si)
	return string(val)
}

// GetServiceKey ETCD注册需要使用
func (si *ServiceInfo) GetServiceKey(prefix string) string {
	return fmt.Sprintf("/%s/%s/%s/%s://%s", prefix, si.Name, si.Kind.String(), si.Scheme, si.Address)
}

// Server ...
type Server interface {
	standard.Component
	GracefulStop(ctx context.Context) error
	Info() *ServiceInfo
}

// OrderServer ...
// Experimental
type OrderServer interface {
	standard.Component
	// Prepare 用于一些准备数据
	// 因为在OrderServer中，也会有invoker操作，需要放这个里面执行，需要区分他和真正server的init操作
	// server的init操作有一些listen，必须先执行，否则有些通信，会有问题
	Prepare() error
	GracefulStop(ctx context.Context) error
	Info() *ServiceInfo
	Health() bool
	Invoker(fns ...func() error) // 用户初始化函数，放在order server里执行
}

// Route ...
type Route struct {
	// 权重组，按照
	WeightGroups []WeightGroup
	// 方法名
	Method string
}

// WeightGroup ...
type WeightGroup struct {
	Group  string
	Weight int
}

// ApplyOptions 设置可选项
func ApplyOptions(options ...Option) ServiceInfo {
	info := defaultServiceInfo()
	for _, option := range options {
		option(&info)
	}
	return info
}

// WithMetaData 设置metadata信息
func WithMetaData(key, value string) Option {
	return func(c *ServiceInfo) {
		c.Metadata[key] = value
	}
}

// WithScheme 设置协议
func WithScheme(scheme string) Option {
	return func(c *ServiceInfo) {
		c.Scheme = scheme
	}
}

// WithAddress 设置地址
func WithAddress(address string) Option {
	return func(c *ServiceInfo) {
		c.Address = address
	}
}

// WithName 定义服务名称
func WithName(name string) Option {
	return func(c *ServiceInfo) {
		c.Name = name
	}
}

// WithKind 设置类型
func WithKind(kind constant.ServiceKind) Option {
	return func(c *ServiceInfo) {
		c.Kind = kind
	}
}

func defaultServiceInfo() ServiceInfo {
	si := ServiceInfo{
		Name:       eapp.Name(),
		Weight:     100,
		Enable:     true,
		Healthy:    true,
		Metadata:   make(map[string]string),
		Region:     eapp.AppRegion(),
		Zone:       eapp.AppZone(),
		Kind:       0,
		Deployment: "",
		Group:      "",
	}
	si.Metadata["appMode"] = eapp.AppMode()
	si.Metadata["appHost"] = eflag.String("host")
	si.Metadata["startTime"] = eapp.StartTime()
	si.Metadata["buildTime"] = eapp.BuildTime()
	si.Metadata["appVersion"] = eapp.AppVersion()
	si.Metadata["egoVersion"] = eapp.EgoVersion()
	si.Metadata["depEnv"] = os.Getenv(constant.EgoDeploymentEnv) // 部署环境
	return si
}
