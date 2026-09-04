<template>
  <nav v-if="items?.length" class="flex items-center gap-2 text-sm text-gray-500" :aria-label="$t('common.breadcrumb', 'Breadcrumb')">
    <template v-for="(item, i) in items" :key="i">
      <svg v-if="i > 0" class="w-4 h-4 shrink-0 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
      </svg>
      <span v-if="i === items.length - 1" class="text-[#0D1B2A] font-medium truncate max-w-[200px]" :aria-current="'page'">
        {{ item.name }}
      </span>
      <NuxtLink v-else :to="item.to" class="hover:text-[#0D1B2A] transition-colors truncate max-w-[160px]">
        {{ item.name }}
      </NuxtLink>
    </template>
  </nav>
</template>

<script setup lang="ts">
import { buildBreadcrumbSchema } from '~/composables/useSeo'

export interface BreadcrumbItemData {
  name: string
  to?: string
}

interface Props {
  items: BreadcrumbItemData[]
}

const props = defineProps<Props>()
const route = useRoute()
const SITE_URL = 'https://sportswear-platform.com'

if (props.items?.length) {
  const schema = buildBreadcrumbSchema(
    props.items.map((item, index) => ({
      name: item.name,
      item: item.to
        ? `${SITE_URL}${item.to.startsWith('/') ? '' : '/'}${item.to}`
        : `${SITE_URL}${route.path}`,
    }))
  )
  useHead({
    script: [
      {
        type: 'application/ld+json',
        children: JSON.stringify(schema),
      },
    ],
  })
}
</script>