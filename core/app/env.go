package app

import (
	"crypto/md5"
	"fmt"
	"os"

	"github.com/gotomicro/ego/core/constant"
)

var (
	appLogDir       string
	appMode         string
	appRegion       string
	appZone         string
	appHost         string
	appInstance     string
	egoDebug      string
	egoConfigPath string
)

func InitEnv() {
	appLogDir = os.Getenv(constant.EnvAppLogDir)
	appMode = os.Getenv(constant.EnvAppMode)
	appRegion = os.Getenv(constant.EnvAppRegion)
	appZone = os.Getenv(constant.EnvAppZone)
	appHost = os.Getenv(constant.EnvAppHost)
	appInstance = os.Getenv(constant.EnvAppInstance)
	if appInstance == "" {
		appInstance = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s:%s", HostName(), AppID()))))
	}
	egoDebug = os.Getenv(constant.EgoDebug)
	egoConfigPath = os.Getenv(constant.EgoConfigPath)
	if egoConfigPath == "" {
		egoConfigPath = "config/local.toml"
	}
}

func AppLogDir() string {
	return appLogDir
}

func SetAppLogDir(logDir string) {
	appLogDir = logDir
}

func AppMode() string {
	return appMode
}

func SetAppMode(mode string) {
	appMode = mode
}

func AppRegion() string {
	return appRegion
}

func SetAppRegion(region string) {
	appRegion = region
}

func AppZone() string {
	return appZone
}

func SetAppZone(zone string) {
	appZone = zone
}

func AppHost() string {
	return appHost
}

func SetAppHost(host string) {
	appHost = host
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
