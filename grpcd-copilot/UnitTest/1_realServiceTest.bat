@echo off

:: check if image golang:10235667 existed，if not, build it.
docker images -q golang:10235667 2> nul
:: check Dockerfile.10235667
if exist "..\src\Dockerfile.10235667" (
    docker build -t golang:10235667 -f ..\src\Dockerfile.10235667 .
) else (
    echo Dockerfile.10235667 not found in ..\src directory.
    exit /b 1
)

:: make sure goUnitTest container is not running or existed
FOR /f "tokens=*" %%i IN ('docker ps -aq --filter "name=goUnitTest"') DO docker rm -f %%i

:: Compile proto documentation
docker run -d --rm --name goUnitTest -v %cd%/..:/go/src/app -w /go/src/app golang:10235667 sh -c "protoc --go_out=. --go-grpc_out=. *.proto"

timeout /t 5
:: if ../canf22g2 folder exist, move here else exist
if exist "..\canf22g2" (
    xcopy /E /I /Y "..\canf22g2" ".\canf22g2"
) else (
    echo No canf22g2 folder found in parent directory
    exit /b 1
)

:: make sure goUnitTest container is not running or existed
FOR /f "tokens=*" %%i IN ('docker ps -aq --filter "name=goUnitTest"') DO docker rm -f %%i

:: if ../src/go.mod exist, copy here else exist
if exist "..\src\go.mod" (
    :: if ../src/go.sum exist, copy here else exist
    if exist "..\src\go.sum" (
        copy /Y "..\src\go.mod" ".\"
        copy /Y "..\src\go.sum" ".\"
    ) else (
        echo No go.sum file found in ../src directory
        exit /b 1
    )
) else (
    echo No go.mod file found in ../src directory
    exit /b 1
)

:: environment setup & Run ./app
docker run -it --name goUnitTest -v %cd%:/go/src/app -w /go/src/app golang:10235667 sh -c "cp -r canf22g2 /usr/local/go/src && go run test_real_service.go"