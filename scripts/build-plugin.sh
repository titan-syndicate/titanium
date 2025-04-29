#!/bin/bash

# Generate GRPC code
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    internal/plugin/test/test.proto

# Build the plugin
go build -o test-plugin ./cmd/test-plugin

# Build the test program
go build -o plugin-test ./cmd/plugin-test