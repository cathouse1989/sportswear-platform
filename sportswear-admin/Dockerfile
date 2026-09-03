# ============ 构建阶段 ============
# BASE_REGISTRY 允许用户指定国内镜像仓库前缀，例如：
#   docker build --build-arg BASE_REGISTRY=docker.mirrors.ustc.edu.cn/library/ .
# 留空则直连 Docker Hub（需在 daemon.json 配置 registry-mirrors）。
ARG BASE_REGISTRY=
FROM ${BASE_REGISTRY}node:22-alpine AS builder

WORKDIR /app

# 设置 npm 国内镜像源，确保依赖安装顺利
RUN npm config set registry https://registry.npmmirror.com

COPY package.json package-lock.json* ./
RUN npm ci || npm install

COPY . .
RUN npm run build

# ============ 运行阶段：nginx 托管静态文件 + API 反代 ============
FROM ${BASE_REGISTRY}nginx:1.27-alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]