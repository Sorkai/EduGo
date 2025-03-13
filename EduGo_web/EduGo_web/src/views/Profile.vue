<template>
  <MainLayout>
    <div class="profile-container">
      <h1>用户信息</h1>
      <div class="profile-content">
        <div class="profile-item">
          <label>用户名：</label>
          <span>{{ user.username }}</span>
        </div>
        <div class="profile-item">
          <label>邮箱：</label>
          <span>{{ user.email }}</span>
        </div>
        <div class="profile-item">
          <label>注册时间：</label>
          <span>{{ user.createdAt }}</span>
        </div>
      </div>
    </div>
  </MainLayout>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted } from 'vue'
import MainLayout from '../components/MainLayout.vue'
import userService from '../services/userService'
import type { UserProfile } from '../services/userService'

export default defineComponent({
  name: 'Profile',
  components: {
    MainLayout
  },
  setup() {
    const user = ref<UserProfile>({
      username: '',
      email: '',
      createdAt: new Date().toISOString()
    })

    const errorMessage = ref('')

    onMounted(async () => {
      try {
        const response = await userService.getUserProfile()
        user.value = response
        errorMessage.value = ''
      } catch (error) {
        console.error('获取用户信息失败:', error)
        errorMessage.value = '获取用户信息失败，请检查网络连接或重新登录'
      }
    })

    return {
      user,
      errorMessage
    }
  }
})
</script>

<style scoped>
.profile-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.profile-content {
  margin-top: 20px;
}

.profile-item {
  margin: 10px 0;
  padding: 10px;
  border-bottom: 1px solid #eee;
}

.profile-item label {
  font-weight: bold;
  margin-right: 10px;
}
</style>
