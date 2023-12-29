#!/bin/bash

# Generate Go code using protoc with paths=source_relative
protoc --proto_path=. --go_opt=paths=source_relative --go_out=. $(find yamcs/protobuf -name '*.proto') $(find yamcs/api -name '*.proto')
