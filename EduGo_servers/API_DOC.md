# EduGo 后端API文档

## 用户类型和角色

EduGo系统支持以下用户角色：

- **超级管理员 (super_admin)**: 拥有最高权限，可以管理所有用户和功能。系统中首个注册的用户默认为超级管理员。
- **管理员 (admin)**: 教学领导，可以管理教师和学生。
- **教师 (teacher)**: 可以管理学生。
- **学生 (student)**: 可以管理与家长的关系。
- **家长 (parent)**: 可以查看关联学生的信息。

用户注册时可以选择角色（教师、学生、家长），超级管理员和管理员角色需要由超级管理员指定。

## 用户管理

### 用户注册
- **URL**: `/api/v1/register`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "username": "string",
    "password": "string",
    "email": "string",
    "firstName": "string",
    "lastName": "string"
  }
  ```
- **Response**:
  ```json
  {
    "message": "string",
    "user": {
      "id": "number",
      "username": "string",
      "email": "string",
      "firstName": "string",
      "lastName": "string"
    }
  }
  ```

### 用户登录
- **URL**: `/api/v1/login`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **Response**:
  ```json
  {
    "message": "string",
    "token": "string"
  }
  ```

### 更新用户信息
- **URL**: `/api/v1/user`
- **Method**: `PUT`
- **Request Body**:
  ```json
  {
    "email": "string",
    "firstName": "string",
    "lastName": "string"
  }
  ```
- **Response**:
  ```json
  {
    "message": "string",
    "user": {
      "id": "number",
      "username": "string",
      "email": "string",
      "firstName": "string",
      "lastName": "string"
    }
  }
  ```

### 重置密码
- **URL**: `/api/v1/user/password`
- **Method**: `PUT`
- **Request Body**:
  ```json
  {
    "oldPassword": "string",
    "newPassword": "string"
  }
  ```
- **Response**:
  ```json
  {
    "message": "string"
  }
  ```

### 用户注销
- **URL**: `/api/v1/logout`
- **Method**: `POST`
- **Response**:
  ```json
  {
    "message": "string"
  }
  ```

### 刷新Token
- **URL**: `/api/v1/refresh`
- **Method**: `POST`
- **Response**:
  ```json
  {
    "message": "string",
    "token": "string"
  }
  ```

### 获取所有用户（超级管理员权限）
- **URL**: `/api/v1/super-admin/users`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <jwt_token>`
- **Response**:
  ```json
  {
    "users": [
      {
        "id": "number",
        "username": "string",
        "email": "string",
        "firstName": "string",
        "lastName": "string",
        "role": "string",
        "status": "string",
        "createdAt": "string"
      }
    ]
  }
  ```

### 根据角色获取用户（管理员及以上权限）
- **URL**: `/api/v1/admin/users/role/:role`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <jwt_token>`
- **URL Parameters**: `role` - 用户角色（super_admin, admin, teacher, student, parent）
- **Response**:
  ```json
  {
    "users": [
      {
        "id": "number",
        "username": "string",
        "email": "string",
        "firstName": "string",
        "lastName": "string",
        "role": "string",
        "status": "string",
        "createdAt": "string"
      }
    ]
  }
  ```

### 根据ID获取用户（管理员及以上权限）
- **URL**: `/api/v1/admin/users/:id`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <jwt_token>`
- **URL Parameters**: `id` - 用户ID
- **Response**:
  ```json
  {
    "user": {
      "id": "number",
      "username": "string",
      "email": "string",
      "firstName": "string",
      "lastName": "string",
      "role": "string",
      "status": "string",
      "createdAt": "string"
    }
  }
  ```

### 更新用户角色（超级管理员权限）
- **URL**: `/api/v1/super-admin/users/:id/role`
- **Method**: `PUT`
- **Headers**: `Authorization: Bearer <jwt_token>`
- **URL Parameters**: `id` - 用户ID
- **Request Body**:
  ```json
  {
    "role": "string" // super_admin, admin, teacher, student, parent
  }
  ```
- **Response**:
  ```json
  {
    "message": "string",
    "user": {
      "id": "number",
      "role": "string"
    }
  }
  ```

### 更新用户状态（管理员及以上权限）
- **URL**: `/api/v1/admin/users/:id/status`
- **Method**: `PUT`
- **Headers**: `Authorization: Bearer <jwt_token>`
- **URL Parameters**: `id` - 用户ID
- **Request Body**:
  ```json
  {
    "status": "string" // active, inactive, blocked
  }
  ```
- **Response**:
  ```json
  {
    "message": "string",
    "user": {
      "id": "number",
      "status": "string"
    }
  }
  ```

## 用户关系管理

### 创建管理员-教师关系（管理员及以上权限）
- **URL**: `/api/v1/admin/relations/teacher`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer <jwt_token>`
- **Request Body**:
  ```json
  {
    "teacher_id": "number",
    "department": "string", // 可选
    "position": "string" // 可选
  }
  ```
- **Response**:
  ```json
  {
    "message": "string",
    "relation": {
      "id": "number",
      "admin_id": "number",
      "teacher_id": "number",
      "department": "string",
      "position": "string"
    }
  }
  ```

### 创建教师-学生关系（教师及以上权限）
- **URL**: `/api/v1/teacher/relations/student`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer <jwt_token>`
- **Request Body**:
  ```json
  {
    "student_id": "number",
    "course_id": "number", // 可选
    "course_name": "string", // 可选
    "semester": "string" // 可选
  }
  ```
