package server

import (
	"context"
	"fmt"
	"github.com/gotomicro/ego/core/app"
	"github.com/gotomicro/ego/core/standard"

	"github.com/gotomicro/ego/core/constant"
)

type Option func(c *ServiceInfo)

// ServiceConfigurator represents service configurator
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

// Server ...
type Server interface {
	standard.Component
	GracefulStop(ctx context.Context) error
	Info() *ServiceInfo
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

func ApplyOptions(options ...Option) ServiceInfo {
	info := defaultServiceInfo()
	for _, option := range options {
		option(&info)
	}
	return info
}

func WithMetaData(key, value string) Option {
	return func(c *ServiceInfo) {
		c.Metadata[key] = value
	}
}

func WithScheme(scheme string) Option {
	return func(c *ServiceInfo) {
		c.Scheme = scheme
	}
}

func WithAddress(address string) Option {
	return func(c *ServiceInfo) {
		c.Address = address
	}
}

func WithKind(kind constant.ServiceKind) Option {
	return func(c *ServiceInfo) {
		c.Kind = kind
	}
}

func defaultServiceInfo() ServiceInfo {
	si := ServiceInfo{
		Name:       app.Name(),
		Weight:     100,
		Enable:     true,
		Healthy:    true,
		Metadata:   make(map[string]string),
		Region:     app.AppRegion(),
		Zone:       app.AppZone(),
		Kind:       0,
		Deployment: "",
		Group:      "",
	}
	si.Metadata["appMode"] = app.AppMode()
	si.Metadata["appHost"] = app.AppHost()
	si.Metadata["startTime"] = app.StartTime()
	si.Metadata["buildTime"] = app.BuildTime()
	si.Metadata["appVersion"] = app.AppVersion()
	si.Metadata["egoVersion"] = app.EgoVersion()
	return si
}
