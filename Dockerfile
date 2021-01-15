FROM alpine:3.9

ARG KRANE_VERSION
ENV KRANE_DOWNLOAD_URL=https://github.com/krane/krane/releases/download/${KRANE_VERSION}/krane_${KRANE_VERSION}_linux_386.tar.gz

RUN apk add curl ca-certificates

WORKDIR /bin
RUN curl -fSL $KRANE_DOWNLOAD_URL | tar xz && chmod +x krane

EXPOSE 8500
VOLUME ["/var/run/docker.sock"]
ENTRYPOINT krane
