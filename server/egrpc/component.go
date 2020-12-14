// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package egrpc

import (
	"context"
	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

const PackageName = "server.egrpc"

// Component ...
type Component struct {
	name   string
	config *Config
	logger *elog.Component
	*grpc.Server
	listener   net.Listener
	serverInfo *server.ServiceInfo
	quit       chan error
}

func newComponent(name string, config *Config, logger *elog.Component) *Component {
	newServer := grpc.NewServer(config.serverOptions...)
	reflection.Register(newServer)

	return &Component{
		name:       name,
		config:     config,
		logger:     logger,
		Server:     newServer,
		listener:   nil,
		serverInfo: nil,
		quit:       make(chan error),
	}
}

func (s *Component) Name() string {
	return s.name
}

func (c *Component) PackageName() string {
	return PackageName
}

func (s *Component) Init() error {
	listener, err := net.Listen(s.config.Network, s.config.Address())
	if err != nil {
		s.logger.Panic("new grpc server err", elog.FieldErrKind("listen err"), elog.FieldErr(err))
	}
	s.config.Port = listener.Addr().(*net.TCPAddr).Port

	info := server.ApplyOptions(
		server.WithScheme("grpc"),
		server.WithAddress(listener.Addr().String()),
		server.WithKind(constant.ServiceProvider),
	)
	s.listener = listener
	s.serverInfo = &info
	return nil
}

// Component implements server.Component interface.
func (s *Component) Start() error {
	err := s.Server.Serve(s.listener)
	return err
}

// Stop implements server.Component interface
// it will terminate echo server immediately
func (s *Component) Stop() error {
	s.Server.Stop()
	return nil
}

// GracefulStop implements server.Component interface
// it will stop echo server gracefully
func (s *Component) GracefulStop(ctx context.Context) error {
	go func() {
		s.Server.GracefulStop()
		close(s.quit)
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-s.quit:
			return nil
		}
	}
	return nil
}

// Info returns server info, used by governor and consumer balancer
func (s *Component) Info() *server.ServiceInfo {
	return s.serverInfo
}

func (c *Component) Address() string {
	return c.config.Address()
}
