# EduGo 前端文档

EduGo是一个基于AI的数智化教育应用平台，前端使用Vue 3 + TypeScript + ArcoVue开发。

## 项目结构

```
src/
├── assets/              # 静态资源
├── components/          # 公共组件
│   ├── MainLayout.vue   # 主布局组件
│   └── ...
├── router/              # 路由配置
├── services/            # API服务
│   └── userService.ts   # 用户相关API服务
├── views/               # 页面视图
│   ├── Login.vue        # 登录页面
│   ├── Register.vue     # 注册页面
│   ├── Profile.vue      # 用户个人资料页面
│   ├── IntelligentAssessment.vue  # 智能测评页面
│   ├── IntelligentTeaching.vue    # 智能教学页面
│   ├── VirtualReality.vue         # 虚拟现实页面
│   ├── EducationalRobot.vue       # 教育机器人页面
│   └── IntelligentEvaluation.vue  # 智能评价页面
├── App.vue              # 根组件
└── main.ts              # 入口文件
```

## 功能模块

### 用户管理模块

- 用户注册：支持用户名、邮箱、密码注册
- 用户登录：支持用户名密码登录，JWT认证
- 个人资料：查看和修改个人资料
- 密码修改：支持修改密码
- 用户注销：支持用户注销
- 登录状态管理：使用JWT令牌管理登录状态，支持记住登录状态

### 智能测评模块

- 在线测评：支持在线答题测评
- AI分析：AI分析测评结果
- 个性化建议：根据测评结果提供个性化学习建议

### 智能教学模块

- 自适应学习：根据学生能力水平提供个性化学习资源
- 学习进度跟踪：跟踪学习进度和成果
- 智能推荐：推荐适合的学习内容

### 虚拟现实模块

- 沉浸式学习：提供沉浸式学习环境
- 3D交互：支持3D交互学习

### 教育机器人模块

- 智能对话：与教育机器人进行对话
- 学习辅助：机器人辅助学习

### 智能评价模块

- 全面评价：对学习过程和结果进行全面评价
- 数据分析：使用数据分析技术评价学习效果

## 技术栈

- Vue 3：前端框架
- TypeScript：类型系统
- Vite：构建工具
- Vue Router：路由管理
- ArcoVue：UI组件库
- Axios：HTTP客户端

## 开发环境设置

### 推荐的IDE

[VSCode](https://code.visualstudio.com/) + [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (禁用Vetur)

### 项目设置

```sh
# 安装依赖
npm install

# 开发环境运行
npm run dev

# 生产环境构建
npm run build
```

### 前后端分离部署

项目支持前后端分离部署，前端可以通过环境变量配置后端API地址。

#### 配置步骤

1. 在前端项目根目录创建环境变量文件：
   - 开发环境：`.env.local`
   - 生产环境：`.env.production`

2. 在环境变量文件中设置API地址：
   ```
   VITE_API_BASE_URL=http://your-api-server.com/api/v1
   ```

3. 构建前端项目：
   ```sh
   npm run build
   ```

4. 将构建后的文件部署到Web服务器（如Nginx、Apache等）

#### 后端部署

后端服务需要配置CORS（跨域资源共享）以允许前端访问：

```go
// 在main.go中添加CORS中间件
router.Use(func(c *gin.Context) {
    c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
    c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

    if c.Request.Method == "OPTIONS" {
        c.AbortWithStatus(204)
        return
    }

    c.Next()
})
```

## API接口

前端通过`services/userService.ts`与后端API交互，主要包括：

- 用户注册：`POST /api/v1/register`
- 用户登录：`POST /api/v1/login`
- 获取用户信息：`GET /api/v1/user`
- 更新用户信息：`PUT /api/v1/user`
- 修改密码：`PUT /api/v1/user/password`
- 用户注销：`POST /api/v1/logout`
- 刷新Token：`POST /api/v1/refresh`

详细API文档请参考后端API文档。

## 组件说明

### MainLayout

主布局组件，包含导航栏、内容区域和页脚。

- 响应式设计：适配不同屏幕尺寸
- 用户状态管理：显示登录/注册按钮或用户头像
- 导航菜单：提供主要功能模块的导航

### 用户相关组件

- Login：用户登录页面
- Register：用户注册页面
- Profile：用户个人资料页面，支持查看和修改个人信息、修改密码

### 功能模块组件

- IntelligentAssessment：智能测评页面
- IntelligentTeaching：智能教学页面
- VirtualReality：虚拟现实页面
- EducationalRobot：教育机器人页面
- IntelligentEvaluation：智能评价页面

## 样式和主题

项目使用ArcoVue的主题系统，支持亮色和暗色模式。主要样式变量定义在CSS变量中，便于统一管理和修改。

## 最近更新

### 2025/3/16 - 修复登录状态不能正常显示的问题

- 修复了JWT令牌中用户ID字段名称不一致导致的登录状态不能正常显示的问题
- 增强了JWT中间件，使其能够处理不同格式的令牌
- 改进了RefreshToken功能，确保在上下文中没有用户名时能够从数据库获取
- 更新了API文档和前端开发文档，添加了关于JWT认证和登录状态管理的详细说明
