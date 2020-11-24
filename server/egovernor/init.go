package egovernor

import (
	"encoding/json"
	"github.com/gotomicro/ego/core/app"
	"github.com/gotomicro/ego/core/conf"
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
		encoder.Encode(conf.Traverse("."))
	})

	HandleFunc("/debug/env", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_ = jsoniter.NewEncoder(w).Encode(os.Environ())
	})

	HandleFunc("/build/info", func(w http.ResponseWriter, r *http.Request) {
		serverStats := map[string]string{
			"name":         app.Name(),
			"appID":        app.AppID(),
			"appMode":      app.AppMode(),
			"appVersion":   app.AppVersion(),
			"egoVersion": app.EgoVersion(),
			"buildUser":    app.BuildUser(),
			"buildHost":    app.BuildHost(),
			"buildTime":    app.BuildTime(),
			"startTime":    app.StartTime(),
			"hostName":     app.HostName(),
			"goVersion":    app.GoVersion(),
		}
		_ = jsoniter.NewEncoder(w).Encode(serverStats)
	})
}
