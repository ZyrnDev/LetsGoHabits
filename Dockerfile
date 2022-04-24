FROM golang:1.18

RUN apt update && apt install protobuf-compiler -y
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN wget "https://github.com/grpc/grpc-web/releases/download/1.3.1/protoc-gen-grpc-web-1.3.1-linux-x86_64" > /usr/local/bin/protoc-gen-grpc-web && chmod +x /usr/local/bin/protoc-gen-grpc-web

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod ./
# go.sum ./
RUN go mod download && go mod verify && go mod tidy

# Super slow to build, so we cache the result
RUN go get "gorm.io/driver/sqlite"
RUN go get "gorm.io/gorm"

ARG APP_DIR=./

COPY . .
RUN protoc -I=proto --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative proto/*.proto 
RUN go build -v -o ./app $APP_DIR
RUN cp ./app /usr/local/bin/

CMD ["app"]