//go:build !windows

package ego

import (
	"os"
	"syscall"
)

var (
	// DefaultShutdownSignals 默认停止信号
	DefaultShutdownSignals = []os.Signal{
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	}

	// ReloadSignal 重载信号
	ReloadSignal = syscall.SIGUSR1
)
