package constant

const (
	// EnvAppName is ego application name.
	EnvAppName = "EGO_NAME"
	// EnvAppMode is ego application running mode.
	EnvAppMode = "EGO_MODE"
	// EnvAppRegion is ego application instance running region
	EnvAppRegion = "EGO_REGION"
	// EnvAppZone is ego application instance running zone
	EnvAppZone = "EGO_ZONE"
	// EnvAppHost is ego application instance host name
	EnvAppHost = "EGO_HOST"
	// EnvAppInstance is ego application unique instance ID
	EnvAppInstance = "EGO_INSTANCE"
	// EgoDebug means if turn on debug mode or not, when set true, verbose log will be print in terminal.
	EgoDebug = "EGO_DEBUG"
	// EgoConfigPath chooses what configuration path will be used to start application.
	EgoConfigPath = "EGO_CONFIG_PATH"
	// EgoLogPath if user config log writer as fileWriter, EGO_LOG_PATH means directory path of log file.
	EgoLogPath = "EGO_LOG_PATH"
	// EgoLogAddApp if set true, all log entries will append application name filed.
	// EgoLogAddApp 应用日志增加应用名环境变量，如果增加该环境变量，日志里会将应用名写入到app字段里
	EgoLogAddApp = "EGO_LOG_ADD_APP"
	// EgoLogExtraKeys 扩展追踪字段，通常用于打印自定义Headers/Metadata。如用户ID(X-Ego-Uid)、订单ID(X-Ego-Order-Id)等。
	// 配置格式 {key1},{key2},{key3}...，多个 key 之间通过 "," 分割。
	// 比如 export EGO_LOG_EXTRA_KEYS=X-Ego-Uid,X-Ego-Order-Id
	// 这些扩展的追踪字段会根据配置的 key1、key2、key3 等键名，从 Headers(HTTP) 或 Metadata(gRPC) 查找对应值并打印到请求日志中
	EgoLogExtraKeys = "EGO_LOG_EXTRA_KEYS"
	// EgoLogWriter writer方式： file | stderr
	EgoLogWriter = "EGO_LOG_WRITER"
	// EgoLogTimeType 记录的时间类型，默认 second，millisecond，%Y-%m-%d %H:%M:%S
	EgoLogTimeType = "EGO_LOG_TIME_TYPE"
	// EgoTraceIDName 应用链路ID环境变量，不配置，默认x-trace-id
	EgoTraceIDName = "EGO_TRACE_ID_NAME"
	// EgoGovernorEnableConfig 是否开启查看config
	EgoGovernorEnableConfig = "EGO_GOVERNOR_ENABLE_CONFIG"
)
