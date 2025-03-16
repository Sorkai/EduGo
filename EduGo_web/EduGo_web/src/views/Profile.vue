<template>
  <MainLayout>
    <a-card class="profile-container">
      <template #title>
        <a-space>
          <a-avatar :size="64" :style="{ backgroundColor: '#3370ff' }">
            {{ user.username ? user.username.charAt(0).toUpperCase() : 'U' }}
          </a-avatar>
          <div class="profile-title">
            <h2>{{ user.username }}</h2>
            <a-tag color="blue">{{ user.role || '普通用户' }}</a-tag>
          </div>
        </a-space>
      </template>
      
      <a-spin :loading="loading">
        <a-alert v-if="errorMessage" type="error" :content="errorMessage" closable class="mb-4" />
        
        <a-tabs default-active-key="1">
          <a-tab-pane key="1" title="个人资料">
            <a-descriptions :data="profileData" layout="vertical" bordered />
          </a-tab-pane>
          
          <a-tab-pane key="2" title="修改资料">
            <a-form :model="form" @submit="handleUpdateProfile" auto-label-width>
              <a-form-item field="email" label="邮箱" :rules="[{ type: 'email', message: '请输入有效的邮箱地址' }]">
                <a-input v-model="form.email" placeholder="请输入邮箱" />
              </a-form-item>
              
              <a-form-item field="firstName" label="名">
                <a-input v-model="form.firstName" placeholder="请输入名" />
              </a-form-item>
              
              <a-form-item field="lastName" label="姓">
                <a-input v-model="form.lastName" placeholder="请输入姓" />
              </a-form-item>
              
              <a-form-item>
                <a-button type="primary" html-type="submit" :loading="updating">保存修改</a-button>
              </a-form-item>
            </a-form>
          </a-tab-pane>
          
          <a-tab-pane key="3" title="修改密码">
            <a-form :model="passwordForm" @submit="handleUpdatePassword" auto-label-width>
              <a-form-item field="oldPassword" label="当前密码" :rules="[{ required: true, message: '请输入当前密码' }]">
                <a-input-password v-model="passwordForm.oldPassword" placeholder="请输入当前密码" />
              </a-form-item>
              
              <a-form-item field="newPassword" label="新密码" :rules="[
                { required: true, message: '请输入新密码' },
                { minLength: 8, message: '密码至少8位' },
                { validator: validatePasswordStrength, message: '密码需包含大小写字母、数字和特殊字符' }
              ]">
                <a-input-password v-model="passwordForm.newPassword" placeholder="请输入新密码" />
              </a-form-item>
              
              <a-form-item field="confirmPassword" label="确认新密码" :rules="[
                { required: true, message: '请确认新密码' },
                { validator: validatePasswordMatch, message: '两次输入的密码不一致' }
              ]">
                <a-input-password v-model="passwordForm.confirmPassword" placeholder="请再次输入新密码" />
              </a-form-item>
              
              <a-form-item>
                <a-button type="primary" html-type="submit" :loading="updatingPassword">更新密码</a-button>
              </a-form-item>
            </a-form>
          </a-tab-pane>
        </a-tabs>
      </a-spin>
    </a-card>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import MainLayout from '@/components/MainLayout.vue'
import userService from '@/services/userService'
import type { UserProfile } from '@/services/userService'
import axios from 'axios'

const router = useRouter()
const loading = ref(true)
const updating = ref(false)
const updatingPassword = ref(false)
const errorMessage = ref('')

const user = ref<UserProfile>({
  id: 0,
  username: '',
  email: '',
  firstName: '',
  lastName: '',
  createdAt: ''
})

const form = reactive({
  email: '',
  firstName: '',
  lastName: ''
})

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const profileData = computed(() => [
  {
    label: '用户名',
    value: user.value.username
  },
  {
    label: '邮箱',
    value: user.value.email
  },
  {
    label: '姓名',
    value: `${user.value.firstName || ''} ${user.value.lastName || ''}`.trim() || '未设置'
  },
  {
    label: '注册时间',
    value: user.value.createdAt ? new Date(user.value.createdAt).toLocaleString() : '未知'
  }
])

const validatePasswordStrength = (value: string) => {
  return /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]).{8,}$/.test(value)
}

const validatePasswordMatch = (value: string) => {
  return value === passwordForm.newPassword
}

onMounted(async () => {
  try {
    loading.value = true
    const userProfile = await userService.getUserProfile()
    user.value = userProfile
    
    // 初始化表单数据
    form.email = userProfile.email
    form.firstName = userProfile.firstName
    form.lastName = userProfile.lastName
    
    errorMessage.value = ''
  } catch (error: any) {
    console.error('获取用户信息失败:', error)
    if (axios.isAxiosError(error) && error.response?.status === 401) {
      Message.error('登录已过期，请重新登录')
      router.push('/login')
    } else {
      errorMessage.value = '获取用户信息失败，请检查网络连接或重新登录'
    }
  } finally {
    loading.value = false
  }
})

const handleUpdateProfile = async () => {
  try {
    updating.value = true
    const updatedUser = await userService.updateUser({
      email: form.email,
      firstName: form.firstName,
      lastName: form.lastName
    })
    
    user.value = updatedUser
    Message.success('个人资料更新成功')
  } catch (error: any) {
    if (axios.isAxiosError(error) && error.response?.data) {
      errorMessage.value = error.response.data.error || '更新个人资料失败'
    } else if (error instanceof Error) {
      errorMessage.value = error.message
    } else {
      errorMessage.value = '更新个人资料失败'
    }
  } finally {
    updating.value = false
  }
}

const handleUpdatePassword = async () => {
  try {
    updatingPassword.value = true
    await userService.resetPassword({
      oldPassword: passwordForm.oldPassword,
      newPassword: passwordForm.newPassword
    })
    
    Message.success('密码更新成功')
    
    // 清空密码表单
    passwordForm.oldPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
  } catch (error: any) {
    if (axios.isAxiosError(error) && error.response?.data) {
      errorMessage.value = error.response.data.error || '更新密码失败'
    } else if (error instanceof Error) {
      errorMessage.value = error.message
    } else {
      errorMessage.value = '更新密码失败'
    }
  } finally {
    updatingPassword.value = false
  }
}
</script>

<style scoped>
.profile-container {
  max-width: 800px;
  margin: 2rem auto;
}

.profile-title {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.profile-title h2 {
  margin: 0 0 8px 0;
}

.mb-4 {
  margin-bottom: 16px;
}
</style>
