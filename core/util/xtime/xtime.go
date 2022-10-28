package xtime

import (
	"os"
	"time"
)

// Duration wraps time.ParseDuration but do panic when parse duration occurs.
func Duration(str string) time.Duration {
	dur, err := time.ParseDuration(str)
	if err != nil {
		panic(err)
	}
	return dur
}

// TS RFC 3339 with seconds
// Deprecated: this function will be moved to internal package, user should not use it any more.
var TS TimeFormat = "2006-01-02 15:04:05"

// TimeFormat ...
// Deprecated: this function will be moved to internal package, user should not use it any more.
type TimeFormat string

// Format 格式化
// Deprecated: this function will be moved to internal package, user should not use it any more.
func (ts TimeFormat) Format(t time.Time) string {
	return t.Format(string(ts))
}

// ParseInLocation parses time with location from env "TZ", if "TZ" hasn't been set then we use UTC by default.
func ParseInLocation(layout, value string) (time.Time, error) {
	loc, err := time.LoadLocation(os.Getenv("TZ"))
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation(layout, value, loc)
}
