FROM golang:1.17

RUN apt update && apt install protobuf-compiler -y
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod ./
# go.sum ./
RUN go mod download && go mod verify && go mod tidy

# Super slow to build, so we cache the result
RUN go get "gorm.io/driver/sqlite"

COPY . .
RUN protoc -I=proto --go_out=proto proto/*.proto
RUN go build -v -o ./app
RUN cp ./app /usr/local/bin/app

CMD ["app"]