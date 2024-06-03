package file

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/econf"
)

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
