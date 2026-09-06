# 门户「路径 × SEO × 导航」闭环联动改造说明

> 变更日期：2026-09；涉及子项目：sportswear-backend / sportswear-admin / sportswear-portal
> 无数据库结构变更（仅新增一份幂等 seed 数据）。

## 一、背景与目标

### 原有痛点
1. 「页面类型 → 路径」映射散落在 5 处硬编码（navigations.vue / pages.vue × 2 / seo.vue / portal default.vue），相互漂移。
2. 改页面 slug 不联动：关联导航 URL 不变、route SEO 不迁移 → 死链 + SEO 失效。
3. SEO 管理下拉为写死的 10 个路由，新建落地页后**不可见、无法配置**。
4. 导航 URL 手填无校验，历史遗留 `/oem /odm /factory` 死链靠 seed 脚本手动清理。
5. 页面/导航/SEO 三套数据靠人肉对齐，运营无「死链/缺导航/缺 SEO」的统一体检入口。

### 改造目标
以 **「路由（route）」为唯一锚点**，形成 `页面(slug) ↔ 导航(url) ↔ SEO(route)` 闭环，并提供可视化健康检查，降低运营使用门槛。

## 二、核心设计

### 1. 路由注册表（单一事实来源）
后端新增路由注册表，聚合「系统预置路由 + 已发布页面派生路由」，前端三处硬编码改为消费该接口。

### 2. 路径派生规则（与前台 Nuxt 路由对齐）
| 页面类型 | 路径 path | 路由 key route |
|---------|-----------|----------------|
| home | `/` | `home` |
| product / product_category | `/products` | `products` |
| blog | `/blog` | `blog` |
| case | `/cases` | `cases` |
| faq | `/faq` | `faq` |
| contact | `/contact` | `contact` |
| 其余（normal/seo_landing/oem/odm/...） | `/{slug}` | `{slug}` |

> 落地页统一走 `/{slug}`（门户 `[slug].vue` 兜底页），修正了原 `seo_landing: '/landing'` 等错误映射。

### 3. slug 联动
页面 slug/type 变更时，后端自动：① 更新「URL 仍等于旧派生路径」的关联导航；② 迁移 route SEO（目标不存在时才迁移，避免覆盖）。

### 4. 健康检查（4 类断点）
`dead_link` 导航死链 / `missing_nav` 缺导航 / `missing_seo` 缺 SEO / `nav_drift` 路径漂移。

## 三、改动清单

### 后端（sportswear-backend）
| 文件 | 改动 |
|------|------|
| `internal/services/portal_health.go`（新增） | `PagePath`/`RouteKey`/`RouteFromPage` 路径派生；`ListPortalRoutes()` 路由注册表；`PortalHealthCheck()` 健康检查；`SyncPagePathEffects()` slug 联动 |
| `internal/handlers/cms_handler.go` | 新增 `ListPortalRoutes`、`PortalHealthCheck` handler |
| `internal/router/router.go` | 注册 `GET /admin/portal/routes`、`GET /admin/portal/health` |
| `internal/services/cms_service.go` | `UpdatePage` 在 slug/type 变更时级联联动 |

### 前端（sportswear-admin）
| 文件 | 改动 |
|------|------|
| `src/constants/routes.ts`（新增） | 统一 `pagePath`/`routeFromPage`/`PRESET_ROUTES`（接口兜底） |
| `src/api/index.ts` | 新增 `portalHealthApi`（listRoutes / healthCheck） |
| `src/types/index.ts` | 新增 `PortalRouteEntry`/`HealthIssue`/`PortalHealthReport` 类型 |
| `src/views/system/seo.vue` | SEO 下拉动态化 + 支持 `?route=xx` 定位 |
| `src/views/cms/navigations.vue` | 统一 `pagePath`；新增「SEO」状态列 |
| `src/views/cms/pages.vue` | 统一 `pagePath`；slug 变更提示 |
| `src/views/system/health.vue`（新增） | 门户闭环健康检查页 |
| `src/router/index.ts` + `src/config/menu.ts` | 恢复「页面管理」入口；新增「门户健康检查」路由与菜单 |

