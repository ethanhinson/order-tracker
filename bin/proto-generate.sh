#!/bin/bash

echo "Generating service stubs..."

for x in $(find "$(pwd)/protos" -name '*.proto'); do
  protoc --proto_path="$(pwd)"/protos \
  -I /usr/local/include \
  -I $GOPATH/src \
  -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --go_out=paths=import,plugins=grpc:./service  "$x"
done

echo "Building descriptors for Envoy proxy..."

protoc --proto_path="$(pwd)"/protos \
  -I /usr/local/include \
  -I $GOPATH/src \
  -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --include_imports \
  --descriptor_set_out "./descriptors/proto.pb"  \
  "$(pwd)/protos/service.proto"
