/// <reference types="vite/client" />

interface ImportMetaEnv {
  /** 门户站点根地址，用于后台「门户预览」iframe，如 http://localhost:3000 */
  readonly VITE_PORTAL_BASE_URL?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<Record<string, never>, Record<string, never>, any>
  export default component
}