- **Response**:
  ```json
  {
    "message": "string",
    "relation": {
      "id": "number",
      "teacher_id": "number",
      "student_id": "number",
      "course_id": "number",
      "course_name": "string",
      "semester": "string"
    }
  }
  ```

### 创建学生-家长关系（学生及以上权限）
- **URL**: `/api/v1/student/relations/parent`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer <jwt_token>`
- **Request Body**:
  ```json
  {
    "parent_id": "number",
    "relationship": "string" // 可选，例如：father, mother, guardian, other
  }
  ```
- **Response**:
  ```json
  {
    "message": "string",
    "relation": {
      "id": "number",
      "student_id": "number",
      "parent_id": "number",
      "relationship": "string"
    }
  }
  ```

### 获取管理员管理的教师列表（管理员及以上权限）
- **URL**: `/api/v1/admin/relations/teachers`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <jwt_token>`
- **Response**:
  ```json
  {
    "teachers": [
      {
        "id": "number",
        "username": "string",
        "email": "string",
        "firstName": "string",
        "lastName": "string",
        "role": "string",
        "status": "string"
      }
    ]
  }
  ```

### 获取教师教授的学生列表（教师及以上权限）
- **URL**: `/api/v1/teacher/relations/students`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <jwt_token>`
- **Response**:
  ```json
  {
    "students": [
      {
        "id": "number",
        "username": "string",
        "email": "string",
        "firstName": "string",
        "lastName": "string",
        "role": "string",
        "status": "string"
      }
    ]
  }
  ```

### 获取学生的家长列表（学生及以上权限）
- **URL**: `/api/v1/student/relations/parents`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <jwt_token>`
- **Response**:
  ```json
  {
    "parents": [
      {
        "id": "number",
        "username": "string",
        "email": "string",
        "firstName": "string",
        "lastName": "string",
        "role": "string",
        "status": "string"
      }
    ]
  }
  ```

## 用户管理页面API

### 获取用户列表（根据当前用户角色返回不同的用户列表）
- **URL**: `/api/v1/user-management/users`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <jwt_token>`
- **Response**:
  ```json
  {
    "users": [
      {
        "id": "number",
        "username": "string",
        "email": "string",
        "firstName": "string",
        "lastName": "string",
        "role": "string",
        "status": "string",
        "createdAt": "string"
      }
    ]
  }
  ```
  注意：返回的用户列表会根据当前用户的角色进行过滤：
  - 超级管理员：返回所有用户
  - 管理员：返回教师和学生
  - 教师：返回学生

## 认证机制

### JWT认证与用户角色

EduGo使用JWT（JSON Web Token）进行用户认证。当用户登录成功后，服务器会返回一个JWT令牌，前端应用需要将此令牌存储在本地（localStorage或sessionStorage），并在后续请求中通过Authorization头部发送给服务器。

#### JWT令牌格式

JWT令牌包含以下信息：

- `user_id`: 用户ID
- `username`: 用户名
- `role`: 用户角色（super_admin, admin, teacher, student, parent）
- `exp`: 令牌过期时间（24小时后）
- `iat`: 令牌签发时间

#### 请求认证

需要认证的API请求应在HTTP头部包含以下字段：

```
Authorization: Bearer <jwt_token>
```

#### 令牌刷新

当令牌即将过期时，前端应用可以调用刷新令牌API获取新的令牌，无需用户重新登录。

## 跨域资源共享 (CORS)

为支持前后端分离部署，API服务器配置了CORS中间件，允许来自不同域的前端应用访问API。

### CORS配置

服务器设置了以下CORS响应头：

- `Access-Control-Allow-Origin`: `*` (允许任何来源访问，生产环境中应限制为特定域名)
- `Access-Control-Allow-Credentials`: `true` (允许携带凭证)
- `Access-Control-Allow-Headers`: 允许的请求头包括 `Content-Type`, `Authorization` 等
- `Access-Control-Allow-Methods`: 允许的HTTP方法包括 `GET`, `POST`, `PUT`, `DELETE`, `OPTIONS`

### 前后端分离部署

在前后端分离部署环境中，前端应用需要配置正确的API基础URL，指向后端API服务器地址。详细配置方法请参考前端开发文档。

## 部署说明

详细部署步骤请参考[部署指南](../all/部署指南.md)，包含以下内容：

1. 服务器要求
2. 环境准备（MySQL、Go安装）
3. 后端服务部署
4. 监控与维护
5. 容器化部署选项
6. 常见问题处理

## 错误处理

API返回的错误格式统一为：

```json
{
  "error": "错误信息"
}
```

常见的HTTP状态码：

- `200 OK`: 请求成功
- `201 Created`: 资源创建成功
- `400 Bad Request`: 请求参数错误
- `401 Unauthorized`: 未认证或认证失败
- `403 Forbidden`: 权限不足
- `404 Not Found`: 资源不存在
- `500 Internal Server Error`: 服务器内部错误
