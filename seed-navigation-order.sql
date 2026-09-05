-- 门户导航排序统一：header / footer 顺序一致，并补齐 footer 导航（幂等）
-- 顺序：首页 / 产品中心 / 案例展示 / 关于我们 / 博客 / 常见问题 / 联系我们
SET client_encoding TO 'UTF8';

-- 1. 隐藏指向不存在页面的历史遗留 header 顶级导航（/oem /odm /factory 无对应路由页面）
UPDATE navigations SET is_visible = false
WHERE type = 'header' AND url IN ('/oem', '/odm', '/factory') AND deleted_at IS NULL;

-- 2. 统一 header 顶级导航排序（按 URL 精确匹配，幂等）
UPDATE navigations SET sort_order = 1 WHERE type = 'header' AND url = '/' AND deleted_at IS NULL;
UPDATE navigations SET sort_order = 2 WHERE type = 'header' AND url = '/products' AND deleted_at IS NULL;
UPDATE navigations SET sort_order = 3 WHERE type = 'header' AND url = '/cases' AND deleted_at IS NULL;
UPDATE navigations SET sort_order = 4 WHERE type = 'header' AND url = '/about' AND deleted_at IS NULL;
UPDATE navigations SET sort_order = 5 WHERE type = 'header' AND url = '/blog' AND deleted_at IS NULL;
UPDATE navigations SET sort_order = 6 WHERE type = 'header' AND url = '/faq' AND deleted_at IS NULL;
UPDATE navigations SET sort_order = 7 WHERE type = 'header' AND url = '/contact' AND deleted_at IS NULL;

-- 3. 补齐 footer 顶级导航（与 header 一致，已存在则跳过）
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
