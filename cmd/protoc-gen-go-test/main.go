package main

import (
	"errors"
	"flag"
	"fmt"

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
	var flags flag.FlagSet
	pkg := flags.String("pkg", "", "specified go package name")
	protogen.Options{ParamFunc: flags.Set}.Run(func(gen *protogen.Plugin) error {
		if *pkg == "" {
			return errors.New("empty pkg! you must set generated go package name with [pkg] flag")
		}
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generateFile(gen, f, *pkg)
		}
		return nil
	})
}
