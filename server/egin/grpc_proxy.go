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
	"github.com/golang/protobuf/ptypes/any"
	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	opts = protojson.MarshalOptions{
		EmitUnpopulated: true,
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
func protoError(c *gin.Context, code int, e error) ([]byte, error) {
	s, ok := status.FromError(e)
	c.Header(HeaderGRPCPROXYError, "true")
	if e != nil {
		c.Header("Error", e.Error())
	}
	if ok {
		if de, ok := statusFromString(s.Message()); ok {
			return protoJSON(c, code, de.Proto())
		}
	}
	return protoJSON(c, code, e)
}

// protoJSON sends a Protobuf JSON response with status code and data.
func protoJSON(c *gin.Context, code int, i interface{}) ([]byte, error) {
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
		return []byte{}, err
	}

	c.Header(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
	c.Writer.WriteHeader(code)

	return opts.Marshal(m)
}

// GRPCProxy experimental
func GRPCProxy(h interface{}) gin.HandlerFunc {
	t := reflect.TypeOf(h)
	if t.Kind() != reflect.Func {
		panic("reflect error: handler must be func")
	}
	return func(c *gin.Context) {
		var req = reflect.New(t.In(1).Elem()).Interface()
		if c.Request.Method == http.MethodGet {
			_ = c.BindUri(req)
		}
		if err := c.Bind(req); err != nil {
			output, _ := protoError(c, http.StatusBadRequest, errBadRequest)
			_, _ = c.Writer.Write(output)
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
			output, _ := protoError(c, http.StatusInternalServerError, errMicroInvoke)
			_, _ = c.Writer.Write(output)
			return
		}
		if len(vs) != 2 {
			output, _ := protoError(c, http.StatusInternalServerError, errMicroInvokeLen)
			_, _ = c.Writer.Write(output)
			return
		}
		repV, errV := vs[0], vs[1]
		if !errV.IsNil() || repV.IsNil() {
			if e, ok := errV.Interface().(error); ok {
				output, _ := protoError(c, http.StatusOK, e)
				_, _ = c.Writer.Write(output)
				return
			}
			output, _ := protoError(c, http.StatusInternalServerError, errMicroInvokeInvalid)
			_, _ = c.Writer.Write(output)
			return
		}
		if !repV.IsValid() {
			output, _ := protoError(c, http.StatusInternalServerError, errMicroResInvalid)
			_, _ = c.Writer.Write(output)
			return
		}
		// todo 根据gRPC状态码转换为HTTP状态码
		output, _ := protoJSON(c, http.StatusOK, repV.Interface())
		_, _ = c.Writer.Write(output)
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
