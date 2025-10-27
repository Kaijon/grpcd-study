# proto definition and docker build/compile proto enviornment guide

## Confluence
- [System Architecture | EPIC(s) and Stories | Confluence](https://jira.getac.com/confluence/pages/viewpage.action?pageId=107819324)

## build environment
- Tool  Chain: https://gitlab.veretos.com/iot/toolchains/kl730BuildEnv
- Docker image: deviceenv.azurecr.io/kl730:v1.0
  ```
    docker login deviceenv.azurecr.io --username deviceenv --password 0wWSmpavQA6X9wjj1TRoqi+qDDcffSWR
    docker pull deviceenv.azurecr.io/kl730:v1.0
  ```

### Run PC simulator
- run gRPC server
    ```
    cd services
    0_pcSimulator.bat
    ``` 
- run bloomRPC App and import which <func>.proto you need.

### Building the firmware image 
**Use docker image** for protoc
- Enter docker
  ```
  docker run --name canf22build --rm -it -v "$(pwd):/coding" -w /coding  --privileged deviceenv.azurecr.io/kl730:v1.0 /bin/bash 
  ```
    - the **/coding** folder in docker will map to the current folder
- example package: /canf22g2/grpc
  - proto file
  ```
    package canf22g2.grpc;
    option go_package = "canf22g2/grpc";
  ```
  - golang
  ```
  import (
	pb "canf22g2/grpc"
    )
  ```
- run proto command directly in /coding folder for building proto file /canf22g2/grpc package
    `protoc --go_out=. --go-grpc_out=. *.proto`

    ::: Note: *you can change output folder path '.' and proto file name '*.proto'* :::


### Building web-gRPC javascript
- run gRPC server
    ```
    cd services
    0_dockerCompile_js.bat
    ```
- You'll get <func>.js files you need in proto folder.

### Generate protobuf go material (Including swagger and grpc-gateway)
- Build go building evironment
    ```
    docker compose build 
    ```
- run gen-proto.sh in the docker container
    ```
    docker compose run --remove-orphans --rm goGo ./gen-proto.sh
    ```

### Generate grpcd
- Build grpcd binary
    ```
    ./gen-grpcd.sh
    ``` 
