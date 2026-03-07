# API 接口文档

## 📋 概述

本项目提供了一套完整的用户认证和权限管理 API 接口，包括用户管理、角色管理、菜单管理、字典管理和文件上传等功能。

### 基础信息

- **基础路径**: `/api`
- **认证方式**: JWT Token（除登录注册外）
- **数据格式**: JSON
- **字符编码**: UTF-8

### 响应格式

#### 成功响应
```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

#### 错误响应
```json
{
  "code": 400,
  "message": "错误信息",
  "data": null
}
```

### 认证说明

除了登录和注册接口外，其他所有接口都需要在请求头中携带 JWT Token：

```
Authorization: Bearer <token>
```

或者使用自定义 header：
```
x-token: <token>
```

---

## 🔐 认证接口 (Auth)

### 1. 用户登录

**接口**: `POST /api/auth/login`

**权限**: 公开（无需认证）

**请求参数**:
```json
{
  "username": "string",
  "password": "string"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**curl 示例**:
```bash
curl -X POST http://localhost:8888/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'
```

---

### 2. 用户注册

**接口**: `POST /api/auth/register`

**权限**: 公开（无需认证）

**请求参数**:
```json
{
  "username": "string",
  "password": "string",
  "nickname": "string",
  "email": "string",
  "phone": "string"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "注册成功",
  "data": null
}
```

**curl 示例**:
```bash
curl -X POST http://localhost:8888/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456",
    "nickname": "测试用户",
    "email": "test@example.com",
    "phone": "13800138000"
  }'
```

---

## 👤 用户接口 (User)

### 1. 获取当前用户信息

**接口**: `GET /api/user/info`

**权限**: 已认证用户

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "admin",
    "nickname": "管理员",
    "avatar": "/uploads/avatar.jpg",
    "email": "admin@example.com",
    "phone": "13800138000",
    "gender": 1,
    "status": 1,
    "roles": [
      {
        "id": 1,
        "name": "超级管理员",
        "code": "admin"
      }
    ]
  }
}
```

**curl 示例**:
```bash
curl -X GET http://localhost:8888/api/user/info \
  -H "Authorization: Bearer <token>"
```

---

### 2. 更新用户信息

**接口**: `PUT /api/user`

**权限**: 已认证用户

**请求参数**:
```json
{
  "id": 1,
  "nickname": "新昵称",
  "avatar": "/uploads/new-avatar.jpg",
  "email": "new@example.com",
  "phone": "13900139000",
  "gender": 1,
  "status": 1
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "更新成功",
  "data": null
}
```

**curl 示例**:
```bash
curl -X PUT http://localhost:8888/api/user \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "nickname": "新昵称",
    "email": "new@example.com"
  }'
```

---

### 3. 修改密码

**接口**: `PUT /api/user/password`

**权限**: 已认证用户

**请求参数**:
```json
{
  "old_password": "旧密码",
  "new_password": "新密码"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "密码修改成功",
  "data": null
}
```

**curl 示例**:
```bash
curl -X PUT http://localhost:8888/api/user/password \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "123456",
    "new_password": "new123456"
  }'
```

---

### 4. 获取用户列表

**接口**: `GET /api/user/list`

**权限**: 已认证用户

**请求参数**:
```json
{
  "page": 1,
  "page_size": 10,
  "username": "可选，用户名模糊搜索",
  "status": 0
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "username": "admin",
        "nickname": "管理员",
        "email": "admin@example.com",
        "status": 1
      }
    ],
    "total": 10
  }
}
```

**curl 示例**:
```bash
curl -X GET "http://localhost:8888/api/user/list?page=1&page_size=10" \
  -H "Authorization: Bearer <token>"
```

---

### 5. 删除用户

**接口**: `DELETE /api/user/:id`

**权限**: 已认证用户

**路径参数**:
- `id`: 用户 ID

**响应示例**:
```json
{
  "code": 200,
  "message": "删除成功",
  "data": null
}
```

**curl 示例**:
```bash
curl -X DELETE http://localhost:8888/api/user/1 \
  -H "Authorization: Bearer <token>"
```

---

### 6. 分配角色

**接口**: `POST /api/user/assign-roles`

**权限**: 已认证用户

**请求参数**:
```json
{
  "user_id": 1,
  "role_ids": [1, 2]
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "分配角色成功",
  "data": null
}
```

**curl 示例**:
```bash
curl -X POST http://localhost:8888/api/user/assign-roles \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "role_ids": [1, 2]
  }'
