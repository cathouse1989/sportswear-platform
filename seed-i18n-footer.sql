-- 补配门户页脚词条 + 产品分类名多语言词条（4 语言），已存在则跳过（key+language 幂等）
-- 说明：产品分类名统一走 product.categories.<slug> 词条，与后台「分类管理」categories.slug 一一对应，
--       门户 footer / 产品列表 / 分类页共用同一词条；历史 footer.cat_* 词条已废弃并清理。
SET client_encoding TO 'UTF8';
DELETE FROM i18n_entries WHERE key = 'footer.about_desc' OR key LIKE 'footer.cat_%' OR key LIKE 'product.categories.%';
INSERT INTO i18n_entries (key, language, value, module, is_active, sort_order)
SELECT t.k, t.l, t.v, t.m, true, 0
FROM (VALUES
  ('footer.about_desc','en','Professional sportswear OEM/ODM manufacturer with 15+ years of experience. ISO 9001, BSCI, OEKO-TEX certified.','footer'),
  ('footer.about_desc','zh',E'专业运动服装 OEM/ODM 制造商，\n深耕行业 15 年以上。\n通过 ISO 9001、BSCI、OEKO-TEX 认证。','footer'),
  ('footer.about_desc','es','Fabricante profesional de ropa deportiva OEM/ODM con más de 15 años de experiencia. Certificado ISO 9001, BSCI y OEKO-TEX.','footer'),
  ('footer.about_desc','fr',E'Fabricant professionnel de vêtements de sport OEM/ODM avec plus de 15 ans d''expérience. Certifié ISO 9001, BSCI et OEKO-TEX.','footer'),
  ('product.categories.yoga-wear','en','Yoga Wear','product'),
  ('product.categories.yoga-wear','zh','瑜伽服','product'),
  ('product.categories.yoga-wear','es','Ropa de Yoga','product'),
  ('product.categories.yoga-wear','fr','Vêtements de Yoga','product'),
  ('product.categories.running-gear','en','Running Gear','product'),
  ('product.categories.running-gear','zh','跑步装备','product'),
  ('product.categories.running-gear','es','Equipo de Correr','product'),
  ('product.categories.running-gear','fr','Équipement de Course','product'),
  ('product.categories.training-apparel','en','Training Apparel','product'),
  ('product.categories.training-apparel','zh','训练服饰','product'),
  ('product.categories.training-apparel','es','Ropa de Entrenamiento','product'),
  ('product.categories.training-apparel','fr','Vêtements d''Entraînement','product'),
  ('product.categories.team-uniforms','en','Team Uniforms','product'),
  ('product.categories.team-uniforms','zh','团队队服','product'),
  ('product.categories.team-uniforms','es','Uniformes de Equipo','product'),
  ('product.categories.team-uniforms','fr','Uniformes d''Équipe','product'),
  ('product.categories.custom-design','en','Custom Design','product'),
  ('product.categories.custom-design','zh','定制设计','product'),
  ('product.categories.custom-design','es','Diseño Personalizado','product'),
  ('product.categories.custom-design','fr','Conception Personnalisée','product')
) AS t(k,l,v,m)
WHERE NOT EXISTS (
  SELECT 1 FROM i18n_entries e WHERE e.key = t.k AND e.language = t.l AND e.deleted_at IS NULL
);
