// +build !windows

package xfile

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"
)

// FileStat return a FileInfo describing the named file.
func FileStat(name string) (fi FileInfo, err error) {
	if Exists(name) {
		f, err := os.Open(name)
		if err != nil {
			return fi, err
		}
		defer f.Close()
		stats, _ := f.Stat()
		fi.Uid = stats.Sys().(*syscall.Stat_t).Uid
		fi.Gid = stats.Sys().(*syscall.Stat_t).Gid
		fi.Mode = stats.Mode()
		h := md5.New()
		_, _ = io.Copy(h, f)
		fi.Md5 = fmt.Sprintf("%x", h.Sum(nil))
		return fi, nil
	}
	return fi, errors.New("file not found")
}