```

---

### 7. 上传头像

**接口**: `POST /api/user/avatar`

**权限**: 已认证用户

**请求参数**: `multipart/form-data`
- `file`: 头像图片文件

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "avatar_url": "/uploads/avatar-123.jpg",
    "file_name": "avatar-123.jpg",
    "file_path": "/path/to/uploads/avatar-123.jpg"
  }
}
```

**curl 示例**:
```bash
curl -X POST http://localhost:8888/api/user/avatar \
  -H "Authorization: Bearer <token>" \
  -F "file=@/path/to/avatar.jpg"
```

---

## 🎭 角色接口 (Role)

### 1. 创建角色

**接口**: `POST /api/role`

**权限**: 管理员

**请求参数**:
```json
{
  "name": "角色名称",
  "code": "role_code",
  "desc": "角色描述",
  "status": 1,
  "menu_ids": [1, 2, 3]
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "创建成功",
  "data": null
}
```

**curl 示例**:
```bash
curl -X POST http://localhost:8888/api/role \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "普通用户",
    "code": "user",
    "desc": "普通用户角色",
    "status": 1,
    "menu_ids": [1, 2]
  }'
```

---

### 2. 更新角色

**接口**: `PUT /api/role`

**权限**: 管理员

**请求参数**:
```json
{
  "id": 1,
  "name": "新角色名称",
  "code": "new_role_code",
  "desc": "新描述",
  "status": 1,
  "menu_ids": [1, 2, 3]
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "更新成功",
  "data": null
}
```

**curl 示例**:
```bash
curl -X PUT http://localhost:8888/api/role \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "name": "新角色名称",
    "menu_ids": [1, 2, 3]
  }'
```

---

### 3. 获取角色详情

**接口**: `GET /api/role/:id`

**权限**: 已认证用户

**路径参数**:
- `id`: 角色 ID

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "超级管理员",
    "code": "admin",
    "desc": "超级管理员角色",
    "status": 1,
    "menus": [
      {
        "id": 1,
        "name": "用户管理",
        "path": "/user"
      }
    ]
  }
}
```

**curl 示例**:
```bash
curl -X GET http://localhost:8888/api/role/1 \
  -H "Authorization: Bearer <token>"
```

---

### 4. 获取角色列表

**接口**: `GET /api/role/list`

**权限**: 已认证用户

**请求参数**:
```json
{
  "page": 1,
  "page_size": 10,
  "name": "可选，角色名称模糊搜索",
  "status": 0
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "超级管理员",
        "code": "admin",
        "status": 1
      }
    ],
    "total": 5
  }
}
```

**curl 示例**:
```bash
curl -X GET "http://localhost:8888/api/role/list?page=1&page_size=10" \
  -H "Authorization: Bearer <token>"
```

---

### 5. 删除角色

**接口**: `DELETE /api/role/:id`

**权限**: 管理员

**路径参数**:
- `id`: 角色 ID

**响应示例**:
```json
{
  "code": 200,
  "message": "删除成功",
  "data": null
}
```

**curl 示例**:
```bash
curl -X DELETE http://localhost:8888/api/role/1 \
  -H "Authorization: Bearer <token>"
```

---

### 6. 获取所有角色

**接口**: `GET /api/role/all`

**权限**: 已认证用户

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "超级管理员",
      "code": "admin"
    },
    {
      "id": 2,
      "name": "普通用户",
      "code": "user"
    }
  ]
}
```

**curl 示例**:
```bash
curl -X GET http://localhost:8888/api/role/all \
  -H "Authorization: Bearer <token>"
```

---

## 📋 菜单接口 (Menu)

### 1. 创建菜单

**接口**: `POST /api/menu`

**权限**: 管理员

**请求参数**:
```json
{
  "name": "菜单名称",
  "path": "/user",
  "component": "user/index",
  "icon": "user",
  "sort": 1,
  "parent_id": 0,
  "type": 1,
  "status": 1
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "创建成功",
  "data": null
}
```

**curl 示例**:
```bash
curl -X POST http://localhost:8888/api/menu \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "用户管理",
    "path": "/user",
    "component": "user/index",
    "icon": "user",
    "sort": 1,
    "parent_id": 0,
    "type": 1,
    "status": 1
  }'
```

---

### 2. 更新菜单

**接口**: `PUT /api/menu`

**权限**: 管理员

