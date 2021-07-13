// +build !linux

package elog

import (
	"os"
)

func chown(_ string, _ os.FileInfo) error {
	return nil
}
