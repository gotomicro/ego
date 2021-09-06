#!/bin/bash


rm -rf ./protoc-gen-go-errors
go build -o ./protoc-gen-go-errors ../../../cmd/protoc-gen-go-errors
export PATH=$PATH:$(pwd)/
protoc --proto_path=. --go_out=paths=source_relative:. --go-errors_out=paths=source_relative,file=errors.proto:. ./errors.proto
