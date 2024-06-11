package econf

import (
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
)

type mockDataSource struct {
	path        string
	enableWatch bool
	changed     chan struct{}
}

func (m *mockDataSource) Parse(path string, watch bool) ConfigType {
	_, err := url.Parse(path)
	if err != nil {
		log.Panicf("Error: %v\n", err)
	}
	return ConfigTypeToml
}

func (m *mockDataSource) ReadConfig() ([]byte, error) {
	return os.ReadFile(m.path)
}

func (m *mockDataSource) IsConfigChanged() <-chan struct{} {
	changed := make(chan struct{})
	if content, err := m.ReadConfig(); err == nil {
		// 创建临时的配置对象
		tempC := &Configuration{}
		if err := toml.Unmarshal(content, tempC); err == nil {
			tempC.mu.RLock()
			defer tempC.mu.RUnlock()
			for _, change := range tempC.onChanges {
				change(tempC)
			}
			close(changed)
		}
	}
	return changed
}

func (m *mockDataSource) Close() error {
	close(m.changed)
	return nil
}

var _ DataSource = (*mockDataSource)(nil)

func newMockDataSource(Addr string, watch bool) DataSource {
	return &mockDataSource{path: Addr, enableWatch: watch}
}

func TestWatchFile(t *testing.T) {
	if runtime.GOOS == "linux" {
		t.Skip("Skip test on Linux ...")
	}
	t.Run("file content changed", func(t *testing.T) {
		// given a `config.toml` file being watched
		v, configFile, cleanup, wg := newWithConfigFile(t)
		defer cleanup()
		_, err := os.Stat(configFile)
		assert.NoError(t, err)
		t.Logf("test config file: %s\n", configFile)
		// when overwriting the file and waiting for the custom change notification handler to be triggered
		err1 := os.WriteFile(configFile, []byte(`foo= "baz"`), 0640)
		assert.NoError(t, err1)
		// wg.Wait()
		wg.Done()
		// then the config value should have changed
		assert.Equal(t, "baz", v.Get("foo"))
	})

	t.Run("link to real file changed (Kubernetes)", func(t *testing.T) {
		// skip if not executed on Linux
		t.Skipf("Skipping test as symlink replacements don't work on non-linux environment...")

		v, watchDir, _, _, wg := newWithSymlinkedConfigFile(t)
		// defer cleanup()
		// when link to another `config.toml` file
		dataDir2 := path.Join(watchDir, "data2")
		err := os.Mkdir(dataDir2, 0777)
		assert.NoError(t, err)
		configFile2 := path.Join(dataDir2, "config.toml")
		err1 := os.WriteFile(configFile2, []byte(`foo= "baz"`), 0640)
		assert.NoError(t, err1)
		// change the symlink using the `ln -sfn` command
		err2 := exec.Command("ln", "-sfn", dataDir2, path.Join(watchDir, "data")).Run()
		assert.NoError(t, err2)
		wg.Wait()
		// then
		assert.Equal(t, "baz", v.Get("foo"))
	})
}

func newWithConfigFile(t *testing.T) (*Configuration, string, func(), *sync.WaitGroup) {
	var watchDir = os.TempDir()
	configFile := path.Join(watchDir, "config.toml")
	err := os.WriteFile(configFile, []byte(`foo= "baz"`), 0640)
	assert.NoError(t, err)
	content, err := os.ReadFile(configFile)
	if err != nil {
		log.Panicf("Error: %v\n", err)
	}
	t.Logf("Content of configFile: %v\n", string(content))

	cleanup := func() {
		os.RemoveAll(watchDir)
	}

	v := New()
	provider := newMockDataSource(configFile, true)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	var init int64
	v.OnChange(func(configuration *Configuration) {
		if atomic.CompareAndSwapInt64(&init, 0, 1) {
			t.Logf("config init")
		} else {
			t.Logf("config file changed")
		}
		wg.Done()
	})

	err1 := v.LoadFromDataSource(provider, toml.Unmarshal)
	assert.NoError(t, err1)
	assert.Equal(t, "baz", v.Get("foo"))
	return v, configFile, cleanup, wg
}

func newWithSymlinkedConfigFile(t *testing.T) (*Configuration, string, string, func(), *sync.WaitGroup) {
	watchDir := os.TempDir()
	dataDir1 := path.Join(watchDir, "data1")
	err := os.Mkdir(dataDir1, 0777)
	assert.NoError(t, err)
	realConfigFile := path.Join(dataDir1, "config.toml")
	t.Logf("Real config file location: %s\n", realConfigFile)
	err1 := os.WriteFile(realConfigFile, []byte(`foo= "baz"`), 0640)
	assert.NoError(t, err1)
	cleanup := func() {
		os.RemoveAll(watchDir)
	}
	// now, symlink the tm `data1` dir to `data` in the baseDir
	err2 := os.Symlink(dataDir1, path.Join(watchDir, "data"))
	assert.NoError(t, err2)

	// and link the `<watchdir>/datadir1/config.toml` to `<watchdir>/config.toml`
	configFile := path.Join(watchDir, "config.toml")
	err3 := os.Symlink(path.Join(watchDir, "data", "config.toml"), configFile)
	assert.NoError(t, err3)
	t.Logf("Config file location: %s\n", path.Join(watchDir, "config.toml"))

	v := New()
	provider := newMockDataSource(configFile, true)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	var init int64
	v.OnChange(func(configuration *Configuration) {
		if atomic.CompareAndSwapInt64(&init, 0, 1) {
			t.Logf("config init")
		} else {
			t.Logf("config file changed")
		}
		wg.Done()
	})
	err4 := v.LoadFromDataSource(provider, toml.Unmarshal)
	assert.NoError(t, err4)
	assert.Equal(t, "bar", v.Get("foo"))
	return v, watchDir, configFile, cleanup, wg
}
