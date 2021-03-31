package econf_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gotomicro/ego/core/econf"
	_ "github.com/gotomicro/ego/core/econf/file"
	"github.com/gotomicro/ego/core/econf/manager"
)

func TestWatchFile(t *testing.T) {
	t.Run("file content changed", func(t *testing.T) {
		// given a `config.yaml` file being watched
		v, configFile, cleanup, wg := newWithConfigFile(t)
		defer cleanup()
		_, err := os.Stat(configFile)
		require.NoError(t, err)
		t.Logf("test config file: %s\n", configFile)
		// when overwriting the file and waiting for the custom change notification handler to be triggered
		err = ioutil.WriteFile(configFile, []byte("foo: baz\n"), 0640)
		require.Nil(t, err)
		wg.Wait()
		// then the config value should have changed
		assert.Equal(t, "baz", v.Get("foo"))
	})

	t.Run("link to real file changed (Kubernetes)", func(t *testing.T) {
		v, watchDir, _, _, wg := newWithSymlinkedConfigFile(t)
		// defer cleanup()
		// when link to another `config.yaml` file
		dataDir2 := path.Join(watchDir, "data2")
		err := os.Mkdir(dataDir2, 0777)
		require.Nil(t, err)
		configFile2 := path.Join(dataDir2, "config.yaml")
		err = ioutil.WriteFile(configFile2, []byte("foo: baz\n"), 0640)
		require.Nil(t, err)
		// change the symlink using the `ln -sfn` command
		err = exec.Command("ln", "-sfn", dataDir2, path.Join(watchDir, "data")).Run()
		require.Nil(t, err)
		wg.Wait()
		// then
		require.Nil(t, err)
		assert.Equal(t, "baz", v.Get("foo"))
	})
}

func newWithConfigFile(t *testing.T) (*econf.Configuration, string, func(), *sync.WaitGroup) {
	watchDir, err := ioutil.TempDir("", "")
	require.Nil(t, err)
	configFile := path.Join(watchDir, "config.yaml")
	err = ioutil.WriteFile(configFile, []byte("foo: bar\n"), 0640)
	require.Nil(t, err)
	cleanup := func() {
		os.RemoveAll(watchDir)
	}
	v := econf.New()
	provider, parser, tag, err := manager.NewDataSource(configFile, true)
	assert.Nil(t, err)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	v.OnChange(func(configuration *econf.Configuration) {
		t.Logf("config file changed")
		wg.Done()
	})
	err = v.LoadFromDataSource(provider, parser, econf.WithTagName(tag))
	assert.Nil(t, err)
	require.Equal(t, "bar", v.Get("foo"))
	return v, configFile, cleanup, wg
}

func newWithSymlinkedConfigFile(t *testing.T) (*econf.Configuration, string, string, func(), *sync.WaitGroup) {
	watchDir, err := ioutil.TempDir("", "")
	require.Nil(t, err)
	dataDir1 := path.Join(watchDir, "data1")
	err = os.Mkdir(dataDir1, 0777)
	require.Nil(t, err)
	realConfigFile := path.Join(dataDir1, "config.yaml")
	t.Logf("Real config file location: %s\n", realConfigFile)
	err = ioutil.WriteFile(realConfigFile, []byte("foo: bar\n"), 0640)
	require.Nil(t, err)
	cleanup := func() {
		os.RemoveAll(watchDir)
	}
	// now, symlink the tm `data1` dir to `data` in the baseDir
	err = os.Symlink(dataDir1, path.Join(watchDir, "data"))
	require.Nil(t, err)

	// and link the `<watchdir>/datadir1/config.yaml` to `<watchdir>/config.yaml`
	configFile := path.Join(watchDir, "config.yaml")
	err = os.Symlink(path.Join(watchDir, "data", "config.yaml"), configFile)
	require.Nil(t, err)

	t.Logf("Config file location: %s\n", path.Join(watchDir, "config.yaml"))

	v := econf.New()
	provider, parser, tag, err := manager.NewDataSource(configFile, true)
	require.Nil(t, err)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	v.OnChange(func(configuration *econf.Configuration) {
		t.Logf("config file changed")
		wg.Done()
	})
	err = v.LoadFromDataSource(provider, parser, econf.WithTagName(tag))
	require.Nil(t, err)
	require.Equal(t, "bar", v.Get("foo"))
	return v, watchDir, configFile, cleanup, wg
}
