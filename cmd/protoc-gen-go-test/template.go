package main

import (
	"bytes"
	"log"
	"text/template"

	"github.com/gotomicro/ego/internal/tools"
)

const (
	baseTpl = `
var svc *egrpc.Component

func init() {
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return svc.Listener().(*bufconn.Listener).Dial()
}
`
	stream2streamTpl = `
func Test{{.Name}}(t *testing.T) {
	cli := {{.Package}}.New{{.Service.Name}}Client(cegrpc.DefaultContainer().Build(cegrpc.WithDialOption(grpc.WithContextDialer(bufDialer))).ClientConn)
    ctx := context.Background()
    req := &{{.Package}}.{{.InType}}{}

    stream, err := cli.{{.Name}}(ctx)
    assert.NoError(t, err)
	err = stream.Send(req)
    assert.NoError(t, err)

    wantRes := &{{.Package}}.{{.OutType}}{}
	res, err := stream.Recv()
    assert.NoError(t, err)
    assert.NotNil(t, res)
	assert.True(t, proto.Equal(wantRes, res))
	t.Logf("res: %+v", res)
}

`
	unary2streamTpl = `func Test{{.Name}}(t *testing.T) {
	cli := {{.Package}}.New{{.Service.Name}}Client(cegrpc.DefaultContainer().Build(cegrpc.WithDialOption(grpc.WithContextDialer(bufDialer))).ClientConn)
    ctx := context.Background()
    req := &{{.Package}}.{{.InType}}{}

    stream, err := cli.{{.Name}}(ctx, req)
    wantRes := &{{.Package}}.{{.OutType}}{}
    res, err := stream.Recv()
    assert.NoError(t, err)

    assert.NoError(t, err)
    assert.NotNil(t, res)
	assert.True(t, proto.Equal(wantRes, res))
	t.Logf("res: %+v", res)
}

`
	stream2unaryTpl = `
func Test{{.Name}}(t *testing.T) {
	cli := {{.Package}}.New{{.Service.Name}}Client(cegrpc.DefaultContainer().Build(cegrpc.WithDialOption(grpc.WithContextDialer(bufDialer))).ClientConn)
    ctx := context.Background()

    stream, err := cli.{{.Name}}(ctx)
    assert.NoError(t, err)

    req := &{{.Package}}.{{.InType}}{}
	err = stream.Send(req)
    assert.NoError(t, err)

    wantRes := &{{.Package}}.{{.OutType}}{}
	res, err := stream.CloseAndRecv()
    assert.NoError(t, err)
	assert.Equal(t, res, nil)
	assert.True(t, proto.Equal(wantRes, res))
	t.Logf("res: %+v", res)
}

`
	unary2unaryTpl = `
func Test{{.Name}}(t *testing.T) {
	cli := {{.Package}}.New{{.Service.Name}}Client(cegrpc.DefaultContainer().Build(cegrpc.WithDialOption(grpc.WithContextDialer(bufDialer))).ClientConn)
    ctx := context.Background()
    req := &{{.Package}}.{{.InType}}{}

    wantRes := &{{.Package}}.{{.OutType}}{}
    res, err := cli.{{.Name}}(ctx, req)
    assert.NoError(t, err)
    assert.NotNil(t, res)
	assert.True(t, proto.Equal(wantRes, res))
	t.Logf("res: %+v", res)
}

`
)

var tmpls = map[string]string{
	"stream_stream": stream2streamTpl,
	"stream_unary":  stream2unaryTpl,
	"unary_stream":  unary2streamTpl,
	"unary_unary":   unary2unaryTpl,
}

type svcData struct {
	Name              string
	InType            string
	OutType           string
	Package           string
	Service           service
	isStreamingServer bool
	isStreamingClient bool
}

type service struct {
	Name string
}

type svcWrapper struct {
	Svcs []*svcData
}

func (s svcData) tmpl() string {
	var srv, cli string
	if s.isStreamingServer {
		srv = "stream"
	} else {
		srv = "unary"
	}

	if s.isStreamingClient {
		cli = "stream"
	} else {
		cli = "unary"
	}

	return tmpls[cli+"_"+srv]
}

func (w *svcWrapper) execute() string {
	buf := new(bytes.Buffer)

	tpl := template.New("test")
	for _, svc := range w.Svcs {
		err := template.Must(tpl.Parse(svc.tmpl())).Execute(buf, svc)
		if err != nil {
			log.Fatal("render tmpl fail", err)
			return ""
		}
	}
	return string(tools.GoFmt(buf.Bytes()))
}
