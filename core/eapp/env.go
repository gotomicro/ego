package eapp

import (
	"crypto/md5"
	"fmt"
	"os"

	"github.com/gotomicro/ego/core/constant"
)

var (
	appMode       string
	appRegion     string
	appZone       string
	appHost       string
	appInstance   string
	egoDebug      string
	egoConfigPath string
	egoLogPath    string
	egoLogAddApp  string
)

func InitEnv() {
	appMode = os.Getenv(constant.EnvAppMode)
	appRegion = os.Getenv(constant.EnvAppRegion)
	appZone = os.Getenv(constant.EnvAppZone)
	appHost = os.Getenv(constant.EnvAppHost)
	appInstance = os.Getenv(constant.EnvAppInstance)
	if appInstance == "" {
		appInstance = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s:%s", HostName(), Name()))))
	}
	egoDebug = os.Getenv(constant.EgoDebug)
	egoConfigPath = os.Getenv(constant.EgoConfigPath)
	if egoConfigPath == "" {
		egoConfigPath = "config/local.toml"
	}
	egoLogPath = os.Getenv(constant.EgoLogPath)
	egoLogAddApp = os.Getenv(constant.EgoLogAddApp)
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
