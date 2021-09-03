package ecode

import (
	"google.golang.org/grpc/codes"
)

// EcodeNum 低于10000均为系统错误码，业务错误码请使用10000以上
const EcodeNum codes.Code = 9999
