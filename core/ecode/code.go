package ecode

import (
	"encoding/json"
	"net/http"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server/egovernor"
	"github.com/golang/protobuf/ptypes/any"
	spb "google.golang.org/genproto/googleapis/rpc/status"
)

// EcodeNum 低于10000均为系统错误码，业务错误码请使用10000以上
const EcodeNum int32 = 9999

var (
	aid              int
	maxCustomizeCode = 9999
	_codes           sync.Map
	// OK ...
	OK = add(int(codes.OK), "OK")
)

func init() {
	// status code list
	egovernor.HandleFunc("/status/code/list", func(w http.ResponseWriter, r *http.Request) {
		var rets = make(map[int]*spbStatus)
		_codes.Range(func(key, val interface{}) bool {
			code := key.(int)
			status := val.(*spbStatus)
			rets[code] = status
			return true
		})
		_ = json.NewEncoder(w).Encode(rets)
	})
}

// Add ...
func Add(code int, message string) *spbStatus {
	if code > maxCustomizeCode {
		elog.Panic("customize code must less than 9999", elog.Any("code", code))
	}

	return add(aid*10000+code, message)
}

func add(code int, message string) *spbStatus {
	status := &spbStatus{
		&spb.Status{
			Code:    int32(code),
			Message: message,
			Details: make([]*any.Any, 0),
		},
	}
	_codes.Store(code, status)
	return status
}

// ExtractCodes cause from error to ecode.
func ExtractCodes(e error) *spbStatus {
	if e == nil {
		return OK
	}
	// todo 不想做code类型转换，所以全部用grpc标准码处理
	// 如果存在标准的grpc的错误，直接返回自定义的ecode编码
	gst, _ := status.FromError(e)
	return &spbStatus{
		&spb.Status{
			Code:    int32(gst.Code()),
			Message: gst.Message(),
			Details: make([]*any.Any, 0),
		},
	}
}
