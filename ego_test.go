package ego

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/server"
)

func TestEgoRun(t *testing.T) {
	t.Run("ego run start error", func(t *testing.T) {
		svc := &testServer{
			ServeErr: fmt.Errorf("when server call start error"),
		}
		app := New()
		app.Serve(svc)
		go func() {
			time.Sleep(time.Millisecond * 100)
			err := app.Stop(context.Background(), false, false)
			assert.Nil(t, err)
		}()
		err := app.Run()
		assert.EqualError(t, err, "when server call start error")
	})

	t.Run("ego run stop error", func(t *testing.T) {
		svc := &testServer{
			StopErr: fmt.Errorf("when server call stop error"),
		}
		app := New()
		app.Serve(svc)
		go func() {
			time.Sleep(time.Millisecond * 100)
			err := app.Stop(context.Background(), false, false)
			assert.Nil(t, err)
		}()
		err := app.Run()
		assert.EqualError(t, err, "when server call stop error")
	})
}
func TestEgoNew(t *testing.T) {
	app := New()
	assert.NotNil(t, app.logger)
	assert.NotNil(t, app.servers)
	assert.NotNil(t, app.jobs)
	assert.NotNil(t, app.logger)
}

func TestEgoEconfRace(t *testing.T) {
	New()
	go func() {
		econf.OnChange(func(c *econf.Configuration) {})
	}()
	econf.OnChange(func(c *econf.Configuration) {})

}

type testServer struct {
	ServeBlockTime time.Duration
	ServeErr       error

	StopBlockTime time.Duration
	StopErr       error

	GstopBlockTime time.Duration
	GstopErr       error
}

func (s *testServer) Name() string {
	return "test_server"
}

func (s *testServer) PackageName() string {
	return "server"
}

func (s *testServer) Init() error {
	time.Sleep(s.ServeBlockTime)
	return s.ServeErr
}

func (s *testServer) Start() error {
	time.Sleep(s.ServeBlockTime)
	return s.ServeErr
}

func (s *testServer) Stop() error {
	time.Sleep(s.StopBlockTime)
	return s.StopErr
}
func (s *testServer) GracefulStop(ctx context.Context) error {
	time.Sleep(s.GstopBlockTime)
	return s.GstopErr
}
func (s *testServer) Info() *server.ServiceInfo {
	return &server.ServiceInfo{}
}
