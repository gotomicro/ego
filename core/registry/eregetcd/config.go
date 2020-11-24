package eregetcd

import (
	"time"
)

// Config ...
type Config struct {
	ReadTimeout time.Duration
	Prefix      string
	ServiceTTL  time.Duration
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		ReadTimeout: time.Second * 3,
		Prefix:      "ego",
		ServiceTTL:  0,
	}
}
