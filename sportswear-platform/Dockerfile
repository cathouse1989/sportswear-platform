# ============ 构建阶段 ============
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 先复制依赖清单，利用 Docker 层缓存
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server

# ============ 运行阶段 ============
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /server ./server

EXPOSE 8080

ENTRYPOINT ["./server"]