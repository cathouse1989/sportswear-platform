// 全局 Schema.org 结构化数据注入
// 注入 Organization + WebSite + WebPage 基础结构化数据到每个页面
export default defineNuxtPlugin(() => {
  const SITE_NAME = 'OEM/ODM Sportswear Manufacturer'
  const SITE_URL = 'https://sportswear-platform.com'

  useHead({
    script: [
      {
        type: 'application/ld+json',
        children: JSON.stringify({
          '@context': 'https://schema.org',
          '@type': 'Organization',
          name: SITE_NAME,
          alternateName: 'Sportswear OEM/ODM Factory',
          description: 'Professional OEM/ODM sportswear manufacturer. Custom sportswear, activewear, and athletic apparel for global brands.',
          url: SITE_URL,
          logo: `${SITE_URL}/favicon.svg`,
          sameAs: [
            'https://www.linkedin.com/company/sportswear-platform',
            'https://www.instagram.com/sportswear_platform',
            'https://www.facebook.com/sportswearplatform',
          ],
          contactPoint: [
            {
              '@type': 'ContactPoint',
              telephone: '+86-123-4567-8900',
              contactType: 'sales',
              availableLanguage: ['English', '中文', 'Español', 'Français'],
            },
          ],
        }),
      },
      {
        type: 'application/ld+json',
        children: JSON.stringify({
          '@context': 'https://schema.org',
          '@type': 'WebSite',
          name: SITE_NAME,
          url: SITE_URL,
          potentialAction: {
            '@type': 'SearchAction',
            target: {
              '@type': 'EntryPoint',
              urlTemplate: `${SITE_URL}/en/products?search={search_term_string}`,
            },
            'query-input': 'required name=search_term_string',
          },
        }),
      },
    ],
  })
})