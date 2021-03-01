package eapp

import (
	"os"

	"github.com/gotomicro/ego/core/constant"
)

var (
	appMode        string
	appRegion      string
	appZone        string
	appHost        string // 应用的ip
	appInstance    string // 通常是实例的机器名
	egoDebug       string
	egoConfigPath  string
	egoLogPath     string
	egoLogAddApp   string
	egoTraceIDName string
)

func InitEnv() {
	appMode = os.Getenv(constant.EnvAppMode)
	appRegion = os.Getenv(constant.EnvAppRegion)
	appZone = os.Getenv(constant.EnvAppZone)
	appHost = os.Getenv(constant.EnvAppHost)
	appInstance = os.Getenv(constant.EnvAppInstance)
	if appInstance == "" {
		appInstance = HostName()
	}
	egoDebug = os.Getenv(constant.EgoDebug)
	egoConfigPath = os.Getenv(constant.EgoConfigPath)
	if egoConfigPath == "" {
		egoConfigPath = "config/local.toml"
	}
	egoLogPath = os.Getenv(constant.EgoLogPath)
	egoLogAddApp = os.Getenv(constant.EgoLogAddApp)
	egoTraceIDName = os.Getenv(constant.EgoTraceIDName)
	if egoTraceIDName == "" {
		egoTraceIDName = "x-trace-id"
	}
}

func AppMode() string {
	return appMode
}

func AppRegion() string {
	return appRegion
}

func AppZone() string {
	return appZone
}

func AppHost() string {
	return appHost
}

func AppInstance() string {
	return appInstance
}

// IsDevelopmentMode 判断是否是生产模式
func IsDevelopmentMode() bool {
	return egoDebug == "true"
}

func EgoConfigPath() string {
	return egoConfigPath
}

func EgoLogPath() string {
	return egoLogPath
}

func EnableLoggerAddApp() bool {
	return egoLogAddApp == "true"
}

func EgoTraceIDName() string {
	return egoTraceIDName
}
