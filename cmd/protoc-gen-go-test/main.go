package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

var versionFlag = flag.Bool("version", false, "print the version and exit")

var version string

func main() {
	flag.Parse()
	if *versionFlag {
		fmt.Printf("Version: %s\n", version)
		return
	}
	var flags flag.FlagSet
	var outFlag = flags.String("out", "", "specified output directory name")
	var modFlag = flags.String("mod", "", "specified pb stub code module path")
	// check args
	if len(os.Args) > 1 {
		exit(fmt.Errorf("unknown argument %q (this program should be run by protoc, not directly)", os.Args[1]))
	}
	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		exit(err)
	}

	// prepare generation
	req := &pluginpb.CodeGeneratorRequest{}
	if err := proto.Unmarshal(in, req); err != nil {
		exit(err)
	}
	gen, err := protogen.Options{ParamFunc: flags.Set}.New(req)
	if err != nil {
		exit(err)
	}

	// execute generation
	if err := run(gen, outFlag, modFlag); err != nil {
		gen.Error(err)
	}
	resp := gen.Response()
	out, err := proto.Marshal(resp)
	if err != nil {
		exit(err)
	}

	// write to stdout
	if _, err := os.Stdout.Write(out); err != nil {
		exit(err)
	}
}

func run(gen *protogen.Plugin, outFlag *string, modFlag *string) error {
	gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	if err := checkOut(*outFlag); err != nil {
		return err
	}
	_, err := generateInitFile(gen, nil, *outFlag, *modFlag)
	if err != nil {
		return fmt.Errorf("generate init test file fail, %w", err)
	}
	for _, f := range gen.Files {
		if !f.Generate {
			continue
		}
		// if proto contains services, we do not generate code.
		if f.Desc.Services().Len() == 0 {
			continue
		}
		_, err := generateFile(gen, f, *outFlag, *modFlag)
		if err != nil {
			return fmt.Errorf("generate test file fail, %w", err)
		}
	}
	return nil
}

func checkOut(out string) error {
	fi, err := os.Stat(out)
	if err != nil {
		return fmt.Errorf("out parameter is invalid, %w", err)
	}
	if mode := fi.Mode(); !mode.IsDir() {
		return fmt.Errorf("out parameter is not a directory, %w", err)
	}
	return nil
}

func exit(err error) {
	if _, err := fmt.Fprintf(os.Stderr, "%s: %v\n", filepath.Base(os.Args[0]), err); err != nil {
		log.Println("fprintf fail, %w", err)
	}
	os.Exit(1)
}
