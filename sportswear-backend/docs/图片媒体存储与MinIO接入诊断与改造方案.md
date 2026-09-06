# 图片媒体存储与 MinIO 接入 · 诊断报告与改造方案

> 状态：已按方案 C 落地核心闭环（DB + MinIO + 展示数据一致性），详见文末「十一、落地实现记录」
> 范围：媒体管理（`/media`、`/media/upload`）、存储源（`/storage-sources`）、门户图片展示链路
> 关联文档：`sportswear-portal/docs/requirements/系统需求文档.md` §4.5 / §4.6、`sportswear-admin/docs/管理后台产品说明书-v1.md`（B-08）

---

## 一、结论（TL;DR）

**目前图片管理并不是读取 MinIO，而是本地磁盘 `uploads/` 目录。**

- 上传链路：`SaveToLocal` 把文件写到 `uploads/<category>/<年月>/<uuid>.<ext>`，URL 生成相对路径 `/uploads/...`，`source` 字段**硬编码为 `local`**。
- 读取链路：后端 `r.Static("/uploads", "./uploads")` 静态服务 + 门户 Nitro 反代 `server/routes/uploads/[...path].ts`，**全部从本地磁盘读取**。
- MinIO 只存在于配置（`MINIO_*` 环境变量）、`docker-compose.yml`、需求文档里；`go.mod` **没有 minio-go SDK**；唯一会拼 MinIO 地址的 `UploadService.PublicURL()` 是**死代码，从未被调用**。

换句话说：需求里设计的「local / minio / external 三种存储源 + 路径解析」这套能力，目前**只做了一半**（枚举、模型、service 都写好了），但**上传/列表读取链路没有接入**，`source`、`path`、`url` 三者与实际存储、存储源配置之间**没有形成闭环**。

---

## 二、现状诊断（逐环节核对）

### 2.1 上传链路（全走本地，未走 MinIO）

`internal/handlers/media_handler.go:95-158` `UploadFile`：

1. `FormFile("file")` 取文件 → `ValidateFile` 白名单 + 大小校验。
2. `GenerateStoragePath` 生成 `category/2006-01/uuid.ext` 路径。
3. `uploadService.SaveToLocal(src, storagePath)` **直接落本地磁盘**（`upload_service.go:108` 写到 `uploads/<path>`）。
4. `URL: h.uploadService.LocalURL(storagePath)` → 相对路径 `/uploads/<path>`（`upload_service.go:117`）。
5. `Source: models.MediaSourceLocal` **硬编码 local**（`media_handler.go:147`）。

`UploadService` 里明明有 `PublicURL()`（`upload_service.go:98`，拼 `http://endpoint/bucket/path` 的 MinIO 风格 URL），但**无任何调用点**；也没有 MinIO 客户端/SDK。

> 询盘附件上传同理：`internal/handlers/public_handler.go:704` 也走 `SaveToLocal` + `LocalURL`，与 MinIO 无关。

### 2.2 读取/展示链路（本地静态服务 + 反代）

| 端 | 位置 | 说明 |
|---|---|---|
| 后端静态服务 | `router.go:27` `r.Static("/uploads", "./uploads")` | 从本地磁盘读 |
| 管理后台代理 | `sportswear-admin/vite.config.ts:33` `/uploads → http://localhost:8080` | 预览/选择器经此代理 |
| 门户反代 | `sportswear-portal/server/routes/uploads/[...path].ts` | 同源反代到后端 `/uploads` |
| 门户缓存头 | `sportswear-portal/nuxt.config.ts:137` `/uploads/**` CDN 缓存 1 天 | 仅针对本地相对路径 |

前端归一化逻辑 `sportswear-portal/composables/useImageUrl.ts` 只认识 `localhost/backend/host.docker.internal` 三种后端地址，**没有 MinIO 地址的归一化规则**。

### 2.3 MinIO 配置与代码完全脱节

