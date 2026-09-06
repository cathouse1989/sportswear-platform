-- ============================================================
--  门户预置路由的初始 route SEO（英文源语言；其他语言可在后台「SEO 管理」补充，未填则回退英文）
--  说明：与后台「SEO 管理」编辑的结果完全一致（seo 表 entity_type='route'），
--        运营后续可直接在后台覆盖。文案取自门户各列表页的代码级默认值。
--  幂等两步，可重复执行：
--    1) 清理「英文源语言中 title/description 均为空」的占位记录（视为未配置）；
--    2) 插入尚不存在的记录（已配置的非空记录不会被覆盖）。
--  entity_id 为确定性 UUID v5：sha1(namespaceURL, "route:"+route)，
--  与后端 services.RouteSEOID 完全一致，请勿手工修改。
-- ============================================================
SET client_encoding TO 'UTF8';

-- 1. 清理空占位记录（视为未配置）
DELETE FROM seos
WHERE entity_type = 'route'
  AND language = 'en'
  AND COALESCE(title, '') = ''
  AND COALESCE(description, '') = '';

-- 2. 插入初始 route SEO
INSERT INTO seos (entity_type, entity_id, language, title, description, keywords)
SELECT v.entity_type, v.entity_id, 'en', v.title, v.description, v.keywords
FROM (VALUES
  ('route', '280451b2-e2ac-5804-b6b6-d65842705595'::uuid,
   'OEM/ODM Sportswear Manufacturer - Custom Activewear Factory',
   'Professional OEM/ODM sportswear manufacturer. Custom sportswear, activewear, and athletic apparel for global brands. 15+ years of experience, ISO certified.',
   'sportswear manufacturer, OEM, ODM, custom sportswear, activewear, athletic apparel, China factory, private label'),

  ('route', 'ab4f05c1-7fc2-5e40-a902-2ecf36d0912e'::uuid,
   'Custom Sportswear Products - OEM/ODM Activewear Collection',
   'Explore custom sportswear manufacturing: yoga wear, running gear, training apparel, team uniforms and more. OEM/ODM production with low MOQ and fast sampling.',
   'custom sportswear, OEM activewear, ODM sportswear, yoga wear manufacturer, running gear factory, team uniforms'),

  ('route', 'a00fe6fd-df6a-5bbb-8e2b-934df1f450f7'::uuid,
   'Case Studies - OEM/ODM Sportswear Manufacturing',
   'Explore our successful OEM/ODM sportswear projects. See how we help global brands bring their activewear visions to life with quality manufacturing and innovative solutions.',
   'sportswear case studies, OEM projects, ODM manufacturing cases, activewear brand partnerships, garment factory portfolio'),

  ('route', '2ec55bc4-185a-565d-bb51-74d19286d6da'::uuid,
   'About Us - OEM/ODM Sportswear Manufacturer Since 2008',
   'Professional OEM/ODM sportswear manufacturer with 15+ years of experience. ISO 9001, BSCI, OEKO-TEX certified. 50,000 sqm facility serving 500+ global brands.',
   'sportswear manufacturer history, OEM factory China, ODM sportswear company, ISO certified activewear factory, garment manufacturing China'),

  ('route', '55e9cbbe-1ec4-5373-baed-04ce88a2ce38'::uuid,
   'Sportswear Manufacturing Blog - Industry Insights & Tips',
   'Expert insights on sportswear manufacturing, OEM/ODM processes, fabric technology, and industry trends. Learn how to choose the right manufacturer for your activewear brand.',
   'sportswear blog, OEM manufacturing insights, activewear industry, garment factory, sportswear manufacturing tips'),

  ('route', 'c0e07dec-1eab-5bda-b838-8ef0be3072b3'::uuid,
   'FAQ - Sportswear OEM/ODM Manufacturing',
   'Find answers to frequently asked questions about OEM/ODM sportswear manufacturing, MOQ, customization processes, pricing, shipping, and quality control.',
   'sportswear FAQ, OEM questions, ODM manufacturing, MOQ sportswear, custom activewear process, garment manufacturing FAQ'),

  ('route', '35e9aa6c-845e-519a-8ee8-e94aa02d0055'::uuid,
   'Contact Us - OEM/ODM Sportswear Manufacturer',
   'Get in touch with our OEM/ODM sportswear manufacturing team. Request a quote, ask about customization options, or discuss your activewear project. We respond within 24 hours.',
   'contact sportswear manufacturer, OEM inquiry, ODM quote, activewear supplier, China garment factory contact'),

  ('route', '96680249-2ec4-5a35-8120-4ef968d24795'::uuid,
   'Privacy Policy - OEM/ODM Sportswear Manufacturer',
   'Privacy policy for our sportswear OEM/ODM manufacturing website. Learn how we collect, use and protect your personal information.',
   'privacy policy, sportswear manufacturer, data protection, OEM factory website')
) AS v(entity_type, entity_id, title, description, keywords)
WHERE NOT EXISTS (
  SELECT 1 FROM seos
  WHERE entity_type = v.entity_type AND entity_id = v.entity_id AND language = 'en'
);
