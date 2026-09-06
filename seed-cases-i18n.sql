-- ============================================================
--  案例展示（/cases）多语言补配：初始化 3 个案例（英文源，status=published），
--  并写入 zh / es / fr 三语言翻译（status=published），门户据此本地化展示。
--  幂等：案例按 slug 判断，翻译按 (case_id, language) 判断，可重复执行。
--  执行方式：
--    docker cp seed-cases-i18n.sql sportswear-postgres:/tmp/seed-cases-i18n.sql
--    docker exec sportswear-postgres psql -U postgres -d sportswear_platform -f /tmp/seed-cases-i18n.sql
--  或：
--    psql "$DATABASE_URL" -f seed-cases-i18n.sql
-- ============================================================
SET client_encoding TO 'UTF8';

-- 1) 初始化案例（英文源内容）
INSERT INTO cases (title, slug, client_industry, project_type, products, client_need, problem, solution, process, result, cover_image, status, published_at, created_at, updated_at)
VALUES
(
  $$Full-Custom Yoga Collection for a North American Activewear Brand$$,
  $$yoga-wear-brand-launch$$,
  $$Activewear / Yoga$$,
  $$odm$$,
  $$Yoga Leggings, Sports Bras, Tank Tops$$,
  $$<p>The client, a fast-growing North American activewear brand, needed a complete yoga collection with their own fabrics, fits and branding — from first sample to finished goods.</p>$$,
  $$<p>They had struggled with previous suppliers on sizing consistency and color matching across multi-fabric styles, which delayed their launch twice.</p>$$,
  $$<p>We proposed a custom blended fabric (nylon/spandex) with a shared color palette and a unified size spec, then produced graded samples for every style in one cycle.</p>$$,
  $$<p>Fabric development → pattern making → fit samples → bulk production → quality inspection → packing.</p>$$,
  $$<p>Launched on schedule with a 98.5% first-pass quality rate, and the client returned for three follow-up seasonal collections.</p>$$,
  $$https://images.unsplash.com/photo-1544367567-0f2fcb009e0b?w=800&h=500&fit=crop$$,
  $$published$$, now(), now(), now()
),
(
  $$Sublimated Football Kits for a European Football Club$$,
  $$sublimated-football-uniforms$$,
  $$Team Sports$$,
  $$oem$$,
  $$Jerseys, Shorts, Socks$$,
  $$<p>The club required match and training kits with full sublimation printing and precise club crest reproduction.</p>$$,
  $$<p>The crest contained fine gradients that previous suppliers could not reproduce without cracking or fading after washing.</p>$$,
  $$<p>We used high-resolution sublimation with heat-sealed crests and tested for wash fastness before bulk production.</p>$$,
  $$<p>Artwork approval → sublimation proofing → sample kits → bulk printing → sewing → QC → delivery.</p>$$,
  $$<p>Delivered 3,500 kits ahead of the season, with the crest remaining crisp after 50+ washes.</p>$$,
  $$https://images.unsplash.com/photo-1517466787929-bc90951d0974?w=800&h=500&fit=crop$$,
  $$published$$, now(), now(), now()
),
(
  $$Private Label Running Gear for an Australian Retailer$$,
  $$private-label-running-gear$$,
  $$Running / Retail$$,
  $$private_label$$,
  $$Running Jackets, Tights, Shorts$$,
  $$<p>The retailer wanted a private-label running range with their own branding, hangtags and packaging.</p>$$,
  $$<p>They needed a low MOQ to test the market first, while keeping premium fabric quality and a consistent fit.</p>$$,
  $$<p>We ran a small-batch production line with shared fabric stock and modular trims, allowing low MOQ without compromising quality.</p>$$,
  $$<p>Fabric selection → sample approval → small-batch production → custom labeling → packing → shipping.</p>$$,
  $$<p>The test batch sold out in six weeks, and the client scaled to a full seasonal program.</p>$$,
  $$https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=800&h=500&fit=crop$$,
  $$published$$, now(), now(), now()
)
ON CONFLICT (slug) DO NOTHING;

