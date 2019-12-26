#!/bin/sh -eu
ACCOUNT=$1
if grep -q $1 "/nfs/passwd"; then
   echo "fail"
   exit 0
fi
if grep -q $1 "/nfs/exports"; then
    echo "fail"
    exit 0
fi
if [ -d "/nfs/$1" ]; then
    echo "fail"
    exit 0
fi
echo "OK"