import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'

export default defineConfig({
  plugins: [
    vue(),
    // Element Plus 按需自动引入（组件 + 样式），避免全量打包。
    // dts 保持关闭：避免生成的精确类型声明（如 ElTable scope.row: DefaultRow）
    // 与各页面显式声明的实体类型（Blog/Product 等）冲突；运行时按需引入不受影响。
    Components({
      resolvers: [ElementPlusResolver()],
      dts: false,
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    host: '0.0.0.0',
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      // 上传文件的静态资源代理：媒体记录存的是相对路径（/uploads/...），
      // 后台预览/媒体库选择器需要经此代理到后端 8080 的静态文件服务
      '/uploads': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  build: {
    // 拆分公共依赖 chunk，利于浏览器长期缓存（依赖版本不变则 hash 不变）
    rollupOptions: {
      output: {
        manualChunks: {
          'vue-vendor': ['vue', 'vue-router', 'pinia'],
          'axios': ['axios'],
        },
      },
    },
  },
})