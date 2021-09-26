#!/bin/bash


ROOT=$(echo "$(cd "$(dirname "../../../../")" && pwd)" )
PROTO=${ROOT}/examples/helloworld
SERVER=${ROOT}/examples/grpc/direct/server

rm -rf ./protoc-gen-go-errors
go build -o ${ROOT}/bin/protoc-gen-go-errors ${ROOT}/cmd/protoc-gen-go-errors
export PATH=$PATH:${ROOT}/bin/
protoc --proto_path=. --go_out=paths=source_relative:. --go-errors_out=paths=source_relative:. ./errors.proto
