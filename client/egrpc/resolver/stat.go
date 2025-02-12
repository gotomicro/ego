package resolver

import (
	"net/http"
	"sync"

	"github.com/gotomicro/ego/server/egovernor"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/grpc/resolver"
)

var instances = sync.Map{}

func init() {
	egovernor.HandleFunc("/debug/client/egrpc/stats", func(w http.ResponseWriter, r *http.Request) {
		_ = jsoniter.NewEncoder(w).Encode(stats())
	})

}

// stats
func stats() (stats map[string][]resolver.Address) {
	stats = make(map[string][]resolver.Address)
	instances.Range(func(key, val interface{}) bool {
		name := key.(string)
		addresses := val.([]resolver.Address)
		stats[name] = addresses
		return true
	})
	return
}
