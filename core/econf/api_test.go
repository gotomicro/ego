package econf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApi(t *testing.T) {
	fn := func(configuration *Configuration) {}
	OnChange(fn)
	Sub("")
	Reset()
	Traverse("")
	RawConfig()
	Debug("")
	Get("")
	assert.NoError(t, nil)
}
