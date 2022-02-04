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
