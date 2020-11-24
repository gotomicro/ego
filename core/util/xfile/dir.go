package xfile

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/multierr"
)

// GetCurrentDirectory ...
func GetCurrentDirectory() string {
	if dir, err := os.Getwd(); err == nil {
		return dir
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal("err", err)
	}

	return strings.Replace(dir, "\\", "/", -1)
}

// GetCurrentPackage ...
func GetCurrentPackage() string {
	if dir, err := os.Getwd(); err == nil {
		return dir
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal("err", err)
	}

	return strings.Replace(dir, "\\", "/", -1)
}

// MakeDirectory ...
func MakeDirectory(dirs ...string) error {
	var errs error
	for _, dir := range dirs {
		if !Exists(dir) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				errs = multierr.Append(errs, err)
			}
		}
	}
	return errs
}
