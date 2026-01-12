FROM golang:1.25-alpine

WORKDIR /usr/local/src

ADD go.* /usr/local/src/ansi-convert/
ADD src/ansi-convert/ /usr/local/src/ansi-convert/src/ansi-convert/

RUN go mod tidy

ADD . .
