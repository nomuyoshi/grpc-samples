FROM golang:1.16 AS build-env

ENV PKG_PATH /go/src/github.com/nomuyoshi/grpc-samples/todo-microservices
ENV SERVICE_PATH ${PKG_PATH}/project

WORKDIR ${PKG_PATH}
COPY go.mod go.sum ./
RUN go mod download
COPY . ${PKG_PATH}

# ビルド
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go install ${SERVICE_PATH}

FROM alpine:3.13
ENV TZ Asia/Tokyo
COPY --from=build-env /go/bin/project /go/bin/project
ENTRYPOINT ["/go/bin/project"]
