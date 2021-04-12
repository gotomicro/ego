package compound

import (
	"context"

	"golang.org/x/sync/errgroup"

	cregistry "github.com/gotomicro/ego/core/eregistry"
	"github.com/gotomicro/ego/server"
)

type compoundRegistry struct {
	registries []cregistry.Registry
}

// ListServices ...
func (c compoundRegistry) ListServices(ctx context.Context, target cregistry.Target) ([]*server.ServiceInfo, error) {
	var eg errgroup.Group
	var services = make([]*server.ServiceInfo, 0)
	for _, registry := range c.registries {
		registry := registry
		eg.Go(func() error {
			infos, err := registry.ListServices(ctx, target)
			if err != nil {
				return err
			}
			services = append(services, infos...)
			return nil
		})
	}
	err := eg.Wait()
	return services, err
}

// WatchServices ...
func (c compoundRegistry) WatchServices(ctx context.Context, target cregistry.Target) (chan cregistry.Endpoints, error) {
	panic("compound registry doesn't support watch services")
}

// SyncServices ...
func (c compoundRegistry) SyncServices(context.Context, cregistry.SyncServicesOptions) error {
	panic("compound registry doesn't support sync services")
}

// RegisterService ...
func (c compoundRegistry) RegisterService(ctx context.Context, bean *server.ServiceInfo) error {
	var eg errgroup.Group
	for _, registry := range c.registries {
		registry := registry
		eg.Go(func() error {
			return registry.RegisterService(ctx, bean)
		})
	}
	return eg.Wait()
}

// UnregisterService ...
func (c compoundRegistry) UnregisterService(ctx context.Context, bean *server.ServiceInfo) error {
	var eg errgroup.Group
	for _, registry := range c.registries {
		registry := registry
		eg.Go(func() error {
			return registry.UnregisterService(ctx, bean)
		})
	}
	return eg.Wait()
}

// Close ...
func (c compoundRegistry) Close() error {
	var eg errgroup.Group
	for _, registry := range c.registries {
		registry := registry
		eg.Go(func() error {
			return registry.Close()
		})
	}
	return eg.Wait()
}

// New ...
func New(registries ...cregistry.Registry) cregistry.Registry {
	return compoundRegistry{
		registries: registries,
	}
}
