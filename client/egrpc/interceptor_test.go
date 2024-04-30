package egrpc

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"

	"github.com/gotomicro/ego/core/util/xtime"
	"github.com/gotomicro/ego/internal/test/helloworld"
	"github.com/gotomicro/ego/internal/tools"
)

func Test_customHeader(t *testing.T) {
	md := metadata.New(map[string]string{
		"X-Ego-Uid": "9527",
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	interceptor := customHeader([]string{"X-Ego-Uid"})

	cc := new(grpc.ClientConn)
	err := interceptor(ctx, "/foo", nil, nil, cc,
		func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			info := tools.GrpcHeaderValue(ctx, "X-Ego-Uid")
			assert.Equal(t, "9527", info)
			return nil
		})
	assert.Nil(t, err)
}

func TestMetric(t *testing.T) {
	cmp := DefaultContainer()
	cmp.name = "test"
	interceptor := cmp.metricUnaryClientInterceptor()
	cc := new(grpc.ClientConn)
	err := interceptor(context.Background(), "/foo", nil, nil, cc,
		func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			return nil
		})
	assert.Nil(t, err)
}

func TestTimeout(t *testing.T) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	helloworld.RegisterGreeterServer(server, &GreeterTimeout{})
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()
	begin := time.Now()
	cmp := DefaultContainer().Build(
		WithBufnetServerListener(listener),
		WithReadTimeout(2*time.Second),
	)
	cli := helloworld.NewGreeterClient(cmp.ClientConn)
	_, _ = cli.SayHello(context.Background(), &helloworld.HelloRequest{})
	cost := time.Since(begin)
	assert.True(t, cost > xtime.Duration("1.8s"))
}

func TestDebugLog(t *testing.T) {
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	helloworld.RegisterGreeterServer(server, &GreeterDebuglog{})
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()
	container := DefaultContainer()
	cmp := container.Build(
		WithBufnetServerListener(listener),
		WithDialOption(grpc.WithChainUnaryInterceptor(container.debugUnaryClientInterceptor())),
	)
	cli := helloworld.NewGreeterClient(cmp.ClientConn)
	_, _ = cli.SayHello(context.Background(), &helloworld.HelloRequest{})
	line := buf.String()
	assert.Contains(t, line, "Hello DebugLog")
	log.SetOutput(os.Stderr)
}

func TestCustomHeaderAppAndCpu(t *testing.T) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	helloworld.RegisterGreeterServer(server, &GreeterHeader{
		t: t,
	})
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()
	container := DefaultContainer()
	cmp := container.Build(
		WithBufnetServerListener(listener),
	)
	cli := helloworld.NewGreeterClient(cmp.ClientConn)
	_, _ = cli.SayHello(context.Background(), &helloworld.HelloRequest{})
}

func TestPrometheusUnary(t *testing.T) {
	// 1 获取prometheus的handler的数据
	ts := httptest.NewServer(promhttp.Handler())
	defer ts.Close()

	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	helloworld.RegisterGreeterServer(server, &GreeterDebuglog{})
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()
	container := DefaultContainer()
	cmp := container.Build(
		WithName("hello"),
		WithAddr("bufnet"),
		WithBufnetServerListener(listener),
	)
	cli := helloworld.NewGreeterClient(cmp.ClientConn)
	_, _ = cli.SayHello(context.Background(), &helloworld.HelloRequest{})

	pc := ts.Client()
	res, err := pc.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	text, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	err = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, string(text), `ego_client_handle_seconds_count{method="/helloworld.Greeter/SayHello",name="hello",peer="bufnet",type="unary"}`)
	assert.Contains(t, string(text), `ego_client_handle_seconds_bucket{method="/helloworld.Greeter/SayHello",name="hello",peer="bufnet",type="unary",le="0.005"}`)
}

// Greeter ...
type GreeterTimeout struct {
	helloworld.UnimplementedGreeterServer
}

// SayHello ...
func (g GreeterTimeout) SayHello(context context.Context, request *helloworld.HelloRequest) (*helloworld.HelloResponse, error) {
	time.Sleep(2 * time.Second)
	return &helloworld.HelloResponse{
		Message: "Hello",
	}, nil
}

// Greeter ...
type GreeterDebuglog struct {
	helloworld.UnimplementedGreeterServer
}

// SayHello ...
func (g GreeterDebuglog) SayHello(context context.Context, request *helloworld.HelloRequest) (*helloworld.HelloResponse, error) {
	return &helloworld.HelloResponse{
		Message: "Hello DebugLog",
	}, nil
}

// Greeter ...
type GreeterHeader struct {
	t *testing.T
	helloworld.UnimplementedGreeterServer
}

// SayHello ...
func (g GreeterHeader) SayHello(context context.Context, request *helloworld.HelloRequest) (*helloworld.HelloResponse, error) {
	appName := tools.GrpcHeaderValue(context, "app")
	// cpu := tools.GrpcHeaderValue(context, "enable-cpu-usage")
	// assert.Equal(g.t, "true", cpu)
	assert.Equal(g.t, "egrpc.test", appName)

	return &helloworld.HelloResponse{
		Message: "Hello",
	}, nil
}
