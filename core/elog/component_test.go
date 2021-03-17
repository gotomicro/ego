package elog

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"

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

func BenchmarkAliWriter(b *testing.B) {
	b.Logf("Logging at a disabled level with some accumulated context.")
	aliLogger := newAliLogger()
	b.Run("ali", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				aliLogger.Info(getMessage(0))
			}
		})
	})
}
