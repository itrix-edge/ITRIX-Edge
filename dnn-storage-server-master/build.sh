#!/bin/bash -eu
docker run -it --rm -v "$PWD":/go/myapp -w /go/myapp -e CC=gcc -e GOARCH=amd64 -e CGO_ENABLED=1 -e GOOS=linux macchiang/mygolang:1.8.3 sh build-within-docker.sh
yes | cp dnn-storage-server installation/
yes | cp vsftpd/repo_config installation/
if [ "$1" == "all" ]; then
    cd vsftpd
    sh build_push.sh
    cd ..
elif [ "$1" == "upgrade" ]; then
    cd installation
    sh upgrade.sh
fi