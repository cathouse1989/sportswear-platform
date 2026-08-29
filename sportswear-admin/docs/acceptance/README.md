# 验收用例库（Acceptance Cases）

> 版本：v0（骨架）　建立日期：2026-08-27
> 依据：《管理后台产品说明书》第四章模块明细 + 第七章质检基线生成；打磨方案第 0 期 W3 产出物。

## 一、用途与执行节奏

1. **发版前**：当期改动涉及的页面用例全量执行，全绿才允许合并 release；
2. **每周五走查**：随机抽 3 个页面对照用例运行（对应打磨方案 §四-4），漂移即修；
3. **回归来源**：说明书第八章缺陷表中状态为 ✅ 的条目，对应页面的用例永久保留回归标记。

## 二、用例文件规范

- 一个菜单页一个文件，命名 = 路由末段（如 `blogs.md`）；
- 用例结构：`前置条件（角色+数据）→ 步骤 → 预期 →（可选）关联缺陷号`；
- 新增页面/功能时先写用例再开发（DoD ②）；修改行为时先更新用例再更新说明书。

## 三、目录索引

| 文件 | 覆盖页面 | 状态 |
|------|---------|------|
| blogs.md | /blogs 博客管理 | ✅ 含 B-01 回归 |
| cases.md | /cases 案例管理 | ✅ 含 B-01 回归 |
| faqs.md | /faqs FAQ 管理 | ✅ 含 B-01 回归 |
| pages.md | /pages 页面管理 | ✅ 含 B-13 回归 |
| trash.md | /trash 回收站 | ✅ 含 B-02/B-11 回归 |
| products.md | /products 产品管理（北极星链路） | ✅ |
| leads.md | /leads 询盘管理 | ✅ |
| media.md | /media 媒体管理 | ✅ 含 B-08 现状标记 |
| dashboard.md | /dashboard 仪表盘 | ✅ |
| analytics.md | /analytics 数据分析 | ✅ |
| users.md | /users 用户管理 | ✅ 含 B-15 回归 |
| roles.md | /roles 角色权限 | ✅ 含 B-15 回归 |
| theme.md | /theme 主题配置 | ✅ |
| i18n.md | /i18n 词条管理 | ✅ |
| storage-sources.md | /storage-sources 存储源配置 | ✅ |
| quotes.md | /quotes 报价管理 | ✅ |
| notifications.md | /notifications 通知列表 | ✅ |
| operation-logs.md | /operation-logs 操作日志 | ✅ |
| factories.md | /factories 工厂管理 | ✅ |
| certifications.md | /certifications 认证管理 | ✅ 含 B-12 回归 |
| production-processes.md | /production-processes 生产流程 | ✅ |

> portal-preview（/portal-preview）用例并入 products/pages 的"前台可见性"断言，独立用例文件待 P1-#7 预览体系调整后补充。

## 四、全局横切用例（每页均适用，不重复写入各文件）

| # | 步骤 | 预期 |
|---|------|------|
| G1 | 以 viewer 角色访问各页面 | 页面可进入但写操作按钮不可见/禁用（按钮级权限实现后生效，见 P0-#3） |
| G2 | 断网/后端停机时点击查询 | 统一错误提示，页面不白屏、不加载数据假象 |
| G3 | 列表接口返回空 | 展示空态而非报错 |
| G4 | 删除类操作 | 一律有二次确认，确认文案含对象名称 |
