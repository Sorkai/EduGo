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

## 认证机制

### JWT认证

EduGo使用JWT（JSON Web Token）进行用户认证。当用户登录成功后，服务器会返回一个JWT令牌，前端应用需要将此令牌存储在本地（localStorage或sessionStorage），并在后续请求中通过Authorization头部发送给服务器。

#### JWT令牌格式

JWT令牌包含以下信息：

- `user_id`: 用户ID
- `username`: 用户名
- `role`: 用户角色（默认为"user"）
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
