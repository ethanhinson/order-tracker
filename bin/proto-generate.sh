#!/bin/bash

echo "Generating service stubs..."

for x in $(find "$(pwd)/protos" -name '*.proto'); do
  protoc --proto_path="$(pwd)"/protos \
  -I /usr/local/include \
  -I $GOPATH/src \
  --go_out=paths=import,plugins=grpc:./service  "$x"
done