#!/bin/bash


ROOT=$(echo "$(cd "$(dirname "../../../../../")" && pwd)" )
PROTO=${ROOT}/examples/helloworld
SERVER=${ROOT}/examples/grpc/direct/server

cd ${PROTO}
protoc --proto_path=. --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. helloworld.proto

cd ${SERVER}
rm -rf ./protoc-gen-go-test
go build -o ${ROOT}/bin/protoc-gen-go-test ${ROOT}/cmd/protoc-gen-go-test
export PATH=$PATH:${ROOT}/bin/
protoc --proto_path=${ROOT}/examples/helloworld  --go-test_out=mod=github.com/gotomicro/ego/examples/helloworld,out=./server,paths=source_relative:. helloworld.proto
