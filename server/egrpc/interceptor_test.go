package egrpc

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

// func Test_enableCPUUsage(t *testing.T) {
//	md := metadata.New(map[string]string{
//		"enable-cpu-usage": "true",
//	})
//	ctx := metadata.NewIncomingContext(context.Background(), md)
//	value := enableCPUUsage(ctx)
//	assert.Equal(t, true, value)
//
//	ctx2 := metadata.NewIncomingContext(context.Background(), nil)
//	value2 := enableCPUUsage(ctx2)
//	assert.Equal(t, false, value2)
//
//	md3 := metadata.New(map[string]string{
//		"enable-cpu-usage": "test",
//	})
//	ctx3 := metadata.NewIncomingContext(context.Background(), md3)
//	value3 := enableCPUUsage(ctx3)
//	assert.Equal(t, false, value3)
// }

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
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return cmp.Listener().(*bufconn.Listener).Dial()
		}))
	assert.Nil(t, err)
	cli := helloworld.NewGreeterClient(client)
	_, err = cli.SayHello(context.Background(), &helloworld.HelloRequest{})
	assert.Nil(t, err)
	logged, err := os.ReadFile(path.Join(logger.ConfigDir(), logger.ConfigName()))
	assert.Nil(t, err)
	assert.Contains(t, string(logged), "/helloworld.Greeter/SayHello")
	_ = os.Remove(path.Join(logger.ConfigDir(), logger.ConfigName()))
}

func Test_ServerPanicAccessLogger(t *testing.T) {
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
	helloworld.RegisterGreeterServer(cmp.Server, &PanicGreeter{})
	_ = cmp.Init()
	go func() {
		_ = cmp.Start()
	}()

	client, err := grpc.Dial("",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return cmp.Listener().(*bufconn.Listener).Dial()
		}))
	assert.Nil(t, err)
	cli := helloworld.NewGreeterClient(client)
	_, err = cli.SayHello(context.Background(), &helloworld.HelloRequest{})
	assert.Contains(t, err.Error(), "panic recover, origin err: we have a panic")

	logged, err := os.ReadFile(path.Join(logger.ConfigDir(), logger.ConfigName()))
	fmt.Printf("logged--------------->"+"%+v\n", string(logged))
	assert.Nil(t, err)
	assert.Contains(t, string(logged), "500")
	assert.Contains(t, string(logged), `"code":13`)
	_ = os.Remove(path.Join(logger.ConfigDir(), logger.ConfigName()))
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
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return cmp.Listener().(*bufconn.Listener).Dial()
		}))
	assert.Nil(t, err)
	cli := helloworld.NewGreeterClient(client)
	ctx := metadata.AppendToOutgoingContext(context.Background(), "app", "ego")
	_, err = cli.SayHello(ctx, &helloworld.HelloRequest{})
	assert.Nil(t, err)
	logged, err := os.ReadFile(path.Join(logger.ConfigDir(), logger.ConfigName()))
	assert.Nil(t, err)
	assert.Contains(t, string(logged), `"peerName":"ego"`)
	_ = os.Remove(path.Join(logger.ConfigDir(), logger.ConfigName()))
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
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return cmp.Listener().(*bufconn.Listener).Dial()
		}))
	assert.Nil(t, err)
	cli := helloworld.NewGreeterClient(client)
	ctx := metadata.AppendToOutgoingContext(context.Background(), "app", "ego")
	_, err = cli.SayHello(ctx, &helloworld.HelloRequest{})
	assert.Nil(t, err)
	logged, err := os.ReadFile(path.Join(logger.ConfigDir(), logger.ConfigName()))
	assert.Nil(t, err)
	assert.Contains(t, string(logged), `"peerName":"ego"`)
	_ = os.Remove(path.Join(logger.ConfigDir(), logger.ConfigName()))

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
	assert.Contains(t, string(text), `ego_server_handle_seconds_count{method="/helloworld.Greeter/SayHello",peer="ego",rpc_service="helloworld.Greeter",type="unary"}`)
	assert.Contains(t, string(text), `ego_server_handle_total{code="OK",method="/helloworld.Greeter/SayHello",peer="ego",rpc_service="helloworld.Greeter",type="unary",uniform_code="200"}`)
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

// Greeter ...
type PanicGreeter struct {
	helloworld.UnimplementedGreeterServer
}

// SayHello ...
func (g PanicGreeter) SayHello(context context.Context, request *helloworld.HelloRequest) (*helloworld.HelloResponse, error) {
	panic("we have a panic")
}

func TestCtxStoreSet(t *testing.T) {
	ctx := context.Background()
	Ctx := context.WithValue(ctx, ctxStoreStruct{}, &ctxStore{kvs: make(map[string]interface{})})
	CtxStore := Ctx.Value(ctxStoreStruct{}).(*ctxStore)
	assert.NotNil(t, CtxStore)
	CtxStoreSet(Ctx, "hello", "world")
	CtxStore.mu.RLock()
	storeValue, out := CtxStore.kvs["hello"]
	CtxStore.mu.RUnlock()
	assert.Equal(t, "world", storeValue)
	assert.True(t, true, out)
}
