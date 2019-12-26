#!/bin/sh -eu
#docker run -it --rm -v "$PWD":/go/myapp -w /go/myapp -e CC=gcc -e GOARCH=amd64 -e CGO_ENABLED=1 -e GOOS=linux mygolang:16.04 go get -v ./...;go build -v -ldflags "-linkmode external -extldflags -static" -a -installsuffix cgo -o main .
#CC=gcc GOARCH=amd64 CGO_ENABLED=1 GOOS=linux go build -v -ldflags "-linkmode external -extldflags -static" -a -installsuffix cgo -o main .
sh build.sh
docker build -t dnn-storage-server:v1 .
REPO=dnn-registry.haha/system
docker tag dnn-storage-server:v1 $REPO/dnn-storage-server:v1
#docker push $REPO/dnn-storage-server:v1

#docker run -it --rm -v /data:/data dnn-storage-server:v1
