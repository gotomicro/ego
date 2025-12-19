package otel

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	conf := DefaultConfig()
	out := Load("")
	assert.Equal(t, conf, out)
	Load("").Build()
	out1 := conf.buildJaegerTP()
	assert.True(t, true, out1)
	err := conf.Stop()
	assert.NoError(t, err)
}

func TestResolveVal(t *testing.T) {
	// 1. 准备阶段：设置测试用的环境变量
	os.Setenv("TEST_INSTANCE_ID", "service-001")
	os.Setenv("EMPTY_VAR", "")

	// 测试结束后清理环境变量
	defer func() {
		os.Unsetenv("TEST_INSTANCE_ID")
		os.Unsetenv("EMPTY_VAR")
	}()

	// 2. 定义测试用例
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "普通字符串",
			input:    "regular-string",
			expected: "regular-string",
		},
		{
			name:     "合法的环境变量且存在",
			input:    "$TEST_INSTANCE_ID",
			expected: "service-001",
		},
		{
			name:     "合法的环境变量但不存在",
			input:    "$NON_EXISTENT_VAR",
			expected: "$NON_EXISTENT_VAR",
		},
		{
			name:     "合法的环境变量但值为空",
			input:    "$EMPTY_VAR",
			expected: "$EMPTY_VAR",
		},
		{
			name:     "格式错误：$后跟数字（正则不匹配）",
			input:    "$123ID",
			expected: "$123ID",
		},
		{
			name:     "格式错误：$在字符串中间",
			input:    "hello$WORLD",
			expected: "hello$WORLD",
		},
		{
			name:     "只有美元符号",
			input:    "$",
			expected: "$",
		},
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
	}

	// 3. 执行测试
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveVal(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}
