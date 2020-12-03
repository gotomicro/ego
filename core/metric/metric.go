package metric

import (
	"github.com/gotomicro/ego/core/app"
	"github.com/gotomicro/ego/server/egovernor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var (
	// TypeHTTP ...
	TypeHTTP = "http"
	// TypeGRPCUnary ...
	TypeGRPCUnary = "unary"
	// TypeGRPCStream ...
	TypeGRPCStream = "stream"
	// TypeRedis ...
	TypeRedis = "redis"
	TypeGorm  = "gorm"
	// TypeRocketMQ ...
	TypeRocketMQ = "rocketmq"
	// TypeWebsocket ...
	TypeWebsocket = "ws"

	// TypeMySQL ...
	TypeMySQL = "mysql"

	// CodeJob
	CodeJobSuccess = "ok"
	// CodeJobFail ...
	CodeJobFail = "fail"
	// CodeJobReentry ...
	CodeJobReentry = "reentry"

	// CodeCache
	CodeCacheMiss = "miss"
	// CodeCacheHit ...
	CodeCacheHit = "hit"

	// Namespace
	DefaultNamespace = "ego"
)

var (
	// ServerHandleCounter ...
	ServerHandleCounter = CounterVecOpts{
		Namespace: DefaultNamespace,
		Name:      "server_handle_total",
		Labels:    []string{"type", "method", "peer", "code"},
	}.Build()

	// ServerHandleHistogram ...
	ServerHandleHistogram = HistogramVecOpts{
		Namespace: DefaultNamespace,
		Name:      "server_handle_seconds",
		Labels:    []string{"type", "method", "peer"},
	}.Build()

	// ClientHandleCounter ...
	ClientHandleCounter = CounterVecOpts{
		Namespace: DefaultNamespace,
		Name:      "client_handle_total",
		Labels:    []string{"type", "name", "method", "peer", "code"},
	}.Build()

	// ClientHandleHistogram ...
	ClientHandleHistogram = HistogramVecOpts{
		Namespace: DefaultNamespace,
		Name:      "client_handle_seconds",
		Labels:    []string{"type", "name", "method", "peer"},
	}.Build()

	// JobHandleCounter ...
	JobHandleCounter = CounterVecOpts{
		Namespace: DefaultNamespace,
		Name:      "job_handle_total",
		Labels:    []string{"type", "name", "code"},
	}.Build()

	// JobHandleHistogram ...
	JobHandleHistogram = HistogramVecOpts{
		Namespace: DefaultNamespace,
		Name:      "job_handle_seconds",
		Labels:    []string{"type", "name"},
	}.Build()

	LibHandleHistogram = HistogramVecOpts{
		Namespace: DefaultNamespace,
		Name:      "lib_handle_seconds",
		Labels:    []string{"type", "method", "address"},
	}.Build()
	// LibHandleCounter ...
	LibHandleCounter = CounterVecOpts{
		Namespace: DefaultNamespace,
		Name:      "lib_handle_total",
		Labels:    []string{"type", "method", "address", "code"},
	}.Build()

	LibHandleSummary = SummaryVecOpts{
		Namespace: DefaultNamespace,
		Name:      "lib_handle_stats",
		Labels:    []string{"name", "status"},
	}.Build()

	// CacheHandleCounter ...
	CacheHandleCounter = CounterVecOpts{
		Namespace: DefaultNamespace,
		Name:      "cache_handle_total",
		Labels:    []string{"type", "name", "action", "code"},
	}.Build()

	// CacheHandleHistogram ...
	CacheHandleHistogram = HistogramVecOpts{
		Namespace: DefaultNamespace,
		Name:      "cache_handle_seconds",
		Labels:    []string{"type", "name", "action"},
	}.Build()

	// BuildInfoGauge ...
	BuildInfoGauge = GaugeVecOpts{
		Namespace: DefaultNamespace,
		Name:      "build_info",
		Labels:    []string{"name", "mode", "region", "zone", "app_version", "ego_version", "start_time", "build_time", "go_version"},
	}.Build()
)

func init() {
	BuildInfoGauge.WithLabelValues(
		app.Name(),
		app.AppMode(),
		app.AppRegion(),
		app.AppZone(),
		app.AppVersion(),
		app.EgoVersion(),
		app.StartTime(),
		app.BuildTime(),
		app.GoVersion(),
	).Set(float64(time.Now().UnixNano() / 1e6))

	egovernor.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})
}
