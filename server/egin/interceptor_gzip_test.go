package egin

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGzip(t *testing.T) {
	opts := func(options *GzipOptions) {}
	Gzip(3, opts)
	assert.NoError(t, nil)

	newGzipHandler(3, opts)
	assert.NoError(t, nil)
}

func TestWithGzipExcludedExtensions(t *testing.T) {
	a := []string{"hello", "world"}
	WithGzipExcludedExtensions(a)
	WithGzipExcludedPaths(a)
	WithGzipExcludedPathsRegexs(a)
	NewExcludedPaths(a)
	NewExcludedPathesRegexs(a)
	assert.NoError(t, nil)
}

func TestWithGzipDecompressFn(t *testing.T) {
	d := func(c *gin.Context) {}
	WithGzipDecompressFn(d)
	assert.NoError(t, nil)
}

func TestContains(t *testing.T) {
	var e = ExcludedPathesRegexs{}
	out := e.Contains("")
	assert.Equal(t, false, out)

	var ee = ExcludedPaths{}
	out1 := ee.Contains("")
	assert.Equal(t, false, out1)

	var a = ExcludedExtensions{}
	out2 := a.Contains("")
	assert.Equal(t, false, out2)
}
