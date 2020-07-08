# start a golang base image, version 1.8
FROM golang:1.14

#switch to our app directory
RUN mkdir -p /go/src/worker
WORKDIR /go/src/worker

#copy the source files
COPY . /go/src/worker

RUN go mod download 

#disable crosscompiling 
ENV CGO_ENABLED=0
ENV GO111MODULE=on

#compile linux only
ENV GOOS=linux

#build the binary with debug information removed
RUN go build  -ldflags '-w -s' -a -installsuffix cgo -o worker