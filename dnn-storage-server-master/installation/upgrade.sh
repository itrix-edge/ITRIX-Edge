#!/bin/sh -eu
 #dnn storage server
 systemctl stop dnn-storage-server.service
 yes | cp dnn-storage-server /usr/lib/dnn-storage-server
 systemctl start dnn-storage-server.service