-- 2) 中文翻译
INSERT INTO case_translations (case_id, language, title, client_need, problem, solution, process, result, status, created_at, updated_at)
SELECT c.id, t.language, t.title, t.client_need, t.problem, t.solution, t.process, t.result, 'published', now(), now()
FROM cases c
JOIN (VALUES
  (
    $$yoga-wear-brand-launch$$, $$zh$$,
    $$北美运动品牌全定制瑜伽系列$$,
    $$<p>客户是一家快速成长的北美运动服饰品牌，需要一套拥有自有面料、版型与品牌元素的完整瑜伽系列，从首版样品到成品交付。</p>$$,
    $$<p>此前供应商在尺码一致性与多面料配色上多次出错，导致其上市计划两次延期。</p>$$,
    $$<p>我们提出定制混纺面料（锦纶/氨纶）+ 统一配色体系 + 统一尺码规格，并在一个周期内完成全款式的放码样品。</p>$$,
    $$<p>面料开发 → 打版制样 → 试身样 → 大货生产 → 质量检验 → 包装。</p>$$,
    $$<p>如期上市，首检合格率达 98.5%，客户随后连续追加三个季度系列订单。</p>$$
  ),
  (
    $$sublimated-football-uniforms$$, $$zh$$,
    $$欧洲足球俱乐部升华印球队服$$,
    $$<p>俱乐部需要比赛与训练服，采用全幅升华印花并精确还原队徽。</p>$$,
    $$<p>队徽含细腻渐变色，此前供应商印刷后易开裂或水洗后褪色。</p>$$,
    $$<p>我们采用高清升华印花 + 热封队徽，并在大货前完成水洗牢度测试。</p>$$,
    $$<p>图案确认 → 升华打样 → 样品套件 → 批量印花 → 缝制 → 质检 → 交付。</p>$$,
    $$<p>赛季前交付 3,500 套队服，队徽经过 50 余次水洗依然清晰。</p>$$
  ),
  (
    $$private-label-running-gear$$, $$zh$$,
    $$澳大利亚零售商自有品牌跑步装备$$,
    $$<p>零售商希望推出自有品牌跑步系列，含自有品牌吊牌与包装。</p>$$,
    $$<p>他们需要低起订量先行测试市场，同时保证面料品质与版型一致。</p>$$,
    $$<p>我们启用小批量产线，共享面料库存与模块化辅料，在低起订量下仍保证品质。</p>$$,
    $$<p>面料选型 → 样品确认 → 小批量生产 → 定制吊牌 → 包装 → 发货。</p>$$,
    $$<p>试销批次六周售罄，客户随即升级为整季计划。</p>$$
  )
) AS t(case_slug, language, title, client_need, problem, solution, process, result)
  ON c.slug = t.case_slug AND c.deleted_at IS NULL
WHERE NOT EXISTS (
  SELECT 1 FROM case_translations e
  WHERE e.case_id = c.id AND e.language = t.language AND e.deleted_at IS NULL
);

-- 3) 西班牙语翻译
INSERT INTO case_translations (case_id, language, title, client_need, problem, solution, process, result, status, created_at, updated_at)
SELECT c.id, t.language, t.title, t.client_need, t.problem, t.solution, t.process, t.result, 'published', now(), now()
FROM cases c
JOIN (VALUES
  (
    $$yoga-wear-brand-launch$$, $$es$$,
    $$Colección de yoga totalmente personalizada para una marca norteamericana$$,
    $$<p>El cliente, una marca de ropa deportiva norteamericana en rápido crecimiento, necesitaba una colección de yoga completa con sus propias telas, ajustes y marca, desde la primera muestra hasta el producto terminado.</p>$$,
    $$<p>Había tenido problemas con proveedores anteriores en cuanto a la consistencia de tallas y la igualación de colores en estilos de múltiples tejidos, lo que retrasó dos veces su lanzamiento.</p>$$,
    $$<p>Propusimos un tejido mezclado personalizado (nailon/spandex) con una paleta de colores compartida y una especificación de tallas unificada, y produjimos muestras escaladas para cada estilo en un solo ciclo.</p>$$,
    $$<p>Desarrollo de tejido → patronaje → muestras de ajuste → producción en serie → control de calidad → embalaje.</p>$$,
    $$<p>Se lanzó a tiempo con una tasa de calidad a la primera del 98,5%, y el cliente volvió para tres colecciones de temporada posteriores.</p>$$
  ),
  (
    $$sublimated-football-uniforms$$, $$es$$,
    $$Equipaciones de fútbol por sublimación para un club europeo$$,
    $$<p>El club necesitaba equipaciones de partido y entrenamiento con impresión por sublimación completa y reproducción precisa del escudo.</p>$$,
    $$<p>El escudo contenía degradados finos que los proveedores anteriores no podían reproducir sin agrietarse ni desteñirse tras el lavado.</p>$$,
    $$<p>Utilizamos sublimación de alta resolución con escudos termosellados y realizamos pruebas de solidez al lavado antes de la producción en serie.</p>$$,
    $$<p>Aprobación del arte → prueba de sublimación → equipaciones de muestra → impresión en serie → confección → control de calidad → entrega.</p>$$,
    $$<p>Entregamos 3.500 equipaciones antes del inicio de temporada, y el escudo se mantiene nítido tras más de 50 lavados.</p>$$
  ),
  (
    $$private-label-running-gear$$, $$es$$,
    $$Equipamiento de running de marca blanca para un minorista australiano$$,
    $$<p>El minorista quería una gama de running de marca blanca con su propia marca, etiquetas colgantes y embalaje.</p>$$,
    $$<p>Necesitaba un MOQ bajo para probar primero el mercado, manteniendo a la vez la calidad premium del tejido y un ajuste uniforme.</p>$$,
    $$<p>Pusimos en marcha una línea de producción en lotes pequeños con stock de tejido compartido y detalles modulares, permitiendo un MOQ bajo sin renunciar a la calidad.</p>$$,
    $$<p>Selección de tejido → aprobación de muestras → producción en lotes pequeños → etiquetado personalizado → embalaje → envío.</p>$$,
    $$<p>El lote de prueba se agotó en seis semanas y el cliente amplió a un programa de temporada completo.</p>$$
  )
) AS t(case_slug, language, title, client_need, problem, solution, process, result)
  ON c.slug = t.case_slug AND c.deleted_at IS NULL
