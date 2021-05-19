// Package rotate provides a rolling logger.
//
// Note that this is v2.0 of rotate, and should be imported using gopkg.in
// thusly:
//
//   import "gopkg.in/natefinch/rotate.v2"
//
// The package name remains simply rotate, and the code resides at
// https://github.com/natefinch/rotate under the v2.0 branch.
//
// rotate is intended to be one part of a logging infrastructure.
// It is not an all-in-one solution, but instead is a pluggable
// component at the bottom of the logging stack that simply controls the files
// to which logs are written.
//
// rotate plays well with any logging package that can write to an
// io.Writer, including the standard library's log package.
//
// rotate assumes that only one process is writing to the output files.
// Using the same rotate configuration from multiple processes on the same
// machine will result in improper behavior.

// +build linux

package rotate

import (
	"os"
	"syscall"
	"time"
)

func ctime(file *os.File) (time.Time, error) {
	fi, err := file.Stat()
	if err != nil {
		return time.Now(), err
	}

	stat := fi.Sys().(*syscall.Stat_t)
	return time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec)), nil
}
