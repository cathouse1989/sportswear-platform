# Sportswear Portal - 客户门户展示站

## 项目定位

面向全球访客的 B2B 运动服饰 OEM/ODM 门户展示站，SSR 服务端渲染，SEO 优先。

- URL 结构：`/en/`, `/zh/`, `/es/`, `/fr/`（多语言子路径）
- 框架：Nuxt 3（SSR）
- 样式：Tailwind CSS
- 多语言：@nuxtjs/i18n（自动 hreflang 标签）

## 页面清单

| 页面 | 路由 | 后端 API |
|------|------|---------|
| 首页 | / | GET /public/home |
| 产品列表 | /products | GET /public/products |
| 产品详情 | /products/:slug | GET /public/products/:slug |
| 博客列表 | /blog | GET /public/blogs |
| 博客详情 | /blog/:slug | GET /public/blogs/:slug |
| 案例列表 | /cases | GET /public/cases |
| FAQ | /faq | GET /public/faqs |
| 联系我们 | /contact | POST /public/leads |

## 快速启动

```bash
# 安装依赖（网络不好时用镜像）
npm install --registry=https://registry.npmmirror.com

# 开发模式
npm run dev

# 构建
npm run build

# 预览构建产物
npm run preview
```

## 环境变量

创建 `.env` 文件：

```
NUXT_PUBLIC_API_BASE=http://localhost:8080/api/v1
```

## 技术栈

- Nuxt 3.15 + Vue 3.5
- @nuxtjs/i18n 9.x（4 种语言：en/zh/es/fr）
- ofetch（HTTP 客户端）
- Tailwind CSS