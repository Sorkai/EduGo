# EduGo 后端API文档

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
