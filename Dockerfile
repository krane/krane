FROM alpine

WORKDIR ~/.krane/

COPY bootstrap.sh .

COPY usr/local/bin/krane-server .

RUN chmod 755 krane-server

CMD ["./krane-server"]