#!/bin/sh
FROM golang:1.17-alpine AS build

WORKDIR /root

RUN apk --no-cache add curl

COPY example-service example-service 
EXPOSE 3000
EXPOSE 3001

CMD ["./example-service","start"]




