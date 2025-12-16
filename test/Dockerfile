FROM golang:1.25-alpine

WORKDIR /usr/local/src

ADD go.* /usr/local/src/
ADD src/ /usr/local/src/src/

RUN go mod tidy

ADD . .