- `config.go:98-104` 读取了 `MINIO_ENDPOINT / MINIO_ACCESS_KEY / MINIO_SECRET_KEY / MINIO_BUCKET / MINIO_USE_SSL`，但 `UploadService` 只把 endpoint/useSSL 存到 struct 字段，**从未用于发起连接或读写**。
- `docker-compose.yml:39` 起了 `minio` 容器（端口 9000/9001），`deploy.ps1`/`dev.ps1` 也提示 MinIO 账号，但**后端进程从不连接它**。
- `go.mod` 无 `github.com/minio/minio-go/v7` 依赖。

### 2.4 已存在但「未接通」的抽象能力

| 抽象 | 位置 | 现状 |
|---|---|---|
| `MediaSource` 枚举 local/minio/external | `models/media.go:7-14` | 已定义，上传时只用到 `local` |
| `StorageSource` 模型 + `BuildBaseURL()` | `models/storage_source.go` | 已定义，`BuildBaseURL` 只被 `ResolveBaseURL` 间接调用，实际无人用 |
| `StorageSourceService.GetDefault/GetByCode/ResolveBaseURL` | `services/storage_source_service.go:38-135` | 已定义，**业务上传/读取未调用** |
| `StorageSourceService.InitDefault()` | `storage_source_service.go:137` + `main.go:100` | 启动时种了 1 条默认 `local` 源（host=localhost:8080, base_path=uploads），但上传逻辑忽略它 |
| `Media.ResolveAccessURL(baseURL)` | `models/media.go:54-65` | **全库零调用**（搜索仅命中定义处与需求文档） |
| `Media.Path / Thumbnail / WebPURL / AVIFURL / Width / Height` | `models/media.go:36-39` | 字段齐全，但上传不填 Width/Height、不生成 WebP/AVIF |

---

## 三、数据现状盘点（存量数据里会踩的坑）

1. **`source` 字段失真**：所有通过后台上传的媒体记录 `source` 都是 `local`，即便日后切到 MinIO，历史数据也全标 `local`。
2. **`url` 与 `path` 口径不统一**：
   - 本地上传：`url = /uploads/...`（相对路径）、`path = products/2026-08/uuid.jpg`（相对存储路径）。
   - `ResolveAccessURL` 的设计是「baseURL + path」拼完整 URL，但列表接口没调用它，前端拿到的是**原始相对 `url`**，靠反代兜底。
   - 若 `path` 为空（手动建记录、external 源），`url` 直接原样返回。
3. **默认存储源与真实访问方式对不上**：`InitDefault` 种的是 `http://localhost:8080/uploads`，但本地开发实际走 vite 代理、容器/隧道下走门户 `/uploads` 反代；`localhost:8080` 在浏览器侧往往不可达。`BuildBaseURL()` 直接拼出的地址在「本地 / Docker / 隧道 / 生产」四种环境表现不一致。
4. **产品封面/图集与媒体库无关联**：`product.cover_image`、`product.images[].url` 存的是**字符串 URL**（seed 数据大量是 `images.unsplash.com` 外链，见 `auth_service.go`、`setup-data.cjs`；上传后存 `/uploads/...` 或完整地址），**没有 media_id 外键**，删除/替换媒体不会联动产品。
5. **删除媒体只删 DB 记录**：`MediaService.DeleteMedia`（`media_service.go:155`）只 `db.Delete`，不清理磁盘文件/对象存储对象，**产生孤儿文件**。
6. **无图片尺寸与格式转换**：`Width/Height` 恒 0；`WebPURL/AVIFURL` 恒空（需求 §4.5 要求 WebP/AVIF 转换）。
7. **seed 大量使用第三方外链**：hero banner、产品封面等直接引用 `images.unsplash.com`，这些**既不在媒体库、也不在本地/MinIO**，属于 `external` 来源，但媒体库里没有对应记录。

---

## 四、问题清单（按优先级）

