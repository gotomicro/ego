package ecode

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/any"
	spb "google.golang.org/genproto/googleapis/rpc/status"
)

// EcodeNum 低于10000均为系统错误码，业务错误码请使用10000以上
const EcodeNum int32 = 9999

var (
	// OK ...
	OK = &SpbStatus{
		&spb.Status{
			Code:    int32(codes.OK),
			Message: "OK",
			Details: make([]*any.Any, 0),
		},
	}
)

// ExtractCodes cause from error to ecode.
func ExtractCodes(e error) *SpbStatus {
	if e == nil {
		return OK
	}
	// 如果存在标准的grpc的错误，直接返回自定义的ecode编码
	gst, _ := status.FromError(e)
	return &SpbStatus{
		&spb.Status{
			Code:    int32(gst.Code()),
			Message: gst.Message(),
			Details: make([]*any.Any, 0),
		},
	}
}
