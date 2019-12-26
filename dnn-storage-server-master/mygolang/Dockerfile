#Build environment for my environment
FROM ubuntu:17.04
MAINTAINER Mac Chiang <macchiang@gmail.com>
RUN apt-get update && apt-get install -y --no-install-recommends \
		g++ \
		gcc \
		libc6-dev \
		make \
		pkg-config \
        wget \
        git \
		tzdata \
		ca-certificates \
	&& rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.8.3

RUN set -eux; \
	\
# this "case" statement is generated via "update.sh"
	dpkgArch="$(dpkg --print-architecture)"; \
	case "${dpkgArch##*-}" in \
		amd64) goRelArch='linux-amd64' ;; \
		armhf) goRelArch='linux-armv6l';; \
		arm64) goRelArch='linux-arm64';; \
		i386) goRelArch='linux-386';; \
		ppc64el) goRelArch='linux-ppc64le';; \
		s390x) goRelArch='linux-s390x';; \
	esac; \
	\
	url="https://golang.org/dl/go${GOLANG_VERSION}.${goRelArch}.tar.gz"; \
	wget -O go.tgz "$url" --no-check-certificate; \
	tar -C /usr/local -xzf go.tgz; \
	rm go.tgz; \
	export PATH="/usr/local/go/bin:$PATH"; \
	ln -fs /usr/share/zoneinfo/Asia/Taipei /etc/localtime; \
	go version

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH

CMD ["bash"]