package egovernor

import (
	"encoding/json"
	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/econf"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"os"
)

func init() {
	HandleFunc("/configs", func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		if r.URL.Query().Get("pretty") == "true" {
			encoder.SetIndent("", "    ")
		}
		encoder.Encode(econf.Traverse("."))
	})

	HandleFunc("/debug/env", func(w http.ResponseWriter, r *http.Request) {
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
