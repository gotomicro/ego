package ecron

import (
	"sync/atomic"
	"time"

	"github.com/gotomicro/ego/core/standard"
)

type immediatelyScheduler struct {
	Schedule
	initOnce uint32
}

// Next ...
func (is *immediatelyScheduler) Next(curr time.Time) (next time.Time) {
	if atomic.CompareAndSwapUint32(&is.initOnce, 0, 1) {
		return curr
	}

	return is.Schedule.Next(curr)
}

// Ecron ...
type Ecron interface {
	standard.Component
}
