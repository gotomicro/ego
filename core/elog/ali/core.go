package ali

import (
	"fmt"
	"sync"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var _ zapcore.Encoder = &mapObjEncoder{}

var encoderPool = sync.Pool{New: func() interface{} {
	return &mapObjEncoder{MapObjectEncoder: zapcore.NewMapObjectEncoder()}
}}

func getEncoder() *mapObjEncoder {
	return encoderPool.Get().(*mapObjEncoder)
}

func putEncoder(enc *mapObjEncoder) {
	enc.MapObjectEncoder = zapcore.NewMapObjectEncoder()
	encoderPool.Put(enc)
}

// mapObjEncoder ...
type mapObjEncoder struct {
	*zapcore.EncoderConfig
	parentFields []zapcore.Field
	*zapcore.MapObjectEncoder
}

// newMapObjEncoder ...
func newMapObjEncoder(cfg zapcore.EncoderConfig) *mapObjEncoder {
	return &mapObjEncoder{
		EncoderConfig:    &cfg,
		MapObjectEncoder: zapcore.NewMapObjectEncoder(),
	}
}

// Clone ...
func (e *mapObjEncoder) Clone() zapcore.Encoder {
	return e.clone()
}

func (e *mapObjEncoder) clone() *mapObjEncoder {
	clone := getEncoder()
	clone.EncoderConfig = e.EncoderConfig
	// copy parentFields
	clone.parentFields = make([]zapcore.Field, 0, len(e.parentFields))
	clone.parentFields = append(clone.parentFields, e.parentFields...)

	// copy fields
	for k, v := range e.MapObjectEncoder.Fields {
		clone.MapObjectEncoder.Fields[k] = v
	}
	return clone
}

// EncodeEntry ...
func (e *mapObjEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// do nothing, just implement zapcore.Encoder.EncodeEntry()
	return nil, nil
}

func (e *mapObjEncoder) encodeEntry(ent zapcore.Entry, fields []zapcore.Field) *mapObjEncoder {
	final := e.clone()
	if final.LevelKey != "" {
		final.AddString(final.LevelKey, ent.Level.String())
	}
	if final.TimeKey != "" {
		final.AddInt64(final.TimeKey, ent.Time.Unix())
	}
	if ent.LoggerName != "" && final.NameKey != "" {
		final.AddString(final.NameKey, ent.LoggerName)
	}
	if ent.Caller.Defined && final.CallerKey != "" {
		final.AddString(final.CallerKey, ent.Caller.String())
	}
	if final.MessageKey != "" {
		final.AddString(final.MessageKey, ent.Message)
	}
	if ent.Stack != "" && final.StacktraceKey != "" {
		final.AddString(final.StacktraceKey, ent.Stack)
	}
	addFields(final, fields)
	return final
}

// NewCore creates a Core that writes logs to a WriteSyncer.
func NewCore(ops ...Option) (zapcore.Core, func() error) {
	var c config
	for _, o := range ops {
		o(&c)
	}
	aliLs, err := newWriter(c)
	if err != nil {
		panic(fmt.Errorf("NewCore fail, %w", err))
	}

	core := &ioCore{
		LevelEnabler: c.levelEnabler,
		enc:          c.encoder,
		writer:       aliLs,
	}
	closeFunc := func() error {
		core.writer.cancel()
		return core.writer.flush()
	}
	return core, closeFunc
}

func addFields(enc zapcore.ObjectEncoder, fields []zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}

type ioCore struct {
	zapcore.LevelEnabler
	enc    zapcore.Encoder
	writer *writer
}

func (c *ioCore) With(fields []zapcore.Field) zapcore.Core {
	c.enc.(*mapObjEncoder).parentFields = fields
	clone := c.clone()
	// NOTICE: we must reset parentFields otherwise parent logger with also print parent fields
	c.enc.(*mapObjEncoder).parentFields = nil
	return clone
}

func (c *ioCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *ioCore) Write(ent zapcore.Entry, fields []zapcore.Field) (err error) {
	clone := c.enc.(*mapObjEncoder).encodeEntry(ent, fields)
	addFields(clone, append(fields, clone.parentFields...))
	if err := c.writer.write(clone.Fields); err != nil {
		return err
	}
	if ent.Level > zapcore.ErrorLevel {
		err = c.Sync()
	}
	putEncoder(clone)
	return
}

func (c *ioCore) Sync() error {
	return c.writer.flush()
}

func (c *ioCore) clone() *ioCore {
	return &ioCore{
		LevelEnabler: c.LevelEnabler,
		enc:          c.enc.Clone(),
		writer:       c.writer,
	}
}