**请求参数**:
```json
{
  "id": 1,
  "name": "新菜单名称",
  "path": "/new-path",
  "component": "new/component",
  "icon": "new-icon",
  "sort": 2,
  "parent_id": 0,
  "type": 1,
  "status": 1
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "更新成功",
  "data": null
}
```

**curl 示例**:
```bash
curl -X PUT http://localhost:8888/api/menu \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "name": "新菜单名称",
    "sort": 2
  }'
```

---

### 3. 获取菜单详情

**接口**: `GET /api/menu/:id`

**权限**: 已认证用户

**路径参数**:
- `id`: 菜单 ID

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "用户管理",
    "path": "/user",
    "component": "user/index",
    "icon": "user",
    "sort": 1,
    "parent_id": 0,
    "type": 1,
    "status": 1
  }
}
```

**curl 示例**:
```bash
curl -X GET http://localhost:8888/api/menu/1 \
  -H "Authorization: Bearer <token>"
```

---

### 4. 获取菜单列表

**接口**: `GET /api/menu/list`

**权限**: 已认证用户

**请求参数**:
```json
{
  "page": 1,
  "page_size": 10
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "用户管理",
        "path": "/user",
        "sort": 1
      }
    ],
    "total": 10
  }
}
```

**curl 示例**:
```bash
curl -X GET "http://localhost:8888/api/menu/list?page=1&page_size=10" \
  -H "Authorization: Bearer <token>"
```

---

### 5. 获取菜单树

**接口**: `GET /api/menu/tree`

**权限**: 已认证用户

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "系统管理",
      "children": [
        {
          "id": 2,
          "name": "用户管理"
        },
        {
          "id": 3,
          "name": "角色管理"
        }
      ]
    }
  ]
}
```

**curl 示例**:
```bash
curl -X GET http://localhost:8888/api/menu/tree \
  -H "Authorization: Bearer <token>"
```

---

### 6. 删除菜单

**接口**: `DELETE /api/menu/:id`

**权限**: 管理员

**路径参数**:
- `id`: 菜单 ID

**响应示例**:
```json
{
  "code": 200,
  "message": "删除成功",
  "data": null
}
```

**curl 示例**:
```bash
curl -X DELETE http://localhost:8888/api/menu/1 \
  -H "Authorization: Bearer <token>"
```

---

### 7. 获取角色的菜单

**接口**: `GET /api/menu/role/:role_id`

**权限**: 已认证用户

**路径参数**:
- `role_id`: 角色 ID

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "用户管理",
      "path": "/user"
    },
    {
      "id": 2,
      "name": "角色管理",
      "path": "/role"
    }
  ]
}
```

**curl 示例**:
```bash
curl -X GET http://localhost:8888/api/menu/role/1 \
  -H "Authorization: Bearer <token>"
```

---

## 📖 字典接口 (Dict)

### 1. 根据编码查询字典（公开）

**接口**: `GET /api/dict/code/:code`

**权限**: 公开（无需认证）

**路径参数**:
- `code`: 字典编码

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "性别字典",
    "code": "gender",
    "items": [
      {
        "label": "男",
        "value": "1",
        "sort": 1
      },
      {
        "label": "女",
        "value": "2",
        "sort": 2
      }
    ]
  }
}
```

**curl 示例**:
```bash
curl -X GET http://localhost:8888/api/dict/code/gender
```

---

### 2. 获取字典详情

**接口**: `GET /api/dict/:id`

**权限**: 已认证用户

**路径参数**:
- `id`: 字典 ID

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "性别字典",
    "code": "gender",
    "desc": "性别字典",
    "status": 1,
    "items": [
      {
        "label": "男",
        "value": "1",
        "sort": 1
      }
    ]
  }
}
```

**curl 示例**:
```bash
curl -X GET http://localhost:8888/api/dict/1 \
  -H "Authorization: Bearer <token>"
```

---

### 3. 获取字典列表

**接口**: `GET /api/dict/list`

**权限**: 已认证用户

**请求参数**:
```json
{
  "page": 1,
  "page_size": 10,
  "name": "可选，字典名称模糊搜索",
  "status": 0
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "性别字典",
        "code": "gender",
        "status": 1
      }
    ],
    "total": 5
  }
}
```

**curl 示例**:
```bash
curl -X GET "http://localhost:8888/api/dict/list?page=1&page_size=10" \
  -H "Authorization: Bearer <token>"
