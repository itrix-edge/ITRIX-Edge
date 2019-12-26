#!/bin/sh -eu
DIR="/nfs"
#REPO="100.86.2.10:32190"
source repo_config
apt-get install nfs-kernel-server nfs-common docker.io -y
mkdir -p $DIR/utils
mkdir -p $DIR/logs
chmod -R 777 $DIR
mv /etc/exports $DIR/exports
echo "$DIR/utils *(ro,no_subtree_check,sync,all_squash,anonuid=0,anongid=0)" >> $DIR/exports
echo "$DIR/logs *(rw,no_subtree_check,sync,all_squash,anonuid=0,anongid=0)" >> $DIR/exports
ln -s $DIR/exports /etc/exports 
/etc/init.d/nfs-kernel-server restart
ufw allow from any to any port nfs

# REPO certs
 mkdir -p /etc/docker/certs.d/$REPO
 cp  ca.crt  /etc/docker/certs.d/$REPO

 # vsftpd
 docker run -d -v $DIR:/nfs -p 21:21 -p 20:20 --name vsftpd --restart always $REPO/vsftpd:v1
 ufw allow ftp

 #dnn storage server
 mkdir -p /usr/lib/dnn-storage-server
 cp dnn-storage-server /usr/lib/dnn-storage-server
 cp dnn-storage-server.service /etc/systemd/system/
 systemctl daemon-reload
 systemctl enable dnn-storage-server.service
 systemctl start dnn-storage-server.service
 #nohup /usr/lib/dnn-storage-server/dnn-storage-server &

 ufw allow from any to any port 8888 tcp