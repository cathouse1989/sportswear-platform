-- 补配门户页脚词条（4 语言），已存在则跳过（key+language 幂等）
SET client_encoding TO 'UTF8';
DELETE FROM i18n_entries WHERE key = 'footer.about_desc' OR key LIKE 'footer.cat_%';
INSERT INTO i18n_entries (key, language, value, module, is_active, sort_order)
SELECT t.k, t.l, t.v, 'footer', true, 0
FROM (VALUES
  ('footer.about_desc','en','Professional sportswear OEM/ODM manufacturer with 15+ years of experience. ISO 9001, BSCI, OEKO-TEX certified.'),
  ('footer.about_desc','zh',E'专业运动服装 OEM/ODM 制造商，\n深耕行业 15 年以上。\n通过 ISO 9001、BSCI、OEKO-TEX 认证。'),
  ('footer.about_desc','es','Fabricante profesional de ropa deportiva OEM/ODM con más de 15 años de experiencia. Certificado ISO 9001, BSCI y OEKO-TEX.'),
  ('footer.about_desc','fr',E'Fabricant professionnel de vêtements de sport OEM/ODM avec plus de 15 ans d''expérience. Certifié ISO 9001, BSCI et OEKO-TEX.'),
  ('footer.cat_yoga_wear','en','Yoga Wear'),
  ('footer.cat_yoga_wear','zh','瑜伽服'),
  ('footer.cat_yoga_wear','es','Ropa de Yoga'),
  ('footer.cat_yoga_wear','fr','Vêtements de Yoga'),
  ('footer.cat_running_gear','en','Running Gear'),
  ('footer.cat_running_gear','zh','跑步装备'),
  ('footer.cat_running_gear','es','Equipo de Correr'),
  ('footer.cat_running_gear','fr','Équipement de Course'),
  ('footer.cat_training_apparel','en','Training Apparel'),
  ('footer.cat_training_apparel','zh','训练服饰'),
  ('footer.cat_training_apparel','es','Ropa de Entrenamiento'),
  ('footer.cat_training_apparel','fr','Vêtements d''Entraînement'),
  ('footer.cat_team_uniforms','en','Team Uniforms'),
  ('footer.cat_team_uniforms','zh','团队队服'),
  ('footer.cat_team_uniforms','es','Uniformes de Equipo'),
  ('footer.cat_team_uniforms','fr','Uniformes d''Équipe'),
  ('footer.cat_custom_design','en','Custom Design'),
  ('footer.cat_custom_design','zh','定制设计'),
  ('footer.cat_custom_design','es','Diseño Personalizado'),
  ('footer.cat_custom_design','fr','Conception Personnalisée')
) AS t(k,l,v)
WHERE NOT EXISTS (
  SELECT 1 FROM i18n_entries e WHERE e.key = t.k AND e.language = t.l AND e.deleted_at IS NULL
);
