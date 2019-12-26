docker build . -t nfs-server:v1;docker rm -f nfs-server ;docker run -d --privileged -v /data:/data -v /data/exports:/etc/exports -p 111:111/udp -p 2049:2049/tcp --name nfs-server nfs-server:v1
