package eapp

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/util/xcolor"
	"github.com/gotomicro/ego/core/util/xtime"
)

var (
	startTime  string
	goVersion  string
	egoVersion string
)

// build info
var (
	appName         string
	hostName        string
	buildAppVersion string
	buildUser       string
	buildHost       string
	buildStatus     string
	buildTime       string
)

func init() {
	if appName == "" {
		appName = os.Getenv(constant.EnvAppName)
		if appName == "" {
			appName = filepath.Base(os.Args[0])
		}
	}

	name, err := os.Hostname()
	if err != nil {
		name = "unknown"
	}
	hostName = name
	startTime = xtime.TS.Format(time.Now())
	setBuildTime(buildTime)
	goVersion = runtime.Version()
	initEnv()

	// ego version
	egoVersion = "unknown version"
	info, ok := debug.ReadBuildInfo()
	if ok {
		for _, value := range info.Deps {
			if value.Path == "github.com/gotomicro/ego" {
				egoVersion = value.Version
			}
		}
	}
}

// Name gets application name.
func Name() string {
	return appName
}

// AppVersion get buildAppVersion
func AppVersion() string {
	return buildAppVersion
}

// EgoVersion get egoVersion
func EgoVersion() string {
	return egoVersion
}

// BuildTime get buildTime
func BuildTime() string {
	return buildTime
}

// BuildUser get buildUser
func BuildUser() string {
	return buildUser
}

// BuildHost get buildHost
func BuildHost() string {
	return buildHost
}

// setBuildTime set buildTime
func setBuildTime(param string) {
	buildTime = strings.Replace(param, "--", " ", 1)
}

// HostName get host name
func HostName() string {
	return hostName
}

// StartTime get start time
func StartTime() string {
	return startTime
}

// GoVersion get go version
func GoVersion() string {
	return goVersion
}

// PrintVersion print formated version info
func PrintVersion() {
	fmt.Printf("%-20s : %s\n", xcolor.Green("EGO"), xcolor.Blue("I am EGO"))
	fmt.Printf("%-20s : %s\n", xcolor.Green("AppName"), xcolor.Blue(appName))
	fmt.Printf("%-20s : %s\n", xcolor.Green("AppHost"), xcolor.Blue(HostName()))
	fmt.Printf("%-20s : %s\n", xcolor.Green("Region"), xcolor.Blue(AppRegion()))
	fmt.Printf("%-20s : %s\n", xcolor.Green("Zone"), xcolor.Blue(AppZone()))
	fmt.Printf("%-20s : %s\n", xcolor.Green("AppVersion"), xcolor.Blue(buildAppVersion))
	fmt.Printf("%-20s : %s\n", xcolor.Green("EgoVersion"), xcolor.Blue(egoVersion))
	fmt.Printf("%-20s : %s\n", xcolor.Green("BuildUser"), xcolor.Blue(buildUser))
	fmt.Printf("%-20s : %s\n", xcolor.Green("BuildHost"), xcolor.Blue(buildHost))
	fmt.Printf("%-20s : %s\n", xcolor.Green("BuildTime"), xcolor.Blue(BuildTime()))
	fmt.Printf("%-20s : %s\n", xcolor.Green("BuildStatus"), xcolor.Blue(buildStatus))
}
