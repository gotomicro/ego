package econf

import (
	"os"
	"path"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
)

func TestWithTagName(t *testing.T) {
	watchDir := os.TempDir()
	configFile := path.Join(watchDir, "config.toml")
	err := os.WriteFile(configFile, []byte(`foo= "baz"`), 0640)
	assert.NoError(t, err)
	defer func() {
		os.RemoveAll(configFile)
	}()
	v := New()
	provider := newMockDataSource(configFile, true)

	err1 := v.LoadFromDataSource(provider, toml.Unmarshal, WithTagName("toml"), WithWeaklyTypedInput(true))
	assert.NoError(t, err1)
	assert.Equal(t, "toml", GetOptionTagName())
	assert.Equal(t, true, GetOptionWeaklyTypedInput())

	err2 := v.LoadFromDataSource(provider, toml.Unmarshal, WithSquash(true))
	assert.NoError(t, err2)
}
