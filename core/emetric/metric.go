package emetric

import (
	"time"

	"github.com/gotomicro/ego/core/eapp"
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
	// TypeGorm ...
	TypeGorm = "gorm"
	// TypeWebsocket ...
	TypeWebsocket = "ws"
	// TypeMySQL ...
	TypeMySQL = "mysql"
	// DefaultNamespace ...
	DefaultNamespace = "ego"
	// Conn 连接信息
	Conn = "conn"
)

var (
	// ServerHandleCounter ...
	ServerHandleCounter = CounterVecOpts{
		Namespace: DefaultNamespace,
		Name:      "server_handle_total",
		Labels:    []string{"type", "method", "peer", "code", "uniform_code", "rpc_service"},
	}.Build()

	// ServerStartedCounter ...
	ServerStartedCounter = CounterVecOpts{
		Namespace: DefaultNamespace,
		Name:      "server_started_total",
		Labels:    []string{"type", "method", "peer", "rpc_service"},
	}.Build()

	// ServerHandleHistogram ...
	ServerHandleHistogram = HistogramVecOpts{
		Namespace: DefaultNamespace,
		Name:      "server_handle_seconds",
		Labels:    []string{"type", "method", "peer", "rpc_service"},
	}.Build()

	// ClientHandleCounter ...
	ClientHandleCounter = CounterVecOpts{
		Namespace: DefaultNamespace,
		Name:      "client_handle_total",
		Labels:    []string{"type", "name", "method", "peer", "code"},
	}.Build()

	// ClientStartedCounter ...
	ClientStartedCounter = CounterVecOpts{
		Namespace: DefaultNamespace,
		Name:      "client_started_total",
		Labels:    []string{"type", "name", "method", "peer"},
	}.Build()

	// ClientHandleHistogram ...
	ClientHandleHistogram = HistogramVecOpts{
		Namespace: DefaultNamespace,
		Name:      "client_handle_seconds",
		Labels:    []string{"type", "name", "method", "peer"},
	}.Build()

	// ClientStatsGauge ...
	ClientStatsGauge = GaugeVecOpts{
		Namespace: DefaultNamespace,
		Name:      "client_stats_gauge",
		Labels:    []string{"type", "name", "index"},
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

	// LibHandleHistogram ...
	// Deprecated LibHandleHistogram
	LibHandleHistogram = HistogramVecOpts{
		Namespace: DefaultNamespace,
		Name:      "lib_handle_seconds",
		Labels:    []string{"type", "method", "address"},
	}.Build()

	// LibHandleCounter ...
	// Deprecated LibHandleCounter
	LibHandleCounter = CounterVecOpts{
		Namespace: DefaultNamespace,
		Name:      "lib_handle_total",
		Labels:    []string{"type", "method", "address", "code"},
	}.Build()

	// LibHandleSummary ...
	// Deprecated LibHandleSummary
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

	// ConnGauge ...
	ConnGauge = GaugeVecOpts{
		Namespace: DefaultNamespace,
		Name:      "connection_states",
		Labels:    []string{"state"},
	}.Build()
)

func init() {
	BuildInfoGauge.WithLabelValues(
		eapp.Name(),
		eapp.AppMode(),
		eapp.AppRegion(),
		eapp.AppZone(),
		eapp.AppVersion(),
		eapp.EgoVersion(),
		eapp.StartTime(),
		eapp.BuildTime(),
		eapp.GoVersion(),
	).Set(float64(time.Now().UnixNano() / 1e6))
}
