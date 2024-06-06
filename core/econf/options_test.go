package econf

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithTagName(t *testing.T) {
	watchDir, err := ioutil.TempDir("", "")
	require.Nil(t, err)
	configFile := path.Join(watchDir, "config.toml")
	err = os.WriteFile(configFile, []byte(`foo= "baz"`), 0640)
	require.Nil(t, err)
	defer func() {
		os.RemoveAll(configFile)
	}()
	v := New()
	provider := newMockDataSource(configFile, true)

	err = v.LoadFromDataSource(provider, toml.Unmarshal, WithTagName("toml"), WithWeaklyTypedInput(true))
	require.Nil(t, err)
	assert.Equal(t, "toml", GetOptionTagName())
	assert.Equal(t, true, GetOptionWeaklyTypedInput())

	err = v.LoadFromDataSource(provider, toml.Unmarshal, WithSquash(true))
	require.Nil(t, err)
}
