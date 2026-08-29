// eslint 扁平配置（P0-#6 工程护栏）
// 策略：结构规则严格、风格规则交由 Prettier、历史遗留的 any 用法暂关闭（P0-#5 契约硬化后再收紧）
import pluginVue from 'eslint-plugin-vue'
import { defineConfigWithVueTs, vueTsConfigs } from '@vue/eslint-config-typescript'
import globals from 'globals'

export default defineConfigWithVueTs(
  { files: ['**/*.{ts,mts,tsx,vue}'] },
  { ignores: ['dist/**', 'node_modules/**', 'coverage/**', '*.log'] },
  pluginVue.configs['flat/essential'],
  vueTsConfigs.recommended,
  {
    languageOptions: {
      globals: { ...globals.browser },
    },
    rules: {
      // 结构与正确性
      'no-console': ['warn', { allow: ['warn', 'error'] }],
      'prefer-const': 'error',
      eqeqeq: ['error', 'smart'],
      '@typescript-eslint/no-unused-vars': ['warn', { argsIgnorePattern: '^_', varsIgnorePattern: '^_' }],
      // 历史债务：P0-#5 契约硬化时移除豁免
      '@typescript-eslint/no-explicit-any': 'off',
      // 风格交给 Prettier
      'vue/multi-word-component-names': 'off',
      'vue/max-attributes-per-line': 'off',
      'vue/singleline-html-element-content-newline': 'off',
      'vue/html-indent': 'off',
      'vue/html-closing-bracket-newline': 'off',
    },
  },
)
