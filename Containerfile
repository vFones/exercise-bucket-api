FROM golang:alpine

RUN apk update && \
    apk upgrade && \
    apk add --no-cache curl

WORKDIR /app

COPY ./main /app/main

EXPOSE 8080

ENTRYPOINT ["./main"]
