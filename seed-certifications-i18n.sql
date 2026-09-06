-- ============================================================
--  认证管理多语言补配：为 4 个默认认证补齐 code / 颁发日期 / 到期日期，
--  并写入 zh / es / fr 三语言翻译（status=published），门户据此本地化展示。
--  幂等：认证按 name 匹配，翻译按 (certification_id, language) 判断，可重复执行。
--  执行方式：
--    docker exec -i sportswear-postgres psql -U postgres -d sportswear_platform < seed-certifications-i18n.sql
--  或：
--    psql "$DATABASE_URL" -f seed-certifications-i18n.sql
-- ============================================================
SET client_encoding TO 'UTF8';

-- 1) 补齐认证主表元数据（编号 / 颁发日期 / 到期日期）
UPDATE certifications
SET code = v.code, issue_date = v.issue_date, expiry_date = v.expiry_date
FROM (VALUES
  ('ISO 9001:2015',         'ISO-9001',     '2023-05-10', '2026-05-09'),
  ('BSCI Compliance',       'BSCI',         '2024-01-15', '2026-01-14'),
  ('OEKO-TEX Standard 100', 'OEKO-TEX-100', '2024-03-01', '2027-02-28'),
  ('SEDEX Registered',      'SEDEX',        '2024-06-20', '2026-06-19')
) AS v(name, code, issue_date, expiry_date)
WHERE certifications.name = v.name AND certifications.deleted_at IS NULL;

-- 2) 写入三语言翻译（已存在则跳过）
INSERT INTO certification_translations (certification_id, language, name, description, status, created_at, updated_at)
SELECT c.id, t.language, t.tr_name, t.tr_desc, 'published', now(), now()
FROM certifications c
JOIN (VALUES
  ('ISO 9001:2015',         'zh', 'ISO 9001:2015 质量管理体系认证', '质量管理体系认证'),
  ('ISO 9001:2015',         'es', 'ISO 9001:2015', 'Sistema de Gestión de Calidad certificado'),
  ('ISO 9001:2015',         'fr', 'ISO 9001:2015', 'Système de management de la qualité certifié'),
  ('BSCI Compliance',       'zh', 'BSCI 社会责任认证', 'Amfori BSCI 认证 — 道德制造'),
  ('BSCI Compliance',       'es', 'Cumplimiento BSCI', 'Amfori BSCI certificado — fabricación ética'),
  ('BSCI Compliance',       'fr', 'Conformité BSCI', 'Amfori BSCI certifié — fabrication éthique'),
  ('OEKO-TEX Standard 100', 'zh', 'OEKO-TEX Standard 100 产品安全认证', '产品安全认证'),
  ('OEKO-TEX Standard 100', 'es', 'OEKO-TEX Standard 100', 'Certificación de seguridad del producto'),
  ('OEKO-TEX Standard 100', 'fr', 'OEKO-TEX Standard 100', 'Certification de sécurité des produits'),
  ('SEDEX Registered',      'zh', 'SEDEX 注册供应商', '供应商道德数据交换（SEDEX）'),
  ('SEDEX Registered',      'es', 'SEDEX Registrado', 'Intercambio de datos éticos de proveedores'),
  ('SEDEX Registered',      'fr', 'SEDEX Enregistré', 'Échange de données éthiques des fournisseurs')
) AS t(cert_name, language, tr_name, tr_desc)
  ON c.name = t.cert_name AND c.deleted_at IS NULL
WHERE NOT EXISTS (
  SELECT 1 FROM certification_translations e
  WHERE e.certification_id = c.id AND e.language = t.language AND e.deleted_at IS NULL
);
