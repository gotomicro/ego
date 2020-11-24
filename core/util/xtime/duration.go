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
