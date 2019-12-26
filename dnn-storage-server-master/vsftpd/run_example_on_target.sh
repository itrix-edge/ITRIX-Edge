docker run --name dnn-ftp-server --restart=unless-stopped -d -p 20:20 -p 21:21 -p 13001-13100 --mount type=bind,src=/var/nfs,dst=/nfs dnn-registry.haha/system/dnn-ftp-server:v1
