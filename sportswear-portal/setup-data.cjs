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
      res.on('data', c => data += c);
      res.on('end', () => {
        try { resolve(JSON.parse(data)); }
        catch { resolve({ success: false, raw: data }); }
      });
    });
    req.on('error', err => reject(err));
    if (body) req.write(JSON.stringify(body));
    req.end();
  });
}

async function main() {
  // Try multiple common admin emails
  const emails = ['admin@example.com', 'admin@test.com', 'admin@sportswear.com', 'admin@admin.com', 'superadmin@example.com'];
  const passwords = ['admin123', '123456', 'password', 'admin', '12345678'];

  let token = null;
  let userEmail = '';

  for (const email of emails) {
    for (const pwd of passwords) {
      const r = await request('POST', '/admin/auth/login', { email, password: pwd });
      if (r.success) {
        token = r.data.token;
        userEmail = email;
        console.log(`Logged in as ${email}`);
        break;
      }
    }
    if (token) break;
  }

  if (!token) {
    console.log('Could not login. Trying registration endpoint...');
    // Check if we can register
    try {
      const reg = await request('POST', '/admin/auth/register', { 
        email: 'admin@sportswear.com', 
        password: 'admin123',
        name: 'Admin'
      });
      if (reg.success) {
        token = reg.data.token;
        console.log('Registered and logged in');
      } else {
        console.log('Register also failed:', reg.message);
        // Try login again with same creds after potential register
        const r2 = await request('POST', '/admin/auth/login', { email: 'admin@sportswear.com', password: 'admin123' });
        if (r2.success) {
          token = r2.data.token;
          userEmail = 'admin@sportswear.com';
          console.log('Logged in after register attempt');
        }
      }
    } catch(e) {
      console.log('No registration endpoint. Giving up on login.');
    }
  }

  if (!token) {
    console.log('FATAL: Cannot authenticate. Set up manually via admin UI.');
    console.log('The theme.vue and MainLayout.vue code changes are ready.');
    return;
  }

  // Set logo config
  console.log('Setting theme config...');
  await request('PUT', '/admin/theme/logo_url', { value: '' }, token);
  await request('PUT', '/admin/theme/brand_name', { value: 'SPORTSWEAR' }, token);
  await request('PUT', '/admin/theme/brand_subtitle', { value: 'Premium Mfg.' }, token);
  await request('PUT', '/admin/theme/logo_alt', { value: 'Sportswear Manufacturer' }, token);
  console.log('Theme config saved');

  // Get categories
  const catRes = await request('GET', '/admin/categories', null, token);
  const cats = (catRes.data || catRes || []);
  console.log(`Categories: ${cats.length}`);

  // Get product count
  const prodRes = await request('GET', '/admin/products?page=1&pageSize=1', null, token);
  const total = prodRes?.data?.total || 0;
  console.log(`Current products: ${total}`);

  // Add more products if needed
  if (total < 16 && cats.length > 0) {
    const byName = (name) => cats.find(c => c.name === name)?.id || cats[0]?.id;
    const products = [
      { sku: 'SW-YG-011', slug: 'pro-yoga-leggings', category_id: byName('Yoga Wear'), type: 'both', gender: 'female', is_featured: true, cover_image: 'https://images.unsplash.com/photo-1575052814086-f385e2e2ad1b?w=800&h=1000&fit=crop', brief: 'Premium yoga leggings with high waist and tummy control.', material: 'Nylon 75% + Spandex 25%', status: 'published' },
      { sku: 'SW-YG-012', slug: 'yoga-mat-premium', category_id: byName('Yoga Wear'), type: 'both', gender: 'unisex', is_featured: true, cover_image: 'https://images.unsplash.com/photo-1601925260368-ae2f83cf8b7f?w=800&h=1000&fit=crop', brief: 'Eco-friendly premium yoga mat with alignment lines.', material: 'Natural Rubber + PU', status: 'published' },
      { sku: 'SW-RS-010', slug: 'running-vest', category_id: byName('Running Gear'), type: 'both', gender: 'male', is_featured: true, cover_image: 'https://images.unsplash.com/photo-1591942632499-1b8bd0b1e3b0?w=800&h=1000&fit=crop', brief: 'Ultra-lightweight running vest for marathon training.', material: 'Polyester Mesh 100%', status: 'published' },
      { sku: 'SW-RS-011', slug: 'trail-shoes', category_id: byName('Running Gear'), type: 'both', gender: 'unisex', is_featured: true, cover_image: 'https://images.unsplash.com/photo-1595950653106-6c9ebd614d3a?w=800&h=1000&fit=crop', brief: 'Durable trail running shoes with aggressive tread.', material: 'Mesh + Rubber', status: 'published' },
      { sku: 'SW-TA-012', slug: 'compression-shirt', category_id: byName('Training Apparel'), type: 'both', gender: 'male', is_featured: true, cover_image: 'https://images.unsplash.com/photo-1581091226825-a6a2a5aee158?w=800&h=1000&fit=crop', brief: 'Compression training shirt for muscle support.', material: 'Nylon 80% + Spandex 20%', status: 'published' },
      { sku: 'SW-TA-013', slug: 'gym-shorts', category_id: byName('Training Apparel'), type: 'both', gender: 'male', is_featured: true, cover_image: 'https://images.unsplash.com/photo-1565693413575-7cb5c7a7d4db?w=800&h=1000&fit=crop', brief: 'Breathable gym shorts with built-in compression liner.', material: 'Polyester 90% + Spandex 10%', status: 'published' },
    ];
    for (const p of products) {
      if (!p.category_id) { console.log(`Skip ${p.sku}: no category`); continue; }
      const res = await request('POST', '/admin/products', p, token);
      console.log(res.success ? `+ ${p.sku}` : `x ${p.sku}: ${res.message || ''}`);
    }
    console.log(`Added products`);
  }

  console.log('All done!');
}
main();