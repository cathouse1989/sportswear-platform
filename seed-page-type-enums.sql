-- ============================================================
--  页面类型枚举补丁：page.type 新增 14 个类型（与后端 PageType 常量一致）
--  说明：后端 enum_service.go 已内置该枚举，重启后会自动 seed（幂等）。
--        本脚本用于【后端未重启时】手动补枚举，让后台「页面管理」类型下拉立即可用。
--  幂等：类型不存在才建；条目 ON CONFLICT (type_id, value) 跳过已存在项，可重复执行。
-- ============================================================
SET client_encoding TO 'UTF8';

INSERT INTO sys_enum_types (code, name, module, i18n_prefix, description, is_system, is_active, sort_order)
SELECT 'page.type', '页面类型', 'content', 'page.type_options', '门户页面类型（决定访问路径派生）', true, true, 190
WHERE NOT EXISTS (SELECT 1 FROM sys_enum_types WHERE code = 'page.type');

INSERT INTO sys_enum_items (type_id, value, label, translations, is_active, sort_order)
SELECT t.id, v.value, v.label, v.translations, true, v.sort_order
FROM sys_enum_types t
CROSS JOIN (VALUES
  ('normal','Normal',jsonb_build_object('zh','普通','es','Normal','fr','Normal'),1),
  ('home','Home',jsonb_build_object('zh','首页','es','Inicio','fr','Accueil'),2),
  ('product','Product',jsonb_build_object('zh','产品','es','Producto','fr','Produit'),3),
  ('product_category','Product Category',jsonb_build_object('zh','产品分类','es','Categoría de producto','fr','Catégorie de produit'),4),
  ('oem','OEM',jsonb_build_object('zh','OEM','es','OEM','fr','OEM'),5),
  ('odm','ODM',jsonb_build_object('zh','ODM','es','ODM','fr','ODM'),6),
  ('private_label','Private Label',jsonb_build_object('zh','贴牌','es','Marca propia','fr','Marque privée'),7),
  ('factory','Factory',jsonb_build_object('zh','工厂','es','Fábrica','fr','Usine'),8),
  ('production','Production',jsonb_build_object('zh','生产流程','es','Producción','fr','Production'),9),
  ('blog','Blog',jsonb_build_object('zh','博客','es','Blog','fr','Blog'),10),
  ('case','Case',jsonb_build_object('zh','案例','es','Caso','fr','Étude de cas'),11),
  ('faq','FAQ',jsonb_build_object('zh','FAQ','es','FAQ','fr','FAQ'),12),
  ('contact','Contact',jsonb_build_object('zh','联系我们','es','Contacto','fr','Contact'),13),
  ('seo_landing','SEO Landing',jsonb_build_object('zh','SEO 落地页','es','Página de destino SEO','fr','Page de destination SEO'),14)
) AS v(value, label, translations, sort_order)
WHERE t.code = 'page.type'
ON CONFLICT (type_id, value) DO NOTHING;