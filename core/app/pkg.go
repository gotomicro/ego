package app

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gotomicro/ego/core/util/xtime"

	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/util/xcolor"
)

const egoVersion = "0.0.1"

var (
	startTime string
	goVersion string
)

// build info
/*

 */
var (
	appName         string
	appID           string
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
	SetBuildTime(buildTime)
	goVersion = runtime.Version()
	InitEnv()
}

// Name gets application name.
func Name() string {
	return appName
}

//SetName set app anme
func SetName(s string) {
	appName = s
}

//AppID get appID
func AppID() string {
	return appID
}

//SetAppID set appID
func SetAppID(s string) {
	appID = s
}

//AppVersion get buildAppVersion
func AppVersion() string {
	return buildAppVersion
}

//appVersion not defined
// func SetAppVersion(s string) {
// 	appVersion = s
// }

// EgoVersion get egoVersion
func EgoVersion() string {
	return egoVersion
}

// todo: egoVersion is const not be set
// func SetegoVersion(s string) {
// 	egoVersion = s
// }

//BuildTime get buildTime
func BuildTime() string {
	return buildTime
}

//BuildUser get buildUser
func BuildUser() string {
	return buildUser
}

//BuildHost get buildHost
func BuildHost() string {
	return buildHost
}

//SetBuildTime set buildTime
func SetBuildTime(param string) {
	buildTime = strings.Replace(param, "--", " ", 1)
}

// HostName get host name
func HostName() string {
	return hostName
}

//StartTime get start time
func StartTime() string {
	return startTime
}

//GoVersion get go version
func GoVersion() string {
	return goVersion
}

// PrintVersion print formated version info
func PrintVersion() {
	fmt.Printf("%-8s]> %-30s => %s\n", "ego", xcolor.Red("name"), xcolor.Blue(appName))
	fmt.Printf("%-8s]> %-30s => %s\n", "ego", xcolor.Red("appID"), xcolor.Blue(appID))
	fmt.Printf("%-8s]> %-30s => %s\n", "ego", xcolor.Red("region"), xcolor.Blue(AppRegion()))
	fmt.Printf("%-8s]> %-30s => %s\n", "ego", xcolor.Red("zone"), xcolor.Blue(AppZone()))
	fmt.Printf("%-8s]> %-30s => %s\n", "ego", xcolor.Red("appVersion"), xcolor.Blue(buildAppVersion))
	fmt.Printf("%-8s]> %-30s => %s\n", "ego", xcolor.Red("egoVersion"), xcolor.Blue(egoVersion))
	fmt.Printf("%-8s]> %-30s => %s\n", "ego", xcolor.Red("buildUser"), xcolor.Blue(buildUser))
	fmt.Printf("%-8s]> %-30s => %s\n", "ego", xcolor.Red("buildHost"), xcolor.Blue(buildHost))
	fmt.Printf("%-8s]> %-30s => %s\n", "ego", xcolor.Red("buildTime"), xcolor.Blue(BuildTime()))
	fmt.Printf("%-8s]> %-30s => %s\n", "ego", xcolor.Red("buildStatus"), xcolor.Blue(buildStatus))
}
