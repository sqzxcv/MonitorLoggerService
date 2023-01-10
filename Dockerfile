FROM golang:alpine as builder

WORKDIR /monitorloggerservice
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.io,direct \
    && go env -w CGO_ENABLED=0 \
#    && go env -w GOOS=linux \
#    && go env -w GOARCH=amd64 \
    && go env \
    && go mod tidy \
    && go build -o monitor_server .

FROM alpine:latest

LABEL MAINTAINER="test@qq.com"

WORKDIR /monitorloggerservice

COPY --from=0 /monitorloggerservice/monitor_server ./
#COPY --from=0 /monitorloggerservice/resource ./resource/
#COPY --from=0 /monitorloggerservice/config.docker.yaml ./

EXPOSE 80
ENTRYPOINT ./monitor_server
