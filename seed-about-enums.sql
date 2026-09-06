-- ============================================================
--  「关于我们」模块枚举补丁：page.module.type 新增 6 个类型
--  说明：后端 enum_service.go 已内置这些枚举，重启后会自动 seed（幂等）。
--        本脚本用于【后端未重启时】手动补枚举，让后台「模块编辑」下拉立即可用。
--  幂等：ON CONFLICT (type_id, value) 跳过已存在项，可重复执行。
-- ============================================================
SET client_encoding TO 'UTF8';

INSERT INTO sys_enum_items (type_id, value, label, translations, is_active, sort_order)
SELECT t.id, v.value, v.label, v.translations::jsonb, true, v.sort_order
FROM sys_enum_types t
CROSS JOIN (VALUES
  ('about_hero', 'About Hero', $json${"zh":"关于我们横幅","es":"Hero de nosotros","fr":"Bannière à propos"}$json$, 14),
  ('about_story', 'About Story', $json${"zh":"公司故事（图文）","es":"Nuestra historia","fr":"Notre histoire"}$json$, 15),
  ('about_stats', 'Stats', $json${"zh":"数据统计","es":"Estadísticas","fr":"Statistiques"}$json$, 16),
  ('about_certifications', 'Certifications', $json${"zh":"认证墙","es":"Certificaciones","fr":"Certifications"}$json$, 17),
  ('about_process', 'Process Steps', $json${"zh":"生产流程步骤","es":"Pasos del proceso","fr":"Étapes de production"}$json$, 18),
  ('about_cta', 'Call To Action', $json${"zh":"行动号召","es":"Llamada a la acción","fr":"Appel à l'action"}$json$, 19)
) AS v(value, label, translations, sort_order)
WHERE t.code = 'page.module.type'
ON CONFLICT (type_id, value) DO NOTHING;
