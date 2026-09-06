-- ============================================================
--  门户「关于我们」页（slug=about）初始化：页面 + 6 个模块 + 四语言翻译
--  说明：与后台「内容管理 → 页面」编辑的结果完全一致，运营可直接在后台修改/发布。
--  幂等：page 按 slug 判断，module 按 (page_id, type) 判断，可重复执行。
--  config 为多语言 JSONB：文案字段存 { en, zh, es, fr }，前台按当前语言取值。
-- ============================================================
SET client_encoding TO 'UTF8';

DO $$
DECLARE
  p_id uuid;
BEGIN
  SELECT id INTO p_id FROM pages WHERE slug = 'about' AND deleted_at IS NULL;
  IF p_id IS NULL THEN
    INSERT INTO pages (title, slug, type, status, template, sort_order, published_at)
    VALUES ('About Us', 'about', 'normal', 'published', '', 0, now())
    RETURNING id INTO p_id;
  END IF;

  -- 页面标题翻译（四语言）
  INSERT INTO page_translations (page_id, language, title, content, status)
  SELECT p_id, v.language, v.title, '', 'published'
  FROM (VALUES
    ('zh', '关于我们'),
    ('es', 'Sobre Nosotros'),
    ('fr', 'À Propos')
  ) AS v(language, title)
  WHERE NOT EXISTS (
    SELECT 1 FROM page_translations t WHERE t.page_id = p_id AND t.language = v.language AND t.deleted_at IS NULL
  );

  -- 页面模块（6 个结构化区块）
  INSERT INTO page_modules (page_id, type, title, sort_order, is_visible, config)
  SELECT p_id, v.type, v.title, v.sort_order, true, v.config::jsonb
  FROM (VALUES
    ('about_hero', 'About Us', 1, $json${"title":{"en":"About Us","zh":"关于我们","es":"Sobre Nosotros","fr":"À Propos"},"subtitle":{"en":"Professional OEM/ODM sportswear manufacturer with 15+ years of experience.","zh":"拥有15年以上经验的专业OEM/ODM运动服装制造商。","es":"Fabricante profesional OEM/ODM de ropa deportiva con más de 15 años de experiencia.","fr":"Fabricant professionnel OEM/ODM de vêtements de sport avec plus de 15 ans d'expérience."},"bg_image":""}$json$),
    ('about_story', 'Our Story', 2, $json${"title":{"en":"Our Story","zh":"我们的故事","es":"Nuestra Historia","fr":"Notre Histoire"},"body":{"en":"<p>Founded in 2008, we have grown from a small workshop into a leading OEM/ODM sportswear manufacturer serving global brands.</p><p>We specialize in high-performance activewear, yoga wear, running gear, team uniforms and custom sportswear solutions for brands worldwide.</p>","zh":"<p>我们成立于 2008 年，已从一间小作坊成长为服务全球品牌的领先 OEM/ODM 运动服装制造商。</p><p>我们专注于高性能运动服、瑜伽服、跑步装备、团队队服及定制运动服装解决方案。</p>","es":"<p>Fundada en 2008, hemos pasado de ser un pequeño taller a un fabricante líder OEM/ODM de ropa deportiva.</p><p>Nos especializamos en ropa deportiva de alto rendimiento, yoga, running, uniformes y soluciones personalizadas.</p>","fr":"<p>Fondée en 2008, nous sommes passés d'un petit atelier à un fabricant leader OEM/ODM de vêtements de sport.</p><p>Nous sommes spécialisés dans les vêtements de sport haute performance, yoga, running et uniformes.</p>"},"image":"","image_side":"left"}$json$),
    ('about_stats', '', 3, $json${"items":[{"value":"15+","label":{"en":"Years Experience","zh":"年行业经验","es":"Años de Experiencia","fr":"Ans d'Expérience"}},{"value":"500+","label":{"en":"Global Clients","zh":"全球客户","es":"Clientes Globales","fr":"Clients Mondiaux"}},{"value":"50K","label":{"en":"sqm Facility","zh":"平方米工厂","es":"m² de Instalaciones","fr":"m² d'Installations"}}]}$json$),
    ('about_certifications', 'Our Certifications', 4, $json${"title":{"en":"Our Certifications","zh":"资质认证","es":"Certificaciones","fr":"Certifications"},"ref_ids":[]}$json$),
    ('about_process', 'Factory Tour', 5, $json${"title":{"en":"Factory Tour","zh":"工厂参观","es":"Visita a la Fábrica","fr":"Visite de l'Usine"},"steps":[{"icon":"✂️","title":{"en":"Cutting","zh":"裁剪","es":"Corte","fr":"Coupe"},"desc":{"en":"Computerized cutting with 0.1mm precision","zh":"电脑裁剪机，精度高达 0.1mm","es":"Corte computarizado con precisión de 0.1mm","fr":"Découpe informatisée avec une précision de 0.1mm"}},{"icon":"🧵","title":{"en":"Sewing","zh":"缝制","es":"Costura","fr":"Couture"},"desc":{"en":"500+ skilled workers on modern lines","zh":"500 余名熟练工人操作现代化流水线","es":"500+ trabajadores calificados","fr":"Plus de 500 ouvriers qualifiés"}},{"icon":"🖨️","title":{"en":"Printing","zh":"印花","es":"Impresión","fr":"Impression"},"desc":{"en":"Sublimation, screen print & embroidery","zh":"升华、丝印与刺绣","es":"Sublimación, serigrafía y bordado","fr":"Sublimation, sérigraphie et broderie"}},{"icon":"🔬","title":{"en":"Quality Control","zh":"质量检验","es":"Control de Calidad","fr":"Contrôle Qualité"},"desc":{"en":"Multi-point inspection at every stage","zh":"每道工序多点检验","es":"Inspección multipunto en cada etapa","fr":"Inspection multi-points à chaque étape"}},{"icon":"📦","title":{"en":"Packaging","zh":"包装","es":"Embalaje","fr":"Emballage"},"desc":{"en":"Professional packaging to your brand","zh":"按品牌要求专业包装","es":"Embalaje profesional según su marca","fr":"Emballage professionnel selon votre marque"}},{"icon":"🚢","title":{"en":"Shipping","zh":"运输","es":"Envío","fr":"Expédition"},"desc":{"en":"Global logistics for on-time delivery","zh":"全球物流，准时交付","es":"Logística global para entrega puntual","fr":"Logistique mondiale pour livraison ponctuelle"}}]}$json$),
    ('about_cta', 'Ready to Start?', 6, $json${"title":{"en":"Ready to Start Your Project?","zh":"准备开始您的项目？","es":"¿Listo para empezar?","fr":"Prêt à démarrer votre projet ?"},"description":{"en":"Let's talk about your custom sportswear needs.","zh":"让我们聊聊您的定制运动服装需求。","es":"Hablemos de sus necesidades de ropa deportiva.","fr":"Parlons de vos besoins en vêtements de sport."},"button_text":{"en":"Get a Quote","zh":"获取报价","es":"Solicitar Cotización","fr":"Obtenir un Devis"},"button_url":"/contact"}$json$)
  ) AS v(type, title, sort_order, config)
  WHERE NOT EXISTS (
    SELECT 1 FROM page_modules m WHERE m.page_id = p_id AND m.type = v.type AND m.deleted_at IS NULL
  );
END $$;
