package eflag

import (
	"strings"

	"github.com/gotomicro/ego/internal/ienv"
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
			set.FlagSet.IntVar(f.Variable, field, ienv.EnvOrInt(f.EnvVar, f.Default), f.Usage)
		} else {
			set.FlagSet.Int(field, ienv.EnvOrInt(f.EnvVar, f.Default), f.Usage)
		}
		set.actions[field] = f.Action
	}
}
