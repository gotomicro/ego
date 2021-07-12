package ecode

import (
	spb "google.golang.org/genproto/googleapis/rpc/status"
)

// SpbStatus ...
type SpbStatus struct {
	*spb.Status
}
