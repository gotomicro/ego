package elog

import (
	"context"
	"strings"
	"time"

	"github.com/gotomicro/ego/core/etrace"
	"go.uber.org/zap"
)

// FieldComponent constructs an elog Field with component type name
func FieldComponent(value string) Field {
	return String("comp", value)
}

// FieldComponentName constructs an elog Field with component name
func FieldComponentName(value string) Field {
	return String("compName", value)
}

// FieldApp constructs an elog Field with ego application name
func FieldApp(value string) Field {
	return String("app", value)
}

// FieldAddr constructs an elog Field with some address
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

// FieldUniformCode uniform code
func FieldUniformCode(value int32) Field {
	return Int32("ucode", value)
}

// FieldTid constructs an elog Field with traceID
func FieldTid(value string) Field {
	return String("tid", value)
}

// FieldCtxTid constructs an elog Field with traceID which extracted from context
func FieldCtxTid(ctx context.Context) Field {
	return String("tid", etrace.ExtractTraceID(ctx))
}

// FieldSize ...
func FieldSize(value int32) Field {
	return Int32("size", value)
}

// FieldCost construct an elog Field with time cost
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

// FieldCustomKeyValue constructs a custom Key and value
func FieldCustomKeyValue(key string, value string) Field {
	return String(strings.ToLower(key), value)
}

// FieldLogName constructs a field log name
func FieldLogName(value string) Field {
	return String("lname", value)
}
