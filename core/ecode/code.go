package ecode

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// EcodeNum 低于10000均为系统错误码，业务错误码请使用10000以上
const EcodeNum codes.Code = 9999

// ExtractCodes cause from error to ecode.
func ExtractCodes(e error) *status.Status {
	if e == nil {
		return status.New(codes.OK, "OK")
	}
	// 如果存在标准的grpc的错误，直接返回自定义的ecode编码
	return status.FromContextError(e)
}
