package eerrors

import (
	"errors"
	"fmt"

	"github.com/gotomicro/ego/internal/ecode"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

//go:generate protoc -I. --go_out=paths=source_relative:. errors.proto

// Error 错误接口
type Error interface {
	error
	WithMetadata(map[string]string) Error
	WithMessage(string) Error
}

const (
	// UnknownReason is unknown reason for error info.
	UnknownReason = ""
	// SupportPackageIsVersion1 this constant should not be referenced by any other code.
	SupportPackageIsVersion1 = true
)

var _ Error = &EgoError{}

type errKey string

var errs = map[errKey]*EgoError{}

// Register 注册错误信息
func Register(egoError *EgoError) {
	errs[errKey(egoError.Reason)] = egoError
}

func (x *EgoError) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s metadata = %v", x.Code, x.Reason, x.Message, x.Metadata)
}

// GRPCStatus returns the Status represented by se.
func (x *EgoError) GRPCStatus() *status.Status {
	s, _ := status.New(codes.Code(x.Code), x.Message).
		WithDetails(&errdetails.ErrorInfo{
			Reason:   x.Reason,
			Metadata: x.Metadata,
		})
	return s
}

// WithMetadata with an MD formed by the mapping of key, value.
func (x *EgoError) WithMetadata(md map[string]string) Error {
	err := proto.Clone(x).(*EgoError)
	err.Metadata = md
	return err
}

// WithMessage set message to current EgoError
func (x *EgoError) WithMessage(msg string) Error {
	err := proto.Clone(x).(*EgoError)
	err.Message = msg
	return err
}

// New returns an error object for the code, message.
func New(code int, reason, message string) *EgoError {
	return &EgoError{
		Code:    int32(code),
		Message: message,
		Reason:  reason,
	}
}

// ToHTTPStatusCode Get equivalent HTTP status code from x.Code
func (x *EgoError) ToHTTPStatusCode() int {
	return ecode.GrpcToHTTPStatusCode(codes.Code(x.Code))
}

// FromError try to convert an error to *Error.
// It supports wrapped errors.
func FromError(err error) *EgoError {
	if err == nil {
		return nil
	}
	if se := new(EgoError); errors.As(err, &se) {
		return se
	}
	gs, ok := status.FromError(err)
	if ok {
		for _, detail := range gs.Details() {
			switch d := detail.(type) {
			case *errdetails.ErrorInfo:
				e, ok := errs[errKey(d.Reason)]
				if ok {
					return e
				}
				return New(
					int(gs.Code()),
					d.Reason,
					gs.Message(),
				).WithMetadata(d.Metadata).(*EgoError)
			}
		}
	}
	return New(int(codes.Unknown), UnknownReason, err.Error())
}
