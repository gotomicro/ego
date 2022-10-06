package egrpc

import (
	"context"
	"io/ioutil"
	"net"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"

	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/internal/test/helloworld"
)

func Test_getPeerName(t *testing.T) {
	md := metadata.New(map[string]string{
		"app": "ego-svc",
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	value := getPeerName(ctx)
	assert.Equal(t, "ego-svc", value)

	ctx2 := metadata.NewIncomingContext(context.Background(), nil)
	value2 := getPeerName(ctx2)
	assert.Equal(t, "", value2)
}

// todo add more unittest
func Test_getPeerIP(t *testing.T) {
	md := metadata.New(map[string]string{
		"client-ip": "127.0.0.1",
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	value := getPeerIP(ctx)
	assert.Equal(t, "127.0.0.1", value)
}

func Test_enableCPUUsage(t *testing.T) {
	md := metadata.New(map[string]string{
		"enable-cpu-usage": "true",
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	value := enableCPUUsage(ctx)
	assert.Equal(t, true, value)

	ctx2 := metadata.NewIncomingContext(context.Background(), nil)
	value2 := enableCPUUsage(ctx2)
	assert.Equal(t, false, value2)

	md3 := metadata.New(map[string]string{
		"enable-cpu-usage": "test",
	})
	ctx3 := metadata.NewIncomingContext(context.Background(), md3)
	value3 := enableCPUUsage(ctx3)
	assert.Equal(t, false, value3)
}

func Test_ServerAccessLogger(t *testing.T) {
	// 使用非异步日志
	logger := elog.DefaultContainer().Build(
		elog.WithDebug(false),
		elog.WithEnableAddCaller(true),
		elog.WithEnableAsync(false),
	)
	cmp := DefaultContainer().Build(
		WithNetwork("bufnet"),
		WithLogger(logger),
	)
	helloworld.RegisterGreeterServer(cmp.Server, &Greeter{})
	_ = cmp.Init()
	go func() {
		_ = cmp.Start()
	}()

	client, err := grpc.Dial("",
		grpc.WithInsecure(),
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return cmp.Listener().(*bufconn.Listener).Dial()
		}))
	assert.Nil(t, err)
	cli := helloworld.NewGreeterClient(client)
	_, err = cli.SayHello(context.Background(), &helloworld.HelloRequest{})
	assert.Nil(t, err)
	logged, err := ioutil.ReadFile(path.Join(logger.ConfigDir(), logger.ConfigName()))
	assert.Nil(t, err)
	assert.Contains(t, string(logged), "/helloworld.Greeter/SayHello")
	os.Remove(path.Join(logger.ConfigDir(), logger.ConfigName()))
}

func Test_ServerAccessAppName(t *testing.T) {
	// 使用非异步日志
	logger := elog.DefaultContainer().Build(
		elog.WithDebug(false),
		elog.WithEnableAddCaller(true),
		elog.WithEnableAsync(false),
	)
	cmp := DefaultContainer().Build(
		WithNetwork("bufnet"),
		WithLogger(logger),
	)
	helloworld.RegisterGreeterServer(cmp.Server, &Greeter{})
	_ = cmp.Init()
	go func() {
		_ = cmp.Start()
	}()

	client, err := grpc.Dial("",
		grpc.WithInsecure(),
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return cmp.Listener().(*bufconn.Listener).Dial()
		}))
	assert.Nil(t, err)
	cli := helloworld.NewGreeterClient(client)
	ctx := metadata.AppendToOutgoingContext(context.Background(), "app", "ego")
	_, err = cli.SayHello(ctx, &helloworld.HelloRequest{})
	assert.Nil(t, err)
	logged, err := ioutil.ReadFile(path.Join(logger.ConfigDir(), logger.ConfigName()))
	assert.Nil(t, err)
	assert.Contains(t, string(logged), `"peerName":"ego"`)
	os.Remove(path.Join(logger.ConfigDir(), logger.ConfigName()))
}

func TestPrometheus(t *testing.T) {
	// 1 获取prometheus的handler的数据
	ts := httptest.NewServer(promhttp.Handler())
	defer ts.Close()

	// 使用非异步日志
	logger := elog.DefaultContainer().Build(
		elog.WithDebug(false),
		elog.WithEnableAddCaller(true),
		elog.WithEnableAsync(false),
	)
	cmp := DefaultContainer().Build(
		WithNetwork("bufnet"),
		WithLogger(logger),
	)
	helloworld.RegisterGreeterServer(cmp.Server, &Greeter{})
	_ = cmp.Init()
	go func() {
		_ = cmp.Start()
	}()

	client, err := grpc.Dial("",
		grpc.WithInsecure(),
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return cmp.Listener().(*bufconn.Listener).Dial()
		}))
	assert.Nil(t, err)
	cli := helloworld.NewGreeterClient(client)
	ctx := metadata.AppendToOutgoingContext(context.Background(), "app", "ego")
	_, err = cli.SayHello(ctx, &helloworld.HelloRequest{})
	assert.Nil(t, err)
	logged, err := ioutil.ReadFile(path.Join(logger.ConfigDir(), logger.ConfigName()))
	assert.Nil(t, err)
	assert.Contains(t, string(logged), `"peerName":"ego"`)
	os.Remove(path.Join(logger.ConfigDir(), logger.ConfigName()))

	pc := ts.Client()
	res, err := pc.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	text, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	err = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, string(text), `ego_server_handle_seconds_count{method="/helloworld.Greeter/SayHello",peer="ego",rpc_service="helloworld.Greeter",type="unary"}`)
	assert.Contains(t, string(text), `ego_server_handle_total{code="OK",method="/helloworld.Greeter/SayHello",peer="ego",rpc_service="helloworld.Greeter",type="unary",uniform_code="OK"}`)
	assert.Contains(t, string(text), `ego_server_started_total{method="/helloworld.Greeter/SayHello",peer="ego",rpc_service="helloworld.Greeter",type="unary"}`)
}

// Greeter ...
type Greeter struct {
	helloworld.UnimplementedGreeterServer
}

// SayHello ...
func (g Greeter) SayHello(context context.Context, request *helloworld.HelloRequest) (*helloworld.HelloResponse, error) {
	return &helloworld.HelloResponse{
		Message: "Hello",
	}, nil
}
