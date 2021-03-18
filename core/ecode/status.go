package ecode

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	spb "google.golang.org/genproto/googleapis/rpc/status"
)

// Status ...
type Status interface {
}

// SpbStatus ...
type SpbStatus struct {
	*spb.Status
}

// GetCodeAsInt ...
func (s *SpbStatus) GetCodeAsInt() int {
	return int(s.Code)
}

// GetCodeAsUint32 ...
func (s *SpbStatus) GetCodeAsUint32() uint32 {
	return uint32(s.Code)
}

// GetCodeAsBool ...
func (s *SpbStatus) GetCodeAsBool() bool {
	return s.CauseCode() == 0
}

// GetMessage ...
func (s *SpbStatus) GetMessage(exts ...interface{}) string {
	if len(exts)%2 != 0 {
		panic("parameter must be odd")
	}

	var buf bytes.Buffer
	buf.WriteString(s.Message)

	if len(exts) > 0 {
		buf.WriteByte(',')
	}
	for i := 0; i < len(exts); i++ {
		buf.WriteString(fmt.Sprintf("%v", exts[i]))
		buf.WriteByte(':')
		buf.WriteString(fmt.Sprintf("%v", exts[i+1]))
		i++
	}
	return buf.String()
}

// GetDetailMessage ...
func (s *SpbStatus) GetDetailMessage(exts ...interface{}) string {
	var buf bytes.Buffer
	buf.WriteString(s.GetMessage(exts...))
	for _, detail := range s.Details {
		buf.WriteByte('\n')
		buf.WriteString(detail.String())
	}
	return buf.String()
}

// String ...
func (s *SpbStatus) String() string {
	bs, _ := json.Marshal(s)
	return string(bs)
}

// CauseCode ...
func (s *SpbStatus) CauseCode() int {
	return int(s.Code) % 10000
}

// Proto ...
func (s *SpbStatus) Proto() *spb.Status {
	if s == nil {
		return nil
	}
	return proto.Clone(s.Status).(*spb.Status)
}

// MustWithDetails ...
func (s *SpbStatus) MustWithDetails(details ...interface{}) *SpbStatus {
	status, err := s.WithDetails(details...)
	if err != nil {
		panic(err)
	}
	return status
}

// WithDetails returns a new status with the provided details messages appended to the status.
// If any errors are encountered, it returns nil and the first error encountered.
func (s *SpbStatus) WithDetails(details ...interface{}) (*SpbStatus, error) {
	if s.CauseCode() == 0 {
		return nil, errors.New("no error details for status with code OK")
	}
	p := s.Proto()
	for _, detail := range details {
		if pmsg, ok := detail.(proto.Message); ok {
			any, err := marshalAnyProtoMessage(pmsg)
			if err != nil {
				return nil, err
			}
			p.Details = append(p.Details, any)
		} else {
			any, err := marshalAny(detail)
			if err != nil {
				return nil, err
			}
			p.Details = append(p.Details, any)
		}
	}
	return &SpbStatus{Status: p}, nil
}

func marshalAny(obj interface{}) (*any.Any, error) {
	typ := reflect.TypeOf(obj)
	val := fmt.Sprintf("%+v", obj)

	return &any.Any{TypeUrl: typ.Name(), Value: []byte(val)}, nil
}

func marshalAnyProtoMessage(pb proto.Message) (*any.Any, error) {
	value, err := proto.Marshal(pb)
	if err != nil {
		return nil, err
	}
	return &any.Any{TypeUrl: proto.MessageName(pb), Value: value}, nil
}
