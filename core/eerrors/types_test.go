package eerrors

import (
	"testing"
)

func TestTypes(t *testing.T) {
	var (
		input = []*EgoError{
			Canceled("reason_canceled", "message_canceled"),
			DeadlineExceeded("reason_deadline_exceeded", "message_deadline_exceeded"),
			NotFound("reason_not_found", "message_not_found"),
		}
		output = []func(egoError *EgoError) bool{
			IsCanceled,
			IsDeadlineExceeded,
			IsNotFound,
		}
	)

	for i, in := range input {
		if !output[i](in) {
			t.Errorf("not expect: %v", in)
		}
	}
}
