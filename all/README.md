# EduGo 项目文档

EduGo是一个基于AI的数智化教育应用平台，前端使用Vue 3 + TypeScript + ArcoVue开发，后端使用Go语言和Gin框架开发。

## 用户类型和角色设计

EduGo系统支持以下用户角色：

- **超级管理员 (super_admin)**: 拥有最高权限，可以管理所有用户和功能。系统中首个注册的用户默认为超级管理员。
- **管理员 (admin)**: 教学领导，可以管理教师和学生。
- **教师 (teacher)**: 可以管理学生。
- **学生 (student)**: 可以管理与家长的关系。
- **家长 (parent)**: 可以查看关联学生的信息。

用户注册时可以选择角色（教师、学生、家长），超级管理员和管理员角色需要由超级管理员指定。

## 用户关系设计

EduGo系统实现了三种用户关系：

1. **管理员-教师关系 (admin_teacher)**
   - 管理员可以管理多个教师，教师可以归属于一个或多个管理员
   - 包含部门和职位等附加信息

2. **教师-学生关系 (teacher_student)**
   - 教师可以教授多个学生，学生可以有多个教师
   - 包含课程ID、课程名称和学期等附加信息

3. **学生-家长关系 (student_parent)**
   - 学生可以有多个家长，家长可以有多个孩子
   - 包含关系类型（父亲、母亲、监护人等）

这种关系设计使得系统能够灵活地管理不同用户之间的关联，并支持教育场景中常见的组织结构。

## 前端项目结构

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

### 2025/4/1 - 智能测评模块实现

- 实现了智能测评模块的后端API，支持教师创建、发布测评，学生参与测评
- 实现了测评题库管理功能，支持添加选择题
- 实现了学生在线答题功能，系统自动评分
- 实现了AI分析功能，根据学生答题情况提供个性化学习建议
- 实现了测评结果查看功能，包括得分、正确答案、解析等
- 前端实现了教师测评管理界面，包括创建测评、添加题目、发布测评、查看学生完成情况等功能
- 前端实现了学生测评参与界面，包括查看可参与测评、开始测评、答题、提交答案、查看结果等功能
- 更新了API文档，添加了智能测评模块的API接口说明

### 2025/3/16 - 用户类型扩充和用户管理功能

- 扩充用户类型为：超级管理员、管理员、教师、学生、家长
- 实现默认首个注册的用户为超级管理员，拥有最高权限
- 添加用户注册时可选教师、学生、家长角色
- 增加用户管理页面，支持超级管理员、管理员、教师使用，并区分不同用户权限
- 实现管理员(教学领导)、教师、学生、家长账号间关系管理
- 修复了JWT令牌中用户ID字段名称不一致导致的登录状态不能正常显示的问题
- 增强了JWT中间件，使其能够处理不同格式的令牌
- 改进了RefreshToken功能，确保在上下文中没有用户名时能够从数据库获取
- 更新了API文档和前端开发文档，添加了关于JWT认证和登录状态管理的详细说明