| 级别 | 问题 | 影响 | 现状位置 |
|---|---|---|---|
| P0 | 上传未接 MinIO，全走本地磁盘 | 与需求「MinIO/S3 存储」不符；本地盘难横向扩展、容器重建需卷挂载否则丢文件 | `media_handler.go:128`、`public_handler.go:704` |
| P0 | `source/path/url` 三字段未闭环，`ResolveAccessURL` 零调用 | 媒体列表返回的 URL 依赖反代兜底，换存储源/上生产会大面积死链 | `media.go:54`、`media_handler.go:143-147` |
| P1 | 删除媒体不删物理文件 | 孤儿文件堆积，磁盘/桶泄漏 | `media_service.go:155` |
| P1 | 默认存储源 host=localhost:8080，跨环境不可达 | 存储源配置形同虚设 | `storage_source_service.go:137-155` |
| P1 | 产品封面/图集与媒体库无外键 | 媒体被删后产品图死链；无法做引用计数/替换 | `models/product.go`（images/cover_image） |
| P2 | 上传不写 Width/Height、不生成 WebP/AVIF | 与需求 §4.5 不符，影响性能与 SEO | `media_handler.go:138-150` |
| P2 | seed 大量 unsplash 外链未入库 | 数据源不可控、无版权/可用性保障 | `auth_service.go`、`setup-data.cjs` |
| P2 | MinIO 配置/容器/文档齐备但 SDK 未引入 | 资源闲置 | `go.mod`、`docker-compose.yml:39` |

---

## 五、改造方案（三选一，推荐 C）

### 方案 A：接入 MinIO 真正上传/读取（引入 minio-go SDK）

**目标**：上传写 MinIO、读取走 MinIO（或 presigned URL），本地 `uploads/` 仅作降级。

**改动**：
1. `go.mod` 引入 `github.com/minio/minio-go/v7`。
2. `UploadService` 增加 MinIO 客户端初始化（懒加载，`minio.New(endpoint, &minio.Options{Creds, Secure})`），新增 `SaveToMinIO(bucket, objectName, reader, size, contentType)`、`PresignedGetURL(bucket, objectName, expiry)`。
3. 新增配置开关（如 `STORAGE_DRIVER=local|minio`，默认 `local` 以兼容现状），`UploadFile`/`UploadLeadAttachment` 按开关走 MinIO 或本地。
4. `EnsureBucket` 启动时自动建桶（幂等）。
5. 前端 `/uploads` 反代仍保留用于 local 模式；minio 模式下返回 MinIO 公网/presigned URL。

**优点**：真正满足需求；对象存储高可用、易扩展。
**缺点**：新增依赖；需要 MinIO 公网可达或额外反代；历史本地文件需迁移或双读兼容。

### 方案 B：只打通「存储源/URL 解析」闭环（不引入 SDK）

**目标**：让 `source/path/url` 与 `storage_sources` 配置一致，媒体列表/详情返回**可访问的最终 URL**。

**改动**：
1. `MediaHandler.ListMedia/GetMedia/CreateMedia` 后置处理：用 `StorageSourceService.GetDefault()`（或按 `source` 匹配 code）取 baseURL，调用 `Media.ResolveAccessURL(baseURL)` 统一输出。
2. 上传时按默认存储源写 `source`（local→`local`，minio→`minio`），URL 规则统一为「相对 path + 由前端/门户按环境解析」，或「完整 baseURL+path」。
3. 删除媒体时同步删除物理文件（local 删盘文件；minio/external 只删记录）。
4. 前端 `useImageUrl.ts` 增加 MinIO/CDN 域名归一化规则；`storage-sources.vue` 补 `is_active`、access_key 字段。

**优点**：改动小、无新依赖；先把「数据正确性」修好。
**缺点**：仍是「本地磁盘实际存储 + 配置层 MinIO」的假闭环，未真正把文件放进 MinIO。

### 方案 C：两者结合（推荐，分两步走）

- **第 1 步（数据闭环，可独立上线）**：落地方案 B，先把 `source/path/url` 口径统一、URL 解析闭环、删库删文件、默认存储源修正——解决「数据混乱」这一核心诉求，且不引入外部依赖、风险低。
- **第 2 步（真正接入 MinIO）**：落地方案 A，引入 SDK + `STORAGE_DRIVER` 开关 + 双读兼容（新文件走 MinIO、老 `/uploads/` 本地文件继续可读），完成「MinIO 落地」。

