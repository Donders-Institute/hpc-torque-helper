#!/bin/bash

# requires packages:
# - protobuf: https://github.com/google/protobuf/releases
# - protoc-gen-go: go get -u github.com/golang/protobuf/protoc-gen-go

# The following command generate client and server gRPC interfaces in ../internal/grpc.pb.go
protoc --go_out=plugins=grpc:../internal grpc.proto
