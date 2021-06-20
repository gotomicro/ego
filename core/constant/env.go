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
	// EgoTraceIDName 应用链路ID环境变量，不配置，默认x-trace-id
	EgoTraceIDName = "EGO_TRACE_ID_NAME"
	// EgoExtraTraceKeys 扩展追踪字段，通常用于追踪自定义业务字段。如用户ID(X-Ego-Uid)、订单ID(X-Ego-Order-Id)等。
	// 配置格式 {key1},{key2},{key3}...，多个 key 之间通过 "," 分割。
	// 比如 export EGO_EXTRA_TRACE_KEYS=X-Ego-Uid,X-Ego-Order-Id
	// 这些扩展的追踪字段会根据配置的 key1、key2、key3 等键名，通过 Headers(HTTP) 或 metadata(gRPC) 进行存取，并往下游传递
	EgoExtraTraceKeys = "EGO_EXTRA_TRACE_KEYS"
	// DefaultDeployment ...
	DefaultDeployment = ""
)
