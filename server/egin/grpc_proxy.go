package egin

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/codegangsta/inject"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	jsonpbMarshaler = jsonpb.Marshaler{
		EmitDefaults: true,
	}
	statusMSDefault *rpcstatus.Status
)

var (
	errBadRequest         = status.Errorf(codes.InvalidArgument, createStatusErr(codeMSInvalidParam, "bad request"))
	errMicroDefault       = status.Errorf(codes.Internal, createStatusErr(codeMS, "micro default"))
	errMicroInvoke        = status.Errorf(codes.Internal, createStatusErr(codeMSInvoke, "invoke failed"))
	errMicroInvokeLen     = status.Errorf(codes.Internal, createStatusErr(codeMSInvokeLen, "invoke result not 2 item"))
	errMicroInvokeInvalid = status.Errorf(codes.Internal, createStatusErr(codeMSSecondItemNotError, "second invoke res not a error"))
	errMicroResInvalid    = status.Errorf(codes.Internal, createStatusErr(codeMSResErr, "response is not valid"))
)

func init() {
	s, _ := status.FromError(errMicroDefault)
	de, _ := statusFromString(s.Message())
	statusMSDefault = de.Proto()
}

// protoError ...
func protoError(c *gin.Context, code int, e error) error {
	s, ok := status.FromError(e)
	c.Header(HeaderGRPCPROXYError, "true")
	if ok {
		if de, ok := statusFromString(s.Message()); ok {
			return protoJSON(c, code, de.Proto())
		}
	}
	return protoJSON(c, code, e)
}

// protoJSON sends a Protobuf JSON response with status code and data.
func protoJSON(c *gin.Context, code int, i interface{}) error {
	var acceptEncoding = c.Request.Header.Get(HeaderAcceptEncoding)
	var ok bool
	var m proto.Message
	if m, ok = i.(proto.Message); !ok {
		c.Header(HeaderGRPCPROXYError, "true")
		m = statusMSDefault
	}
	// protobuf output
	if strings.Contains(acceptEncoding, MIMEApplicationProtobuf) {
		c.Header(HeaderContentType, MIMEApplicationProtobuf)
		c.Writer.WriteHeader(code)
		bs, _ := proto.Marshal(m)
		_, err := c.Writer.Write(bs)
		return err
	}

	c.Header(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
	c.Writer.WriteHeader(code)

	return jsonpbMarshaler.Marshal(c.Writer, m)
}

// GRPCProxy experimental
func GRPCProxy(h interface{}) gin.HandlerFunc {
	t := reflect.TypeOf(h)
	if t.Kind() != reflect.Func {
		panic("reflect error: handler must be func")
	}
	return func(c *gin.Context) {
		var req = reflect.New(t.In(1).Elem()).Interface()
		if err := c.Bind(req); err != nil {
			protoError(c, http.StatusBadRequest, errBadRequest)
			return
		}
		var md = metadata.MD{}
		for k, vs := range c.Request.Header {
			for _, v := range vs {
				bs := bytes.TrimFunc([]byte(v), func(r rune) bool {
					return r == '\n' || r == '\r' || r == '\000'
				})
				md.Append(k, string(bs))
			}
		}
		ctx := metadata.NewOutgoingContext(context.TODO(), md)
		var inj = inject.New()
		inj.Map(ctx)
		inj.Map(req)
		vs, err := inj.Invoke(h)
		if err != nil {
			protoError(c, http.StatusInternalServerError, errMicroInvoke)
			return
		}
		if len(vs) != 2 {
			protoError(c, http.StatusInternalServerError, errMicroInvokeLen)
			return
		}
		repV, errV := vs[0], vs[1]
		if !errV.IsNil() || repV.IsNil() {
			if e, ok := errV.Interface().(error); ok {
				protoError(c, http.StatusOK, e)
				return
			}
			protoError(c, http.StatusInternalServerError, errMicroInvokeInvalid)
			return
		}
		if !repV.IsValid() {
			protoError(c, http.StatusInternalServerError, errMicroResInvalid)
			return
		}
		// todo 根据gRPC状态码转换为HTTP状态码
		protoJSON(c, http.StatusOK, repV.Interface())
		return
	}
}

type statusErr struct {
	s *rpcstatus.Status
}

// Proto ...
func (e *statusErr) Proto() *rpcstatus.Status {
	if e.s == nil {
		return nil
	}
	return proto.Clone(e.s).(*rpcstatus.Status)
}

func statusFromString(s string) (*statusErr, bool) {
	i := strings.Index(s, ":")
	if i == -1 {
		return nil, false
	}
	u64, err := strconv.ParseInt(s[:i], 10, 32)
	if err != nil {
		return nil, false
	}

	return &statusErr{
		&rpcstatus.Status{
			Code:    int32(u64),
			Message: s[i:],
			Details: []*any.Any{},
		},
	}, true
}

func createStatusErr(code uint32, msg string) string {
	return fmt.Sprintf("%d:%s", code, msg)
}
