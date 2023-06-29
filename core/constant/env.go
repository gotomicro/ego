package constant

const (
	// EnvAppName defines your application name.
	EnvAppName = "EGO_NAME"
	// EnvAppMode defines your application running mode, you can set EGO_MODE with words such as "development/testing/production".
	EnvAppMode = "EGO_MODE"
	// EnvAppRegion defines your application running region, such as "ASIA".
	EnvAppRegion = "EGO_REGION"
	// EnvAppZone defines your application running zone "ShangHai2".
	EnvAppZone = "EGO_ZONE"
	// EnvAppHost defines your application running HOST name.
	EnvAppHost = "EGO_HOST"
	// EnvAppInstance defines your application replication unique ID.
	EnvAppInstance = "EGO_INSTANCE"
	// EgoDebug when set to true, ego will print verbose logs.
	EgoDebug = "EGO_DEBUG"
	// EgoConfigPath defines your application configuration dsn.
	EgoConfigPath = "EGO_CONFIG_PATH"
	// EgoLogPath if your application used file writer logger to print logs, EgoLogPath is logs directory.
	EgoLogPath = "EGO_LOG_PATH"
	// EgoLogAddApp defines if we should append application name to every logger entries.
	EgoLogAddApp = "EGO_LOG_ADD_APP"
	// EgoLogExtraKeys used to append extra tracing keys to every access logger entries, the keys usually comes from HTTP Headers or gRPC Metadata.
	// you can trace you custom business clues, such as "X-Biz-Uid"(your application user ID) or "X-Biz-Order-Id"(your application order ID).
	// each keys separated with ",". For example, export EGO_LOG_EXTRA_KEYS=X-Ego-Uid,X-Ego-Order-Id
	EgoLogExtraKeys = "EGO_LOG_EXTRA_KEYS"
	// EgoLogWriter defines your log writer, available types are: "file/stderr"
	EgoLogWriter = "EGO_LOG_WRITER"
	// EgoLogTimeType defines time format on your logger entries, available types are "second/millisecond/%Y-%m-%d %H:%M:%S"
	EgoLogTimeType = "EGO_LOG_TIME_TYPE"
	// EgoTraceIDName defines tracing ID NAME, default value is "x-trace-id"
	EgoTraceIDName = "EGO_TRACE_ID_NAME"
	// EgoGovernorEnableConfig defines if you can query current configuration with governor APIs.
	EgoGovernorEnableConfig = "EGO_GOVERNOR_ENABLE_CONFIG"
	// EgoLogEnableAddCaller when set to true, your log will show caller, default value is false
	EgoLogEnableAddCaller = "EGO_LOG_ENABLE_ADD_CALLER"
	// EgoDefaultConfigExt defines default config file extension, support ".toml"，".yaml"，".json",
	// EgoDefaultConfigExt effective only the configuration file path without extension name
	EgoDefaultConfigExt = "EGO_DEFAULT_CONFIG_EXT"
	// EgoDeploymentEnv defines deployment environment, such as "k8s", "ecs"
	EgoDeploymentEnv = "EGO_DEPLOYMENT_ENV"
)
