package elog

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/gotomicro/ego/core/etrace"
)

func TestFieldAddr(t *testing.T) {
	value := zap.Field{Key: "addr", Type: zapcore.StringType, String: "127.0.0.1"}
	assert.Equal(t, value, FieldAddr("127.0.0.1"))
}

func TestFieldApp(t *testing.T) {
	value := zap.Field{Key: "app", Type: zapcore.StringType, String: "ego-svc"}
	assert.Equal(t, value, FieldApp("ego-svc"))
}

func TestFieldCode(t *testing.T) {
	value := zap.Field{Key: "code", Type: zapcore.Int32Type, Integer: int64(1)}
	assert.Equal(t, value, FieldCode(1))
}

func TestFieldComponent(t *testing.T) {
	value := zap.Field{Key: "comp", Type: zapcore.StringType, String: "server"}
	assert.Equal(t, value, FieldComponent("server"))
}

func TestFieldComponentName(t *testing.T) {
	value := zap.Field{Key: "compName", Type: zapcore.StringType, String: "ego"}
	assert.Equal(t, value, FieldComponentName("ego"))
}

func TestFieldName(t *testing.T) {
	value := zap.Field{Key: "name", Type: zapcore.StringType, String: "ego"}
	assert.Equal(t, value, FieldName("ego"))
}

func TestFieldType(t *testing.T) {
	value := zap.Field{Key: "type", Type: zapcore.StringType, String: "ego"}
	assert.Equal(t, value, FieldType("ego"))
}

func TestFieldKind(t *testing.T) {
	value := zap.Field{Key: "kind", Type: zapcore.StringType, String: "ego"}
	assert.Equal(t, value, FieldKind("ego"))
}

func TestFieldUniformCode(t *testing.T) {
	value := zap.Field{Key: "ucode", Type: zapcore.Int32Type, Integer: int64(20)}
	assert.Equal(t, value, FieldUniformCode(20))
}

func TestFieldTid(t *testing.T) {
	value := zap.Field{Key: "tid", Type: zapcore.StringType, String: "111"}
	assert.Equal(t, value, FieldTid("111"))
}

func TestFieldCtxTid(t *testing.T) {
	var ctx context.Context
	value := zap.Field{Key: "tid", Type: zapcore.StringType, String: etrace.ExtractTraceID(ctx)}
	assert.Equal(t, value, FieldCtxTid(ctx))
}

func TestFieldSize(t *testing.T) {
	value := zap.Field{Key: "size", Type: zapcore.Int32Type, Integer: int64(1)}
	assert.Equal(t, value, FieldSize(1))
}

func TestFieldKey(t *testing.T) {
	value := zap.Field{Key: "key", Type: zapcore.StringType, String: "ego"}
	assert.Equal(t, value, FieldKey("ego"))
}

func TestFieldValue(t *testing.T) {
	value := zap.Field{Key: "value", Type: zapcore.StringType, String: "server"}
	assert.Equal(t, value, FieldValue("server"))
}

func TestFieldErrKind(t *testing.T) {
	value := zap.Field{Key: "errKind", Type: zapcore.StringType, String: "ego-err"}
	assert.Equal(t, value, FieldErrKind("ego-err"))
}

func TestFieldDescription(t *testing.T) {
	value := zap.Field{Key: "desc", Type: zapcore.StringType, String: "server-ego"}
	assert.Equal(t, value, FieldDescription("server-ego"))
}

func TestFieldMethod(t *testing.T) {
	value := zap.Field{Key: "method", Type: zapcore.StringType, String: "ego"}
	assert.Equal(t, value, FieldMethod("ego"))
}

func TestFieldEvent(t *testing.T) {
	value := zap.Field{Key: "event", Type: zapcore.StringType, String: "ego--service"}
	assert.Equal(t, value, FieldEvent("ego--service"))
}

func TestFieldIP(t *testing.T) {
	value := zap.Field{Key: "ip", Type: zapcore.StringType, String: "127.162.1.1"}
	assert.Equal(t, value, FieldIP("127.162.1.1"))
}

func TestFieldPeerIP(t *testing.T) {
	value := zap.Field{Key: "peerIp", Type: zapcore.StringType, String: "197.162.1.1"}
	assert.Equal(t, value, FieldPeerIP("197.162.1.1"))
}

func TestFieldPeerName(t *testing.T) {
	value := zap.Field{Key: "peerName", Type: zapcore.StringType, String: "ego-peer"}
	assert.Equal(t, value, FieldPeerName("ego-peer"))
}

func TestFieldLogName(t *testing.T) {
	value := zap.Field{Key: "lname", Type: zapcore.StringType, String: "logger"}
	assert.Equal(t, value, FieldLogName("logger"))
}
