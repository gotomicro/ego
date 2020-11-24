// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris,!windows

package xfile

import (
	"fmt"
	"io"
	"runtime"
)

// Lock ...
func Lock(name string) (io.Closer, error) {
	return nil, fmt.Errorf("file locking is not implemented on %s/%s", runtime.GOOS, runtime.GOARCH)
}
