package elog

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/gotomicro/ego/core/econf"
)

func TestRotateLogger(t *testing.T) {
	err := os.Setenv("EGO_DEBUG", "false")
	assert.NoError(t, err)
	conf := `
[default]
debug = false
level = "info"
enableAsync = false
`
	err = econf.LoadFromReader(strings.NewReader(conf), toml.Unmarshal)
	assert.NoError(t, err)
	logger := Load("default").Build()
	defer logger.Flush()
	logger.Error("1")

	childLogger := logger.With(String("prefix", "PREFIX"))
	childLogger.Error("2")
	defer childLogger.Flush()

	logger.Error("3")
	logger.With(String("prefix2", "PREFIX2"))
	logger.Error("4")
}

var messages = fakeMessages(1000)

func newFileLogger(path string) *Component {
	conf := `
[file]
level = "info"
name = "%s"
`
	conf = fmt.Sprintf(conf, path)
	var err error
	if err = econf.LoadFromReader(strings.NewReader(conf), toml.Unmarshal); err != nil {
		log.Println("load conf fail", err)
		return nil
	}
	log.Println("start to send logs to file")
	return Load("file").Build()
}

func newStderrLogger() *Component {
	conf := `
[stderr]
level = "info"
writer = "stderr"
`
	var err error
	if err = econf.LoadFromReader(strings.NewReader(conf), toml.Unmarshal); err != nil {
		log.Println("load conf fail", err)
		return nil
	}
	log.Println("start to send logs to stderr")
	return Load("stderr").Build()
}

func newStdoutLogger() *Component {
	conf := `
[stdout]
level = "info"
writer = "stdout"
`
	var err error
	if err = econf.LoadFromReader(strings.NewReader(conf), toml.Unmarshal); err != nil {
		log.Println("load conf fail", err)
		return nil
	}
	log.Println("start to send logs to stdout")
	return Load("stdout").Build()
}

func newAliLogger() *Component {
	conf := `
[ali]
level = "info"
enableAsync = false
writer = "ali"
flushBufferSize = 2097152     # flushBufferSize set to 2MB
flushBufferInterval = "2s"
aliEndpoint = "%s"            # your ali sls endpoint
aliAccessKeyID = "%s"         # your ali sls AK ID
aliAccessKeySecret = "%s"     # your ali sls AK Secret
aliProject = "%s"             # your ali sls project
aliLogstore = "%s"            # your ali logstore
aliApiBulkSize = 512          # al api bulk size
aliApiTimeout = "3s"          # ali api timeout
aliApiRetryCount = 3          # ali api retry
aliApiRetryWaitTime = "1s"    # ali api retry wait time
aliApiRetryMaxWaitTime = "3s" # ali api retry wait max wait time
aliApiMaxIdleConnsPerHost = 20
aliApiMaxIdleConns = 25
`
	conf = fmt.Sprintf(conf,
		os.Getenv("ALI_ENDPOINT"),
		os.Getenv("ALI_AK_ID"),
		os.Getenv("ALI_AK_SECRET"),
		os.Getenv("ALI_PROJECT"),
		os.Getenv("ALI_LOGSTORE"),
	)
	var err error
	if err = econf.LoadFromReader(strings.NewReader(conf), toml.Unmarshal); err != nil {
		log.Println("load conf fail", err)
		return nil
	}
	log.Println("start to send logs to ali sls")

	return Load("ali").Build()
}

func newZapLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	return logger
}

func fakeMessages(n int) []string {
	messages := make([]string, n)
	for i := range messages {
		messages[i] = fmt.Sprintf("Test logging, but use a somewhat realistic message length. (#%v)", i)
	}
	return messages
}

func getMessage(iter int) string {
	return messages[iter%1000]
}

