package egovernor

import (
	"encoding/json"
	"net/http"
	"net/http/pprof"
	"runtime/debug"
)

var (
	// DefaultServeMux ...
	DefaultServeMux = http.NewServeMux()
	routes          = []string{}
)

func init() {
	// 获取全部治理路由
	HandleFunc("/routes", func(resp http.ResponseWriter, req *http.Request) {
		json.NewEncoder(resp).Encode(routes)
	})

	HandleFunc("/debug/pprof/", pprof.Index)
	HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	HandleFunc("/debug/pprof/profile", pprof.Profile)
	HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	HandleFunc("/debug/pprof/trace", pprof.Trace)

	if info, ok := debug.ReadBuildInfo(); ok {
		HandleFunc("/modInfo", func(w http.ResponseWriter, r *http.Request) {
			encoder := json.NewEncoder(w)
			if r.URL.Query().Get("pretty") == "true" {
				encoder.SetIndent("", "    ")
			}
			_ = encoder.Encode(info)
		})
	}
}

// HandleFunc ...
func HandleFunc(pattern string, handler http.HandlerFunc) {
	// todo: 增加安全管控
	DefaultServeMux.HandleFunc(pattern, handler)
	routes = append(routes, pattern)
}
