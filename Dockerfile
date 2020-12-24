FROM alpine:3.9

ARG VERSION

MAINTAINER biensupernice

ENV KRANE_DOWNLOAD_URL=https://github.com/biensupernice/krane/releases/download/${VERSION}/krane_${VERSION}_linux_386.tar.gz

RUN apk add curl ca-certificates

WORKDIR /bin

# Download binary
RUN curl -fSL $KRANE_DOWNLOAD_URL | tar xz && chmod +x krane

EXPOSE 8500
VOLUME ["/var/run/docker.sock"]

ENTRYPOINT krane
