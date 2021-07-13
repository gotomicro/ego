package elog

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Example of how to rotate in response to SIGHUP.
func TestLoggerRotate(t *testing.T) {
	l := &rLogger{}
	log.SetOutput(l)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	go func() {
		for {
			<-c
			err := l.Rotate()
			assert.NoError(t, err)
		}
	}()

	t.Log("test done")
}
