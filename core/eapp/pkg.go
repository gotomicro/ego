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
	"github.com/gotomicro/ego/internal/ienv"
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
		appName = ienv.EnvOrStr(constant.EnvAppName, filepath.Base(os.Args[0]))
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

// BuildUser get buildUser
func BuildUser() string {
	return buildUser
}

// BuildHost get buildHost
func BuildHost() string {
	return buildHost
}

// BuildStatus get buildStatus
func BuildStatus() string {
	return buildStatus
}

// BuildTime get buildTime
func BuildTime() string {
	return buildTime
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

// PrintVersion print formatted version info
func PrintVersion() {
	fmt.Printf("%-20s : %s\n", xcolor.Green("EGO"), xcolor.Blue("I am EGO"))
	fmt.Printf("%-20s : %s\n", xcolor.Green("AppName"), xcolor.Blue(appName))
	fmt.Printf("%-20s : %s\n", xcolor.Green("AppHost"), xcolor.Blue(HostName()))
	fmt.Printf("%-20s : %s\n", xcolor.Green("Region"), xcolor.Blue(AppRegion()))
	fmt.Printf("%-20s : %s\n", xcolor.Green("Zone"), xcolor.Blue(AppZone()))
	fmt.Printf("%-20s : %s\n", xcolor.Green("AppVersion"), xcolor.Blue(AppVersion()))
	fmt.Printf("%-20s : %s\n", xcolor.Green("EgoVersion"), xcolor.Blue(EgoVersion()))
	fmt.Printf("%-20s : %s\n", xcolor.Green("BuildUser"), xcolor.Blue(BuildUser()))
	fmt.Printf("%-20s : %s\n", xcolor.Green("BuildHost"), xcolor.Blue(BuildHost()))
	fmt.Printf("%-20s : %s\n", xcolor.Green("BuildTime"), xcolor.Blue(BuildTime()))
	fmt.Printf("%-20s : %s\n", xcolor.Green("BuildStatus"), xcolor.Blue(BuildStatus()))
}
