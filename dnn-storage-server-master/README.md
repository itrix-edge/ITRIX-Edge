# dnn-storage-server

## Overview

dnn-storage-server project contain total three of docker iamges:

* **macchiang/mygolang:1.8.3**: Go lang build docker image, used for storage server build environment. Sources live under `mygolang/`.
* **dnn-ftp-server:v1**: VSFTPd container image, for containerized FTP server. Sources live under `vsftpd/`.
* **dnn-storage-server:v1**: Containerized stroage node controller server written in Go. Sources are under project root.

Currently, those images are required to run on top of storage node, with NFS server pre-installed. Kubernetes integration are currently unsupported.

## Install Requirement

Storage node (accessible from target system environment) with following conditions:

1. Docker with version 18.06 and up.
2. NFS server configurated. See NFS server setting for detailed information.
3. For online container image dispatcher during installation, accessible container registry with installer container image preload is required. You can use public container registry or use private registry like `registry` or `Harbor` if installation under private environment.


## Firewall information
This project uses following port settings with incoming connections(all of ports are TCP connecitons):

* 8888 for storage server
* 20,21,13001-13100 for FTP protocol. Port range during 13001~13100 are used for FTP passive mode.


## Development & Build Instruction

Entrypoint is `build.sh`, run with following command:

```=shell
$ ./build.sh all
```

If system raise `macchiang/mygolang:1.8.3` was not found, build it from source under `mygolang/`.
If error `ld: connot find -lpthread` occurred, install c static library `glibc-static` and try again. On RHEL platform, use following command:

```=shell
$ yum install glibc-static
```



### Advanced: Installation dnn storage server without container
<pre>
cd installation
sh install.sh
</pre>
