package eflag

import (
	"strings"

	"github.com/gotomicro/ego/internal/ienv"
)

// Float64Flag is a float flag implements of Flag interface.
type Float64Flag struct {
	Name     string
	Usage    string
	EnvVar   string
	Default  float64
	Variable *float64
	Action   func(string, *FlagSet)
}

// Apply implements of Flag Apply function.
func (f *Float64Flag) Apply(set *FlagSet) {
	for _, field := range strings.Split(f.Name, ",") {
		field = strings.TrimSpace(field)
		if f.Variable != nil {
			set.FlagSet.Float64Var(f.Variable, field, ienv.EnvOrFloat64(f.EnvVar, f.Default), f.Usage)
		} else {
			set.FlagSet.Float64(field, ienv.EnvOrFloat64(f.EnvVar, f.Default), f.Usage)
		}
		set.actions[field] = f.Action
	}
}
