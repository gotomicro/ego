package econf_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gotomicro/ego/core/econf"
	_ "github.com/gotomicro/ego/core/econf/file"
	"github.com/gotomicro/ego/core/econf/manager"
)

func TestWatchFile(t *testing.T) {
	if runtime.GOOS == "linux" {
		t.Skip("Skip test on Linux ...")
	}
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
		// skip if not executed on Linux
		if runtime.GOOS != "linux" {
			t.Skipf("Skipping test as symlink replacements don't work on non-linux environment...")
		}

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
	wg.Add(2)
	var init int64
	v.OnChange(func(configuration *econf.Configuration) {
		if atomic.CompareAndSwapInt64(&init, 0, 1) {
			t.Logf("config init")
		} else {
			t.Logf("config file changed")
		}
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
	wg.Add(2)
	var init int64
	v.OnChange(func(configuration *econf.Configuration) {
		if atomic.CompareAndSwapInt64(&init, 0, 1) {
			t.Logf("config init")
		} else {
			t.Logf("config file changed")
		}
		wg.Done()
	})
	err = v.LoadFromDataSource(provider, parser, econf.WithTagName(tag))
	require.Nil(t, err)
	require.Equal(t, "bar", v.Get("foo"))
	return v, watchDir, configFile, cleanup, wg
}

func TestSetKeyDelim(t *testing.T) {
	v := econf.New()
	v.SetKeyDelim("")
	assert.NoError(t, nil)

	err := v.WriteConfig()
	assert.NoError(t, err)
}

func TestGet(t *testing.T) {
	econf.GetString("")
	econf.GetInt("")
	econf.GetInt64("")
	econf.GetFloat64("")
	econf.GetTime("")
	econf.GetDuration("")
	econf.GetStringSlice("")
	econf.GetSlice("")
	econf.GetStringMap("")
	econf.GetStringMapString("")
	econf.GetStringMapStringSlice("")
	assert.NoError(t, nil)
	assert.Equal(t, false, econf.GetBool(""))
}
