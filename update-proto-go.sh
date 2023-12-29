#!/bin/bash

GOPKG="github.com/danieldiamont/go-yamcs-cli"

rm -rf yamcs/api/ yamcs/protobuf/
cp -r ~/dev/yamcs/yamcs-api/src/main/proto/yamcs .

# Set go_package for each .proto file
for proto_file in $(find yamcs -name '*.proto'); do
    echo "option go_package = \"$GOPKG/${proto_file%/*}\";" >> "$proto_file"
done

# Generate Go code using protoc with paths=source_relative
protoc --proto_path=. --go_opt=paths=source_relative --go_out=. $(find yamcs/protobuf -name '*.proto') $(find yamcs/api -name '*.proto')
