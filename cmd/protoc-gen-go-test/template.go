// MIT License
//
// Copyright (c) 2020 go-kratos
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
	err := template.Must(tpl.Parse(baseTpl)).Execute(buf, nil)
	if err != nil {
		log.Fatal("render tmpl fail", err)
	}
	for _, svc := range w.Svcs {
		err := template.Must(tpl.Parse(svc.tmpl())).Execute(buf, svc)
		if err != nil {
			log.Fatal("render tmpl fail", err)
			return ""
		}
	}
	return string(tools.GoFmt(buf.Bytes()))
}
