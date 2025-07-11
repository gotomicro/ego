package eapp

import (
	"os"
	"strings"

	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/internal/ienv"
)

var (
	appMode                 string
	appRegion               string
	appZone                 string
	appInstance             string
	egoDebug                string
	egoLogPath              string
	egoLogAddApp            string
	egoTraceIDName          string
	egoLogExtraKeys         []string
	egoLogWriter            string
	egoGovernorEnableConfig string
	egoLogTimeType          string
	egoLogEnableAddCaller   bool
	egoHeaderExpose         string
)

func initEnv() {
	appMode = os.Getenv(constant.EnvAppMode)
	appRegion = os.Getenv(constant.EnvAppRegion)
	appZone = os.Getenv(constant.EnvAppZone)
	appInstance = ienv.EnvOrStr(constant.EnvAppInstance, HostName())
	egoDebug = os.Getenv(constant.EgoDebug)
	egoLogPath = os.Getenv(constant.EgoLogPath)
	egoLogAddApp = os.Getenv(constant.EgoLogAddApp)
	egoTraceIDName = ienv.EnvOrStr(constant.EgoTraceIDName, "x-trace-id")
	egoGovernorEnableConfig = os.Getenv(constant.EgoGovernorEnableConfig)
	if envEgoLogExtraKeys := strings.TrimSpace(os.Getenv(constant.EgoLogExtraKeys)); envEgoLogExtraKeys != "" {
		egoLogExtraKeys = strings.Split(envEgoLogExtraKeys, ",")
	}
	egoLogWriter = ienv.EnvOrStr(constant.EgoLogWriter, "file")
	egoLogTimeType = ienv.EnvOrStr(constant.EgoLogTimeType, "second")
	if IsDevelopmentMode() {
		egoLogTimeType = "%Y-%m-%d %H:%M:%S"
	}
	egoLogEnableAddCaller = ienv.EnvOrBool(constant.EgoLogEnableAddCaller, false)
	egoHeaderExpose = ienv.EnvOrStr(constant.EgoHeaderExpose, "x-expose")
}

// AppMode returns application running mode.
func AppMode() string {
	return appMode
}

// AppRegion returns application running region.
func AppRegion() string {
	return appRegion
}

// AppZone returns application running zone.
func AppZone() string {
	return appZone
}

// AppInstance returns application instance ID.
func AppInstance() string {
	return appInstance
}

// IsDevelopmentMode returns flag if application is in debug mode.
func IsDevelopmentMode() bool {
	return egoDebug == "true"
}

// EgoLogPath returns application log file directory path when user choose to write log fo file.
func EgoLogPath() string {
	return egoLogPath
}

// EnableLoggerAddApp returns flag if logger has append app Field to log entry.
func EnableLoggerAddApp() bool {
	return egoLogAddApp == "true"
}

// EgoTraceIDName returns the key in Metadata for storing traceID
func EgoTraceIDName() string {
	return egoTraceIDName
}

// EgoLogExtraKeys returns custom extra log keys.
func EgoLogExtraKeys() []string {
	return egoLogExtraKeys
}

// EgoLogWriter ...
func EgoLogWriter() string {
	return egoLogWriter
}

// EgoGovernorEnableConfig ...
func EgoGovernorEnableConfig() bool {
	return egoGovernorEnableConfig == "true"
}

// EgoLogTimeType ...
func EgoLogTimeType() string {
	return egoLogTimeType
}

// SetEgoDebug returns the flag if debug mode has been triggered
func SetEgoDebug(flag string) {
	egoDebug = flag
}

// EgoLogEnableAddCaller ...
func EgoLogEnableAddCaller() bool {
	return egoLogEnableAddCaller
}

func EgoHeaderExpose() string {
	return egoHeaderExpose
}