func BenchmarkFileWriter(b *testing.B) {
	b.Logf("Logging at a disabled level with some accumulated context.")
	logger := newFileLogger("./benchmark-file-writer.log")
	b.Run("file", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
}

func BenchmarkStderrWriter(b *testing.B) {
	b.Logf("Logging at a disabled level with some accumulated context.")
	logger := newStderrLogger()
	b.Run("stderr\n", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
}

func BenchmarkStdoutWriter(b *testing.B) {
	b.Logf("Logging at a disabled level with some accumulated context.")
	logger := newStdoutLogger()
	b.Run("stdout\n", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
}

func BenchmarkZapWriter(b *testing.B) {
	b.Logf("Logging at a disabled level with some accumulated context.")
	logger := newZapLogger()
	b.Run("file", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
}

func BenchmarkAliWriter(b *testing.B) {
	b.Logf("Logging at a disabled level with some accumulated context.")
	logger := newAliLogger()
	b.Run("ali", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
}

func TestMultiLogger(t *testing.T) {
	os.RemoveAll("./logs")
	fileMultiLogger := newFileLogger("multi.log")
	fileMultiLoggers := []*Component{}
	for i := 0; i < 10; i++ {
		fileMultiLoggers = append(fileMultiLoggers, fileMultiLogger.With())
	}
	for i := 0; i < 10; i++ {
		for j := 0; j < 100000; j++ {
			fileMultiLoggers[i].Info(getMessage(0))
		}
	}

	for i := 0; i < 10; i++ {
		_ = fileMultiLoggers[i].Flush()
	}
	log.Println(`done--------------->`)
}

func BenchmarkMultiLogger(b *testing.B) {
	b.Logf("Logging at a single logger and multi child logger.")
	os.RemoveAll("./logs")
	fileMultiLogger := newFileLogger("child.log")
	fileMultiLoggers := []*Component{}
	for i := 0; i < 10; i++ {
		fileMultiLoggers = append(fileMultiLoggers, fileMultiLogger.With())
	}
	b.Run("child-file-logger", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				for i := 0; i < 10; i++ {
					fileMultiLoggers[i].Info(getMessage(0))
				}
			}
		})
	})

	fileSingleLoggers := []*Component{}
	for i := 0; i < 10; i++ {
		fileSingleLoggers = append(fileSingleLoggers, newFileLogger(fmt.Sprintf("independent-%d.log", i)))
	}
	b.Run("independent-file-logger", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				for i := 0; i < 10; i++ {
					fileSingleLoggers[i].Info(getMessage(0))
				}
			}
		})
	})

}

func Test_IsDebugMode(t *testing.T) {
	cmp := &Component{
		config: defaultConfig(),
	}
	cmp.config.Debug = true
	assert.True(t, cmp.IsDebugMode())
}

func TestSetLevel(t *testing.T) {
	logger := &Component{
		config: defaultConfig(),
	}
	logger.lv = &logger.config.al
	logger.SetLevel(zapcore.ErrorLevel)
	assert.Equal(t, "error", logger.lv.String())
}

func TestConfigDir(t *testing.T) {
	logger := &Component{
		config: defaultConfig(),
	}
	assert.Equal(t, "./logs", logger.ConfigDir())
}

func TestConfigName(t *testing.T) {
	logger := &Component{
		config: defaultConfig(),
	}
	assert.Equal(t, "default.log", logger.ConfigName())
}

