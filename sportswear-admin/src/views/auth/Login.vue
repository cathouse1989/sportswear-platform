<template>
  <div class="login-page">
    <div class="login-card">
      <h1 class="title">OEM/ODM 管理后台</h1>
      <p class="subtitle">国际化运动服饰 OEM/ODM 平台</p>

      <el-form ref="formRef" :model="form" :rules="rules" @keyup.enter="handleLogin">
        <el-form-item prop="email">
          <el-input v-model="form.email" placeholder="邮箱" size="large" @input="onEmailInput">
            <template #prefix><el-icon><Message /></el-icon></template>
          </el-input>
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="form.password" type="password" placeholder="密码" size="large" show-password>
            <template #prefix><el-icon><Lock /></el-icon></template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="large" class="login-btn" :loading="loading" @click="handleLogin">
            登录
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { Message, Lock } from '@element-plus/icons-vue'

const router = useRouter()
const authStore = useAuthStore()
const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  email: '',
  password: '',
})

const rules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { whitespace: true, message: '邮箱不能为空格', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { whitespace: true, message: '密码不能为空格', trigger: 'blur' },
  ],
}

// 邮箱实时去除所有空白字符（邮箱本身不含空格，从源头避免复制粘贴带入）
function onEmailInput(val: string) {
  form.email = val.replace(/\s+/g, '')
}

async function handleLogin() {
  await formRef.value?.validate()
  // 兜底：提交前去除首尾空格（密码保留内部字符，仅清理首尾误输入的空格）
  form.email = form.email.trim().replace(/\s+/g, '')
  form.password = form.password.trim()
  loading.value = true
  try {
    await authStore.login(form.email, form.password)
    ElMessage.success('登录成功')
    router.push('/dashboard')
  } catch {
    // 错误已在拦截器提示
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1e3a8a 0%, #3b82f6 100%);
}
.login-card {
  width: 400px;
  padding: 40px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
}
.title {
  text-align: center;
  font-size: 22px;
  margin-bottom: 8px;
  color: #1e3a8a;
}
.subtitle {
  text-align: center;
  font-size: 13px;
  color: #6b7280;
  margin-bottom: 30px;
}
.login-btn {
  width: 100%;
}
</style>