# ============ 构建阶段 ============
FROM node:22-alpine AS builder

WORKDIR /app

# 先安装依赖（利用 Docker 层缓存加速重复构建）
COPY package.json package-lock.json* ./
RUN npm ci 2>/dev/null || npm install

# 复制全部源码并构建
COPY . .
RUN npm run build

# ============ 运行阶段 ============
FROM node:22-alpine

WORKDIR /app

# 从构建阶段复制输出
COPY --from=builder /app/.output ./.output

ENV NODE_ENV=production
ENV HOST=0.0.0.0
ENV PORT=3000

EXPOSE 3000

CMD ["node", ".output/server/index.mjs"]