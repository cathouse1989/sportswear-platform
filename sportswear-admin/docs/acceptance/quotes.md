# 验收用例 · 报价管理（/quotes）

> 前置角色：sales_manager / sales（lead:view，权限码解耦见 B-07）；关联：说明书 §4.3

| # | 类型 | 步骤 | 预期 |
|---|------|------|------|
| 1 | 正向 | 新建报价（Lead ID/数量/单价/总价/币种/MOQ/交期） | 列表出现，金额列显示"USD 1234.00"格式 |
| 2 | 反向 | Lead ID 留空保存 | 表单 required 提示（现状为 el-form required 属性，无 JS 校验拦截，记录为待办：应改 rules 校验） |
| 3 | 正向 | 按状态筛选（draft/sent/viewed/accepted/rejected/expired） | 列表正确过滤 |
| 4 | 正向 | 编辑报价 | 字段完整回填；修改总价保存生效 |
| 5 | 正向 | 删除报价并确认 | 二次确认后移除 |
| 6 | 标记 | PDF 导出 / 审批流 | 无此能力（P2-#15 范围） |
