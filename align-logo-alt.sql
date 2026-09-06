-- ============================================================
--  logo_alt 存量对齐脚本（幂等，可重复执行）
--
--  背景：历史部署脚本 sportswear-portal/setup-data.cjs 曾将 logo_alt 写成
--        'Sportswear Manufacturer'，与后端默认 seed 不一致：
--          - 后端默认（models/portal.go DefaultThemeConfigs）：logo_alt = "Sportswear"
--          - 门户组件 SportswearLogo.vue 兜底：'Sportswear'
--
--  本脚本仅修正该历史错误值，不覆盖运营在后台「主题配置」自定义的其它 alt 文案。
--
--  执行方式：
--    docker exec -i sportswear-postgres psql -U postgres -d sportswear_platform < align-logo-alt.sql
--  或：
--    psql "$DATABASE_URL" -f align-logo-alt.sql
-- ============================================================
SET client_encoding TO 'UTF8';

UPDATE theme_configs
SET value = '"Sportswear"'
WHERE "key" = 'logo_alt'
  AND value IN ('Sportswear Manufacturer', '"Sportswear Manufacturer"');
