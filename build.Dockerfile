# My Go Version
ARG GO_VERSION=1.18.1
ARG GO_OS=bullseye
FROM golang:$GO_VERSION-$GO_OS as go

ARG CGO_ENABLED
ENV CGO_ENABLED=$CGO_ENABLED

ARG GOOS
ENV GOOS=$GOOS

ARG GOARCH
ENV GOARCH=$GOARCH

WORKDIR /usr/src/app


# Add protoc
FROM go as golang_with_protoc
RUN apt update && apt install protobuf-compiler -y
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest


# Build Protobuf Specs
FROM golang_with_protoc as build_proto
COPY proto proto
RUN protoc -I=proto --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative proto/*.proto 


# Build the binary executables
FROM go as build

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go install github.com/mattn/go-sqlite3
RUN go mod verify
RUN go mod tidy

COPY . .
COPY --from=build_proto /usr/src/app/proto proto

RUN mkdir bin

RUN go build -tags osusergo,netgo -ldflags="-extldflags=-static" -v -buildvcs=false -o bin/engine  ./cmd/engine
RUN go build -tags osusergo,netgo -ldflags="-extldflags=-static" -v -buildvcs=false -o bin/handler ./cmd/handler