# ============ 构建阶段 ============
# BASE_REGISTRY 允许用户指定国内镜像仓库前缀，例如：
#   docker build --build-arg BASE_REGISTRY=docker.mirrors.ustc.edu.cn/library/ .
# 留空则直连 Docker Hub（需在 daemon.json 配置 registry-mirrors）。
ARG BASE_REGISTRY=
FROM ${BASE_REGISTRY}node:22-alpine AS builder

WORKDIR /app

# 先安装依赖（利用 Docker 层缓存加速重复构建）
# 设置 npm 国内镜像源，确保依赖安装顺利
RUN npm config set registry https://registry.npmmirror.com

COPY package.json package-lock.json* ./
RUN npm ci 2>/dev/null || npm install

# 复制全部源码并构建
COPY . .
RUN npm run build

# ============ 运行阶段 ============
FROM ${BASE_REGISTRY}node:22-alpine

WORKDIR /app

# 从构建阶段复制输出
COPY --from=builder /app/.output ./.output

ENV NODE_ENV=production
ENV HOST=0.0.0.0
ENV PORT=3000

EXPOSE 3000

CMD ["node", ".output/server/index.mjs"]