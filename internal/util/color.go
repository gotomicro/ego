//go:build darwin || linux

package util

import (
	"fmt"
)

// Yellow prints message with yellow.
// Deprecated: this function will be moved to internal package, user should not use it any more.
func Yellow(msg string) string {
	return fmt.Sprintf("\x1b[33m%s\x1b[0m", msg)
}

// Red prints message with red.
// Deprecated: this function will be moved to internal package, user should not use it any more.
func Red(msg string) string {
	return fmt.Sprintf("\x1b[31m%s\x1b[0m", msg)
}

// Blue prints message with blue.
// Deprecated: this function will be moved to internal package, user should not use it any more.
func Blue(msg string) string {
	return fmt.Sprintf("\x1b[34m%s\x1b[0m", msg)
}

// Green prints message with green.
// Deprecated: this function will be moved to internal package, user should not use it any more.
func Green(msg string) string {
	return fmt.Sprintf("\x1b[32m%s\x1b[0m", msg)
}
