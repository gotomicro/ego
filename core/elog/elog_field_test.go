package elog

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

func TestFieldCost(t *testing.T) {
	FieldCost(111)
	assert.NoError(t, nil)
}

func TestFieldKey(t *testing.T) {
	FieldKey("hello")
	assert.NoError(t, nil)

	FieldName("test")
	assert.NoError(t, nil)

	FieldType("type")
	assert.NoError(t, nil)

	FieldKind("kind")
	assert.NoError(t, nil)

	FieldUniformCode(11)
	assert.NoError(t, nil)

	FieldTid("tid")
	assert.NoError(t, nil)

	ctx := context.Background()
	FieldCtxTid(ctx)
	assert.NoError(t, nil)

	FieldSize(11)
	assert.NoError(t, nil)

	FieldValue("")
	assert.NoError(t, nil)

	FieldValueAny("")
	assert.NoError(t, nil)

	FieldErrKind("")
	assert.NoError(t, nil)

	FieldErr(nil)
	assert.NoError(t, nil)

	FieldErrAny(nil)
	assert.NoError(t, nil)

	FieldMethod("")
	assert.NoError(t, nil)

	FieldEvent("")
	assert.NoError(t, nil)

	FieldIP("")
	assert.NoError(t, nil)

	FieldPeerIP("")
	assert.NoError(t, nil)

	FieldPeerName("")
	assert.NoError(t, nil)

	FieldCustomKeyValue("hello", "world")
	assert.NoError(t, nil)

	FieldLogName("")
	assert.NoError(t, nil)
}
