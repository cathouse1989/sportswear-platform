-- ============================================================
--  FAQ 多语言化数据迁移（幂等，可重复执行）
--  背景：历史 FAQ 每个语言一条主表记录（4 语言 = 4 条），本次改造为
--        「主表存英文源 + faq_translations 存其他语言」，一条 FAQ 一条主表记录。
--  迁移：按 category + sort_order 分组，英文记录保留为主表，
--        其他语言记录转入 faq_translations。
--
--  执行方式：
--    docker exec -i sportswear-postgres psql -U postgres -d sportswear_platform < migrate-faq-i18n.sql
-- ============================================================
SET client_encoding TO 'UTF8';

-- 1. 将非英文记录转入翻译表（关联到同组英文记录）
INSERT INTO faq_translations (faq_id, language, question, answer, status)
SELECT en.id, f.language, f.question, f.answer, 'published'
FROM faqs f
JOIN faqs en ON en.category = f.category AND en.sort_order = f.sort_order
  AND en.language = 'en' AND en.deleted_at IS NULL
WHERE f.language <> 'en' AND f.deleted_at IS NULL
  AND NOT EXISTS (
    SELECT 1 FROM faq_translations t
    WHERE t.faq_id = en.id AND t.language = f.language AND t.deleted_at IS NULL
  );

-- 2. 删除非英文主表记录（保留英文源）
DELETE FROM faqs WHERE language <> 'en' AND deleted_at IS NULL;

-- 3. 主表 language 统一为 en（幂等）
UPDATE faqs SET language = 'en' WHERE deleted_at IS NULL;
