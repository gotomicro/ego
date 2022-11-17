//go:build darwin || linux

package xcolor

import (
	"github.com/gotomicro/ego/internal/util"
)

// Yellow prints message with yellow.
// Deprecated: this function will be moved to internal package, user should not use it anymore.
var Yellow = util.Yellow

// Red prints message with red.
// Deprecated: this function will be moved to internal package, user should not use it anymore.
var Red = util.Red

// Blue prints message with blue.
// Deprecated: this function will be moved to internal package, user should not use it anymore.
var Blue = util.Blue

// Green prints message with green.
// Deprecated: this function will be moved to internal package, user should not use it anymore.
var Green = util.Green
