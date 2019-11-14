FROM golang:latest as builder
ENV GO111MODULE=on GOPROXY=https://mirrors.aliyun.com/goproxy/
WORKDIR /build
COPY ./ /build
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X 'demo/api/controller/v1.BuildTime=`date +"%Y-%m-%d %H:%M:%S"`' -X demo/api/controller/v1.BuildVersion=1.0.1" -tags=jsoniter -o demo .

FROM alpine:latest
LABEL maintainer="dingdayu <614422099@qq.com>"
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk --no-cache add ca-certificates tzdata curl
ENV DREAMENV=TEST DEPLOY_TYPE=DOCKER
WORKDIR /opt/app/demo
COPY --from=0 /build/demo .
COPY --from=0 /build/config/config.release.yaml /opt/app/demo/config/config.yaml
EXPOSE 8000
HEALTHCHECK --interval=5s --timeout=3s \
  CMD curl -fs -X HEAD http://127.0.0.1:8000/health || exit 1
ENTRYPOINT ["/opt/app/demo/demo", "server"]
