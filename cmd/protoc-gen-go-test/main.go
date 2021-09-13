package main

import (
	"flag"
	"fmt"
	"os"

	"google.golang.org/protobuf/compiler/protogen"
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
	var (
		flags flag.FlagSet
		out   = flags.String("out", "", "specified output directory name")
		mod   = flags.String("mod", "", "specified pb stub code module path")
	)
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		var outDir string
		if *out != "" {
			outDir = *out
		}
		if err := checkOut(outDir); err != nil {
			return err
		}
		generateInitFile(gen, *out, *mod)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			if _, err := generateFile(gen, f, *out, *mod); err != nil {
				return fmt.Errorf("generate file fail, %w", err)
			}
		}
		return nil
	})
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
