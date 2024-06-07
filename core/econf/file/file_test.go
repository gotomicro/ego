package file

import (
	"os"
	"os/exec"
	"path"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/econf"
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
		assert.NoError(t, err)
		t.Logf("test config file: %s\n", configFile)
		// when overwriting the file and waiting for the custom change notification handler to be triggered
		err1 := os.WriteFile(configFile, []byte("foo: baz\n"), 0640)
		assert.NoError(t, err1)
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
		assert.NoError(t, err)
		configFile2 := path.Join(dataDir2, "config.yaml")
		err1 := os.WriteFile(configFile2, []byte("foo: baz\n"), 0640)
		assert.NoError(t, err1)
		// change the symlink using the `ln -sfn` command
		err3 := exec.Command("ln", "-sfn", dataDir2, path.Join(watchDir, "data")).Run()
		assert.NoError(t, err3)
		wg.Wait()
		assert.Equal(t, "baz", v.Get("foo"))
	})
}

func newWithConfigFile(t *testing.T) (*econf.Configuration, string, func(), *sync.WaitGroup) {
	watchDir := os.TempDir()
	configFile := path.Join(watchDir, "config.yaml")
	err := os.WriteFile(configFile, []byte("foo: bar\n"), 0640)
	assert.NoError(t, err)
	cleanup := func() {
		os.RemoveAll(watchDir)
	}
	v := econf.New()
	provider, parser, tag, err1 := manager.NewDataSource(configFile, true)
	assert.NoError(t, err1)

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
	err2 := v.LoadFromDataSource(provider, parser, econf.WithTagName(tag))
	assert.NoError(t, err2)
	assert.Equal(t, "bar", v.Get("foo"))
	return v, configFile, cleanup, wg
}

func newWithSymlinkedConfigFile(t *testing.T) (*econf.Configuration, string, string, func(), *sync.WaitGroup) {
	watchDir := os.TempDir()
	dataDir1 := path.Join(watchDir, "data1")
	err := os.Mkdir(dataDir1, 0777)
	assert.NoError(t, err)
	realConfigFile := path.Join(dataDir1, "config.yaml")
	t.Logf("Real config file location: %s\n", realConfigFile)
	err1 := os.WriteFile(realConfigFile, []byte("foo: bar\n"), 0640)
	assert.NoError(t, err1)
	cleanup := func() {
		os.RemoveAll(watchDir)
	}
	// now, symlink the tm `data1` dir to `data` in the baseDir
	err2 := os.Symlink(dataDir1, path.Join(watchDir, "data"))
	assert.NoError(t, err2)

	// and link the `<watchdir>/datadir1/config.yaml` to `<watchdir>/config.yaml`
	configFile := path.Join(watchDir, "config.yaml")
	err3 := os.Symlink(path.Join(watchDir, "data", "config.yaml"), configFile)
	assert.NoError(t, err3)

	t.Logf("Config file location: %s\n", path.Join(watchDir, "config.yaml"))

	v := econf.New()
	provider, parser, tag, err4 := manager.NewDataSource(configFile, true)
	assert.NoError(t, err4)

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
	err5 := v.LoadFromDataSource(provider, parser, econf.WithTagName(tag))
	assert.NoError(t, err5)
	assert.Equal(t, "bar", v.Get("foo"))
	return v, watchDir, configFile, cleanup, wg
}

func TestParse(t *testing.T) {
	cases := []struct {
		in       string
		expected econf.ConfigType
	}{
		{in: "./conf_test/conf.json", expected: "json"},
		{in: "./conf_test/conf.toml", expected: "toml"},
		{in: "./conf_test/conf.yaml", expected: "yaml"},
		{in: "./conf_test/cfg.json", expected: "json"},
	}

	for _, c := range cases {
		fp := &fileDataSource{}
		out := fp.Parse(c.in, true)
		assert.Equal(t, c.expected, out)
	}
}

func TestReadConfig(t *testing.T) {
	cases := []struct {
		in       string
		expected []byte
	}{
		{in: "./conf_test/conf.json", expected: []byte(`{
    "test1": "hello",
    "test2": "world"
}`)},
		{in: "./conf_test/conf.toml", expected: []byte(`[test]
name1 = "hello"
name2 = "world"`)},
		{in: "./conf_test/conf.yaml", expected: []byte(`Test:
  hello: world`)},
		{in: "./conf_test/cfg.json", expected: []byte(``)},
	}

	for _, c := range cases {
		fp := &fileDataSource{path: c.in}
		out, _ := fp.ReadConfig()
		assert.Equal(t, c.expected, out)
	}
}

func TestClose(t *testing.T) {
	c := make(chan struct{})
	fp := &fileDataSource{changed: c}
	out := fp.Close()
	assert.Equal(t, nil, out)
}

func TestIsConfigChanged(t *testing.T) {
	c := make(chan struct{})
	exp := (<-chan struct{})(c)
	fp := &fileDataSource{changed: c}
	out := fp.IsConfigChanged()
	assert.Equal(t, exp, out)
}
