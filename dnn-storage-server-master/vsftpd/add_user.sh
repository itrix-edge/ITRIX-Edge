USER=$1
PASSWD=$2
mkdir -p /nfs/$1
chown -R vsftpd:vsftpd /nfs/$1
echo "$1:$(openssl passwd -1 $2)" >> /nfs/passwd