### 门户（sportswear-portal）
| 文件 | 改动 |
|------|------|
| `layouts/default.vue` | 清理 `NAV_KEY_BY_PATH` 中已废弃死链 `/oem /odm /factory` |

## 四、新增接口

| 方法 | 接口 | 权限（any-of） | 说明 |
|------|------|----------------|------|
| GET | `/api/v1/admin/portal/routes` | navigation:manage / seo:manage / page:view | 路由注册表：预置 + 已发布页面派生，含 `seo_configured`、`nav_count` |
| GET | `/api/v1/admin/portal/health` | navigation:manage / seo:manage / page:view | 健康报告：`total_issues` + 4 类 `issues` |

响应示例（routes 条目）：
```json
{ "route": "products", "path": "/products", "label": "产品中心", "source": "system", "seo_configured": true, "nav_count": 2 }
```
响应示例（health）：
```json
{ "total_issues": 3, "routes": 12, "nav_count": 14, "page_count": 8,
  "issues": [ { "type": "missing_seo", "title": "路由未配置 SEO", "route": "privacy", "path": "/privacy-policy" } ] }
```

## 五、数据库 seed（幂等，可重复执行）

`seed-route-seo.sql`（项目根目录）为 8 个预置路由插入英文源语言初始 route SEO，`entity_id` 为确定性 UUID v5（与后端 `RouteSEOID` 一致，**勿手工修改**）。

执行方式（任选其一）：
```bash
# 方式一：psql
psql "$DATABASE_URL" -f seed-route-seo.sql

# 方式二：docker
docker exec -i <postgres容器> psql -U <user> -d <db> < seed-route-seo.sql
```

执行后「门户健康检查」的「缺 SEO」问题初始为 0，运营可在「SEO 管理」补充中文/西语/法语。

## 六、上线步骤

1. **后端**：重启（使 `/admin/portal/routes`、`/admin/portal/health` 生效）。
2. **前端**：重新构建 / 重启 dev server。
3. **数据库**：执行 `seed-route-seo.sql`（可选，用于初始 route SEO）。
4. **冒烟**：见下方验证清单。

## 七、验证清单

- [ ] 后端 `go build ./...` 无错误。
- [ ] 前端 `npm run typecheck`（vue-tsc）通过。
- [ ] 前端 `npm run lint` 无新增 error。
- [ ] 「系统配置 → 门户健康检查」可打开，统计卡片与 4 类问题正常。
- [ ] 「内容管理 → 页面管理」可进入（入口已恢复），改 slug 保存后提示「已联动更新关联导航与 SEO 路由」。
- [ ] 「SEO 管理」下拉展示预置 + 已发布页面派生的全部路由；从健康检查点「去配置 SEO」可跳转并定位到对应路由。
- [ ] 「导航管理」新增「SEO」状态列，正确显示已配置/未配置。

## 八、回滚说明

- 无数据库结构变更，回滚代码即可；`seed-route-seo.sql` 写入的 route SEO 数据可保留（不影响旧版），如需清理：
  ```sql
  DELETE FROM seos WHERE entity_type = 'route';
  ```

## 九、使用指南（运营视角）

1. **上线新落地页**：页面管理 → 新建（选「SEO 落地页」类型、填 slug）→ 发布 → 导航管理「一键同步」或页面列表「加入导航」→ SEO 管理下拉已自动出现该 slug 路由。
2. **日常巡检**：打开「门户健康检查」点「重新检查」，按「死链/缺导航/缺SEO/路径漂移」逐类点「去修复」即可，无需理解三层数据关系。
3. **改 slug**：页面编辑里直接改 slug 保存，系统自动联动导航 URL 与 SEO 路由。
