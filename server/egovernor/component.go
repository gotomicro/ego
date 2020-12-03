package egovernor

import (
	"context"
	"net"
	"net/http"

	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server"
)

const PackageName = "server.egin"

// Component ...
type Component struct {
	name   string
	config *Config
	logger *elog.Component
	*http.Server
	listener net.Listener
}

func newComponent(name string, config *Config, logger *elog.Component) *Component {
	return &Component{
		name:   name,
		logger: logger,
		Server: &http.Server{
			Addr:    config.Address(),
			Handler: DefaultServeMux,
		},
		listener: nil,
		config:   config,
	}
}

func (c *Component) Name() string {
	return c.name
}

func (c *Component) PackageName() string {
	return PackageName
}

func (c *Component) Init() error {
	var listener, err = net.Listen("tcp4", c.config.Address())
	if err != nil {
		elog.Panic("governor start error", elog.FieldErr(err))
	}
	c.listener = listener
	return nil
}

//Serve ..
func (s *Component) Start() error {
	err := s.Server.Serve(s.listener)
	if err == http.ErrServerClosed {
		return nil
	}
	return err

}

//Stop ..
func (s *Component) Stop() error {
	return s.Server.Close()
}

//GracefulStop ..
func (s *Component) GracefulStop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

//Info ..
func (s *Component) Info() *server.ServiceInfo {
	info := server.ApplyOptions(
		server.WithScheme("http"),
		server.WithAddress(s.listener.Addr().String()),
		server.WithKind(constant.ServiceGovernor),
	)
	// info.Name = info.Name + "." + ModName
	return &info
}