WHERE NOT EXISTS (
  SELECT 1 FROM case_translations e
  WHERE e.case_id = c.id AND e.language = t.language AND e.deleted_at IS NULL
);

-- 4) 法语翻译
INSERT INTO case_translations (case_id, language, title, client_need, problem, solution, process, result, status, created_at, updated_at)
SELECT c.id, t.language, t.title, t.client_need, t.problem, t.solution, t.process, t.result, 'published', now(), now()
FROM cases c
JOIN (VALUES
  (
    $$yoga-wear-brand-launch$$, $$fr$$,
    $$Collection de yoga entièrement personnalisée pour une marque nord-américaine$$,
    $$<p>Le client, une marque de vêtements de sport nord-américaine en forte croissance, avait besoin d'une collection de yoga complète avec ses propres tissus, coupes et image de marque, du premier échantillon au produit fini.</p>$$,
    $$<p>Il avait rencontré des difficultés avec ses précédents fournisseurs sur la régularité des tailles et la correspondance des couleurs entre styles multi-tissus, ce qui avait retardé deux fois son lancement.</p>$$,
    $$<p>Nous avons proposé un tissu mélangé sur mesure (nylon/élasthanne) avec une palette de couleurs partagée et un tableau de tailles unifié, puis produit des échantillons gradés pour chaque style en un seul cycle.</p>$$,
    $$<p>Développement du tissu → patronage → échantillons d'ajustement → production en série → contrôle qualité → emballage.</p>$$,
    $$<p>Lancement dans les délais avec un taux de qualité au premier contrôle de 98,5 %, et le client est revenu pour trois collections saisonnières suivantes.</p>$$
  ),
  (
    $$sublimated-football-uniforms$$, $$fr$$,
    $$Maillots de football sublimés pour un club européen$$,
    $$<p>Le club avait besoin de tenues de match et d'entraînement avec impression par sublimation intégrale et reproduction précise de l'écusson.</p>$$,
    $$<p>L'écusson contenait des dégradés fins que les fournisseurs précédents ne pouvaient pas reproduire sans fissures ni décoloration après lavage.</p>$$,
    $$<p>Nous avons utilisé la sublimation haute résolution avec écussons thermocollés et testé la solidité au lavage avant la production en série.</p>$$,
    $$<p>Validation du graphisme → épreuve de sublimation → tenues d'échantillon → impression en série → confection → contrôle qualité → livraison.</p>$$,
    $$<p>Livraison de 3 500 tenues avant le début de saison, l'écusson restant net après plus de 50 lavages.</p>$$
  ),
  (
    $$private-label-running-gear$$, $$fr$$,
    $$Équipement de course en marque blanche pour un détaillant australien$$,
    $$<p>Le détaillant souhaitait une gamme de course en marque blanche avec sa propre marque, ses étiquettes et son emballage.</p>$$,
    $$<p>Il avait besoin d'un MOQ faible pour tester d'abord le marché, tout en conservant une qualité de tissu premium et une coupe régulière.</p>$$,
    $$<p>Nous avons mis en place une ligne de production en petites séries avec un stock de tissu partagé et des finitions modulaires, permettant un MOQ faible sans compromettre la qualité.</p>$$,
    $$<p>Sélection du tissu → validation des échantillons → production en petites séries → étiquetage personnalisé → emballage → expédition.</p>$$,
    $$<p>Le lot test s'est vendu en six semaines et le client est passé à un programme saisonnier complet.</p>$$
  )
) AS t(case_slug, language, title, client_need, problem, solution, process, result)
  ON c.slug = t.case_slug AND c.deleted_at IS NULL
WHERE NOT EXISTS (
  SELECT 1 FROM case_translations e
  WHERE e.case_id = c.id AND e.language = t.language AND e.deleted_at IS NULL
);

