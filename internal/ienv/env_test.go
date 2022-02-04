package ienv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvOrBoolNoEnv(t *testing.T) {
	flag := EnvOrBool("ego-env-test1", true)
	assert.Equal(t, true, flag)
}

func TestEnvOrBoolHaveEnv(t *testing.T) {
	os.Setenv("ego-env-test1", "false")
	defer os.Unsetenv("ego-env-test1")

	flag := EnvOrBool("ego-env-test1", true)
	assert.Equal(t, false, flag)
}

func TestEnvOrIntNoEnv(t *testing.T) {
	flag := EnvOrInt("ego-env-test1", 1)
	assert.Equal(t, 1, flag)
}

func TestEnvOrIntHaveEnv(t *testing.T) {
	os.Setenv("ego-env-test1", "2")
	defer os.Unsetenv("ego-env-test1")

	flag := EnvOrInt("ego-env-test1", 1)
	assert.Equal(t, 2, flag)
}

func TestEnvOrUintNoEnv(t *testing.T) {
	flag := EnvOrUint("ego-env-test1", 1)
	assert.Equal(t, uint(1), flag)
}

func TestEnvOrUintHaveEnv(t *testing.T) {
	os.Setenv("ego-env-test1", "2")
	defer os.Unsetenv("ego-env-test1")

	flag := EnvOrUint("ego-env-test1", 1)
	assert.Equal(t, uint(2), flag)
}

func TestEnvOrFloat64NoEnv(t *testing.T) {
	flag := EnvOrFloat64("ego-env-test1", 1.1)
	assert.Equal(t, 1.1, flag)
}

func TestEnvOrFloat64HaveEnv(t *testing.T) {
	os.Setenv("ego-env-test1", "1.2")
	defer os.Unsetenv("ego-env-test1")

	flag := EnvOrFloat64("ego-env-test1", 1.1)
	assert.Equal(t, 1.2, flag)
}

func TestEnvOrStrNoEnv(t *testing.T) {
	flag := EnvOrStr("ego-env-test1", "test1")
	assert.Equal(t, "test1", flag)
}

func TestEnvOrStrHaveEnv(t *testing.T) {
	os.Setenv("ego-env-test1", "test2")
	defer os.Unsetenv("ego-env-test1")

	flag := EnvOrStr("ego-env-test1", "test1")
	assert.Equal(t, "test2", flag)
}
