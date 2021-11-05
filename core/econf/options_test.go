package econf_test

import (
	"io/ioutil"
	"path"
	"sync"
	"testing"

	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/econf/manager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithTagName(t *testing.T) {
	watchDir, err := ioutil.TempDir("", "")
	require.Nil(t, err)
	configFile := path.Join(watchDir, "config.yaml")
	err = ioutil.WriteFile(configFile, []byte("foo: bar\n"), 0640)
	require.Nil(t, err)
	// defer func() {
	// 	os.RemoveAll(configFile)
	// }()
	v := econf.New()
	provider, parser, tag, err := manager.NewDataSource(configFile, true)
	require.Nil(t, err)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	v.OnChange(func(configuration *econf.Configuration) {
		t.Logf("config file changed")
		wg.Done()
	})
	err = v.LoadFromDataSource(provider, parser, econf.WithTagName(tag), econf.WithWeaklyTypedInput(true))
	require.Nil(t, err)
	assert.Equal(t, "yaml", econf.GetOptionTagName())
	assert.Equal(t, true, econf.GetOptionWeaklyTypedInput())
}
