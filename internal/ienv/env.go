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

// EnvOrUint ...
func EnvOrUint(envVar string, defaultValue uint) uint {
	if v, ok := os.LookupEnv(envVar); ok && v != "" {
		intValue, _ := strconv.ParseUint(v, 10, 64)
		return uint(intValue)
	}
	return defaultValue
}

// EnvOrFloat64 ...
func EnvOrFloat64(envVar string, defaultValue float64) float64 {
	if v, ok := os.LookupEnv(envVar); ok && v != "" {
		floatValue, _ := strconv.ParseFloat(v, 64)
		return floatValue
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
