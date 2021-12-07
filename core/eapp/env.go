package eapp

import (
	"os"
	"strings"

	"github.com/gotomicro/ego/core/constant"
)

var (
	appMode         string
	appRegion       string
	appZone         string
	appInstance     string // 通常是实例的机器名
	egoDebug        string
	egoLogPath      string
	egoLogAddApp    string
	egoTraceIDName  string
	egoLogExtraKeys []string
	egoLogWriter    string
)

func initEnv() {
	appMode = os.Getenv(constant.EnvAppMode)
	appRegion = os.Getenv(constant.EnvAppRegion)
	appZone = os.Getenv(constant.EnvAppZone)
	appInstance = os.Getenv(constant.EnvAppInstance)
	if appInstance == "" {
		appInstance = HostName()
	}
	egoDebug = os.Getenv(constant.EgoDebug)
	egoLogPath = os.Getenv(constant.EgoLogPath)
	egoLogAddApp = os.Getenv(constant.EgoLogAddApp)
	egoTraceIDName = os.Getenv(constant.EgoTraceIDName)
	if egoTraceIDName == "" {
		egoTraceIDName = "x-trace-id"
	}
	egoLogExtraKeys = strings.Split(os.Getenv(constant.EgoLogExtraKeys), ",")
	egoLogWriter = os.Getenv(constant.EgoLogWriter)
	if egoLogWriter == "" {
		egoLogWriter = "file"
	}
}

// AppMode 获取应用运行的环境
func AppMode() string {
	return appMode
}

// AppRegion 获取APP运行的地区
func AppRegion() string {
	return appRegion
}

// AppZone 获取应用运行的可用区
func AppZone() string {
	return appZone
}

// AppInstance 获取应用实例，通常是实例的机器名
func AppInstance() string {
	return appInstance
}

// IsDevelopmentMode 判断是否是生产模式
func IsDevelopmentMode() bool {
	return egoDebug == "true"
}

// EgoLogPath 获取应用日志路径
func EgoLogPath() string {
	return egoLogPath
}

// EnableLoggerAddApp 日志是否记录应用名信息
func EnableLoggerAddApp() bool {
	return egoLogAddApp == "true"
}

// EgoTraceIDName 获取链路名称
func EgoTraceIDName() string {
	return egoTraceIDName
}

// EgoLogExtraKeys 获取自定义loggerKeys
func EgoLogExtraKeys() []string {
	return egoLogExtraKeys
}

func EgoLogWriter() string {
	return egoLogWriter
}
