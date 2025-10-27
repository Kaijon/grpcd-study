#!/bin/bash

# Remove original file *.go
find . -maxdepth 1 -type f -name "*.go" ! -name "*_test.go" -exec rm -f {} +

# Remove old xml
rm -f grpcd_unitTest.xml

pushd ../src/
./9_dockerCompile.sh

# move ../canf22g2 to current folder if the folder exists else break
sleep 5 # Wait for 5 seconds to give time for any concurrent operations to complete
if [ -d "../canf22g2" ]; then
    cp -r ../canf22g2 ../UnitTest
else
    echo "No canf22g2 folder found in parent directory"
    exit 1
fi

# copy ../services/prebuild_pkg.sh to current folder if the file exists else break
if [ -f "prebuild_pkg.sh" ]; then
    cp -r prebuild_pkg.sh ../UnitTest
else
    echo "No prebuild_pkg.sh file found in services directory"
    exit 1
fi


# copy ../src/config.go to current folder if the file exists else break
if [ -f "config.go" ]; then
    #find . -maxdepth 1 -type f -name "*.go" ! -name "main.go" -exec cp {} ../UnitTest \;
    find . -maxdepth 1 -type f -name "config.go" -exec cp {} ../UnitTest \;
    find . -maxdepth 1 -type f -name "ioctrl.go" -exec cp {} ../UnitTest \;
    find . -maxdepth 1 -type f -name "system.go" -exec cp {} ../UnitTest \;
    find . -maxdepth 1 -type f -name "network.go" -exec cp {} ../UnitTest \;
    find . -maxdepth 1 -type f -name "watermark.go" -exec cp {} ../UnitTest \;
    find . -maxdepth 1 -type f -name "video.go" -exec cp {} ../UnitTest \;
    find . -maxdepth 1 -type f -name "mqtt.go" -exec cp {} ../UnitTest \;
    find . -maxdepth 1 -type f -name "mqttCmd.go" -exec cp {} ../UnitTest \;
else
    echo "No config.go file found in services directory"
    exit 1
fi
popd #../src/

# don't publish mqtt
sed -i 's/MqttClient\.Publish/\/\/MqttClient.Publish/g' mqttCmd.go

# add copy to grpcd for path error
[ -d ./grpcd ] && rm -rf ./grpcd
mkdir -p ./grpcd
cp -rf canf22g2 grpcd/

# environment setup & Run ./app
docker run --rm  --name goUnitTest -v $(pwd):/go/src/app -w /go/src/app --privileged deviceenv.azurecr.io/golang:10235667 sh -c 'cp -r grpcd /usr/local/go/src && ./prebuild_pkg.sh && go test -v 2>&1 ./... | go-junit-report -set-exit-code > grpcd_unitTest.xml'
