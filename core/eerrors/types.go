package eerrors

import (
	"google.golang.org/grpc/codes"
)

// Canceled new Canceled error that is mapped to context.Canceled.
func Canceled(reason, message string) *EgoError {
	return New(int(codes.Canceled), reason, message)
}

// IsCanceled determines if err is an error which indicates a Canceled error.
func IsCanceled(err *EgoError) bool {
	return err.Code == int32(codes.Canceled)
}

// DeadlineExceeded new DeadlineExceeded error that is mapped to context.DeadlineExceeded.
func DeadlineExceeded(reason, message string) *EgoError {
	return New(int(codes.DeadlineExceeded), reason, message)
}

// IsDeadlineExceeded determines if err is an error which indicates a DeadlineExceeded error.
func IsDeadlineExceeded(err *EgoError) bool {
	return err.Code == int32(codes.DeadlineExceeded)
}

// NotFound new NotFound error that is mapped to context.NotFound.
func NotFound(reason, message string) *EgoError {
	return New(int(codes.NotFound), reason, message)
}

// IsNotFound determines if err is an error which indicates a NotFound error.
func IsNotFound(err *EgoError) bool {
	return err.Code == int32(codes.NotFound)
}
