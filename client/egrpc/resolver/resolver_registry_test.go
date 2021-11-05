package resolver

import (
	"context"
	"net/url"
	"testing"

	"github.com/gotomicro/ego/core/eregistry"
	"github.com/gotomicro/ego/server"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"
)

func TestResolver(t *testing.T) {
	builder := &baseBuilder{
		name: "test",
		reg: &testRegistry{
			t: t,
		},
	}
	targetName := "test:///hello"
	target, err := parseTarget(targetName)
	assert.NoError(t, err)
	resolve, err := builder.Build(target, &testRegistry{}, resolver.BuildOptions{})
	assert.NoError(t, err)
	resolve.ResolveNow(resolver.ResolveNowOptions{})
}

// parseTarget uses RFC 3986 semantics to parse the given target into a
// resolver.Target struct containing scheme, authority and endpoint. Query
// params are stripped from the endpoint.
func parseTarget(target string) (resolver.Target, error) {
	u, err := url.Parse(target)
	if err != nil {
		return resolver.Target{}, err
	}
	return resolver.Target{
		URL: *u,
	}, nil
}

type testRegistry struct {
	resolver.ClientConn
	t *testing.T
}

// ListServices ...
func (n testRegistry) ListServices(ctx context.Context, target eregistry.Target) ([]*server.ServiceInfo, error) {
	return nil, nil
}

// WatchServices ...
func (n testRegistry) WatchServices(ctx context.Context, target eregistry.Target) (chan eregistry.Endpoints, error) {
	assert.Equal(n.t, "hello", target.Endpoint)
	assert.Equal(n.t, "test", target.Scheme)
	assert.Equal(n.t, "grpc", target.Protocol)
	return nil, nil
}

// RegisterService ...
func (n testRegistry) RegisterService(context.Context, *server.ServiceInfo) error { return nil }

// UnregisterService ...
func (n testRegistry) UnregisterService(context.Context, *server.ServiceInfo) error { return nil }

// SyncServices 同步所有服务
func (n testRegistry) SyncServices(context.Context, eregistry.SyncServicesOptions) error { return nil }

// Close ...
func (n testRegistry) Close() error { return nil }
