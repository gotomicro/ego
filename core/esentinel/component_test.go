package esentinel

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/elog"
)

var logger = &elog.Component{}

func TestNewComponent(t *testing.T) {
	conf := &Config{}
	newComponent(conf, logger)
	assert.NoError(t, nil)
}

func TestSyncFlowRules(t *testing.T) {
	filePath := "./config_test/sentinel.json"
	err := syncFlowRules(filePath, logger)
	assert.NoError(t, err)
}

func TestIsResMap(t *testing.T) {
	res := "test"
	assert.Equal(t, false, IsResExist(res))
}