```

---

### 4. 创建字典

**接口**: `POST /api/dict`

**权限**: 管理员

**请求参数**:
```json
{
  "name": "字典名称",
  "code": "dict_code",
  "desc": "字典描述",
  "status": 1,
  "items": [
    {
      "label": "选项 1",
      "value": "1",
      "sort": 1,
      "status": 1
    }
  ]
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "创建成功",
  "data": null
}
```

**curl 示例**:
```bash
curl -X POST http://localhost:8888/api/dict \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "性别字典",
    "code": "gender",
    "desc": "性别字典",
    "status": 1,
    "items": [
      {
        "label": "男",
        "value": "1",
        "sort": 1,
        "status": 1
      },
      {
        "label": "女",
        "value": "2",
        "sort": 2,
        "status": 1
      }
    ]
  }'
```

---

### 5. 更新字典

**接口**: `PUT /api/dict`

**权限**: 管理员

**请求参数**:
```json
{
  "id": 1,
  "name": "新字典名称",
  "code": "new_code",
  "desc": "新描述",
  "status": 1,
  "items": [
    {
      "label": "选项 1",
      "value": "1",
      "sort": 1
    }
  ]
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "更新成功",
  "data": null
}
```

**curl 示例**:
```bash
curl -X PUT http://localhost:8888/api/dict \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "name": "新字典名称",
    "items": [
      {
        "label": "男",
        "value": "1",
        "sort": 1
      }
    ]
  }'
```

---

### 6. 删除字典

**接口**: `DELETE /api/dict/:id`

**权限**: 管理员

**路径参数**:
- `id`: 字典 ID

**响应示例**:
```json
{
  "code": 200,
  "message": "删除成功",
  "data": null
}
```

**curl 示例**:
```bash
curl -X DELETE http://localhost:8888/api/dict/1 \
  -H "Authorization: Bearer <token>"
```

---

## 📤 上传接口 (Upload)

### 1. 上传文件

**接口**: `POST /api/upload`

**权限**: 已认证用户

**请求参数**: `multipart/form-data`
- `file`: 要上传的文件

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "file_name": "example-123.jpg",
    "file_path": "/path/to/uploads/example-123.jpg",
    "url": "/uploads/example-123.jpg"
  }
}
```

**curl 示例**:
```bash
curl -X POST http://localhost:8888/api/upload \
  -H "Authorization: Bearer <token>" \
  -F "file=@/path/to/file.jpg"
```

---

### 2. 删除文件

**接口**: `DELETE /api/upload/:file_name`

**权限**: 已认证用户

**路径参数**:
- `file_name`: 文件名

**响应示例**:
```json
{
  "code": 200,
  "message": "删除成功",
  "data": null
}
```

**curl 示例**:
```bash
curl -X DELETE http://localhost:8888/api/upload/example-123.jpg \
  -H "Authorization: Bearer <token>"
```

---

## 🔒 权限说明

### 权限级别

| 权限级别 | 说明 | 接口示例 |
|---------|------|---------|
| 公开 | 无需认证即可访问 | `/api/auth/login`, `/api/auth/register`, `/api/dict/code/:code` |
| 已认证用户 | 需要 JWT Token | `/api/user/info`, `/api/user/list` |
| 管理员 | 需要管理员角色 | `/api/role`, `/api/menu`, `/api/dict`（写操作） |

### 管理员接口列表

以下接口需要管理员权限：
- `POST /api/role` - 创建角色
- `PUT /api/role` - 更新角色
- `DELETE /api/role/:id` - 删除角色
- `POST /api/menu` - 创建菜单
- `PUT /api/menu` - 更新菜单
- `DELETE /api/menu/:id` - 删除菜单
- `POST /api/dict` - 创建字典
- `PUT /api/dict` - 更新字典
- `DELETE /api/dict/:id` - 删除字典

---

## 📊 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证或 Token 过期 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 💡 使用示例

### 完整流程示例

```bash
# 1. 登录获取 Token
TOKEN=$(curl -s -X POST http://localhost:8888/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}' \
  | jq -r '.data.token')

# 2. 使用 Token 访问受保护的接口
curl -X GET http://localhost:8888/api/user/info \
  -H "Authorization: Bearer $TOKEN"

# 3. 获取用户列表
curl -X GET "http://localhost:8888/api/user/list?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN"

# 4. 获取字典数据（公开接口，无需 Token）
curl -X GET http://localhost:8888/api/dict/code/gender
```

---

## 📝 更新日志

- **v1.0.0** - 初始版本
  - 用户认证和管理
  - 角色和权限管理
  - 菜单管理
  - 字典管理
  - 文件上传

---

## 📞 技术支持

如有问题，请查看项目文档或联系开发团队。
