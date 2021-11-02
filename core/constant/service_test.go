package constant

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceKind_String(t *testing.T) {
	assert.Equal(t, ServiceProvider.String(), "providers")
	assert.Equal(t, ServiceGovernor.String(), "governors")
	assert.Equal(t, ServiceConsumer.String(), "consumers")
}
