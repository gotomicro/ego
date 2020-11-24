// +build windows
package xfile

import (
	"io"
	"syscall"
)

// lockCloser hides all of an syscall.Handle's methods, except for Close.
type lockCloser struct {
	fd syscall.Handle
}

// Close ...
func (l lockCloser) Close() error {
	return syscall.Close(l.fd)
}

// Lock ...
func Lock(name string) (io.Closer, error) {
	p, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return nil, err
	}
	fd, err := syscall.CreateFile(p,
		syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		0, nil, syscall.CREATE_ALWAYS,
		syscall.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != nil {
		return nil, err
	}
	return lockCloser{fd: fd}, nil
}
