# 国际化运动服饰 OEM/ODM 平台

基于 Go + Gin + PostgreSQL + Redis 构建的国际化 B2B 运动服饰数字化平台后端系统。

## 功能模块

### 商家管理后台（Admin API）
- 用户与权限管理（RBAC）
- 产品管理（产品、分类、系列、面料、MOQ、定制能力）
- 内容管理（页面、页面模块、导航、博客、案例、FAQ）
- 工厂管理、认证管理、生产流程管理
- 媒体资源管理
- 多语言管理
- SEO 管理
- 询盘管理（Lead 评分、跟进、分配销售）
- 数据统计（Dashboard）
- 操作日志

### 客户前台 API（Public API）
- 首页数据
- 产品列表/详情
- 分类、系列、面料
- 博客、案例、FAQ
- 工厂、认证、生产流程
- 导航
- 询盘提交

## 技术栈

| 组件 | 技术 |
|------|------|
| 语言 | Go 1.25+ |
| Web 框架 | Gin |
| ORM | GORM |
| 数据库 | PostgreSQL 15+ |
| 缓存 | Redis 7+ |
| 对象存储 | MinIO / S3 |
| 认证 | JWT + RBAC |

## 快速开始

### 1. 启动基础设施（Docker）

```bash
docker-compose up -d
```

### 2. 配置环境变量

```bash
cp .env.example .env
# 根据需要修改 .env 配置
```

### 3. 运行项目

```bash
go mod tidy
go run cmd/server/main.go
```

### 4. 默认管理员账号

- 邮箱：`admin@sportswear.com`
- 密码：`Admin@123456`

## API 文档

### 公开 API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/public/home | 首页数据 |
| GET | /api/v1/public/pages/:slug | 页面详情 |
| GET | /api/v1/public/products | 产品列表 |
| GET | /api/v1/public/products/:slug | 产品详情 |
| GET | /api/v1/public/categories | 分类列表 |
| GET | /api/v1/public/blogs | 博客列表 |
| GET | /api/v1/public/blogs/:slug | 博客详情 |
| GET | /api/v1/public/cases | 案例列表 |
| GET | /api/v1/public/faqs | FAQ 列表 |
| GET | /api/v1/public/fabrics | 面料列表 |
| GET | /api/v1/public/factories | 工厂列表 |
| GET | /api/v1/public/certifications | 认证列表 |
| GET | /api/v1/public/production-processes | 生产流程 |
| GET | /api/v1/public/navigations | 导航列表 |
| POST | /api/v1/public/leads | 提交询盘 |

### 后台 API

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/admin/auth/login | 登录 |
| GET | /api/v1/admin/auth/profile | 当前用户 |
| GET/POST | /api/v1/admin/users | 用户管理 |
| GET/POST | /api/v1/admin/roles | 角色管理 |
| GET | /api/v1/admin/permissions | 权限列表 |
| GET/POST/PUT/DELETE | /api/v1/admin/products | 产品管理 |
| GET/POST/PUT/DELETE | /api/v1/admin/categories | 分类管理 |
| GET/POST | /api/v1/admin/series | 系列管理 |
| GET/POST | /api/v1/admin/fabrics | 面料管理 |
| GET/POST/PUT/DELETE | /api/v1/admin/pages | 页面管理 |
| GET/POST/PUT/DELETE | /api/v1/admin/navigations | 导航管理 |
| GET/POST/PUT/DELETE | /api/v1/admin/blogs | 博客管理 |
| GET/POST/PUT/DELETE | /api/v1/admin/cases | 案例管理 |
| GET/POST/PUT/DELETE | /api/v1/admin/faqs | FAQ 管理 |
| GET/POST/PUT/DELETE | /api/v1/admin/factories | 工厂管理 |
| GET/POST/PUT/DELETE | /api/v1/admin/certifications | 认证管理 |
| GET/POST/PUT/DELETE | /api/v1/admin/production-processes | 生产流程 |
| GET/POST/PUT/DELETE | /api/v1/admin/media | 媒体管理 |
| GET/POST/PUT/DELETE | /api/v1/admin/leads | 询盘管理 |
| POST | /api/v1/admin/leads/:id/followups | 询盘跟进 |
| GET | /api/v1/admin/dashboard | 数据统计 |

## 项目结构

```
sportswear-platform/
├── cmd/
│   └── server/
│       └── main.go              # 入口
├── internal/
│   ├── config/                  # 配置加载
│   ├── database/                # 数据库连接
│   ├── models/                  # GORM 模型
│   ├── services/                # 业务逻辑层
│   ├── handlers/                # HTTP 处理器
│   ├── middleware/              # 中间件
│   ├── router/                  # 路由注册
│   └── utils/                   # 工具函数
├── docs/                        # 文档
├── docker-compose.yml           # Docker 基础设施
├── .env.example                 # 环境变量示例
└── README.md
```

## 开发计划

### 第一阶段（已完成）
- ✅ 项目骨架
- ✅ 数据库模型
- ✅ 用户认证与 RBAC
- ✅ 产品管理
- ✅ CMS 内容管理
- ✅ 询盘管理
- ✅ 媒体管理
- ✅ 公开 API
- ✅ 后台 API

### 第二阶段（待开发）
- CRM 客户管理
- 报价系统
- 高级搜索
- 数据统计增强
- AI 翻译
- 邮件自动化

### 第三阶段（待开发）
- 项目配置器
- 价格估算器
- AI 产品推荐
- 客户中心
- 订单管理