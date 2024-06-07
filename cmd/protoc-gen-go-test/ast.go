package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	orderedmap "github.com/wk8/go-ordered-map"
	"golang.org/x/tools/imports"

	"github.com/gotomicro/ego/internal/tools"
)

func checkAndMerge(f *file) ([]byte, error) {
	origBytes, err := os.ReadFile(f.orig)
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

func getDecls(dstFile *dst.File) (*orderedmap.OrderedMap, *orderedmap.OrderedMap, *orderedmap.OrderedMap, *orderedmap.OrderedMap) {
	var importDecls = orderedmap.New()
	var typeDecls = orderedmap.New()
	var valueDecls = orderedmap.New()
	var fnDecls = orderedmap.New()

	for _, decl := range dstFile.Decls {
		switch t := decl.(type) {
		case *dst.FuncDecl:
			fnDecls.Set(t.Name.Name, t)
		case *dst.GenDecl:
			if len(t.Specs) == 0 {
				continue
			}
			switch t.Specs[0].(type) {
			case *dst.ValueSpec:
				for _, sp := range t.Specs {
					val := sp.(*dst.ValueSpec)
					for _, v := range val.Names {
						val.Names = []*dst.Ident{v}
						decl := *t
						decl.Specs = []dst.Spec{val}
						valueDecls.Set(v.Name, &decl)
					}
				}
			case *dst.TypeSpec:
				for _, sp := range t.Specs {
					decl := *t
					decl.Specs = []dst.Spec{sp}
					typeDecls.Set(sp.(*dst.TypeSpec).Name.String(), &decl)
				}
			case *dst.ImportSpec:
				for _, sp := range t.Specs {
					decl := *t
					decl.Specs = []dst.Spec{sp}
					importDecls.Set(sp.(*dst.ImportSpec).Path.Value, &decl)
				}
			}
		default:
			panic("bad decls")
		}
	}
	return importDecls, typeDecls, valueDecls, fnDecls
}

func castStrToBool(str string) bool {
	return str == "true"
}

type decls []dst.Decl

func (decls *decls) rebuildGenDecls(opDecls *orderedmap.OrderedMap, tpDecls *orderedmap.OrderedMap) *decls {
	for op := opDecls.Oldest(); op != nil; op = op.Next() {
		k := op.Key.(string)
		d := op.Value.(*dst.GenDecl)
		anno := getAnnotations(strings.Join(d.Decs.Start, "\n"))[annotationOverride]
		tp := tpDecls.GetPair(op.Key)
		if tp != nil {
			if castStrToBool(anno.val) {
				*decls = append(*decls, d)
			} else {
				*decls = append(*decls, tp.Value.(*dst.GenDecl))
			}
			tpDecls.Delete(k)
		} else {
			*decls = append(*decls, d)
		}
	}
	for tp := tpDecls.Oldest(); tp != nil; tp = tp.Next() {
		*decls = append(*decls, tp.Value.(*dst.GenDecl))
	}
	return decls
}

func (decls *decls) rebuildFnDecls(opDecls *orderedmap.OrderedMap, tpDecls *orderedmap.OrderedMap) *decls {
	for op := opDecls.Oldest(); op != nil; op = op.Next() {
		k := op.Key.(string)
		d := op.Value.(*dst.FuncDecl)
		anno := getAnnotations(strings.Join(d.Decs.Start, "\n"))[annotationOverride]
		tp := tpDecls.GetPair(op.Key)
		if tp != nil {
			if castStrToBool(anno.val) {
				*decls = append(*decls, d)
			} else {
				*decls = append(*decls, tp.Value.(*dst.FuncDecl))
			}
			tpDecls.Delete(k)
		} else {
			*decls = append(*decls, d)
		}
	}
	for tp := tpDecls.Oldest(); tp != nil; tp = tp.Next() {
		*decls = append(*decls, tp.Value.(*dst.FuncDecl))
	}
	return decls
}

func genFinalContent(origContent []byte, tempContent []byte) ([]byte, error) {
	if len(origContent) == 0 {
		return tempContent, nil
	}
	origf, err := decorator.Parse(origContent)
	if err != nil {
		return nil, fmt.Errorf("decorator parse fail, %w", err)
	}
	tempf, err := decorator.Parse(tempContent)
	if err != nil {
		return nil, fmt.Errorf("decorator parse fail, %w", err)
	}

	var decls decls
	origImportDecls, origTypeDecls, origValueDecls, origFnDecls := getDecls(origf)
	tempImportDecls, tempTypeDecls, tempValueDecls, tempFnDecls := getDecls(tempf)
	// 1. rebuild import decls
	// 2. rebuild type decls
	// 3. rebuild value decls
	// 4. rebuild function decls
	decls.rebuildGenDecls(origImportDecls, tempImportDecls).
		rebuildGenDecls(origTypeDecls, tempTypeDecls).
		rebuildGenDecls(origValueDecls, tempValueDecls).
		rebuildFnDecls(origFnDecls, tempFnDecls)
	origf.Decls = decls
	buf := bytes.NewBuffer([]byte{})
	if err := decorator.Fprint(buf, origf); err != nil {
		return nil, fmt.Errorf("fprint fail, %w", err)
	}
	return imports.Process("", tools.GoFmt(buf.Bytes()), &imports.Options{Comments: true, TabIndent: true, TabWidth: 8, FormatOnly: true})
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