func TestDebug(t *testing.T) {
	logger := DefaultContainer().Build(
		WithDebug(false),
		WithLevel("debug"),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	filePath := path.Join(logger.ConfigDir(), logger.ConfigName())
	logger.Debug("some")
	logged, err := os.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Contains(t, string(logged), `"lv":"debug"`)
	os.Remove(filePath)

	logger2 := DefaultContainer().Build(
		WithDebug(true),
		WithLevel("debug"),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	logger2.Debug("some2")
	logged2, err2 := os.ReadFile(filePath)
	assert.Nil(t, err2)
	assert.Contains(t, string(logged2), "\x1b[34mDEBUG\x1b")
	os.Remove(filePath)
}

func TestDebugW(t *testing.T) {
	logger := DefaultContainer().Build(
		WithDebug(false),
		WithLevel("debug"),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	filePath := path.Join(logger.ConfigDir(), logger.ConfigName())
	logger.Debugw("some")
	logged, err := os.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Contains(t, string(logged), `"lv":"debug"`)
	os.Remove(filePath)

	logger2 := DefaultContainer().Build(
		WithDebug(true),
		WithLevel("debug"),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	logger2.Debugw("some2")
	logged2, err2 := os.ReadFile(filePath)
	assert.Nil(t, err2)
	assert.Contains(t, string(logged2), "\x1b[34mDEBUG\x1b")
	os.Remove(filePath)
}

func TestDebugf(t *testing.T) {
	logger := DefaultContainer().Build(
		WithDebug(false),
		WithLevel("debug"),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	filePath := path.Join(logger.ConfigDir(), logger.ConfigName())
	logger.Debugf("hello,%s", "debug")
	logged, err := os.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Contains(t, string(logged), `"lv":"debug"`)
	assert.Contains(t, string(logged), `"msg":"hello,debug"`)
	os.Remove(filePath)
}

func TestInfo(t *testing.T) {
	logger := DefaultContainer().Build(
		WithDebug(false),
		WithLevel("info"),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	filePath := path.Join(logger.ConfigDir(), logger.ConfigName())
	logger.Info("some")
	logged, err := os.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Contains(t, string(logged), `"lv":"info"`)
	os.Remove(filePath)

	logger2 := DefaultContainer().Build(
		WithDebug(true),
		WithLevel("info"),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	logger2.Info("some2")
	logged2, err2 := os.ReadFile(filePath)
	assert.Nil(t, err2)
	assert.Contains(t, string(logged2), "\x1b[32mINFO\x1b")
	os.Remove(filePath)
}

func TestInfow(t *testing.T) {
	logger := DefaultContainer().Build(
		WithDebug(false),
		WithLevel("info"),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	filePath := path.Join(logger.ConfigDir(), logger.ConfigName())
	logger.Infow("some")
	logged, err := os.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Contains(t, string(logged), `"lv":"info"`)
	os.Remove(filePath)

	logger2 := DefaultContainer().Build(
		WithDebug(true),
		WithLevel("info"),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	logger2.Infow("some2")
	logged2, err2 := os.ReadFile(filePath)
	assert.Nil(t, err2)
	assert.Contains(t, string(logged2), "\x1b[32mINFO\x1b")
	os.Remove(filePath)
}

func TestInfof(t *testing.T) {
	logger := DefaultContainer().Build(
		WithDebug(false),
		WithLevel("info"),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	filePath := path.Join(logger.ConfigDir(), logger.ConfigName())
	logger.Infof("hello,%s", "info")
	logged, err := os.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Contains(t, string(logged), `"lv":"info"`)
	assert.Contains(t, string(logged), `"msg":"hello,info"`)
	os.Remove(filePath)
}

func TestWarn(t *testing.T) {
	logger := DefaultContainer().Build(
		WithDebug(false),
		WithLevel("warn"),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	filePath := path.Join(logger.ConfigDir(), logger.ConfigName())
	logger.Warn("some")
	logged, err := os.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Contains(t, string(logged), `"lv":"warn"`)
	os.Remove(filePath)

	logger2 := DefaultContainer().Build(
		WithDebug(true),
		WithLevel("warn"),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	logger2.Warn("some2")
	logged2, err2 := os.ReadFile(filePath)
	assert.Nil(t, err2)
	assert.Contains(t, string(logged2), "\x1b[33mWARN\x1b")
	os.Remove(filePath)
}

func TestWarnw(t *testing.T) {
	logger := DefaultContainer().Build(
		WithDebug(false),
		WithLevel("warn"),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	filePath := path.Join(logger.ConfigDir(), logger.ConfigName())
	logger.Warnw("some")
	logged, err := os.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Contains(t, string(logged), `"lv":"warn"`)
	os.Remove(filePath)

	logger2 := DefaultContainer().Build(
		WithDebug(true),
		WithLevel("warn"),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	logger2.Warnw("some2")
	logged2, err2 := os.ReadFile(filePath)
	assert.Nil(t, err2)
	assert.Contains(t, string(logged2), "\x1b[33mWARN\x1b")
	os.Remove(filePath)
}

func TestWarnf(t *testing.T) {
	logger := DefaultContainer().Build(
		WithDebug(false),
		WithLevel("warn"),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	filePath := path.Join(logger.ConfigDir(), logger.ConfigName())
	logger.Warnf("hello,%s", "warn")
	logged, err := os.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Contains(t, string(logged), `"lv":"warn"`)
	assert.Contains(t, string(logged), `"msg":"hello,warn"`)
	os.Remove(filePath)
}

func TestError(t *testing.T) {
	logger := DefaultContainer().Build(
		WithDebug(false),
		WithLevel("error"),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	filePath := path.Join(logger.ConfigDir(), logger.ConfigName())
	logger.Error("some")
	logged, err := os.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Contains(t, string(logged), `"lv":"error"`)
	os.Remove(filePath)

	logger2 := DefaultContainer().Build(
		WithDebug(true),
		WithLevel("error"),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	logger2.Error("some2")
	logged2, err2 := os.ReadFile(filePath)
	assert.Nil(t, err2)
	assert.Contains(t, string(logged2), "\x1b[31mERROR\x1b")
	os.Remove(filePath)
}

func TestErrorw(t *testing.T) {
	logger := DefaultContainer().Build(
		WithDebug(false),
		WithLevel("error"),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	filePath := path.Join(logger.ConfigDir(), logger.ConfigName())
	logger.Errorw("some")
	logged, err := os.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Contains(t, string(logged), `"lv":"error"`)
	os.Remove(filePath)

	logger2 := DefaultContainer().Build(
		WithDebug(true),
		WithLevel("error"),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	logger2.Errorw("some2")
	logged2, err2 := os.ReadFile(filePath)
	assert.Nil(t, err2)
	assert.Contains(t, string(logged2), "\x1b[31mERROR\x1b")
	os.Remove(filePath)
}

func TestErrorf(t *testing.T) {
	logger := DefaultContainer().Build(
		WithDebug(false),
		WithLevel("error"),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	filePath := path.Join(logger.ConfigDir(), logger.ConfigName())
	logger.Errorf("hello,%s", "error")
	logged, err := os.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Contains(t, string(logged), `"lv":"error"`)
	assert.Contains(t, string(logged), `"msg":"hello,error"`)
	os.Remove(filePath)
}

func TestPanic(t *testing.T) {
	var logger *Component
	assert.Panics(t, func() {
		logger = DefaultContainer().Build(
			WithDebug(false),
			WithLevel("error"),
			WithEnableAddCaller(true),
			WithEnableAsync(false),
		)
		logger.Panic("some")

	})
	filePath := path.Join(logger.ConfigDir(), logger.ConfigName())
	logged, err := os.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Contains(t, string(logged), `"lv":"panic"`)
	os.Remove(filePath)

	var logger2 *Component
	assert.Panics(t, func() {
		logger2 = DefaultContainer().Build(
			WithDebug(true),
			WithLevel("error"),
			WithEnableAddCaller(true),
			WithEnableAsync(false),
		)
		logger2.Panic("some")

	})
	filePath2 := path.Join(logger2.ConfigDir(), logger2.ConfigName())
	logged2, err2 := os.ReadFile(filePath2)
	assert.Nil(t, err2)
	assert.Contains(t, string(logged2), "\x1b[31mPANIC\x1b")
	os.Remove(filePath)
}

func TestPanicw(t *testing.T) {
	var logger *Component
	assert.Panics(t, func() {
		logger = DefaultContainer().Build(
			WithDebug(false),
			WithLevel("error"),
			WithEnableAddCaller(true),
			WithEnableAsync(false),
		)
		logger.Panicw("some")

	})
	filePath := path.Join(logger.ConfigDir(), logger.ConfigName())
	logged, err := os.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Contains(t, string(logged), `"lv":"panic"`)
	os.Remove(filePath)

	var logger2 *Component
	assert.Panics(t, func() {
		logger2 = DefaultContainer().Build(
			WithDebug(true),
			WithLevel("error"),
			WithEnableAddCaller(true),
			WithEnableAsync(false),
		)
		logger2.Panicw("some")

	})
	filePath2 := path.Join(logger2.ConfigDir(), logger2.ConfigName())
	logged2, err2 := os.ReadFile(filePath2)
	assert.Nil(t, err2)
	assert.Contains(t, string(logged2), "\x1b[31mPANIC\x1b")
	os.Remove(filePath)
}

func TestPanicf(t *testing.T) {
	var logger *Component
	assert.Panics(t, func() {
		logger = DefaultContainer().Build(
			WithDebug(false),
			WithLevel("panic"),
			WithEnableAddCaller(true),
			WithEnableAsync(false),
		)
		logger.Panicf("hello,%s", "panic")
	})
	filePath := path.Join(logger.ConfigDir(), logger.ConfigName())
	logged, err := os.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Contains(t, string(logged), `"lv":"panic"`)
	assert.Contains(t, string(logged), `"msg":"hello,panic"`)
	os.Remove(filePath)
}

func TestDPanic(t *testing.T) {
	logger := DefaultContainer().Build(
		WithDebug(false),
		WithLevel("dpanic"),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	filePath := path.Join(logger.ConfigDir(), logger.ConfigName())
	logger.DPanic("some")
	logged, err := os.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Contains(t, string(logged), `"lv":"dpanic"`)
	os.Remove(filePath)

	logger2 := DefaultContainer().Build(
		WithDebug(true),
		WithLevel("error"),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	logger2.DPanic("some2")
	logged2, err2 := os.ReadFile(filePath)
	assert.Nil(t, err2)
	assert.Contains(t, string(logged2), "\x1b[31mDPANIC\x1b")
	os.Remove(filePath)
}

func TestDPanicw(t *testing.T) {
	logger := DefaultContainer().Build(
		WithDebug(false),
		WithLevel("dpanic"),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	filePath := path.Join(logger.ConfigDir(), logger.ConfigName())
	logger.DPanicw("some")
	logged, err := os.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Contains(t, string(logged), `"lv":"dpanic"`)
	os.Remove(filePath)

	logger2 := DefaultContainer().Build(
		WithDebug(true),
		WithLevel("error"),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	logger2.DPanicw("some2")
	logged2, err2 := os.ReadFile(filePath)
	assert.Nil(t, err2)
	assert.Contains(t, string(logged2), "\x1b[31mDPANIC\x1b")
	os.Remove(filePath)
}

func TestDPanicf(t *testing.T) {
	logger := DefaultContainer().Build(
		WithDebug(false),
		WithLevel("dpanic"),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	logger.DPanicf("hello,%s", "dpanic")
	filePath := path.Join(logger.ConfigDir(), logger.ConfigName())
	logged, err := os.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Contains(t, string(logged), `"lv":"dpanic"`)
	assert.Contains(t, string(logged), `"msg":"hello,dpanic"`)
	os.Remove(filePath)
}

func TestWithZapConfig(t *testing.T) {
	testConfig := zap.NewDevelopmentEncoderConfig()
	logger := DefaultContainer().Build(
		WithDebug(false),
		WithEncoderConfig(&testConfig),
		WithEnableAsync(false),
	)
	logger.Info("hello")
	filePath := path.Join(logger.ConfigDir(), logger.ConfigName())
	logged, err := os.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Contains(t, string(logged), `"L":"INFO"`)
	os.Remove(filePath)
}

func TestConfig_AddCaller(t *testing.T) {
	logger := DefaultContainer().Build(
		WithDebug(true),
		WithEnableAddCaller(true),
		WithEnableAsync(false),
	)
	logger.Info("hello")
	filePath := path.Join(logger.ConfigDir(), logger.ConfigName())
	logged, err := os.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Contains(t, string(logged), `elog/component_test.go:`)
	os.Remove(filePath)
}
