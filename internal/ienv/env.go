package ienv

import (
	"os"
	"strconv"
)

// EnvOrBool ...
func EnvOrBool(envVar string, defaultValue bool) bool {
	if v, ok := os.LookupEnv(envVar); ok && v != "" {
		flag, _ := strconv.ParseBool(v)
		return flag
	}
	return defaultValue
}

// EnvOrInt ...
func EnvOrInt(envVar string, defaultValue int) int {
	if v, ok := os.LookupEnv(envVar); ok && v != "" {
		intValue, _ := strconv.ParseInt(v, 10, 64)
		return int(intValue)
	}
	return defaultValue
}

// EnvOrStr returns an env variable's value if it is exists or the default if not
func EnvOrStr(key, defaultValue string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return defaultValue
}
