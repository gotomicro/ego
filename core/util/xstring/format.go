package xstring

import "fmt"

// Formatter ...
type Formatter string

// Format ...
func (fm Formatter) Format(args ...interface{}) string {
	return fmt.Sprintf(string(fm), args...)
}
