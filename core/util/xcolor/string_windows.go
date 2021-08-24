// +build windows

package xcolor

import (
	"fmt"
)

//
//var _ = Random()

// Yellow ...
func Yellow(msg string) string {
	return fmt.Sprintf("%s", msg)
}

// Red ...
func Red(msg string) string {
	return fmt.Sprintf("%s", msg)
}

// Redf ...
func Redf(msg string, arg interface{}) string {
	return fmt.Sprintf("%s %+v\n", msg, arg)
}

// Blue ...
func Blue(msg string) string {
	return fmt.Sprintf("%s", msg)
}

// Green ...
func Green(msg string) string {
	return fmt.Sprintf("%s", msg)
}

// Greenf ...
func Greenf(msg string, arg interface{}) string {
	return fmt.Sprintf("%s %+v\n", msg, arg)
}
