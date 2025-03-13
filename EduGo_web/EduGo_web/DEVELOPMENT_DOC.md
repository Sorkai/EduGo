# EduGo 前端开发文档

## 项目结构

```
src/
├── assets/              # 静态资源
├── components/          # 公共组件
├── router/              # 路由配置
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
