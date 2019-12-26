#!/bin/sh -eu
source ./repo_config
#docker build . -t vsftpd:v1
docker build -t dnn-ftp-server:v1 .
docker tag dnn-ftp-server:v1 $REPO/dnn-ftp-server:v1
#docker push $REPO/vsftpd:v1
#REPO=100.86.2.10:32190
#docker rm -f vsftpd
#docker run -d -v /nfs:/nfs -p 21:21 -p 20:20 --name vsftpd --restart always 100.86.2.10:32190/vsftpd:v1

