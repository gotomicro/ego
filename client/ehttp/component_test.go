package ehttp

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/elog"
)

func TestNewComponent(t *testing.T) {
	// Normal case
	t.Run("Normal case", func(t *testing.T) {
		config := &Config{Addr: "http://hello.com"}
		logger := elog.DefaultLogger
		component := newComponent("hello", config, logger)
		assert.Equal(t, "hello", component.name)
	})

	// Other case...
}
