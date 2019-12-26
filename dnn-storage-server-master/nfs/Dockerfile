FROM alpine

RUN apk add --no-cache nfs-utils

ADD startup.sh /startup.sh

VOLUME /data

EXPOSE 111/udp 2049/tcp

CMD /startup.sh