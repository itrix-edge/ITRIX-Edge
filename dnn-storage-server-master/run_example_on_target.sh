docker run --name dnn-storage-server --restart=unless-stopped -d -p 8888:8888 --mount type=bind,src=/var/nfs,dst=/nfs dnn-registry.haha/system/dnn-storage-server:v1
