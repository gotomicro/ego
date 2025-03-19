//go:build windows

package ego

import (
	"os"
	"syscall"
)

var (
	// DefaultShutdownSignals Windows下的默认停止信号
	DefaultShutdownSignals = []os.Signal{
		syscall.SIGTERM,
		syscall.SIGINT,
	}

	// ReloadSignal Windows下不支持SIGUSR1，使用SIGTERM代替
	ReloadSignal = syscall.SIGTERM
)
