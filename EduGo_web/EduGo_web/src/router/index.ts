import { createRouter, createWebHistory } from 'vue-router'
import IntelligentAssessment from '@/views/IntelligentAssessment.vue'
import IntelligentTeaching from '@/views/IntelligentTeaching.vue'
import VirtualReality from '@/views/VirtualReality.vue'
import EducationalRobot from '@/views/EducationalRobot.vue'

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
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

export default router
