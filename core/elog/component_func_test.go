package elog

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func Test_timeEncoder(t *testing.T) {
	te, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 10:00:00")
	enc := &sliceArrayEncoder{}
	timeEncoder(te, enc)
	assert.Equal(t, int64(1609495200), enc.elems[0].(int64))
}

//func Test_timeDebugEncoder(t *testing.T) {
//	te, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 10:00:00")
//	enc := &sliceArrayEncoder{}
//	timeDebugEncoder(te, enc)
//	assert.Equal(t, "2021-01-01 10:00:00", enc.elems[0].(string))
//}

func Test_debugEncodeLevel1(t *testing.T) {
	type args struct {
		lv  zapcore.Level
		enc *sliceArrayEncoder
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				lv:  zapcore.DebugLevel,
				enc: &sliceArrayEncoder{},
			},
			want: "\x1b[34mDEBUG\x1b[0m",
		},
		{
			args: args{
				lv:  zapcore.InfoLevel,
				enc: &sliceArrayEncoder{},
			},
			want: "\x1b[32mINFO\x1b[0m",
		},
		{
			args: args{
				lv:  zapcore.WarnLevel,
				enc: &sliceArrayEncoder{},
			},
			want: "\x1b[33mWARN\x1b[0m",
		},
		{
			args: args{
				lv:  zapcore.ErrorLevel,
				enc: &sliceArrayEncoder{},
			},
			want: "\x1b[31mERROR\x1b[0m",
		},
		{
			args: args{
				lv:  zapcore.DPanicLevel,
				enc: &sliceArrayEncoder{},
			},
			want: "\x1b[31mDPANIC\x1b[0m",
		},
		{
			args: args{
				lv:  zapcore.PanicLevel,
				enc: &sliceArrayEncoder{},
			},
			want: "\x1b[31mPANIC\x1b[0m",
		},
		{
			args: args{
				lv:  zapcore.FatalLevel,
				enc: &sliceArrayEncoder{},
			},
			want: "\x1b[31mFATAL\x1b[0m",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debugEncodeLevel(tt.args.lv, tt.args.enc)
			assert.Equal(t, tt.want, tt.args.enc.elems[0].(string))
		})
	}
}

func Test_defaultDebugConfig(t *testing.T) {
	cfg := defaultDebugConfig()
	assert.Equal(t, "ts", cfg.TimeKey)
	assert.Equal(t, "lv", cfg.LevelKey)
	assert.Equal(t, "logger", cfg.NameKey)
	assert.Equal(t, "caller", cfg.CallerKey)
	assert.Equal(t, "msg", cfg.MessageKey)
	assert.Equal(t, "stack", cfg.StacktraceKey)
	assert.Equal(t, zapcore.DefaultLineEnding, cfg.LineEnding)
}

func Test_normalizeMessage(t *testing.T) {
	assert.Equal(t, fmt.Sprintf("%-32s", "hello"), normalizeMessage("hello"))
}
