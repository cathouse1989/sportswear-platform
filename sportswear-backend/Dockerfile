# ============ 构建阶段 ============
# BASE_REGISTRY 允许用户指定国内镜像仓库前缀，例如：
#   docker build --build-arg BASE_REGISTRY=docker.mirrors.ustc.edu.cn/library/ .
# 留空则直连 Docker Hub（需在 daemon.json 配置 registry-mirrors）。
ARG BASE_REGISTRY=
FROM ${BASE_REGISTRY}golang:1.25-alpine AS builder

WORKDIR /app

# 先复制依赖清单，利用 Docker 层缓存
COPY go.mod go.sum ./

# 设置 Go 模块代理（国内镜像），确保 `go mod download` 能正常获取依赖
ENV GOPROXY=https://goproxy.cn,direct \
    GOSUMDB=sum.golang.org \
    GONOSUMCHECK=1

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server

# ============ 运行阶段 ============
FROM ${BASE_REGISTRY}alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /server ./server

EXPOSE 8080

ENTRYPOINT ["./server"]