FROM golang:1.17

RUN apt update && apt install protobuf-compiler -y
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod ./
# go.sum ./
RUN go mod download && go mod verify && go mod tidy

# Super slow to build, so we cache the result
RUN go get "gorm.io/driver/sqlite"
RUN go get "gorm.io/gorm"

COPY . .
RUN protoc -I=proto --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative proto/*.proto 
RUN go build -v -o ./app
RUN cp ./app /usr/local/bin/app

CMD ["app"]