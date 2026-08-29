# ============ 构建阶段 ============
FROM node:22-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm ci || npm install

COPY . .

# 浏览器端 API 地址在运行时通过 NUXT_PUBLIC_API_BASE 注入
# SSR 服务端地址通过 NUXT_API_SERVER 注入（见 useApi.ts）
RUN npm run build

# ============ 运行阶段：Nitro SSR 服务 ============
FROM node:22-alpine

WORKDIR /app
COPY --from=builder /app/.output ./.output

ENV NODE_ENV=production
ENV HOST=0.0.0.0
ENV PORT=3000

EXPOSE 3000

CMD ["node", ".output/server/index.mjs"]