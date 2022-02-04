package eflag

import (
	"strings"

	"github.com/gotomicro/ego/internal/ienv"
)

// StringFlag is a string flag implements of Flag interface.
type StringFlag struct {
	Name     string
	Usage    string
	EnvVar   string
	Default  string
	Variable *string
	Action   func(string, *FlagSet)
}

// Apply implements of Flag Apply function.
func (f *StringFlag) Apply(set *FlagSet) {
	for _, field := range strings.Split(f.Name, ",") {
		field = strings.TrimSpace(field)
		if f.Variable != nil {
			set.FlagSet.StringVar(f.Variable, field, ienv.EnvOrStr(f.EnvVar, f.Default), f.Usage)
		} else {
			set.FlagSet.String(field, ienv.EnvOrStr(f.EnvVar, f.Default), f.Usage)
		}
		set.actions[field] = f.Action
	}
}