这样既先修好数据正确性，又为后续切换 MinIO 铺平道路，避免一次性大改带来的回归风险。

---

## 六、各层改动清单（以方案 C 为例）

### 6.1 后端（sportswear-backend）

| 文件 | 改动 |
|---|---|
| `internal/services/upload_service.go` | 新增 `minio` 客户端 + `SaveToMinIO`/`EnsureBucket`/`PresignedGetURL`；`PublicURL` 修正为按存储源解析；保留 `SaveToLocal`/`LocalURL` 作降级 |
| `internal/config/config.go` | 新增 `StorageDriver`（local/minio）、`MinIO` 增加 `PublicEndpoint`（对外可达域名，区分容器内 endpoint） |
| `internal/handlers/media_handler.go` | `UploadFile` 按 driver 分流；`source` 不再硬编码；填 Width/Height（可选）；URL 统一走 `ResolveAccessURL` |
| `internal/handlers/public_handler.go` | `UploadLeadAttachment` 同样按 driver 分流 |
| `internal/services/media_service.go` | `DeleteMedia` 前删除物理文件；`ListMedia` 返回前批量解析 URL |
| `internal/services/storage_source_service.go` | `InitDefault` 修正默认源（相对路径语义，避免写死 localhost:8080）；提供 `ResolveMediaURL(media)` 统一入口 |
| `internal/models/media.go` | `ResolveAccessURL` 支持 local（相对 `/uploads`）与 minio（baseURL+path）两种语义 |

### 6.2 管理后台（sportswear-admin）

| 文件 | 改动 |
|---|---|
| `src/views/media/index.vue` | 表格「来源」列展示 `source`；预览 `src` 用归一化后的 URL；上传成功后刷新列表 |
| `src/views/system/storage-sources.vue` | 表单补 `is_active`、`access_key` 字段；展示 `type` 标签 |
| `src/components/media/MediaPicker.vue` | 上传/选择统一取 `res.url`（已解析），无需改 |
| `vite.config.ts` | 增加 MinIO 域名代理（可选，用于预览 minio 模式图片） |

### 6.3 门户（sportswear-portal）

| 文件 | 改动 |
|---|---|
| `composables/useImageUrl.ts` | `BACKEND_PATTERNS` 增加 MinIO/自定义 CDN 域名归一化；增加 `STORAGE_PUBLIC_ENDPOINT` 识别 |
| `server/routes/uploads/[...path].ts` | 保留（local 模式）；可选新增 minio 反代路由 |
| `nuxt.config.ts` | `/uploads/**` 缓存规则在 minio 模式下调整 |

### 6.4 部署/运维

| 文件 | 改动 |
|---|---|
| `docker-compose.yml` | 后端 env 增加 `STORAGE_DRIVER`、`MINIO_PUBLIC_ENDPOINT`；可选自动建桶 init 容器 |
| `.env.example`（如有） | 补齐 MinIO 与 driver 变量 |
| `dev.ps1`/`deploy.ps1` | 输出提示 MinIO 是否启用 |

---

## 七、数据清洗与迁移（一次性脚本，幂等）

1. **source 对齐**：`UPDATE media SET source='local' WHERE url LIKE '/uploads/%' OR path IS NOT NULL;`（external 记录按 `url LIKE 'http%'` 且不在本域判定）。
2. **url 归一化**：本地上传记录 `url` 若为绝对地址（旧数据），统一改相对 `/uploads/<path>`；minio 记录统一 `url = <baseURL>/<path>`。
3. **孤儿文件清理**：扫描 `uploads/` 磁盘文件，与 `media.path` 比对，删除无记录对应的文件（或反向补建记录）。
4. **产品封面/图集**：可选方案——新增 `media_id` 外键字段并回填，或维持字符串 URL 但接入统一的 URL 解析函数（推荐后者，改动小、风险低）。
5. **seed 外链**：评估将 unsplash 外链替换为本地/MinIO 自有资源，或保留为 `external` 来源并入库登记。

> 建议以 SQL 脚本形式沉淀到仓库根目录（参照现有 `align-logo-alt.sql` 的幂等风格），并在 `deploy.ps1`/文档中给出执行方式。

