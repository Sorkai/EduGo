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
      },
      {
        path: 'user-management',
        component: UserManagement,
        meta: {
          requiresAuth: true,
          roles: ['super_admin', 'admin', 'teacher']
        }
      }
    ]
  }
]
```

## 用户类型和角色

EduGo系统支持以下用户角色：

- **超级管理员 (super_admin)**: 拥有最高权限，可以管理所有用户和功能。
- **管理员 (admin)**: 教学领导，可以管理教师和学生。
- **教师 (teacher)**: 可以管理学生。
- **学生 (student)**: 可以管理与家长的关系。
- **家长 (parent)**: 可以查看关联学生的信息。

用户注册时可以选择角色（教师、学生、家长），超级管理员和管理员角色需要由超级管理员指定。

## 用户管理页面

用户管理页面（UserManagement.vue）提供了用户列表和用户关系管理功能，根据当前用户的角色显示不同的内容和权限。

### 页面结构

用户管理页面分为两个主要标签页：

1. **用户列表**：显示用户列表，支持按角色筛选和关键词搜索。
2. **用户关系**：管理用户之间的关系，如管理员-教师、教师-学生、学生-家长关系。

### 用户列表功能

- 显示用户基本信息：ID、用户名、邮箱、姓名、角色、状态、注册时间
- 按角色筛选用户
- 搜索用户（按用户名或邮箱）
- 查看用户详情
- 修改用户角色（仅超级管理员可用）
- 启用/禁用用户（超级管理员和管理员可用）

### 用户关系管理功能

根据当前用户角色显示不同的关系管理界面：

- **超级管理员和管理员**：
  - 查看和管理教师列表
  - 添加教师关系（指定部门和职位）

- **超级管理员、管理员和教师**：
  - 查看和管理学生列表
  - 添加学生关系（指定课程和学期）

- **学生**：
  - 查看和管理家长列表
  - 添加家长关系（指定关系类型）

### 权限控制

页面根据当前用户角色控制可见内容和操作权限：

- 超级管理员可以查看所有用户，修改任何用户的角色和状态（除了其他超级管理员）
- 管理员可以查看教师和学生，修改其状态（但不能修改角色）
- 教师只能查看学生
- 学生只能管理与家长的关系

### 组件交互

用户管理页面包含多个模态框组件：

- 用户详情模态框
- 修改角色模态框
- 添加教师关系模态框
- 添加学生关系模态框
- 添加家长关系模态框

这些模态框通过v-model绑定控制显示状态，并通过事件通信进行数据交互。

### 数据加载和刷新

页面在加载时会根据当前用户角色获取相应的数据：

```typescript
// 页面加载时初始化
onMounted(async () => {
  await getCurrentUser()
  await fetchUsers()
  await fetchUserRelations()
})
```

当执行添加、修改等操作后，会自动刷新相关数据列表。

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

## 前端部署

详细部署步骤请参考[部署指南](../../all/部署指南.md)，包含以下内容：

1. 服务器要求
2. 环境准备（Node.js、npm安装）
3. 前端项目构建
4. 静态资源部署
5. 容器化部署选项
6. 常见问题处理

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
