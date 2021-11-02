package eapp

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	appName = "test-app"
	out := Name()
	assert.Equal(t, appName, out)
}

func TestAppVersion(t *testing.T) {
	buildAppVersion = "master-1"
	out := AppVersion()
	assert.Equal(t, buildAppVersion, out)
}

func TestEgoVersion(t *testing.T) {
	out := EgoVersion()
	assert.Equal(t, "unknown version", out)
}

func TestBuildTime(t *testing.T) {
	buildTime = time.Now().String()
	out := BuildTime()
	assert.Equal(t, buildTime, out)
}

func TestBuildUser(t *testing.T) {
	buildUser = "unknown"
	out := BuildUser()
	assert.Equal(t, buildUser, out)
}

func TestBuildHost(t *testing.T) {
	buildHost = "localhost"
	out := BuildHost()
	assert.Equal(t, buildHost, out)
}

func TestHostName(t *testing.T) {
	out := HostName()
	assert.NotEmpty(t, out)
}

func TestStartTime(t *testing.T) {
	out := StartTime()
	assert.NotEmpty(t, out)
}

func TestGoVersion(t *testing.T) {
	out := GoVersion()
	assert.Equal(t, runtime.Version(), out)
}

func TestPrintVersion(t *testing.T) {
	appName = "test-app"
	PrintVersion()
}

func TestSetBuildTime(t *testing.T) {
	buildTime = "2021-10-28--12:00"
	setBuildTime(buildTime)
	out := BuildTime()
	assert.Equal(t, "2021-10-28 12:00", out)
}
