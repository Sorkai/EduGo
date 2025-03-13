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
