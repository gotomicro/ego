package esentinel

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/elog"
)

var logger = &elog.Component{}

func TestNewComponent(t *testing.T) {
	conf := &Config{
		AppName:       "APP_NAME",
		LogPath:       "./logs",
		FlowRulesFile: "./config_test/sentinel.json",
	}
	err := newComponent(conf, logger)
	assert.NoError(t, err)
}

func TestSyncFlowRules(t *testing.T) {
	filePath := "./config_test/sentinel.json"
	err := syncFlowRules(filePath, logger)
	assert.NoError(t, err)
}