---

## 八、风险与注意事项

1. **双读兼容**：切换 MinIO 后，存量 `/uploads/` 本地文件必须继续可读，否则历史图片全挂。方案 C 的「分两步走」正是为此。
2. **环境可达性**：MinIO 的 `endpoint`（容器内 `minio:9000`）与浏览器可达地址（`localhost:9000` 或公网域名/隧道）不是一回事，必须拆成两个配置项，否则后台预览/门户图片会 404。
3. **presigned URL 有效期**：若采用签名 URL，需处理过期与缓存（CDN 缓存 1 天会与短时效签名冲突）。
4. **删除的级联**：删媒体若改删物理文件，需先确认无产品/内容引用，避免误删仍被引用的图片（可先做引用计数或软删除）。
5. **SVG 安全**：管理后台上传允许 `.svg`（有 XSS 风险），公开询盘附件已禁用 `.svg`；切 MinIO 时建议保持「管理端可传、公开端禁传」的差异策略不变。
6. **回滚**：`STORAGE_DRIVER` 开关应支持一键切回 `local`，确保 MinIO 故障时降级不丢上传能力。

---

## 九、分阶段实施与验收

### 阶段 1（数据闭环，方案 B）——建议先行
- 验收：媒体列表/详情返回的 URL 与 `source/path` 一致；删除媒体后磁盘文件同步删除；默认存储源不再依赖 `localhost:8080`；跨 local/docker/隧道环境图片均可显示。

### 阶段 2（MinIO 落地，方案 A）
- 验收：设置 `STORAGE_DRIVER=minio` 后，上传文件进入 MinIO 桶，MinIO 控制台可见对象；媒体 URL 指向 MinIO 公网地址且可访问；本地老文件仍可读（双读兼容）；切回 `local` 上传恢复本地。

### 阶段 3（体验与性能增强，可选）
- 图片尺寸写入、WebP/AVIF 转换、缩略图生成、产品封面 media 关联与引用计数。

---

## 十、待确认事项

1. 是否需要在**当前阶段**就真正引入 MinIO SDK（方案 A/C），还是先只做数据闭环（方案 B）？
2. MinIO 是否有**公网可达地址/域名/隧道**？这决定 `MINIO_PUBLIC_ENDPOINT` 如何配置、预览与门户如何取图。
3. 存量 `uploads/` 本地文件规模？是否需要提供「本地 → MinIO 批量迁移」脚本。
4. 产品封面/图集是否要改为「媒体库引用（media_id）」模型，还是继续沿用字符串 URL + 统一解析？


---

## 十一、落地实现记录（方案 C 核心闭环）

> 已完成以下改动，编译 / vet / 测试全部通过。

### 一致性闭环规则

| 层 | 规则 |
|---|---|
| DB `path` | 统一存对象相对路径 `category/2006-01/uuid.ext`（local/minio 通用） |
| DB `source` | 由 `UploadService.MediaSource()` 按 `STORAGE_DRIVER` 决定，与真实存储一致 |
| DB `url` | 统一 `/uploads/<path>` 相对路径（`UploadService.URL()`），跨存储/跨环境一致 |
| 存储 | `STORAGE_DRIVER=local` 写本地磁盘；`=minio` 写 MinIO 桶 |
| 展示 | 浏览器 → 门户/后台反代 → 后端 `/uploads/**` → `ServeUpload` 按驱动读磁盘或 MinIO 流式返回 |
| 删除 | 删 DB 记录 + 同步删物理对象（local 删盘文件 / minio 删对象，双删兜底） |

### 改动文件清单

