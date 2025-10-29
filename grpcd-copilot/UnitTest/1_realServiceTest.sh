#!/bin/bash

# if no docker image golang:10235667, build it
if [ "$(docker images -q golang:10235667 2> /dev/null)" == "" ]; then
    docker build -t golang:10235667 -f Dockerfile.10235667 .
fi

# Ensure no container named "goUnitTest" is running or exists
docker ps -aq -f name=goUnitTest | xargs -r docker rm -f

# compile proto files
docker run -d --rm --name goUnitTest  -v $(pwd)/..:/go/src/app -w /go/src/app golang:10235667 sh -c 'protoc --go_out=. --go-grpc_out=. *.proto'

# move ../canf22g2 to current folder if the folder exists else break
sleep 5 # Wait for 5 seconds to give time for any concurrent operations to complete
if [ -d "../canf22g2" ]; then
    cp -r ../canf22g2 .
else
    echo "No canf22g2 folder found in parent directory"
    exit 1
fi

# Ensure no container named "goUnitTest" is running or exists before running the app
docker ps -aq -f name=goUnitTest | xargs -r docker rm -f

# copy ../services/prebuild_pkg.sh to current folder if the file exists else break
if [ -f "../services/prebuild_pkg.sh" ]; then
    cp -r ../services/prebuild_pkg.sh .
else
    echo "No prebuild_pkg.sh file found in services directory"
    exit 1
fi

# environment setup & Run ./app
docker run -it --name goUnitTest -v $(pwd):/go/src/app -w /go/src/app golang:10235667 sh -c 'cp -r canf22g2 /usr/local/go/src && ./prebuild_pkg.sh && go run test_real_service.go'