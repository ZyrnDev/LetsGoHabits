FROM habits_go_base:latest

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go mod verify
RUN go mod tidy

# # Super slow to build, so we cache the result
RUN go get "gorm.io/driver/sqlite" && go install "gorm.io/driver/sqlite"
RUN go get "gorm.io/gorm" && go install "gorm.io/gorm"

COPY . .
RUN protoc -I=proto --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative proto/*.proto 

RUN mkdir bin
RUN go build -v -buildvcs=false -o bin/development  ./cmd/development
RUN go build -v -buildvcs=false -o bin/engine       ./cmd/engine
RUN go build -v -buildvcs=false -o bin/handler ./cmd/handler
RUN cp bin/* /usr/local/bin/

# Remove the transitive dependencies from build process
ARG EXEC_APP_NAME
ENV EXEC_APP_NAME=${EXEC_APP_NAME:-development}

CMD $EXEC_APP_NAME