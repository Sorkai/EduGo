<template>
  <MainLayout>
    <a-card class="login-card">
      <h1>登录</h1>
      <a-form
        :model="form"
        @submit="handleLogin"
        auto-label-width
      >
        <a-form-item
          field="username"
          label="用户名"
          :rules="[{ required: true, message: '请输入用户名' }]"
          :validate-trigger="['change', 'input']"
        >
          <a-input
            v-model="form.username"
            placeholder="请输入用户名"
          />
        </a-form-item>

        <a-form-item
          field="password"
          label="密码"
          :rules="[{ required: true, message: '请输入密码' }]"
          :validate-trigger="['change', 'input']"
        >
          <a-input-password
            v-model="form.password"
            placeholder="请输入密码"
            allow-clear
          />
        </a-form-item>

        <a-form-item>
          <a-space direction="vertical" size="large">
            <a-checkbox v-model="rememberMe">记住我</a-checkbox>
            <a-link href="/forgot-password">忘记密码？</a-link>
          </a-space>
        </a-form-item>

        <a-divider>或</a-divider>

        <a-form-item>
          <a-space direction="vertical" size="large" fill>
            <a-button type="outline" long @click="handleOAuthLogin('github')">
              <template #icon>
                <icon-github />
              </template>
              使用 GitHub 登录
            </a-button>
            <a-button type="outline" long @click="handleOAuthLogin('google')">
              <template #icon>
                <icon-google />
              </template>
              使用 Google 登录
            </a-button>
          </a-space>
        </a-form-item>

        <a-form-item>
          <a-button
            type="primary"
            html-type="submit"
            :loading="loading"
          >
            登录
          </a-button>
        </a-form-item>

        <a-alert
          v-if="errorMessage"
          type="error"
          :title="errorMessage"
          closable
        />
      </a-form>
    </a-card>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import MainLayout from '@/components/MainLayout.vue'
import { Message } from '@arco-design/web-vue'
import { IconGithub, IconGoogle } from '@arco-design/web-vue/es/icon'

interface LoginForm {
  username: string
  password: string
}

const router = useRouter()
const form = ref<LoginForm>({
  username: '',
  password: ''
})
const rememberMe = ref(false)
const loading = ref(false)
const errorMessage = ref('')

const handleOAuthLogin = (provider: string) => {
  window.location.href = `/api/oauth/${provider}`
}

import userService from '@/services/userService'
import axios from 'axios'

const handleLogin = async () => {
  loading.value = true
  errorMessage.value = ''

  try {
    const data = await userService.login(form.value.username, form.value.password)
    
    // 保存token
    if (rememberMe.value) {
      localStorage.setItem('token', data.token)
    } else {
      sessionStorage.setItem('token', data.token)
    }

    Message.success('登录成功')
    router.push('/')
  } catch (error: any) {
    if (axios.isAxiosError(error) && error.response?.data) {
      errorMessage.value = error.response.data.error || '登录失败'
    } else if (error instanceof Error) {
      errorMessage.value = error.message
    } else {
      errorMessage.value = '登录过程中出现错误'
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-card {
  max-width: 400px;
  margin: 2rem auto;
  padding: 2rem;
}

.arco-alert {
  margin-bottom: 1rem;
}

.arco-divider {
  margin: 1rem 0;
}
</style>
