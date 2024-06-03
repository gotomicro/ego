package econf

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetKeyDelim(t *testing.T) {
	c := &Configuration{}
	c.SetKeyDelim(";")
	if c.keyDelim != ";" {
		t.Errorf("Expected key delimiter to be ';', but got %s", c.keyDelim)
	}
	err := c.WriteConfig()
	assert.NoError(t, err)
}

func TestSub(t *testing.T) {
	c := &Configuration{
		keyDelim: defaultKeyDelim,
		override: map[string]interface{}{
			"key1": "hello",
			"key2": "world",
		},
		keyMap: &sync.Map{},
	}
	t.Run("When key is empty string", func(t *testing.T) {
		out := c.Sub("")
		in := &Configuration{
			keyDelim: defaultKeyDelim,
			override: map[string]interface{}{},
			keyMap:   &sync.Map{},
		}
		assert.Equal(t, in, out)
	})
}

func TestSet(t *testing.T) {
	v := New()
	key := "a.b.c"
	val := 42
	assert.NoError(t, v.Set(key, val))
	assert.Equal(t, "42", v.GetString(key))
	assert.Equal(t, "", GetString(key))
	assert.Equal(t, 42, v.GetInt(key))
	assert.Equal(t, int64(42), v.GetInt64(key))
	assert.Equal(t, float64(42), v.GetFloat64(key))
	assert.Equal(t, []string{"42"}, v.GetStringSlice(key))
}
