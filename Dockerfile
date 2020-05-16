FROM golang:alpine3.11

# api env variables
ENV KRANE_PATH "/.krane"

# rest api port 
ENV PORT 9000
EXPOSE 9000

WORKDIR /app

COPY . .

# Cleanup/Install dependencies
RUN cd cmd/krane-server && go mod tidy && go get

RUN chmod +x "build.sh" 

ENTRYPOINT ["sh", "build.sh"]

CMD ["start"]
