-- ============================================================
--  FAQ category 脏值修正（幂等，可重复执行）
--  背景：早期 seed-faq-i18n.sql 使用了错误分类名（Orders/Production/
--        Customization/Quality/Services），与标准值域不一致。
--  标准值域（后端字典 faq.category + 门户 locale）：
--        moq/oem/odm/sample/payment/production/logistics/fabric/quality/certification
--  映射：
--        Orders        -> moq           （最低起订量）
--        Production    -> production    （生产周期，仅大小写归一）
--        Customization -> fabric        （定制面料）
--        Quality       -> certification （工厂认证）
--        Services      -> oem           （自有品牌/贴牌）
-- ============================================================
SET client_encoding TO 'UTF8';

UPDATE faqs SET category = 'moq'          WHERE category = 'Orders'        AND deleted_at IS NULL;
UPDATE faqs SET category = 'production'    WHERE category = 'Production'    AND deleted_at IS NULL;
UPDATE faqs SET category = 'fabric'        WHERE category = 'Customization' AND deleted_at IS NULL;
UPDATE faqs SET category = 'certification' WHERE category = 'Quality'       AND deleted_at IS NULL;
UPDATE faqs SET category = 'oem'           WHERE category = 'Services'      AND deleted_at IS NULL;
