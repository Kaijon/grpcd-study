#!/bin/sh

rm -f src/canf22g2/grpc/wrappers.pb.go

cd src
GOOS=linux GOARCH=arm64 go build . 
