#!/bin/sh -eu

/usr/sbin/exportfs -r
/sbin/rpcbind --
/usr/sbin/rpc.statd
/usr/sbin/rpc.nfsd |:
/usr/sbin/rpc.mountd -F
