-- ============================================================
--  门户「联系我们」询盘问题模板词条（4 语言，幂等）
--  背景：联系页新增「询盘问题模板」区域，帮助访客组织想问的问题；
--        点击模板即填入留言框。文案走 i18n_entries，运营可在后台
--        「词条管理」（/i18n）直接维护/覆盖，多语言即时生效。
--  兜底：locales/*.json 的 contact.templates.* 为静态兜底包。
--  执行方式：psql -U <user> -d <db> -f seed-contact-inquiry-templates.sql
--  幂等：按 (key, language) 判断，已存在则跳过。
-- ============================================================
SET client_encoding TO 'UTF8';

INSERT INTO i18n_entries (key, language, value, module, is_active, sort_order)
SELECT t.k, t.l, t.v, 'contact', true, 0
FROM (VALUES
  ('contact.templates.title','en',$$What would you like to ask?$$),
  ('contact.templates.hint','en',$$Tap a question to pre-fill the inquiry form below.$$),
  ('contact.templates.q_moq','en',$$What is your minimum order quantity (MOQ)? Can I start with a small order?$$),
  ('contact.templates.q_oem','en',$$Can you produce custom designs with our logo, colors and labels (OEM)?$$),
  ('contact.templates.q_odm','en',$$Can you help us develop new designs and styles (ODM)?$$),
  ('contact.templates.q_sample','en',$$Can I order samples first? What is the sample cost and lead time?$$),
  ('contact.templates.q_price','en',$$Can you send me a quotation with unit price, MOQ and lead time?$$),
  ('contact.templates.q_fabric','en',$$Can I choose specific fabrics or request custom fabric options?$$),
  ('contact.templates.q_payment','en',$$What are your payment terms and which currencies do you accept?$$),
  ('contact.templates.q_shipping','en',$$Do you ship to my country? What shipping options are available?$$),
  ('contact.templates.q_leadtime','en',$$What is the production lead time for a bulk order?$$),
  ('contact.templates.q_quality','en',$$What quality standards and certifications do you have?$$),

  ('contact.templates.title','zh',$$您想咨询什么问题？$$),
  ('contact.templates.hint','zh',$$点击下方问题，自动填入询盘表单。$$),
  ('contact.templates.q_moq','zh',$$请问最低起订量（MOQ）是多少？可以从小批量开始吗？$$),
  ('contact.templates.q_oem','zh',$$可以按我们的设计生产，加我们的 Logo、颜色和标签吗（OEM 贴牌）？$$),
  ('contact.templates.q_odm','zh',$$可以帮我们开发新的设计款式吗（ODM）？$$),
  ('contact.templates.q_sample','zh',$$可以先订购样品吗？样品费和打样周期是多少？$$),
  ('contact.templates.q_price','zh',$$可以给我发一份报价吗（单价、起订量、交期）？$$),
  ('contact.templates.q_fabric','zh',$$可以选择或定制特定面料吗？$$),
  ('contact.templates.q_payment','zh',$$你们的付款方式和接受的币种是什么？$$),
  ('contact.templates.q_shipping','zh',$$可以发货到我们国家吗？有哪些物流方式？$$),
  ('contact.templates.q_leadtime','zh',$$大货的生产周期是多久？$$),
  ('contact.templates.q_quality','zh',$$你们有哪些质量标准和认证？$$),

  ('contact.templates.title','es',$$¿Qué le gustaría preguntar?$$),
  ('contact.templates.hint','es',$$Toque una pregunta para completar el formulario de consulta.$$),
  ('contact.templates.q_moq','es',$$¿Cuál es su cantidad mínima de pedido (MOQ)? ¿Puedo empezar con un pedido pequeño?$$),
  ('contact.templates.q_oem','es',$$¿Pueden fabricar diseños personalizados con nuestro logo, colores y etiquetas (OEM)?$$),
  ('contact.templates.q_odm','es',$$¿Pueden ayudarnos a desarrollar nuevos diseños y estilos (ODM)?$$),
  ('contact.templates.q_sample','es',$$¿Puedo pedir muestras primero? ¿Cuál es el costo y plazo de la muestra?$$),
  ('contact.templates.q_price','es',$$¿Pueden enviarme una cotización con precio unitario, MOQ y plazo de entrega?$$),
  ('contact.templates.q_fabric','es',$$¿Puedo elegir telas específicas o solicitar opciones de tela personalizadas?$$),
  ('contact.templates.q_payment','es',$$¿Cuáles son sus condiciones de pago y qué monedas aceptan?$$),
  ('contact.templates.q_shipping','es',$$¿Envían a mi país? ¿Qué opciones de envío hay disponibles?$$),
  ('contact.templates.q_leadtime','es',$$¿Cuál es el plazo de producción para un pedido al por mayor?$$),
  ('contact.templates.q_quality','es',$$¿Qué estándares de calidad y certificaciones tienen?$$),

  ('contact.templates.title','fr',$$Que souhaitez-vous demander ?$$),
  ('contact.templates.hint','fr',$$Touchez une question pour préremplir le formulaire de demande.$$),
  ('contact.templates.q_moq','fr',$$Quelle est votre quantité minimum de commande (MOQ) ? Puis-je commencer par une petite commande ?$$),
  ('contact.templates.q_oem','fr',$$Pouvez-vous fabriquer des designs personnalisés avec notre logo, couleurs et étiquettes (OEM) ?$$),
  ('contact.templates.q_odm','fr',$$Pouvez-vous nous aider à développer de nouveaux designs et styles (ODM) ?$$),
  ('contact.templates.q_sample','fr',$$Puis-je commander des échantillons d'abord ? Quel est le coût et le délai de l'échantillon ?$$),
  ('contact.templates.q_price','fr',$$Pouvez-vous m'envoyer un devis avec prix unitaire, MOQ et délai de livraison ?$$),
  ('contact.templates.q_fabric','fr',$$Puis-je choisir des tissus spécifiques ou demander des options de tissu personnalisées ?$$),
  ('contact.templates.q_payment','fr',$$Quelles sont vos conditions de paiement et quelles devises acceptez-vous ?$$),
  ('contact.templates.q_shipping','fr',$$Livrez-vous dans mon pays ? Quelles options d'expédition sont disponibles ?$$),
  ('contact.templates.q_leadtime','fr',$$Quel est le délai de production pour une commande en gros ?$$),
  ('contact.templates.q_quality','fr',$$Quelles normes de qualité et certifications avez-vous ?$$)
) AS t(k,l,v)
WHERE NOT EXISTS (
  SELECT 1 FROM i18n_entries e WHERE e.key = t.k AND e.language = t.l AND e.deleted_at IS NULL
);
