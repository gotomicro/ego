package eflag

import (
	"os"
	"strconv"
	"strings"
)

// BoolFlag is a bool flag implements of Flag interface.
type BoolFlag struct {
	Name     string
	Usage    string
	EnvVar   string
	Default  bool
	Variable *bool
	Action   func(string, *FlagSet)
}

// Apply implements of Flag Apply function.
func (f *BoolFlag) Apply(set *FlagSet) {
	for _, field := range strings.Split(f.Name, ",") {
		field = strings.TrimSpace(field)
		if f.Variable != nil {
			set.FlagSet.BoolVar(f.Variable, field, getValueByEnvAndDefaultBoolValue(f.EnvVar, f.Default), f.Usage)
		}

		set.FlagSet.Bool(field, getValueByEnvAndDefaultBoolValue(f.EnvVar, f.Default), f.Usage)
		set.actions[field] = f.Action
	}
}

func getValueByEnvAndDefaultBoolValue(envVar string, defaultValue bool) bool {
	env := os.Getenv(envVar)
	if env != "" {
		flag, _ := strconv.ParseBool(env)
		return flag
	}
	return defaultValue
}
