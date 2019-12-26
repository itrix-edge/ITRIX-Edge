FROM alpine
MAINTAINER Stevennick Ciou <yfciou@itri.org.tw>

RUN set -xe \
    && apk add -U build-base \
    curl \
    linux-pam-dev \
    tar \
    vsftpd \
    openssl \
    db-utils \
    && mkdir pam_pwdfile \
    && cd pam_pwdfile \
    && curl -sSL https://github.com/tiwe-de/libpam-pwdfile/archive/v1.0.tar.gz | tar xz --strip 1 \
    && make install \
    && cd .. \
    && rm -rf pam_pwdfile \
    && apk del build-base \
    curl \
    linux-pam-dev \
    tar \
#    && passwd -l root \
    && rm -rf /var/cache/apk/* \
    && mkdir /data \
    && adduser -s /sbin/nologin -h /data vsftpd -D \
    #&&  addgroup vsftpd \
    && chown -R vsftpd:vsftpd /data \
    && touch /etc/vsftpd.user_list

ADD vsftpd /etc/pam.d/vsftpd	
ADD vsftpd.conf /etc/vsftpd/vsftpd.conf
#ADD vsftpd.pem /etc/pki/tls/certs/vsftpd.pem
#ADD add_user.sh /var/lib/ftp/

VOLUME /nfs

WORKDIR /var/lib/ftp

EXPOSE 21
EXPOSE 20
EXPOSE 13001-13100

CMD ["vsftpd", "/etc/vsftpd/vsftpd.conf"]
