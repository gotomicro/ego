package egrpclog

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/elog"
)

func TestBuild(t *testing.T) {
	exp := elog.EgoLogger.With(elog.FieldComponentName("component.grpc"))
	assert.Equal(t, exp, Build())
}
