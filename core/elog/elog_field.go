package elog

import (
	"context"
	"strings"
	"time"

	"github.com/gotomicro/ego/core/etrace"
	"go.uber.org/zap"
)

// FieldComponent 设置组件
func FieldComponent(value string) Field {
	return String("comp", value)
}

// FieldComponentName 设置组件配置名
func FieldComponentName(value string) Field {
	return String("compName", value)
}

// FieldApp 设置应用名
func FieldApp(value string) Field {
	return String("app", value)
}

// FieldAddr 设置地址
func FieldAddr(value string) Field {
	return String("addr", value)
}

// FieldName ...
func FieldName(value string) Field {
	return String("name", value)
}

// FieldType ... level 1
func FieldType(value string) Field {
	return String("type", value)
}

// FieldKind ... level 2
func FieldKind(value string) Field {
	return String("kind", value)
}

// FieldCode ...
func FieldCode(value int32) Field {
	return Int32("code", value)
}

// FieldTid 设置链路id
func FieldTid(value string) Field {
	return String("tid", value)
}

// FieldCtxTid 设置链路id
func FieldCtxTid(ctx context.Context) Field {
	return String("tid", etrace.ExtractTraceID(ctx))
}

// FieldSize ...
func FieldSize(value int32) Field {
	return Int32("size", value)
}

// FieldCost 耗时时间
func FieldCost(value time.Duration) Field {
	return zap.Float64("cost", float64(value.Microseconds())/1000)
}

// FieldKey ...
func FieldKey(value string) Field {
	return String("key", value)
}

// FieldValue ...
func FieldValue(value string) Field {
	return String("value", value)
}

// FieldValueAny ...
func FieldValueAny(value interface{}) Field {
	return Any("value", value)
}

// FieldErrKind ...
func FieldErrKind(value string) Field {
	return String("errKind", value)
}

// FieldErr ...
func FieldErr(err error) Field {
	return zap.Error(err)
}

// FieldErrAny ...
func FieldErrAny(err interface{}) Field {
	return zap.Any("error", err)
}

// FieldDescription ...
func FieldDescription(value string) Field {
	return String("desc", value)
}

// FieldExtMessage ...
func FieldExtMessage(vals ...interface{}) Field {
	return zap.Any("ext", vals)
}

// FieldStack ...
func FieldStack(value []byte) Field {
	return ByteString("stack", value)
}

// FieldMethod ...
func FieldMethod(value string) Field {
	return String("method", value)
}

// FieldEvent ...
func FieldEvent(value string) Field {
	return String("event", value)
}

// FieldIP ...
func FieldIP(value string) Field {
	return String("ip", value)
}

// FieldPeerIP ...
func FieldPeerIP(value string) Field {
	return String("peerIp", value)
}

// FieldPeerName ...
func FieldPeerName(value string) Field {
	return String("peerName", value)
}

// FieldCustomKeyValue 设置自定义日志
func FieldCustomKeyValue(key string, value string) Field {
	return String(strings.ToLower(key), value)
}
