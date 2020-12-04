package egovernor

import (
	"context"
	"encoding/json"
	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server"
	jsoniter "github.com/json-iterator/go"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"runtime/debug"
)

var (
	// DefaultServeMux ...
	DefaultServeMux = http.NewServeMux()
	routes          = []string{}
)

const PackageName = "server.egin"

func init() {
	// 获取全部治理路由
	HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		json.NewEncoder(resp).Encode(routes)
	})
	HandleFunc("/debug/pprof/", pprof.Index)
	HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	HandleFunc("/debug/pprof/profile", pprof.Profile)
	HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	HandleFunc("/debug/pprof/trace", pprof.Trace)
	if info, ok := debug.ReadBuildInfo(); ok {
		HandleFunc("/module/info", func(w http.ResponseWriter, r *http.Request) {
			encoder := json.NewEncoder(w)
			if r.URL.Query().Get("pretty") == "true" {
				encoder.SetIndent("", "    ")
			}
			_ = encoder.Encode(info)
		})
	}
	HandleFunc("/config/json", func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		if r.URL.Query().Get("pretty") == "true" {
			encoder.SetIndent("", "    ")
		}
		encoder.Encode(econf.Traverse("."))
	})
	HandleFunc("/config/raw", func(w http.ResponseWriter, r *http.Request) {
		w.Write(econf.RawConfig())
	})
	HandleFunc("/env/info", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_ = jsoniter.NewEncoder(w).Encode(os.Environ())
	})
	HandleFunc("/build/info", func(w http.ResponseWriter, r *http.Request) {
		serverStats := map[string]string{
			"name":       eapp.Name(),
			"appMode":    eapp.AppMode(),
			"appVersion": eapp.AppVersion(),
			"egoVersion": eapp.EgoVersion(),
			"buildUser":  eapp.BuildUser(),
			"buildHost":  eapp.BuildHost(),
			"buildTime":  eapp.BuildTime(),
			"startTime":  eapp.StartTime(),
			"hostName":   eapp.HostName(),
			"goVersion":  eapp.GoVersion(),
		}
		_ = jsoniter.NewEncoder(w).Encode(serverStats)
	})
}

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

// HandleFunc ...
func HandleFunc(pattern string, handler http.HandlerFunc) {
	// todo: 增加安全管控
	DefaultServeMux.HandleFunc(pattern, handler)
	routes = append(routes, pattern)
}
