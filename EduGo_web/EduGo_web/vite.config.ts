import { fileURLToPath, URL } from 'node:url'

import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  // 加载环境变量
  const env = loadEnv(mode, process.cwd(), '')
  
  // 获取API基础URL，用于代理配置
  const apiBaseUrl = env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
  
  // 从URL中提取主机和路径
  let target = apiBaseUrl
  let apiPath = '/api/v1'
  
  try {
    const url = new URL(apiBaseUrl)
    target = `${url.protocol}//${url.host}`
    apiPath = url.pathname
  } catch (e) {
    console.error('Invalid API base URL format:', apiBaseUrl)
  }
  
  return {
    plugins: [
      vue(),
      vueDevTools(),
    ],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url))
      },
    },
    server: {
      proxy: {
        // 将/api/v1开头的请求代理到目标服务器
        '/api/v1': {
          target,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/api\/v1/, apiPath)
        }
      }
    }
  }
})
