# This dockerfile builds the krane-server in multi-stages resulting in a lighter image
# The final image is alpine based containing only the krane-server executable and its working env

FROM golang:1.12-alpine AS base

# Dont cache locally, useful for keeping containers small.
RUN apk add --no-cache git

# Set current working directory inside the `golang-alpine` container
WORKDIR /tmp/krane-server

# Cache dependencies before building
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

COPY . .

# Build Go app
RUN go build -o /usr/local/bin/krane-server ./cmd/krane-server

# Start fresh from smaller image image
FROM alpine:3.9

# Install certs to establish secure communitcation via SSL/TLS
RUN apk add ca-certificates

# New working directory inside alpine container
WORKDIR /bin

# Copy executable from previous layer into this new layer which is smaller
# and contains only this executable. Since working directory is /bin
# we copy executable from previous location to working dir `.` 
COPY --from=base /usr/local/bin/krane-server .

# ENV inside the container, can also be passed in as flags with docker run
# ie. docker run -e KRANE_REST_PORT=9292 -p 9292:9292 krane-server
ENV KRANE_PATH "/.krane"
ENV KRANE_REST_PORT 8080
ENV KRANE_LOG_LEVEL "release"
ENV KRANE_PRIVATE_KEY "KbVHZLjpM3IUprwTSRvteRx+d8kmVecnEKvwAuJIaaw="

EXPOSE ${KRANE_REST_PORT}

# Create directory for krane server details
RUN mkdir $KRANE_PATH

ENTRYPOINT ["krane-server"]