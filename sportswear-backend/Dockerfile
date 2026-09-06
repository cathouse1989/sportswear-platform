# ============ 构建阶段 ============
# BASE_REGISTRY 允许用户指定国内镜像仓库前缀，例如：
#   docker build --build-arg BASE_REGISTRY=docker.mirrors.ustc.edu.cn/library/ .
# 留空则直连 Docker Hub（需在 daemon.json 配置 registry-mirrors）。
ARG BASE_REGISTRY=
FROM ${BASE_REGISTRY}golang:1.25-alpine AS builder

WORKDIR /app

# 先复制依赖清单与 vendor 目录，使用 vendor 构建（避免容器内下载依赖受网络波动影响）
COPY go.mod go.sum ./
COPY vendor ./vendor

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -ldflags="-s -w" -o /server ./cmd/server

# ============ 运行阶段 ============
FROM ${BASE_REGISTRY}alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /server ./server

EXPOSE 8080

ENTRYPOINT ["./server"]