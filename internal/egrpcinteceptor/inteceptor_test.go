package egrpcinteceptor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitMethodName(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		f := "/hello.service/GET"
		service, method := SplitMethodName(f)
		assert.Equal(t, "hello.service", service)
		assert.Equal(t, "GET", method)
	})
	t.Run("case 2", func(t *testing.T) {
		f := ""
		service, method := SplitMethodName(f)
		assert.Equal(t, "unknown", service)
		assert.Equal(t, "unknown", method)
	})
}
