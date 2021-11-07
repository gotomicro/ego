package eregistry

import (
	"context"

	"github.com/gotomicro/ego/server"
)

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

// SyncServices 同步所有服务
func (n Nop) SyncServices(context.Context, SyncServicesOptions) error { return nil }

// Close ...
func (n Nop) Close() error { return nil }
