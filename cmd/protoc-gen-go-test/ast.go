package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/gotomicro/ego/internal/tools"
)

func checkAndMerge(f *file) ([]byte, error) {
	origBytes, err := ioutil.ReadFile(f.orig)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("read origFile fail, %w", err)
	}
	tempBytes, err := f.g.Content()
	if err != nil {
		return nil, fmt.Errorf("read tempBytes fail, %w", err)
	}
	finaBytes, err := genFinalContent(origBytes, tempBytes)
	if err != nil {
		return nil, fmt.Errorf("generate finaBytes fail, %w", err)
	}
	return finaBytes, nil
}

func genFinalContent(origContent []byte, tempContent []byte) ([]byte, error) {
	if len(origContent) == 0 {
		return tempContent, nil
	}
	origf, err := decorator.Parse(origContent)
	if err != nil {
		return nil, fmt.Errorf("decorator parse fail, %w", err)
	}
	// cannotOverrideFns store functions we can't override now
	cannotOverrideFns := map[string]*dst.FuncDecl{}
	for _, decl := range origf.Decls {
		fn, ok := decl.(*dst.FuncDecl)
		if ok {
			comments := strings.Join(fn.Decs.Start, "\n")
			annos := getAnnotations(comments)
			if anno, ok := annos[annotationOverride]; ok && anno.val == "true" {
				cannotOverrideFns[fn.Name.Name] = fn
			}
		}
	}

	tempf, err := decorator.Parse(tempContent)
	if err != nil {
		return nil, fmt.Errorf("decorator parse fail, %w", err)
	}
	decls := []dst.Decl{}
	for _, decl := range tempf.Decls {
		fn, ok := decl.(*dst.FuncDecl)
		if !ok {
			decls = append(decls, decl)
			continue
		}
		// if function can't be override now, we ignore it
		oldFn, ok := cannotOverrideFns[fn.Name.Name]
		if ok {
			fn = oldFn
		}
		decls = append(decls, fn)
	}
	tempf.Decls = decls
	buf := bytes.NewBuffer([]byte{})
	if err := decorator.Fprint(buf, tempf); err != nil {
		return nil, fmt.Errorf("fprint fail, %w", err)
	}

	return tools.GoFmt(buf.Bytes()), nil
}

const (
	// annotationOverride is "Override" annotation, Override=true means user have already overridden it.
	annotationOverride = "Override"
)

var commentRgx, _ = regexp.Compile(`@(\w+)=([_a-zA-Z0-9-:]+)`)

type annotation struct {
	name string
	val  string
}

// getAnnotations parse annotations from comments
func getAnnotations(comment string) map[string]annotation {
	matches := commentRgx.FindAllStringSubmatch(comment, -1)
	annotations := make(map[string]annotation)
	for _, v := range matches {
		annotations[v[1]] = annotation{
			name: v[1],
			val:  v[2],
		}
	}
	return annotations
}
