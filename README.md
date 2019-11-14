# Demo

Golang demo project

## First

1. 下载代码后，请执行，`go mod tidy` 检查依赖包情况，会自动下载符合版本的依赖包。
2. 执行 `make run` 启动 `Server` 服务

## Docker

编译镜像

```bash
make docker
```

运行镜像

```bash
docker run --rm -p 8000:8000 demo 
```

## Jaeger

开发环境中启用 Jaeger 服务

```bash
docker run --rm --name jaeger \
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.14
```

> 注意: 添加 `--rm` 容器停止后即删除