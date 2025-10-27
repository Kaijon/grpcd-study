#!/bin/sh

rm -rf gen 

mkdir -p gen/grpc
mkdir -p gen/swagger

protoc -I ./proto \
    --go_out ./gen/grpc --go_opt paths=source_relative \
    --go-grpc_out ./gen/grpc --go-grpc_opt paths=source_relative \
    --grpc-gateway_out ./gen/grpc --grpc-gateway_opt paths=source_relative \
    --openapiv2_out ./gen/swagger \
    ./proto/*.proto

cp -rf gen/grpc/* ./src/canf22g2/grpc
