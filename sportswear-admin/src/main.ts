import { createApp } from 'vue'
import { createPinia } from 'pinia'

// Element Plus 按需引入：
// - 模板中的 <el-xxx> 组件由 unplugin-vue-components 自动按需引入（含样式）
// - 这里只需全局引入基础样式 + 函数式组件（ElMessage/ElMessageBox）的样式
// - 中文语言包由 App.vue 的 <el-config-provider> 提供（组件树与函数式调用均生效）
import 'element-plus/theme-chalk/base.css'
import 'element-plus/theme-chalk/el-message.css'
import 'element-plus/theme-chalk/el-message-box.css'

import App from './App.vue'
import router from './router'
import { vPermission } from './directives/permission'

const app = createApp(App)

app.use(createPinia())
app.use(router)
// 按钮级权限指令：<el-button v-permission="'product:publish'">
app.directive('permission', vPermission)

app.mount('#app')