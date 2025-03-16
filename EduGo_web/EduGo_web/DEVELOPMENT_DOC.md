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

## 用户认证与登录状态管理

### JWT认证

EduGo前端使用JWT（JSON Web Token）进行用户认证。当用户登录成功后，后端会返回一个JWT令牌，前端将此令牌存储在本地（localStorage或sessionStorage），并在后续请求中通过Authorization头部发送给服务器。

### 认证流程

1. 用户登录：
   - 用户输入用户名和密码
   - 前端发送登录请求到后端
   - 后端验证用户名和密码，生成JWT令牌并返回
   - 前端将令牌存储在localStorage（如果选择"记住我"）或sessionStorage中

2. 请求认证：
   - 前端在每个需要认证的请求中添加Authorization头部
   - 使用axios拦截器自动添加令牌到请求头中

3. 登录状态检查：
   - 在MainLayout组件中，使用checkLoginStatus方法检查用户是否已登录
   - 如果本地存储中有令牌，尝试获取用户信息
   - 如果获取成功，设置isLoggedIn为true
   - 如果获取失败（如令牌过期），清除令牌并设置isLoggedIn为false

4. 令牌刷新：
   - 当令牌即将过期时，可以调用刷新令牌API获取新的令牌
   - 无需用户重新登录

### 代码示例

#### 登录并存储令牌

```typescript
const handleLogin = async () => {
  try {
    const data = await userService.login(username, password);
    
    // 保存token
    if (rememberMe) {
      localStorage.setItem('token', data.token);
    } else {
      sessionStorage.setItem('token', data.token);
    }
    
    // 登录成功后的操作
    router.push('/');
  } catch (error) {
    // 处理错误
  }
};
```

#### 添加认证头部

```typescript
// 请求拦截器，添加token
apiClient.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token') || sessionStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  error => {
    return Promise.reject(error);
  }
);
```

#### 检查登录状态

```typescript
const checkLoginStatus = async () => {
  const token = localStorage.getItem('token') || sessionStorage.getItem('token');
  if (token) {
    isLoggedIn.value = true;
    try {
      const userProfile = await userService.getUserProfile();
      username.value = userProfile.username;
    } catch (error) {
      // 如果获取用户信息失败，可能是token过期，清除token
      localStorage.removeItem('token');
      sessionStorage.removeItem('token');
      isLoggedIn.value = false;
    }
  } else {
    isLoggedIn.value = false;
  }
};
```

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