| 文件 | 改动 |
|---|---|
| `internal/config/config.go` | 新增 `StorageDriver`（local/minio）、`MinIO.PublicEndpoint` |
| `internal/services/upload_service.go` | 引入 minio-go；新增 `Driver/MediaSource/Save/SaveToMinIO/URL/EnsureBucket/Delete/Open/getClient/openLocal`；双读（minio 未命中回退本地）+ 双删（对象与磁盘都清理） |
| `internal/services/image_probe.go` | 新增 `ProbeImageDimensions`（JPEG/PNG/GIF/WebP 宽高解析） |
| `internal/services/thumbnail.go` | 新增 `GenerateThumbnail`（缩略图，最长边 320px，输出 JPEG）+ `ThumbPath` |
| `internal/services/media_service.go` | 新增 `CountProductReferences`（统计产品封面/图集引用，防死链） |
| `internal/handlers/media_handler.go` | `UploadFile` 走 `Save`+`URL`+`MediaSource`；图片写入 `Width/Height` 与 `Thumbnail`；`DeleteMedia` 删除前做产品引用检查 + 同步删物理对象与缩略图；新增 `ServeUpload` 读取路由 |
| `internal/handlers/public_handler.go` | 询盘附件上传走统一 `Save`+`URL` |
| `internal/router/router.go` | 删除 `r.Static`，改为 `r.GET("/uploads/*filepath", mediaHandler.ServeUpload)`；服务初始化后 `EnsureBucket` |
| `go.mod` / `go.sum` | 新增 `github.com/minio/minio-go/v7 v7.3.0`、`golang.org/x/image v0.45.0` |
| `docker-compose.yml` | 后端 env 新增 `STORAGE_DRIVER`（默认 local，可切 minio） |
| `cmd/migrate-media/main.go` | 新增本地 → MinIO 批量迁移命令（含 `-dry-run`、`-only-orphans`） |
| `align-media-data.sql`（仓库根） | 存量数据一致性修复（source/url 归一化、path 反推、校验查询，幂等） |

### 使用方式

- **MinIO（默认，图片完全从 MinIO 读取）**：`STORAGE_DRIVER` 不设或 `=minio`（docker-compose 已默认 minio）；启动自动建桶，上传写桶，`/uploads` 反代读桶。
- **本地磁盘（可选）**：`STORAGE_DRIVER=local` 回退本地磁盘。
- **数据清洗**：执行 `psql "$DATABASE_URL" -f align-media-data.sql`（幂等）。
- **存量迁移**：`go run ./cmd/migrate-media -dry-run` 预览 → `go run ./cmd/migrate-media` 实迁；`-only-orphans` 扫描孤儿文件。
- **构建说明**：Dockerfile 已改用 `vendor` 目录构建（`go mod vendor` 后提交 vendor/），避免容器内依赖下载受网络波动影响。

### 边界说明（留待后续）

- WebP/AVIF 格式转换（当前仅做缩略图生成与宽高解析，AVIF 无 Go 解码器跳过）。
- 产品封面/图集与媒体库的 `media_id` 外键关联（当前沿用字符串 URL + 统一反代解析；已用「删除前引用检查」兜底防死链）。
- MinIO 公网直连（`MINIO_PUBLIC_ENDPOINT` + presigned URL + CDN），当前展示统一走后端反代、无需公网暴露 MinIO。

---

## 十二、MinIO 模式联调验证（2026-09-06，实测通过）

在真实 Docker 环境（postgres/redis/minio/backend）完成「完全读取 MinIO」落地与验证：

| 步骤 | 结果 |
|---|---|
| 重建 backend 镜像（vendor 构建）并重启为 `STORAGE_DRIVER=minio` | ✅ 容器 env 确认 `STORAGE_DRIVER=minio`、health ok |
| 本地编译 linux 版 `migrate-media` 放入容器，在容器网络内迁移 | ✅ 成功 21、跳过 0、失败 0 |
| media 表 source 分布 | ✅ `minio | 21`（无 local 残留） |
| MinIO 桶对象确认（`mc stat`） | ✅ 对象存在（111 KiB file） |
| 通过 8080 `/uploads` 反代读回迁移图片 | ✅ HTTP 200、`image/png`、113KB |

> 过程中一并定位并处理：容器内 `goproxy.cn` 下载新依赖网络中断 → Dockerfile 改 vendor 构建绕开；本地 `.env` 的 `DB_PASSWORD` 指向本地 postgres（非 docker）属本地开发配置，不影响 docker 部署。

