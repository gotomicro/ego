package eflag

import (
	"strings"

	"github.com/gotomicro/ego/internal/ienv"
)

// UintFlag is an uint flag implements of Flag interface.
type UintFlag struct {
	Name     string
	Usage    string
	EnvVar   string
	Default  uint
	Variable *uint
	Action   func(string, *FlagSet)
}

// Apply implements of Flag Apply function.
func (f *UintFlag) Apply(set *FlagSet) {
	for _, field := range strings.Split(f.Name, ",") {
		field = strings.TrimSpace(field)
		if f.Variable != nil {
			set.FlagSet.UintVar(f.Variable, field, ienv.EnvOrUint(f.EnvVar, f.Default), f.Usage)
		} else {
			set.FlagSet.Uint(field, ienv.EnvOrUint(f.EnvVar, f.Default), f.Usage)
		}
		set.actions[field] = f.Action
	}
}
