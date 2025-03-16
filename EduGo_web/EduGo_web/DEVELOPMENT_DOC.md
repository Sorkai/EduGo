# EduGo 前端开发文档

## 项目结构

```
src/
├── assets/              # 静态资源
├── components/          # 公共组件
├── config/              # 配置文件
├── router/              # 路由配置
├── services/            # API服务
├── views/               # 页面视图
├── App.vue              # 根组件
└── main.ts              # 入口文件
```

## 主要组件说明

### MainLayout.vue
- 主布局组件，包含导航栏和内容区域
- 使用ArcoVue组件库构建
- 包含响应式设计

### Icon*.vue
- 图标组件集合
- 使用SVG图标
- 支持自定义颜色和大小

## 路由配置

```typescript
const routes = [
  {
    path: '/',
    component: MainLayout,
    children: [
      {
        path: 'assessment',
        component: IntelligentAssessment
      },
      {
        path: 'teaching',
        component: IntelligentTeaching
      },
      {
        path: 'vr',
        component: VirtualReality
      },
      {
        path: 'robot',
        component: EducationalRobot
      },
      {
        path: 'evaluation',
        component: IntelligentEvaluation
      }
    ]
  }
]
```

## 开发规范

1. 使用TypeScript进行开发
2. 遵循ArcoVue组件库设计规范
3. 使用ESLint进行代码检查
4. 使用Prettier进行代码格式化
5. 组件命名采用大驼峰式
6. 变量命名采用小驼峰式

## 开发环境

- Node.js v18+
- npm v9+
- Vite 4+
- TypeScript 5+

## 环境变量配置

项目使用环境变量来配置不同环境下的参数，主要用于前后端分离部署时配置API地址。

### 环境变量文件

- `.env.example`: 环境变量示例文件，包含所有可配置的环境变量
- `.env.local`: 本地开发环境配置，不会被提交到版本控制
- `.env.production`: 生产环境配置

### 可用的环境变量

- `VITE_API_BASE_URL`: API基础URL，例如 `http://api.edugo.com/api/v1`

### 使用方法

1. 复制 `.env.example` 为 `.env.local`（开发环境）或 `.env.production`（生产环境）
2. 修改环境变量值为实际环境的值
3. 重启开发服务器或重新构建项目

### 开发环境中的API代理

在开发环境中，前端开发服务器会自动将API请求代理到配置的后端服务器地址。这是通过Vite的代理功能实现的，配置在`vite.config.ts`文件中。

例如，如果您在`.env.local`中设置：
```
VITE_API_BASE_URL=http://localhost:10086/api/v1
```

当前端代码请求`/api/v1/login`时，Vite开发服务器会将请求代理到`http://localhost:10086/api/v1/login`。

这样，您可以在开发环境中使用相对路径（如`/api/v1/login`）进行API调用，而在生产环境中，这些请求会被发送到环境变量中配置的完整URL。

### 在代码中使用环境变量

环境变量通过 `import.meta.env` 对象访问：

```typescript
// 示例：访问API基础URL
const apiBaseUrl = import.meta.env.VITE_API_BASE_URL
```

项目中已经封装了配置模块，可以通过以下方式使用：

```typescript
import config from '@/config'

// 使用API基础URL
const apiBaseUrl = config.apiBaseUrl
```
