-- ============================================================
--  门户「页面 / 导航 / SEO」闭环脏数据清理脚本（幂等，可重复执行）
--
--  背景：历史版本允许后台自由增删页面/导航，产生了以下脏数据：
--    1) header 导航重复（固定骨架 + 「一键同步」产生的 page_id 关联项并存）
--    2) footer 导航缺失 / header 排序与 seed 定义不一致
--    3) contact 页面标题拼写错误（Contacst Us -> Contact Us）
--    4) seos 空占位记录（language 为空 / title、description 均空）
--    5) 历史已发布产品/博客缺英文源语言 SEO（自动生成仅对新建生效）
--
--  执行方式：
--    docker exec -i sportswear-postgres psql -U postgres -d sportswear_platform < cleanup-dirty-data.sql
--  或：
--    psql "$DATABASE_URL" -f cleanup-dirty-data.sql
-- ============================================================
SET client_encoding TO 'UTF8';

-- 1. 删除 header 重复导航：保留无 page_id 的固定骨架，删除「一键同步」产生的关联项
DELETE FROM navigations
WHERE type = 'header'
  AND page_id IS NOT NULL
  AND deleted_at IS NULL;

-- 2. 修正 contact 页面标题拼写错误
UPDATE pages
SET title = 'Contact Us'
WHERE slug = 'contact' AND title = 'Contacst Us' AND deleted_at IS NULL;

-- 3. 统一 header 顶级导航排序（与 footer 一致：首页/产品/案例/关于/博客/FAQ/联系）
UPDATE navigations SET sort_order = 1 WHERE type = 'header' AND url = '/' AND deleted_at IS NULL;
UPDATE navigations SET sort_order = 2 WHERE type = 'header' AND url = '/products' AND deleted_at IS NULL;
UPDATE navigations SET sort_order = 3 WHERE type = 'header' AND url = '/cases' AND deleted_at IS NULL;
UPDATE navigations SET sort_order = 4 WHERE type = 'header' AND url = '/about' AND deleted_at IS NULL;
UPDATE navigations SET sort_order = 5 WHERE type = 'header' AND url = '/blog' AND deleted_at IS NULL;
UPDATE navigations SET sort_order = 6 WHERE type = 'header' AND url = '/faq' AND deleted_at IS NULL;
UPDATE navigations SET sort_order = 7 WHERE type = 'header' AND url = '/contact' AND deleted_at IS NULL;

-- 4. 补齐 footer 固定骨架导航（幂等）
INSERT INTO navigations (name, type, url, sort_order, is_visible)
SELECT t.n, 'footer', t.u, t.s, true
FROM (VALUES
  ('Home', '/', 1),
  ('Products', '/products', 2),
  ('Cases', '/cases', 3),
  ('About Us', '/about', 4),
  ('Blog', '/blog', 5),
  ('FAQ', '/faq', 6),
  ('Contact Us', '/contact', 7)
) AS t(n, u, s)
WHERE NOT EXISTS (
  SELECT 1 FROM navigations WHERE type = 'footer' AND url = t.u AND deleted_at IS NULL
);

-- 5. 清理 seos 空占位记录：language 为空，或 title/description 均为空（视为未配置，回退英文）
DELETE FROM seos
WHERE deleted_at IS NULL
  AND (
    COALESCE(language, '') = ''
    OR (COALESCE(title, '') = '' AND COALESCE(description, '') = '')
  );

-- 6. 为缺英文源语言 SEO 的已发布产品批量补默认值（title=英文名，description=brief）
INSERT INTO seos (entity_type, entity_id, language, title, description)
SELECT 'product', p.id, 'en',
       COALESCE(
         (SELECT t.name FROM product_translations t
          WHERE t.product_id = p.id AND t.language = 'en' AND t.deleted_at IS NULL
          LIMIT 1),
         p.sku
       ),
       COALESCE(p.brief, '')
FROM products p
WHERE p.deleted_at IS NULL
  AND p.status = 'published'
  AND NOT EXISTS (
    SELECT 1 FROM seos s
    WHERE s.entity_type = 'product' AND s.entity_id = p.id AND s.language = 'en' AND s.deleted_at IS NULL
  );

-- 7. 为缺英文源语言 SEO 的已发布博客批量补默认值（title=标题，description=正文摘要）
INSERT INTO seos (entity_type, entity_id, language, title, description)
SELECT 'blog', b.id, 'en', b.title, COALESCE(left(b.content, 200), '')
FROM blogs b
WHERE b.deleted_at IS NULL
  AND b.status = 'published'
  AND NOT EXISTS (
    SELECT 1 FROM seos s
    WHERE s.entity_type = 'blog' AND s.entity_id = b.id AND s.language = 'en' AND s.deleted_at IS NULL
  );
