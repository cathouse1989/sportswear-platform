-- ============================================================
--  门户「询盘问题模板」默认数据（英文源 + zh/es/fr 翻译，幂等）
--  对应后台「内容管理 → 询盘模板」（/admin/inquiry-templates）与
--  门户「联系我们」询盘引导区块（/public/inquiry-templates）。
--  说明：运营可在后台增删/排序/上下架/维护多语言；本脚本仅提供初始样例。
--  执行方式：psql -U <user> -d <db> -f seed-inquiry-templates.sql
-- ============================================================
SET client_encoding TO 'UTF8';

-- 1) 主表（英文源问题，按 question 幂等）
INSERT INTO inquiry_templates (id, question, category, project_type, sort_order, is_active, created_at, updated_at)
SELECT gen_random_uuid(), t.q, t.cat, t.pt, t.sort, true, now(), now()
FROM (VALUES
  ($$What is your minimum order quantity (MOQ)? Can I start with a small order?$$, 'moq', '', 1),
  ($$Can you produce custom designs with our logo, colors and labels (OEM)?$$, 'oem', 'oem', 2),
  ($$Can you help us develop new designs and styles (ODM)?$$, 'odm', 'odm', 3),
  ($$Can I order samples first? What is the sample cost and lead time?$$, 'sample', '', 4),
  ($$Can you send me a quotation with unit price, MOQ and lead time?$$, 'payment', '', 5),
  ($$Can I choose specific fabrics or request custom fabric options?$$, 'fabric', '', 6),
  ($$What are your payment terms and which currencies do you accept?$$, 'payment', '', 7),
  ($$Do you ship to my country? What shipping options are available?$$, 'logistics', '', 8),
  ($$What is the production lead time for a bulk order?$$, 'production', '', 9),
  ($$What quality standards and certifications do you have?$$, 'quality', '', 10)
) AS t(q, cat, pt, sort)
WHERE NOT EXISTS (
  SELECT 1 FROM inquiry_templates x WHERE x.question = t.q AND x.deleted_at IS NULL
);

-- 2) 翻译表（按 源问题 + 语言 幂等）
INSERT INTO inquiry_template_translations (id, template_id, language, question, status, created_at, updated_at)
SELECT gen_random_uuid(), m.id, t.l, t.q, 'published', now(), now()
FROM (VALUES
  ($$What is your minimum order quantity (MOQ)? Can I start with a small order?$$, 'zh', $$请问最低起订量（MOQ）是多少？可以从小批量开始吗？$$),
  ($$Can you produce custom designs with our logo, colors and labels (OEM)?$$, 'zh', $$可以按我们的设计生产，加我们的 Logo、颜色和标签吗（OEM 贴牌）？$$),
  ($$Can you help us develop new designs and styles (ODM)?$$, 'zh', $$可以帮我们开发新的设计款式吗（ODM）？$$),
  ($$Can I order samples first? What is the sample cost and lead time?$$, 'zh', $$可以先订购样品吗？样品费和打样周期是多少？$$),
  ($$Can you send me a quotation with unit price, MOQ and lead time?$$, 'zh', $$可以给我发一份报价吗（单价、起订量、交期）？$$),
  ($$Can I choose specific fabrics or request custom fabric options?$$, 'zh', $$可以选择或定制特定面料吗？$$),
  ($$What are your payment terms and which currencies do you accept?$$, 'zh', $$你们的付款方式和接受的币种是什么？$$),
  ($$Do you ship to my country? What shipping options are available?$$, 'zh', $$可以发货到我们国家吗？有哪些物流方式？$$),
  ($$What is the production lead time for a bulk order?$$, 'zh', $$大货的生产周期是多久？$$),
  ($$What quality standards and certifications do you have?$$, 'zh', $$你们有哪些质量标准和认证？$$),

  ($$What is your minimum order quantity (MOQ)? Can I start with a small order?$$, 'es', $$¿Cuál es su cantidad mínima de pedido (MOQ)? ¿Puedo empezar con un pedido pequeño?$$),
  ($$Can you produce custom designs with our logo, colors and labels (OEM)?$$, 'es', $$¿Pueden fabricar diseños personalizados con nuestro logo, colores y etiquetas (OEM)?$$),
  ($$Can you help us develop new designs and styles (ODM)?$$, 'es', $$¿Pueden ayudarnos a desarrollar nuevos diseños y estilos (ODM)?$$),
  ($$Can I order samples first? What is the sample cost and lead time?$$, 'es', $$¿Puedo pedir muestras primero? ¿Cuál es el costo y plazo de la muestra?$$),
  ($$Can you send me a quotation with unit price, MOQ and lead time?$$, 'es', $$¿Pueden enviarme una cotización con precio unitario, MOQ y plazo de entrega?$$),
  ($$Can I choose specific fabrics or request custom fabric options?$$, 'es', $$¿Puedo elegir telas específicas o solicitar opciones de tela personalizadas?$$),
  ($$What are your payment terms and which currencies do you accept?$$, 'es', $$¿Cuáles son sus condiciones de pago y qué monedas aceptan?$$),
  ($$Do you ship to my country? What shipping options are available?$$, 'es', $$¿Envían a mi país? ¿Qué opciones de envío hay disponibles?$$),
  ($$What is the production lead time for a bulk order?$$, 'es', $$¿Cuál es el plazo de producción para un pedido al por mayor?$$),
  ($$What quality standards and certifications do you have?$$, 'es', $$¿Qué estándares de calidad y certificaciones tienen?$$),

  ($$What is your minimum order quantity (MOQ)? Can I start with a small order?$$, 'fr', $$Quelle est votre quantité minimum de commande (MOQ) ? Puis-je commencer par une petite commande ?$$),
  ($$Can you produce custom designs with our logo, colors and labels (OEM)?$$, 'fr', $$Pouvez-vous fabriquer des designs personnalisés avec notre logo, couleurs et étiquettes (OEM) ?$$),
  ($$Can you help us develop new designs and styles (ODM)?$$, 'fr', $$Pouvez-vous nous aider à développer de nouveaux designs et styles (ODM) ?$$),
  ($$Can I order samples first? What is the sample cost and lead time?$$, 'fr', $$Puis-je commander des échantillons d'abord ? Quel est le coût et le délai de l'échantillon ?$$),
  ($$Can you send me a quotation with unit price, MOQ and lead time?$$, 'fr', $$Pouvez-vous m'envoyer un devis avec prix unitaire, MOQ et délai de livraison ?$$),
  ($$Can I choose specific fabrics or request custom fabric options?$$, 'fr', $$Puis-je choisir des tissus spécifiques ou demander des options de tissu personnalisées ?$$),
  ($$What are your payment terms and which currencies do you accept?$$, 'fr', $$Quelles sont vos conditions de paiement et quelles devises acceptez-vous ?$$),
  ($$Do you ship to my country? What shipping options are available?$$, 'fr', $$Livrez-vous dans mon pays ? Quelles options d'expédition sont disponibles ?$$),
  ($$What is the production lead time for a bulk order?$$, 'fr', $$Quel est le délai de production pour une commande en gros ?$$),
  ($$What quality standards and certifications do you have?$$, 'fr', $$Quelles normes de qualité et certifications avez-vous ?$$)
) AS t(src, l, q)
JOIN inquiry_templates m ON m.question = t.src AND m.deleted_at IS NULL
WHERE NOT EXISTS (
  SELECT 1 FROM inquiry_template_translations x
  WHERE x.template_id = m.id AND x.language = t.l AND x.deleted_at IS NULL
);
