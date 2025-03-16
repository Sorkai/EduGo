import { createRouter, createWebHistory } from 'vue-router'
import IntelligentAssessment from '@/views/IntelligentAssessment.vue'
import IntelligentTeaching from '@/views/IntelligentTeaching.vue'
import VirtualReality from '@/views/VirtualReality.vue'
import EducationalRobot from '@/views/EducationalRobot.vue'
import UserManagement from '@/views/UserManagement.vue'

// 路由守卫 - 检查用户角色权限
const checkRolePermission = (requiredRoles: string[]) => {
  return (to: any, from: any, next: any) => {
    const userRole = localStorage.getItem('userRole') || sessionStorage.getItem('userRole')
    
    if (!userRole) {
      // 未登录，重定向到登录页
      next({ path: '/login', query: { redirect: to.fullPath } })
    } else if (requiredRoles.includes(userRole)) {
      // 有权限，允许访问
      next()
    } else {
      // 无权限，重定向到首页
      next({ path: '/' })
    }
  }
}

const routes = [
  {
    path: '/',
    redirect: '/assessment'
  },
  {
    path: '/assessment',
    name: 'IntelligentAssessment',
    component: IntelligentAssessment
  },
  {
    path: '/teaching',
    name: 'IntelligentTeaching',
    component: IntelligentTeaching
  },
  {
    path: '/vr',
    name: 'VirtualReality',
    component: VirtualReality
  },
  {
    path: '/robot',
    name: 'EducationalRobot',
    component: EducationalRobot
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue')
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/Profile.vue')
  },
  {
    path: '/evaluation',
    name: 'IntelligentEvaluation',
    component: () => import('@/views/IntelligentEvaluation.vue')
  },
  {
    path: '/user-management',
    name: 'UserManagement',
    component: UserManagement,
    beforeEnter: checkRolePermission(['super_admin', 'admin', 'teacher']),
    meta: {
      requiresAuth: true,
      roles: ['super_admin', 'admin', 'teacher']
    }
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

export default router
