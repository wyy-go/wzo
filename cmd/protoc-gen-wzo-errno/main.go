package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

const version = "0.0.1"

var args = &struct {
	ShowVersion   bool   // version
	ErrorsPackage string // error package
}{

	ShowVersion:   false,
	ErrorsPackage: "",
}

func init() {
	flag.BoolVar(&args.ShowVersion, "version", false, "print the version and exit")
	flag.StringVar(&args.ErrorsPackage, "epk", "github.com/wyy-go/wzo/core/errors", "errors core package in your project")
}

func main() {
	flag.Parse()
	if args.ShowVersion {
		fmt.Printf("protoc-gen-wzo-errno %v\n", version)
		return
	}
	protogen.Options{ParamFunc: flag.CommandLine.Set}.Run(runProtoGen)
}
