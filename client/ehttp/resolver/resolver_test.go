package resolver

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"

	"github.com/gotomicro/ego/core/eregistry"
	"github.com/gotomicro/ego/server"
)

func TestResolver(t *testing.T) {
	builder := &baseBuilder{
		name: "test",
		reg: &testRegistry{
			t: t,
		},
	}
	assert.Equal(t, "test", builder.Scheme())
	addr := "test:///hello"
	_, err := parseTarget(addr)
	assert.NoError(t, err)

	resolve, err := builder.Build(addr)
	assert.NoError(t, err)
	resolve.GetAddr()
}

func TestResolver_http(t *testing.T) {
	var b = &baseHttpBuilder{}
	addr := "test:///hello"
	_, err := b.Build(addr)
	assert.Equal(t, nil, err)
	assert.Equal(t, "", b.Scheme())
	var b1 = &baseHttpResolver{}
	assert.Equal(t, "", b1.GetAddr())
}

// parseTarget uses RFC 3986 semantics to parse the given target into a
// resolver.Target struct containing scheme, authority and endpoint. Query
// params are stripped from the endpoint.
func parseTarget(addr string) (resolver.Target, error) {
	u, err := url.Parse(addr)
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
	assert.Equal(n.t, "http", target.Protocol)
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
