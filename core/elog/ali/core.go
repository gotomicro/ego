package ali

import (
	"errors"
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

type mapObjEncoder struct {
	*zapcore.MapObjectEncoder
}

func (e *mapObjEncoder) Clone() zapcore.Encoder {
	clone := getEncoder()
	// copy fields
	for k, v := range e.MapObjectEncoder.Fields {
		clone.MapObjectEncoder.Fields[k] = v
	}
	return clone
}

func (e *mapObjEncoder) EncodeEntry(zapcore.Entry, []zapcore.Field) (*buffer.Buffer, error) {
	// do nothing, just implement zapcore.Encoder.EncodeEntry()
	return nil, nil
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
		enc:          &mapObjEncoder{MapObjectEncoder: zapcore.NewMapObjectEncoder()},
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
	addFields(c.enc, fields)
	return c.clone()
}

func (c *ioCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *ioCore) Write(ent zapcore.Entry, fields []zapcore.Field) (err error) {
	enc, ok := c.enc.(*mapObjEncoder)
	if !ok {
		return errors.New("type assertion fail")
	}
	addFields(c.enc, fields)
	if err := c.writer.write(enc.Fields); err != nil {
		return err
	}
	if ent.Level > zapcore.ErrorLevel {
		err = c.Sync()
	}
	putEncoder(enc)
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
