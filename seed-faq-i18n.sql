-- 门户 FAQ：补充 header 导航入口 + 中文/西语/法语 FAQ 数据（幂等）
SET client_encoding TO 'UTF8';

-- 1. header 导航加 FAQ（Contact Us 顺延一位）
UPDATE navigations SET sort_order = 7 WHERE type = 'header' AND url = '/contact' AND sort_order = 6;
INSERT INTO navigations (name, type, url, sort_order, is_visible)
SELECT 'FAQ', 'header', '/faq', 6, true
WHERE NOT EXISTS (SELECT 1 FROM navigations WHERE type = 'header' AND url = '/faq' AND deleted_at IS NULL);

-- 2. FAQ 多语言数据（question+language 幂等）
INSERT INTO faqs (question, answer, category, language, sort_order, is_active)
SELECT t.q, t.a, t.c, t.l, t.s, true
FROM (VALUES
  ('最低起订量（MOQ）是多少？',$$我们每个款式的标准起订量是 300 件，打样订单 50 件起。$$,'Orders','zh',1),
  ('生产需要多长时间？',$$打样生产 7-10 天，大货生产 4-6 周，可提供加急订单。$$,'Production','zh',2),
  ('可以定制面料和颜色吗？',$$可以，我们提供全面的面料定制，包括定制颜色、成分和后整理。$$,'Customization','zh',3),
  ('你们的工厂有哪些认证？',$$已通过 ISO 9001:2015、BSCI、OEKO-TEX Standard 100 和 SEDEX 认证。$$,'Quality','zh',4),
  ('你们提供自有品牌（贴牌）服务吗？',$$当然。我们提供全面的自有品牌服务，包括定制品牌、包装和标签。$$,'Services','zh',5),
  ('¿Cuál es la cantidad mínima de pedido (MOQ)?',$$Nuestro MOQ estándar es de 300 unidades por diseño. Los pedidos de muestra comienzan en 50 unidades.$$,'Orders','es',1),
  ('¿Cuánto tiempo tarda la producción?',$$Producción de muestras: 7-10 días. Producción en volumen: 4-6 semanas. Pedidos urgentes disponibles.$$,'Production','es',2),
  ('¿Pueden personalizar telas y colores?',$$Sí, ofrecemos personalización completa de telas, incluyendo colores, composiciones y acabados personalizados.$$,'Customization','es',3),
  ('¿Qué certificaciones tiene su fábrica?',$$Certificados ISO 9001:2015, BSCI, OEKO-TEX Standard 100 y SEDEX.$$,'Quality','es',4),
  ('¿Ofrecen servicios de marca privada?',$$Absolutamente. Servicios completos de marca privada que incluyen marca, embalaje y etiquetado personalizados.$$,'Services','es',5),
  ('Quelle est la quantité minimale de commande (MOQ) ?',$$Notre MOQ standard est de 300 unités par design. Les commandes d''échantillons commencent à 50 unités.$$,'Orders','fr',1),
  ('Combien de temps dure la production ?',$$Production d''échantillons : 7-10 jours. Production en volume : 4-6 semaines. Commandes urgentes disponibles.$$,'Production','fr',2),
  ('Pouvez-vous personnaliser les tissus et les couleurs ?',$$Oui, nous offrons une personnalisation complète des tissus, y compris les couleurs, compositions et finitions personnalisées.$$,'Customization','fr',3),
  ('Quelles certifications votre usine possède-t-elle ?',$$Certifié ISO 9001:2015, BSCI, OEKO-TEX Standard 100 et SEDEX.$$,'Quality','fr',4),
  ('Proposez-vous des services de marque privée ?',$$Absolument. Services complets de marque privée incluant marque, emballage et étiquetage personnalisés.$$,'Services','fr',5)
) AS t(q,a,c,l,s)
WHERE NOT EXISTS (
  SELECT 1 FROM faqs e WHERE e.question = t.q AND e.language = t.l AND e.deleted_at IS NULL
);
