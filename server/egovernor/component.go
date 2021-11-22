package egovernor

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"runtime/debug"

	"github.com/felixge/fgprof"
	jsoniter "github.com/json-iterator/go"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server"
	"github.com/gotomicro/ego/task/ejob"
)

var (
	// DefaultServeMux ...
	DefaultServeMux = http.NewServeMux()
	routes          = []string{}
)

// PackageName 包名
const PackageName = "server.egovernor"

func init() {
	// 获取全部治理路由
	HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		_ = json.NewEncoder(resp).Encode(routes)
	})
	HandleFunc("/debug/pprof/", pprof.Index)
	HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	HandleFunc("/debug/pprof/profile", pprof.Profile)
	HandleFunc("/debug/pprof/fgprofile", fgprof.Handler().ServeHTTP)
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

	// 调试模式开启配置输出
	if eapp.IsDevelopmentMode() {
		HandleFunc("/config/json", func(w http.ResponseWriter, r *http.Request) {
			encoder := json.NewEncoder(w)
			if r.URL.Query().Get("pretty") == "true" {
				encoder.SetIndent("", "    ")
			}
			_ = encoder.Encode(econf.Traverse("."))
		})

		HandleFunc("/config/raw", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(econf.RawConfig())
		})
	}

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
	HandleFunc("/jobs", ejob.Handle)
	HandleFunc("/job/list", ejob.HandleJobList)
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

// Name 配置名称
func (c *Component) Name() string {
	return c.name
}

// PackageName 包名
func (c *Component) PackageName() string {
	return PackageName
}

// Init 初始化
func (c *Component) Init() error {
	var listener, err = net.Listen("tcp4", c.config.Address())
	if err != nil {
		elog.Panic("governor start error", elog.FieldErr(err))
	}
	c.listener = listener
	return nil
}

// Start 开始
func (c *Component) Start() error {
	HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})
	err := c.Server.Serve(c.listener)
	if err == http.ErrServerClosed {
		return nil
	}
	return err

}

//Stop ..
func (c *Component) Stop() error {
	return c.Server.Close()
}

//GracefulStop ..
func (c *Component) GracefulStop(ctx context.Context) error {
	return c.Server.Shutdown(ctx)
}

//Info ..
func (c *Component) Info() *server.ServiceInfo {
	info := server.ApplyOptions(
		server.WithScheme("http"),
		server.WithAddress(c.listener.Addr().String()),
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
