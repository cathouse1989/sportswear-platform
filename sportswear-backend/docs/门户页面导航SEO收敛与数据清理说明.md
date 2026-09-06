# 门户「页面 / 导航 / SEO」收敛与数据清理说明

> 变更日期：2026-09；涉及子项目：sportswear-backend / sportswear-admin / sportswear-portal
> 无数据库结构变更；含一份幂等数据清理脚本。

## 一、背景与目标

历史版本允许后台自由新增/删除页面与导航，导致门户出现以下问题：

1. 新增/删除页面会产生死链或空页面（门户 `[slug].vue` 兜底路由 + 导航 URL 直出）。
2. 导航重复、拼写错误、footer 缺失、指向不存在路径等脏数据。
3. 产品/博客/案例等实体创建时若不手动填 SEO 就无 SEO 记录。
4. 工厂管理/生产流程与门户零数据交互，属"悬空功能"却占菜单位。

**目标**：把「结构（路由/导航骨架）」与「内容（产品/博客/案例/落地页）」彻底分开——路由与导航骨架由代码/seed 维护；内容型数据新增时自动派生 SEO；后台不再暴露能破坏门户的自由增删能力。

## 二、核心方案（三层）

| 层 | 内容 | 维护方式 | 后台开放 |
|---|---|---|---|
| 系统预置 | 固定路由 + 固定导航骨架 | 代码常量 + 幂等 seed | 只读，可编辑文案/排序/显隐 |
| 详情实体 | 产品/博客/案例/分类详情 | Nuxt 固定路由 + 自动 SEO | 不生成 Page/导航；SEO 集中维护 |
| 落地页 | 自定义 slug 落地页 | 发布时自动派生 Page+导航+route SEO | 按需新建，派生全自动 |

## 三、改动清单

### 后端（sportswear-backend）

| 文件 | 改动 |
|---|---|
| `internal/services/seo_service.go` | 新增 `EnsureEntitySEO`（幂等自动生成实体英文 SEO）、`truncateText`、`ListEntitySEO` |
| `internal/services/product_service.go` | `CreateProduct` 未传 SEO 时自动生成默认值 |
| `internal/services/cms_service.go` | `CreatePage`/`CreateBlog`/`CreateCase` 未传 SEO 时自动生成默认值；新增 `ListPublishedLandingPages` |
| `internal/services/portal_health.go` | 新增 `missing_entity_seo` 断点；`missing_nav` 支持 URL 匹配（固定页面经骨架导航匹配） |
| `internal/handlers/public_handler.go` | 新增 `ListLandingPages` |
| `internal/handlers/localization_handler.go` | 新增 `ListEntitySEO`/`UpsertEntitySEO` |
| `internal/router/router.go` | 移除页面/导航 `POST/DELETE` 增删路由；新增 `GET /public/pages`、`GET/POST /admin/seo/entity` |

### 前端（sportswear-admin）

| 文件 | 改动 |
|---|---|
| `src/views/cms/pages.vue` | 隐藏「新建/删除/加入导航」，清理死代码 |
| `src/views/cms/navigations.vue` | 隐藏「新建/添加子级/删除/一键同步/同步页面状态」，清理死代码 |
| `src/views/system/seo.vue` | 新增「路由 / 实体详情」双维度 SEO 管理 |
| `src/api/index.ts` | `seoApi` 扩展 `listEntities`/`saveEntities` |
| `src/config/menu.ts` / `src/router/index.ts` | 下架「工厂管理 / 生产流程」（保留「认证管理」） |

### 门户（sportswear-portal）

| 文件 | 改动 |
|---|---|
| `server/routes/sitemap.xml.ts` | 补落地页进 sitemap |

## 四、数据清理

### 已清理的脏数据（`cleanup-dirty-data.sql`，项目根目录，幂等）

| # | 脏数据 | 清理动作 |
|---|---|---|
| 1 | header 导航重复（固定骨架 + 「一键同步」page_id 关联项并存，14 条） | 删除 7 条 page_id 关联项，保留固定骨架 |
| 2 | footer 导航缺失 | 补齐 7 条 |
| 3 | contact 页面标题拼写错误 `Contacst Us` | 修正为 `Contact Us` |
| 4 | seos 空占位记录（空 language / 空 title） | 删除 8 条 |
| 5 | 9 个已发布产品缺英文 SEO | 批量补默认（title=英文名，description=brief） |
| 6 | 4 个已发布博客缺英文 SEO | 批量补默认（title=标题，description=正文摘要） |

### 执行方式

```bash
docker exec -i sportswear-postgres psql -U postgres -d sportswear_platform < cleanup-dirty-data.sql
# 或
psql "$DATABASE_URL" -f cleanup-dirty-data.sql
```

## 五、上线步骤

1. **数据库**：执行 `cleanup-dirty-data.sql`（幂等，可重复执行）。
2. **后端**：重新构建镜像并重启，使健康检查 `missing_nav` URL 匹配逻辑生效：
   ```bash
   docker compose up -d --no-deps --build backend
   ```
3. **前端**：重新构建 / 重启 dev server。
4. **门户**：重新构建（sitemap 落地页生效）。

## 六、验收清单

- [ ] 后台「页面管理」「导航管理」无「新建/删除」入口。
- [ ] 后台「内容管理」菜单无「工厂管理 / 生产流程」，仍有「认证管理」。
- [ ] 新增产品/博客/案例 → 「SEO 管理」自动出现对应实体 SEO。
- [ ] 「SEO 管理」支持「路由」与「实体详情」两维度编辑。
- [ ] 「门户健康检查」无死链/缺导航/缺 SEO（含 `missing_entity_seo`）。
- [ ] 门户 header/footer 各 7 项、顺序一致、无重复、拼写正确。
- [ ] sitemap 含落地页。
