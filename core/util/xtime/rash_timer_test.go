package xtime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimer(t *testing.T) {
	var testWheel = NewRashTimer(1 * time.Millisecond)
	t1 := testWheel.NewTimer(500 * time.Millisecond)

	before := time.Now()
	<-t1.C
	after := time.Now()

	assert.True(t, after.Sub(before) < time.Millisecond*600)
}
