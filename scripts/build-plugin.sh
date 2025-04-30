#!/bin/bash

# Generate GRPC code
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    internal/plugin/test/test.proto

# Build the plugin
go build -o ti-example-plugin ./cmd/ti-example-plugin

# Create the plugins directory if it doesn't exist
mkdir -p ~/.titanium/plugins

# Copy the plugin to the plugins directory
cp ti-example-plugin ~/.titanium/plugins/

# Build the test program
go build -o plugin-test ./cmd/plugin-test