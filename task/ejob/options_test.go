package ejob

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithName(t *testing.T) {
	container := DefaultContainer()
	container.Build(WithName("test"))
	assert.Equal(t, "test", container.config.Name)
}

func TestWithStartFunc(t *testing.T) {
	container := DefaultContainer()
	fc := func(ctx Context) error {
		return nil
	}
	one := fmt.Sprintf("%p", fc)
	container.Build(WithStartFunc(fc))
	two := fmt.Sprintf("%p", container.config.startFunc)
	assert.Equal(t, one, two)
}
