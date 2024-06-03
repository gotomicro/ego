package elog

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/gotomicro/ego/core/etrace"
)

func TestFieldAddr(t *testing.T) {
	value := zap.Field{Key: "addr", Type: zapcore.StringType, String: "127.0.0.1"}
	assert.True(t, reflect.DeepEqual(value, FieldAddr("127.0.0.1")))
}

func TestFieldApp(t *testing.T) {
	value := zap.Field{Key: "app", Type: zapcore.StringType, String: "ego-svc"}
	assert.True(t, reflect.DeepEqual(value, FieldApp("ego-svc")))
}

func TestFieldCode(t *testing.T) {
	value := zap.Field{Key: "code", Type: zapcore.Int32Type, Integer: int64(1)}
	assert.True(t, reflect.DeepEqual(value, FieldCode(1)))
}

func TestFieldComponent(t *testing.T) {
	value := zap.Field{Key: "comp", Type: zapcore.StringType, String: "server"}
	assert.True(t, reflect.DeepEqual(value, FieldComponent("server")))
}

func TestFieldComponentName(t *testing.T) {
	value := zap.Field{Key: "compName", Type: zapcore.StringType, String: "ego"}
	assert.True(t, reflect.DeepEqual(value, FieldComponentName("ego")))
}

func TestFieldName(t *testing.T) {
	value := zap.Field{Key: "name", Type: zapcore.StringType, String: "ego"}
	assert.True(t, reflect.DeepEqual(value, FieldName("ego")))
}

func TestFieldType(t *testing.T) {
	value := zap.Field{Key: "type", Type: zapcore.StringType, String: "ego"}
	assert.True(t, reflect.DeepEqual(value, FieldType("ego")))
}

func TestFieldKind(t *testing.T) {
	value := zap.Field{Key: "kind", Type: zapcore.StringType, String: "ego"}
	assert.True(t, reflect.DeepEqual(value, FieldKind("ego")))
}

func TestFieldUniformCode(t *testing.T) {
	value := zap.Field{Key: "ucode", Type: zapcore.Int32Type, Integer: int64(20)}
	assert.True(t, reflect.DeepEqual(value, FieldUniformCode(20)))
}

func TestFieldTid(t *testing.T) {
	value := zap.Field{Key: "tid", Type: zapcore.StringType, String: "111"}
	assert.True(t, reflect.DeepEqual(value, FieldTid("111")))
}

func TestFieldCtxTid(t *testing.T) {
	var ctx context.Context
	value := zap.Field{Key: "tid", Type: zapcore.StringType, String: etrace.ExtractTraceID(ctx)}
	assert.True(t, reflect.DeepEqual(value, FieldCtxTid(ctx)))
}

func TestFieldSize(t *testing.T) {
	value := zap.Field{Key: "size", Type: zapcore.Int32Type, Integer: int64(1)}
	assert.True(t, reflect.DeepEqual(value, FieldSize(1)))
}

func TestFieldKey(t *testing.T) {
	value := zap.Field{Key: "key", Type: zapcore.StringType, String: "ego"}
	assert.True(t, reflect.DeepEqual(value, FieldKey("ego")))
}

func TestFieldValue(t *testing.T) {
	value := zap.Field{Key: "value", Type: zapcore.StringType, String: "server"}
	assert.True(t, reflect.DeepEqual(value, FieldValue("server")))
}

func TestFieldErrKind(t *testing.T) {
	value := zap.Field{Key: "errKind", Type: zapcore.StringType, String: "ego-err"}
	assert.True(t, reflect.DeepEqual(value, FieldErrKind("ego-err")))
}

func TestFieldDescription(t *testing.T) {
	value := zap.Field{Key: "desc", Type: zapcore.StringType, String: "server-ego"}
	assert.True(t, reflect.DeepEqual(value, FieldDescription("server-ego")))
}

func TestFieldMethod(t *testing.T) {
	value := zap.Field{Key: "method", Type: zapcore.StringType, String: "ego"}
	assert.True(t, reflect.DeepEqual(value, FieldMethod("ego")))
}

func TestFieldEvent(t *testing.T) {
	value := zap.Field{Key: "event", Type: zapcore.StringType, String: "ego--service"}
	assert.True(t, reflect.DeepEqual(value, FieldEvent("ego--service")))
}

func TestFieldIP(t *testing.T) {
	value := zap.Field{Key: "ip", Type: zapcore.StringType, String: "127.162.1.1"}
	assert.True(t, reflect.DeepEqual(value, FieldIP("127.162.1.1")))
}

func TestFieldPeerIP(t *testing.T) {
	value := zap.Field{Key: "peerIp", Type: zapcore.StringType, String: "197.162.1.1"}
	assert.True(t, reflect.DeepEqual(value, FieldPeerIP("197.162.1.1")))
}

func TestFieldPeerName(t *testing.T) {
	value := zap.Field{Key: "peerName", Type: zapcore.StringType, String: "ego-peer"}
	assert.True(t, reflect.DeepEqual(value, FieldPeerName("ego-peer")))
}

func TestFieldLogName(t *testing.T) {
	value := zap.Field{Key: "lname", Type: zapcore.StringType, String: "logger"}
	assert.True(t, reflect.DeepEqual(value, FieldLogName("logger")))
}
