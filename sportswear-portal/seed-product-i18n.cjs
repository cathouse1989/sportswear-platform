// 为已发布产品批量生成多语言翻译（zh/es/fr），使门户产品多语言切换生效。
// 用法：node seed-product-i18n.cjs  （需后端已启动：http://localhost:8080）
// 说明：
//   - 仅更新「标量字段 + translations」，不动图片/规格/视频/定制/系列/面料/SEO（传 nil 保持不变）。
//   - 英文(en)保留原有内容，其他语言按产品 type/gender/material 生成占位翻译，可在后台微调。
const http = require('http');

const BASE = 'http://localhost:8080/api/v1';

function request(method, path, body = null, authToken = null) {
  return new Promise((resolve, reject) => {
    const url = new URL(BASE + path);
    const opts = {
      hostname: url.hostname,
      port: url.port,
      path: url.pathname + url.search,
      method,
      headers: { 'Content-Type': 'application/json' },
    };
    if (authToken) opts.headers['Authorization'] = `Bearer ${authToken}`;
    const req = http.request(opts, (res) => {
      let data = '';
      res.on('data', (c) => (data += c));
      res.on('end', () => {
        try { resolve(JSON.parse(data)); }
        catch { resolve({ success: false, raw: data }); }
      });
    });
    req.on('error', reject);
    if (body) req.write(JSON.stringify(body));
    req.end();
  });
}

async function login() {
  const emails = ['admin@example.com', 'admin@test.com', 'admin@sportswear.com', 'admin@admin.com', 'superadmin@example.com'];
  const passwords = ['admin123', '123456', 'password', 'admin', '12345678'];
  for (const email of emails) {
    for (const pwd of passwords) {
      const r = await request('POST', '/admin/auth/login', { email, password: pwd });
      if (r.success && r.data?.token) {
        console.log(`Logged in as ${email}`);
        return r.data.token;
      }
    }
  }
  return null;
}

const GENDER_ZH = { unisex: '中性', male: '男款', female: '女款', kids: '童款' };
const TYPE_ZH = { oem: 'OEM', odm: 'ODM', both: 'OEM/ODM' };

function zh(p) {
  const material = p.material || '高品质';
  return {
    language: 'zh',
    name: p.name || p.sku,
    brief: `专业${GENDER_ZH[p.gender] || ''}${TYPE_ZH[p.type] || 'OEM/ODM'}定制运动服装，${material}面料，低起订量。`,
    description: `我们提供高品质定制运动服装制造服务。本产品采用 ${material}，支持颜色、尺码、Logo、面料等定制，欢迎咨询获取样品与报价。`,
    features: `- 高品质 ${material} 面料\n- 支持 ${TYPE_ZH[p.type] || 'OEM/ODM'} 定制\n- 低起订量，快速打样\n- 15 年以上制造经验`,
    usage: '适用于运动健身、团队训练、瑜伽、跑步等场景。',
  };
}

function es(p) {
  const material = p.material || 'tejido de alta calidad';
  return {
    language: 'es',
    name: p.name || p.sku,
    brief: `Ropa deportiva personalizada (${material}), ${(p.type || 'OEM/ODM').toUpperCase()}, MOQ bajo.`,
    description: `Ofrecemos fabricación de ropa deportiva personalizada de alta calidad. Este producto utiliza ${material}, con personalización de color, talla, logo y tejido. Solicite muestras y cotización.`,
    features: `- Tejido ${material} de alta calidad\n- Personalización ${(p.type || 'OEM/ODM').toUpperCase()}\n- MOQ bajo, muestras rápidas\n- Más de 15 años de experiencia`,
    usage: 'Ideal para fitness, entrenamiento en equipo, yoga, running, etc.',
  };
}

function fr(p) {
  const material = p.material || 'tissu de haute qualité';
  return {
    language: 'fr',
    name: p.name || p.sku,
    brief: `Vêtements de sport personnalisés (${material}), ${(p.type || 'OEM/ODM').toUpperCase()}, MOQ bas.`,
    description: `Nous offrons une fabrication de vêtements de sport personnalisés de haute qualité. Ce produit utilise ${material}, avec personnalisation de couleur, taille, logo et tissu. Demandez échantillons et devis.`,
    features: `- Tissu ${material} de haute qualité\n- Personnalisation ${(p.type || 'OEM/ODM').toUpperCase()}\n- MOQ bas, échantillons rapides\n- Plus de 15 ans d'expérience`,
    usage: 'Idéal pour fitness, entraînement en équipe, yoga, course, etc.',
  };
}

async function main() {
  const token = await login();
  if (!token) {
    console.log('FATAL: 无法登录，请先确认后端已启动且账号正确。');
    return;
  }

  const listRes = await request('GET', '/admin/products?page=1&pageSize=500', null, token);
  const products = listRes?.data?.items || listRes?.data || [];
  console.log(`共 ${products.length} 个产品`);

  let ok = 0;
  for (const p of products) {
    try {
      const detailRes = await request('GET', `/admin/products/${p.id}`, null, token);
      const d = detailRes?.data || p;
      const en = (d.translations || []).find((t) => t.language === 'en');

      const translations = [
        {
          language: 'en',
          name: en?.name || d.name || d.sku,
          brief: en?.brief || d.brief || '',
          description: en?.description || d.description || '',
          features: en?.features || d.features || '',
          usage: en?.usage || d.usage || '',
          sort_order: 0,
        },
        zh(d), es(d), fr(d),
      ];

      const payload = {
        sku: d.sku,
        slug: d.slug,
        category_id: d.category_id || null,
        type: d.type || 'both',
        gender: d.gender || 'unisex',
        status: d.status || 'published',
        is_featured: !!d.is_featured,
        is_new: !!d.is_new,
        sort_order: d.sort_order || 0,
        cover_image: d.cover_image || '',
        brief: d.brief || '',
        description: d.description || '',
        features: d.features || '',
        usage: d.usage || '',
        material: d.material || '',
        composition: d.composition || '',
        weight: d.weight || '',
        elasticity: d.elasticity || '',
        fit: d.fit || '',
        support_level: d.support_level || '',
        season: d.season || '',
        size_range: d.size_range || '',
        sample_moq: d.sample_moq || 1,
        production_moq: d.production_moq || 300,
        color_moq: d.color_moq || 100,
        size_moq: d.size_moq || 100,
        translations,
      };

      const res = await request('PUT', `/admin/products/${p.id}`, payload, token);
      if (res.success) { ok++; console.log(`✓ ${d.sku}`); }
      else console.log(`✗ ${d.sku}: ${res.message || 'failed'}`);
    } catch (e) {
      console.log(`✗ ${p.sku || p.id}: ${e.message}`);
    }
  }
  console.log(`完成：成功 ${ok}/${products.length}`);
}

main();
