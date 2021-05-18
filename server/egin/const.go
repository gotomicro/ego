package egin

const (
	// HeaderAcceptEncoding ...
	HeaderAcceptEncoding = "Accept-Encoding"
	// HeaderContentType ...
	HeaderContentType = "Content-Type"
	// HeaderGRPCPROXYError ...
	HeaderGRPCPROXYError = "GRPC-Proxy-Error"
	charsetUTF8          = "charset=utf-8"

	// MIMEApplicationJSON ...
	MIMEApplicationJSON = "application/json"
	// MIMEApplicationJSONCharsetUTF8 ...
	MIMEApplicationJSONCharsetUTF8 = MIMEApplicationJSON + "; " + charsetUTF8
	// MIMEApplicationProtobuf ...
	MIMEApplicationProtobuf = "application/protobuf"
)

const (
	codeMS                   = 1000
	codeMSInvalidParam       = 1001
	codeMSInvoke             = 1002
	codeMSInvokeLen          = 1003
	codeMSSecondItemNotError = 1004
	codeMSResErr             = 1005
)
