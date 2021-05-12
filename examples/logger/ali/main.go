package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
)

func main() {
	var err error
	conf := `
[ali]
level = "info"
enableAsync = false
writer = "ali"
flushBufferSize = 2097152     # flushBufferSize set to 2MB
aliEndpoint = "%s"            # your ali sls endpoint
aliAccessKeyID = "%s"         # your ali sls AK ID
aliAccessKeySecret = "%s"     # your ali sls AK Secret
aliProject = "%s"             # your ali sls project
aliLogstore = "%s"            # your ali logstore
aliApiBulkSize = 128          # al api bulk size
aliApiTimeout = "3s"          # ali api timeout
aliApiRetryCount = 3          # ali api retry
aliApiRetryWaitTime = "1s"    # ali api retry wait time
aliApiRetryMaxWaitTime = "3s" # ali api retry wait max wait time
`
	conf = fmt.Sprintf(conf,
		os.Getenv("ALI_ENDPOINT"),
		os.Getenv("ALI_AK_ID"),
		os.Getenv("ALI_AK_SECRET"),
		os.Getenv("ALI_PROJECT"),
		os.Getenv("ALI_LOGSTORE"),
	)
	if err = econf.LoadFromReader(strings.NewReader(conf), toml.Unmarshal); err != nil {
		log.Println("load conf fail", err)
		return
	}
	log.Println("start to send logs to ali sls")

	logger := elog.Load("ali").Build()
	logger.Info("aaaaaaaaa", elog.Any("map", map[string]interface{}{"aaa": "AAA", "bbb": "BBB"}), elog.Any("slice", []string{"ccc", "ddd"}))
	logger.Info("aaaaaaaaa", elog.Any("map", map[string]interface{}{"ccc": "CCC"}))
	defer logger.Flush()

	childLogger := logger.With(elog.String("prefix", "PREFIX"))
	defer childLogger.Flush()

	childLogger.Error("childLogger1", elog.String("name", "lee"), elog.Int("age", 18))
	childLogger.Error("childLogger2", elog.String("name", "lee"), elog.Int("age", 19))

	logger.Error("parentLogger1")
	logger.With(elog.String("prefix2", "PREFIX2"))
	logger.Error("parentLogger2")

	childLogger.Error("childLogger3", elog.String("name", "lee"), elog.Int("age", 20))

	log.Println("send logs to ali sls success")
}
