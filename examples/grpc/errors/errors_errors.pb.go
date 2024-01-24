// Code generated by protoc-gen-go-errors. DO NOT EDIT.

package bizv1

import (
	eerrors "github.com/gotomicro/ego/core/eerrors"
	codes "google.golang.org/grpc/codes"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the ego package it is being compiled against.
const _ = eerrors.SupportPackageIsVersion1

var errOk *eerrors.EgoError
var errUnknown *eerrors.EgoError
var errUserNotFound *eerrors.EgoError
var errUserIdNotValid *eerrors.EgoError

var i18n = map[string]map[string]string{
	"biz.v1.ERR_OK": map[string]string{
		"cn": "请求成功",
		"en": "OK",
	},
	"biz.v1.ERR_UNKNOWN": map[string]string{
		"cn": "服务内部未知错误",
		"en": "unknown error",
	},
	"biz.v1.ERR_USER_NOT_FOUND": map[string]string{
		"cn": "找不到指定用户",
		"en": "user not found",
	},
	"biz.v1.ERR_USER_ID_NOT_VALID": map[string]string{
		"cn": "用户ID不合法",
		"en": "invalid user id",
	},
}

// ReasonI18n provides error messages in a specified language.
// For instance, to get an error message in Chinese for "@i18n.cn", you can use ReasonI18n(e, "cn").
func ReasonI18n(e eerrors.Error, lan string) string {
	return i18n[eerrors.FromError(e).Reason][lan]
}

func init() {
	errOk = eerrors.New(int(codes.OK), "biz.v1.ERR_OK", Err_ERR_OK.String())
	eerrors.Register(errOk)
	errUnknown = eerrors.New(int(codes.Unknown), "biz.v1.ERR_UNKNOWN", Err_ERR_UNKNOWN.String())
	eerrors.Register(errUnknown)
	errUserNotFound = eerrors.New(int(codes.NotFound), "biz.v1.ERR_USER_NOT_FOUND", Err_ERR_USER_NOT_FOUND.String())
	eerrors.Register(errUserNotFound)
	errUserIdNotValid = eerrors.New(int(codes.InvalidArgument), "biz.v1.ERR_USER_ID_NOT_VALID", Err_ERR_USER_ID_NOT_VALID.String())
	eerrors.Register(errUserIdNotValid)
}

// ErrOk  请求正常，实际上不算是一个错误
// @code=OK
// @i18n.cn="请求成功"
// @i18n.en="OK"
func ErrOk() eerrors.Error {
	return errOk
}

// ErrUnknown  未知错误，比如业务panic了
// @code=UNKNOWN             # 定义了这个错误关联的gRPC Code为：UNKNOWN
// @i18n.cn="服务内部未知错误"        # 定义了一个中文错误文案
// @i18n.en="unknown error"  # 定义了一个英文错误文案
func ErrUnknown() eerrors.Error {
	return errUnknown
}

// ErrUserNotFound  找不到指定用户
// @code=NOT_FOUND
// @i18n.cn="找不到指定用户"
// @i18n.en="user not found"
func ErrUserNotFound() eerrors.Error {
	return errUserNotFound
}

// ErrUserIdNotValid  用户ID不合法
// @code=INVALID_ARGUMENT
// @i18n.cn="用户ID不合法"
// @i18n.en="invalid user id"
func ErrUserIdNotValid() eerrors.Error {
	return errUserIdNotValid
}
