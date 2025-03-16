<template>
  <a-layout class="layout">
    <a-layout-header class="header">
      <div class="logo">
        <img src="@/assets/logo.svg" alt="EduGo Logo" />
        <h1>EduGo</h1>
      </div>
      <div class="menu-wrapper">
        <a-menu
          mode="horizontal"
          :selected-keys="[activeKey]"
          @menu-item-click="handleMenuClick"
        >
          <a-menu-item key="assessment">
            <template #icon><icon-apps /></template>
            智能测评
          </a-menu-item>
          <a-menu-item key="teaching">
            <template #icon><icon-book /></template>
            智能教学
          </a-menu-item>
          <a-menu-item key="vr">
            <template #icon><icon-common /></template>
            虚拟现实
          </a-menu-item>
          <a-menu-item key="robot">
            <template #icon><icon-robot /></template>
            教育机器人
          </a-menu-item>
          <a-menu-item key="evaluation">
            <template #icon><icon-dashboard /></template>
            智能评价
          </a-menu-item>
        </a-menu>
      </div>
      <div class="user-actions">
        <template v-if="isLoggedIn">
          <a-dropdown trigger="click">
            <a-avatar :style="{ backgroundColor: '#3370ff' }">
              {{ userInitial }}
            </a-avatar>
            <template #content>
              <a-doption @click="navigateTo('/profile')">
                <template #icon><icon-user /></template>
                个人中心
              </a-doption>
              <a-doption @click="handleLogout">
                <template #icon><icon-export /></template>
                退出登录
              </a-doption>
            </template>
          </a-dropdown>
        </template>
        <template v-else>
          <a-space>
            <a-button type="text" @click="navigateTo('/login')">登录</a-button>
            <a-button type="primary" @click="navigateTo('/register')">注册</a-button>
          </a-space>
        </template>
      </div>
    </a-layout-header>

    <a-layout>
      <a-layout-content class="content">
        <slot></slot>
      </a-layout-content>
    </a-layout>

    <a-layout-footer class="footer">
      <p>© {{ currentYear }} EduGo 数智化教育平台</p>
      <div class="footer-links">
        <a-link href="#">关于我们</a-link>
        <a-link href="#">使用条款</a-link>
        <a-link href="#">隐私政策</a-link>
        <a-link href="#">联系我们</a-link>
      </div>
    </a-layout-footer>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import {
  IconApps,
  IconBook,
  IconCommon,
  IconRobot,
  IconDashboard,
  IconUser,
  IconExport
} from '@arco-design/web-vue/es/icon'
import userService from '@/services/userService'

const router = useRouter()
const route = useRoute()
const isLoggedIn = ref(false)
const username = ref('')
const currentYear = new Date().getFullYear()

const activeKey = computed(() => {
  const path = route.path
  if (path.startsWith('/assessment')) return 'assessment'
  if (path.startsWith('/teaching')) return 'teaching'
  if (path.startsWith('/vr')) return 'vr'
  if (path.startsWith('/robot')) return 'robot'
  if (path.startsWith('/evaluation')) return 'evaluation'
  return ''
})

const userInitial = computed(() => {
  return username.value ? username.value.charAt(0).toUpperCase() : 'U'
})

onMounted(async () => {
  checkLoginStatus()
})

const checkLoginStatus = async () => {
  const token = localStorage.getItem('token') || sessionStorage.getItem('token')
  if (token) {
    isLoggedIn.value = true
    try {
      const userProfile = await userService.getUserProfile()
      username.value = userProfile.username
    } catch (error) {
      console.error('获取用户信息失败:', error)
      // 如果获取用户信息失败，可能是token过期，清除token
      localStorage.removeItem('token')
      sessionStorage.removeItem('token')
      isLoggedIn.value = false
    }
  } else {
    isLoggedIn.value = false
  }
}

const handleMenuClick = (key: string) => {
  navigateTo(`/${key}`)
}

const navigateTo = (path: string) => {
  router.push(path)
}

const handleLogout = async () => {
  try {
    await userService.logout()
    Message.success('退出登录成功')
    isLoggedIn.value = false
    username.value = ''
    router.push('/login')
  } catch (error) {
    console.error('退出登录失败:', error)
    Message.error('退出登录失败')
  }
}
</script>

<style scoped>
.layout {
  min-height: 100vh;
}

.header {
  display: flex;
  align-items: center;
  padding: 0 20px;
  background-color: var(--color-bg-2);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.logo {
  display: flex;
  align-items: center;
  margin-right: 40px;
}

.logo img {
  height: 32px;
  margin-right: 10px;
}

.logo h1 {
  margin: 0;
  color: var(--color-text-1);
  font-size: 18px;
  font-weight: 600;
}

.menu-wrapper {
  flex: 1;
}

.user-actions {
  display: flex;
  align-items: center;
}

.content {
  padding: 24px;
  background-color: var(--color-bg-1);
}

.footer {
  text-align: center;
  padding: 24px 50px;
  background-color: var(--color-bg-2);
  color: var(--color-text-2);
}

.footer-links {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-top: 10px;
}

@media (max-width: 768px) {
  .header {
    flex-direction: column;
    padding: 10px;
    height: auto;
  }

  .logo {
    margin-right: 0;
    margin-bottom: 10px;
  }

  .menu-wrapper {
    width: 100%;
    margin-bottom: 10px;
  }

  .user-actions {
    width: 100%;
    justify-content: center;
    margin-top: 10px;
  }

  .content {
    padding: 16px;
  }

  .footer {
    padding: 16px;
  }

  .footer-links {
    flex-direction: column;
    gap: 10px;
  }
}
</style>
