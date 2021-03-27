package eflag

import (
	"os"
	"strconv"
	"strings"
)

// IntFlag is an int flag implements of Flag interface.
type IntFlag struct {
	Name     string
	Usage    string
	EnvVar   string
	Default  int
	Variable *int
	Action   func(string, *FlagSet)
}

// Apply implements of Flag Apply function.
func (f *IntFlag) Apply(set *FlagSet) {
	for _, field := range strings.Split(f.Name, ",") {
		field = strings.TrimSpace(field)
		if f.Variable != nil {
			set.FlagSet.IntVar(f.Variable, field, getValueByEnvAndDefaultIntValue(f.EnvVar, f.Default), f.Usage)
		}
		set.FlagSet.Int(field, getValueByEnvAndDefaultIntValue(f.EnvVar, f.Default), f.Usage)
		set.actions[field] = f.Action
	}
}

func getValueByEnvAndDefaultIntValue(envVar string, defaultValue int) int {
	env := os.Getenv(envVar)
	if env != "" {
		intValue, _ := strconv.ParseInt(env, 10, 64)
		return int(intValue)
	}
	return defaultValue
}
