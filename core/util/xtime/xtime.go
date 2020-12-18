package xtime

import "time"

// Duration ...
// panic if parse duration failed
func Duration(str string) time.Duration {
	dur, err := time.ParseDuration(str)
	if err != nil {
		panic(err)
	}
	return dur
}

var TS TimeFormat = "2006-01-02 15:04:05"

type TimeFormat string

func (ts TimeFormat) Format(t time.Time) string {
	return t.Format(string(ts))
}
