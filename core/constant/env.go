package constant

const (
	// EnvAppName 应用名环境变量
	EnvAppName = "EGO_NAME"
	// EnvAppMode 应用模式环境变量
	EnvAppMode = "EGO_MODE"
	// EnvAppRegion ...
	EnvAppRegion = "EGO_REGION"
	// EnvAppZone ...
	EnvAppZone = "EGO_ZONE"
	// EnvAppHost ...
	EnvAppHost = "EGO_HOST"
	// EnvAppInstance 应用实例ID环境变量
	EnvAppInstance = "EGO_INSTANCE"
	// EgoDebug 调试环境变量，export EGO_DEBUG=true，开启应用的调试模式
	EgoDebug = "EGO_DEBUG"
	// EgoConfigPath 应用配置环境变量
	EgoConfigPath = "EGO_CONFIG_PATH"
	// EgoLogPath 应用日志环境变量
	EgoLogPath = "EGO_LOG_PATH"
	// EgoLogAddApp 应用日志增加应用名环境变量，如果增加该环境变量，日志里会将应用名写入到app字段里
	EgoLogAddApp = "EGO_LOG_ADD_APP"
	// EgoLogExtraKeys 扩展追踪字段，通常用于打印自定义Headers/Metadata。如用户ID(X-Ego-Uid)、订单ID(X-Ego-Order-Id)等。
	// 配置格式 {key1},{key2},{key3}...，多个 key 之间通过 "," 分割。
	// 比如 export EGO_LOG_EXTRA_KEYS=X-Ego-Uid,X-Ego-Order-Id
	// 这些扩展的追踪字段会根据配置的 key1、key2、key3 等键名，从 Headers(HTTP) 或 Metadata(gRPC) 查找对应值并打印到请求日志中
	EgoLogExtraKeys = "EGO_LOG_EXTRA_KEYS"
	// EgoLogWriter writer方式： file | stderr
	EgoLogWriter = "EGO_LOG_WRITER"
	// EgoTraceIDName 应用链路ID环境变量，不配置，默认x-trace-id
	EgoTraceIDName = "EGO_TRACE_ID_NAME"
	// DefaultDeployment ...
	DefaultDeployment = ""
)
