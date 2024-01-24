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
	"text/template"

	"github.com/gotomicro/ego/internal/tools"
)

const errorsTpl = `
{{ range .Errors }}
var {{.LowerCamelValue}} *eerrors.EgoError
{{- end }}

var i18n = map[string]map[string]string{
{{- range .Errors }}
	"{{.Key}}": map[string]string{
		{{- range $k,$v :=  .I18n }}
			"{{$k}}": "{{$v}}",
		{{- end }}
	},
{{- end }}
}

// ReasonI18n provides error messages in a specified language. 
// For instance, to get an error message in Chinese for "@i18n.cn", you can use ReasonI18n(e, "cn").
func ReasonI18n(e eerrors.Error, lan string) string {
	return i18n[eerrors.FromError(e).Reason][lan]
}

func init() {
{{- range .Errors }}
{{.LowerCamelValue}} = eerrors.New(int(codes.{{.Code}}), "{{.Key}}", {{.Name}}_{{.Value}}.String())
eerrors.Register({{.LowerCamelValue}})
{{- end }}
}

{{ range .Errors }}
{{if .HasComment}}{{.Comment}}{{end}}func {{.UpperCamelValue}}() eerrors.Error {
	 return {{.LowerCamelValue}}
}
{{ end }}
`

type errorInfo struct {
	Name            string
	Value           string
	Code            string
	UpperCamelValue string
	LowerCamelValue string
	Key             string
	Comment         string
	HasComment      bool
	I18n            map[string]string
}

type errorWrapper struct {
	Errors []*errorInfo
}

func (e *errorWrapper) execute() string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("errors").Parse(errorsTpl)
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, e); err != nil {
		panic(err)
	}
	return string(tools.GoFmt(buf.Bytes()))
}
