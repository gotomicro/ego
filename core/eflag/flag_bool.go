package eflag

import (
	"strings"

	"github.com/gotomicro/ego/internal/ienv"
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
			set.FlagSet.BoolVar(f.Variable, field, ienv.EnvOrBool(f.EnvVar, f.Default), f.Usage)
		} else {
			set.FlagSet.Bool(field, ienv.EnvOrBool(f.EnvVar, f.Default), f.Usage)
		}
		set.actions[field] = f.Action
	}
}
