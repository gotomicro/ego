package xtime

import (
	"os"
	"time"
)

// Duration ...
// panic if parse duration failed
func Duration(str string) time.Duration {
	dur, err := time.ParseDuration(str)
	if err != nil {
		panic(err)
	}
	return dur
}

// TS ...
var TS TimeFormat = "2006-01-02 15:04:05"

// TimeFormat ...
type TimeFormat string

// Format 格式化
func (ts TimeFormat) Format(t time.Time) string {
	return t.Format(string(ts))
}

// ParseInLocation parse time with location from env "TZ", if "TZ" hasn't been set then we use UTC by default.
func ParseInLocation(layout, value string) (time.Time, error) {
	loc, err := time.LoadLocation(os.Getenv("TZ"))
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation(layout, value, loc)
}
