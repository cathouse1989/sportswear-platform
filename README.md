# sportswear
外贸运动服饰 OEM/ODM 数字化平台

```
sportswear-platform/
├── sportswear-backend/    ← 后端 API（Go/Gin）
├── sportswear-admin/      ← 管理后台（Vue3/Element Plus）
├── sportswear-portal/     ← 门户前端（Nuxt3 SSR）
├── docker-compose.yml     ← 统一编排
├── deploy.ps1             ← 一键部署（首次部署）
├── update.ps1             ← 一键更新（后续更新）
└── dev.ps1                ← 一键开发
```

## Docker 部署指南

### 首次部署
```powershell
powershell -ExecutionPolicy Bypass -File deploy.ps1
```

### 后续更新（推荐）
```powershell
powershell -ExecutionPolicy Bypass -File update.ps1
```

### 常用 Docker 命令

```powershell
# 查看所有容器状态
docker ps --filter name=sportswear

# 查看日志
docker compose -p sportswear -f ./docker-compose.yml logs -f

# 重启单个服务
docker compose -p sportswear -f ./docker-compose.yml restart backend

# 停止所有服务
docker compose -p sportswear -f ./docker-compose.yml stop

# 启动所有服务（不重新构建）
docker compose -p sportswear -f ./docker-compose.yml up -d

# 完全清理（包括数据卷）
docker compose -p sportswear -f ./docker-compose.yml down -v
```

## 访问地址

| 服务 | 地址 |
|------|------|
| 后端 API | http://localhost:8080 |
| 管理后台 | http://localhost:8090 |
| 门户前端 | http://localhost:3000 |
| MinIO | http://localhost:9001 |

## 默认账号

- **管理员**: admin@sportswear.com / Admin@123456
- **MinIO**: minioadmin / minioadmin
