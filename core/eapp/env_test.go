package eapp

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/constant"
)

func TestAppMode(t *testing.T) {
	os.Setenv(constant.EnvAppMode, "test-mode")
	defer os.Unsetenv(constant.EnvAppMode)

	initEnv()
	out := AppMode()
	assert.Equal(t, "test-mode", out)
}
func TestAppRegion(t *testing.T) {
	os.Setenv(constant.EnvAppRegion, "test-region")
	defer os.Unsetenv(constant.EnvAppRegion)

	initEnv()
	out := AppRegion()
	assert.Equal(t, "test-region", out)
}

func TestAppZone(t *testing.T) {
	os.Setenv(constant.EnvAppZone, "test-zone")
	defer os.Unsetenv(constant.EnvAppZone)

	initEnv()
	out := AppZone()
	assert.Equal(t, "test-zone", out)
}

func TestAppInstance(t *testing.T) {
	os.Setenv(constant.EnvAppInstance, "test-instance-1")
	defer os.Unsetenv(constant.EnvAppInstance)

	initEnv()
	out := AppInstance()
	assert.Equal(t, "test-instance-1", out)
}

func TestIsDevelopmentMode(t *testing.T) {
	os.Setenv(constant.EgoDebug, "true")
	defer os.Unsetenv(constant.EgoDebug)

	initEnv()
	out := IsDevelopmentMode()
	assert.Equal(t, true, out)
}

func TestEgoLogPath(t *testing.T) {
	os.Setenv(constant.EgoLogPath, "test-ego.log")
	defer os.Unsetenv(constant.EgoLogPath)

	initEnv()
	out := EgoLogPath()
	assert.Equal(t, "test-ego.log", out)
}

func TestEnableLoggerAddApp(t *testing.T) {
	os.Setenv(constant.EgoLogAddApp, "true")
	defer os.Unsetenv(constant.EgoLogAddApp)

	initEnv()
	out := EnableLoggerAddApp()
	assert.Equal(t, true, out)
}

func TestEgoTraceIDName(t *testing.T) {
	os.Setenv(constant.EgoTraceIDName, "x-trace-id")
	defer os.Unsetenv(constant.EgoTraceIDName)

	initEnv()
	out := EgoTraceIDName()
	assert.Equal(t, "x-trace-id", out)
}

func TestEgoLogExtraKeys(t *testing.T) {
	os.Setenv(constant.EgoLogExtraKeys, "x-ego-uid")
	defer os.Unsetenv(constant.EgoLogExtraKeys)

	initEnv()
	out := EgoLogExtraKeys()
	assert.Equal(t, []string{"x-ego-uid"}, out)
}
