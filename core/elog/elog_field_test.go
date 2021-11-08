package elog

import (
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

//func TestFieldCost(t *testing.T) {
//	value := zap.Field{Key: "compName", Type: zapcore.Float64Type, Integer: int64(math.Float64bits(0.16))}
//	assert.True(t, reflect.DeepEqual(value, FieldCost(0.16)))
//}
