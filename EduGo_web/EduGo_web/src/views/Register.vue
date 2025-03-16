<template>
  <MainLayout>
    <a-card class="register-card">
      <h1>注册</h1>
      <a-form
        :model="form"
        @submit="handleRegister"
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
          field="email"
          label="邮箱"
          :rules="[
            { required: true, message: '请输入邮箱' },
            { type: 'email', message: '请输入有效的邮箱地址' }
          ]"
          :validate-trigger="['change', 'input']"
        >
          <a-input
            v-model="form.email"
            placeholder="请输入邮箱"
          />
        </a-form-item>

        <a-form-item
          field="password"
          label="密码"
          :rules="[
            { required: true, message: '请输入密码' },
            { minLength: 8, message: '密码至少8位' },
            { 
              validator: validatePasswordStrength,
              message: '密码需包含大小写字母和数字'
            }
          ]"
          :validate-trigger="['change', 'input']"
        >
          <a-input-password
            v-model="form.password"
            placeholder="请输入密码"
          />
        </a-form-item>

        <a-form-item
          field="confirmPassword"
          label="确认密码"
          :rules="[
            { required: true, message: '请确认密码' },
            { validator: validatePasswordMatch, message: '两次输入的密码不一致' }
          ]"
          :validate-trigger="['change', 'input']"
        >
          <a-input-password
            v-model="form.confirmPassword"
            placeholder="请再次输入密码"
          />
        </a-form-item>

        <a-form-item
          field="role"
          label="用户类型"
          :rules="[{ required: true, message: '请选择用户类型' }]"
        >
          <a-radio-group v-model="form.role">
            <a-radio 
              v-for="option in roleOptions" 
              :key="option.value" 
              :value="option.value"
            >
              {{ option.label }}
            </a-radio>
          </a-radio-group>
        </a-form-item>

        <a-form-item>
          <a-button
            type="primary"
            html-type="submit"
            :loading="loading"
          >
            注册
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
import userService, { USER_ROLES, USER_ROLE_NAMES } from '@/services/userService'
import type { UserRole } from '@/services/userService'
import axios from 'axios'

interface RegisterForm {
  username: string
  email: string
  password: string
  confirmPassword: string
  role: UserRole
}

const router = useRouter()
const form = ref<RegisterForm>({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  role: USER_ROLES.STUDENT // 默认为学生
})

// 可选的用户角色
const roleOptions = [
  { label: USER_ROLE_NAMES[USER_ROLES.TEACHER], value: USER_ROLES.TEACHER },
  { label: USER_ROLE_NAMES[USER_ROLES.STUDENT], value: USER_ROLES.STUDENT },
  { label: USER_ROLE_NAMES[USER_ROLES.PARENT], value: USER_ROLES.PARENT }
]

const loading = ref(false)
const errorMessage = ref('')

const validatePasswordStrength = (value: string) => {
  return /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).+$/.test(value)
}

const validatePasswordMatch = (value: string) => {
  return value === form.value.password
}

const handleRegister = async () => {
  loading.value = true
  errorMessage.value = ''

  try {
    await userService.register({
      username: form.value.username,
      email: form.value.email,
      password: form.value.password,
      firstName: '',
      lastName: '',
      role: form.value.role
    })
    
    Message.success('注册成功')
    router.push('/login')
  } catch (error: any) {
    if (axios.isAxiosError(error) && error.response?.data) {
      errorMessage.value = error.response.data.error || '注册失败'
    } else if (error instanceof Error) {
      errorMessage.value = error.message
    } else {
      errorMessage.value = '注册过程中出现错误'
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-card {
  max-width: 400px;
  margin: 2rem auto;
  padding: 2rem;
}
</style>
