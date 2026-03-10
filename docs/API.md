# API Documentation

## 📋 Overview

This project provides a comprehensive set of user authentication and permission management APIs, including user management, role management, menu management, dictionary management, and file upload functionality.

### Basic Information

- **Base Path**: `/api`
- **Authentication**: JWT Token (except login and registration)
- **Data Format**: JSON
- **Character Encoding**: UTF-8

### Response Format

#### Success Response
```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

#### Error Response
```json
{
  "code": 400,
  "message": "Error message",
  "data": null
}
```

### Authentication

Except for login and registration endpoints, all other endpoints require a JWT Token in the request header:

```
Authorization: Bearer <token>
```

Or use custom header:
```
x-token: <token>
```

---

## 🔐 Authentication APIs (Auth)

### 1. User Login

**Endpoint**: `POST /api/auth/login`

**Permission**: Public (no authentication required)

**Request Parameters**:
```json
{
  "username": "string",
  "password": "string"
}
```

**Response Example**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**curl Example**:
```bash
curl -X POST http://localhost:8888/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'
```

---

### 2. User Registration

**Endpoint**: `POST /api/auth/register`

**Permission**: Public (no authentication required)

**Request Parameters**:
```json
{
  "username": "string",
  "password": "string",
  "nickname": "string",
  "email": "string",
  "phone": "string"
}
```

**Response Example**:
```json
{
  "code": 200,
  "message": "Registration successful",
  "data": null
}
```

**curl Example**:
```bash
curl -X POST http://localhost:8888/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456",
    "nickname": "Test User",
    "email": "test@example.com",
    "phone": "13800138000"
  }'
```

---

## 👤 User APIs (User)

### 1. Get Current User Info

**Endpoint**: `GET /api/user/info`

**Permission**: Authenticated users

**Response Example**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "admin",
    "nickname": "Administrator",
    "avatar": "/uploads/avatar.jpg",
    "email": "admin@example.com",
    "phone": "13800138000",
    "gender": 1,
    "status": 1,
    "roles": [
      {
        "id": 1,
        "name": "Super Administrator",
        "code": "admin"
      }
    ]
  }
}
```

**curl Example**:
```bash
curl -X GET http://localhost:8888/api/user/info \
  -H "Authorization: Bearer <token>"
```

---

### 2. Update User Info

**Endpoint**: `PUT /api/user`

**Permission**: Authenticated users

**Request Parameters**:
```json
{
  "id": 1,
  "nickname": "New Nickname",
  "avatar": "/uploads/new-avatar.jpg",
  "email": "new@example.com",
  "phone": "13900139000",
  "gender": 1,
  "status": 1
}
```

**Response Example**:
```json
{
  "code": 200,
  "message": "Update successful",
  "data": null
}
```

**curl Example**:
```bash
curl -X PUT http://localhost:8888/api/user \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "nickname": "New Nickname",
    "email": "new@example.com"
  }'
```

---

### 3. Change Password

**Endpoint**: `PUT /api/user/password`

**Permission**: Authenticated users

**Request Parameters**:
```json
{
  "old_password": "Old password",
  "new_password": "New password"
}
```

**Response Example**:
```json
{
  "code": 200,
  "message": "Password changed successfully",
  "data": null
}
```

**curl Example**:
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

### 4. Get User List

**Endpoint**: `GET /api/user/list`

**Permission**: Authenticated users

**Request Parameters**:
```json
{
  "page": 1,
  "page_size": 10,
  "username": "Optional, username fuzzy search",
  "status": 0
}
```

**Response Example**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "username": "admin",
        "nickname": "Administrator",
        "email": "admin@example.com",
        "status": 1
      }
    ],
    "total": 10
  }
}
```

**curl Example**:
```bash
curl -X GET "http://localhost:8888/api/user/list?page=1&page_size=10" \
  -H "Authorization: Bearer <token>"
```

---

### 5. Delete User

**Endpoint**: `DELETE /api/user/:id`

**Permission**: Administrator

**Path Parameters**:
- `id`: User ID

**Response Example**:
```json
{
  "code": 200,
  "message": "Delete successful",
  "data": null
}
```

**curl Example**:
```bash
curl -X DELETE http://localhost:8888/api/user/1 \
  -H "Authorization: Bearer <token>"
```

---

### 6. Assign Roles

**Endpoint**: `POST /api/user/assign-roles`

**Permission**: Administrator

**Request Parameters**:
```json
{
  "user_id": 1,
  "role_ids": [1, 2]
}
```

**Response Example**:
```json
{
  "code": 200,
  "message": "Roles assigned successfully",
  "data": null
}
```

**curl Example**:
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

### 7. Upload Avatar

**Endpoint**: `POST /api/user/avatar`

**Permission**: Authenticated users

**Request Parameters**: `multipart/form-data`
- `file`: Avatar image file

**Response Example**:
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

**curl Example**:
```bash
curl -X POST http://localhost:8888/api/user/avatar \
  -H "Authorization: Bearer <token>" \
  -F "file=@/path/to/avatar.jpg"
```

---

## 🎭 Role APIs (Role)

### 1. Create Role

**Endpoint**: `POST /api/role`

**Permission**: Administrator

**Request Parameters**:
```json
{
  "name": "Role Name",
  "code": "role_code",
  "desc": "Role Description",
  "status": 1,
  "menu_ids": [1, 2, 3]
}
```

**Response Example**:
```json
{
  "code": 200,
  "message": "Create successful",
  "data": null
}
```

**curl Example**:
```bash
curl -X POST http://localhost:8888/api/role \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Regular User",
    "code": "user",
    "desc": "Regular user role",
    "status": 1,
    "menu_ids": [1, 2]
  }'
```

---

### 2. Update Role

**Endpoint**: `PUT /api/role`

**Permission**: Administrator

**Request Parameters**:
```json
{
  "id": 1,
  "name": "New Role Name",
  "code": "new_role_code",
  "desc": "New Description",
  "status": 1,
  "menu_ids": [1, 2, 3]
}
```

**Response Example**:
```json
{
  "code": 200,
  "message": "Update successful",
  "data": null
}
```

**curl Example**:
```bash
curl -X PUT http://localhost:8888/api/role \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "name": "New Role Name",
    "menu_ids": [1, 2, 3]
  }'
```

---

### 3. Get Role Details

**Endpoint**: `GET /api/role/:id`

**Permission**: Authenticated users

**Path Parameters**:
- `id`: Role ID

**Response Example**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "Super Administrator",
    "code": "admin",
    "desc": "Super Administrator Role",
    "status": 1,
    "menus": [
      {
        "id": 1,
        "name": "User Management",
        "path": "/user"
      }
    ]
  }
}
```

**curl Example**:
```bash
curl -X GET http://localhost:8888/api/role/1 \
  -H "Authorization: Bearer <token>"
```

---

### 4. Get Role List

**Endpoint**: `GET /api/role/list`

**Permission**: Authenticated users

**Request Parameters**:
```json
{
  "page": 1,
  "page_size": 10,
  "name": "Optional, role name fuzzy search",
  "status": 0
}
```

**Response Example**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "Super Administrator",
        "code": "admin",
        "status": 1
      }
    ],
    "total": 5
  }
}
```

**curl Example**:
```bash
curl -X GET "http://localhost:8888/api/role/list?page=1&page_size=10" \
  -H "Authorization: Bearer <token>"
```

---

### 5. Delete Role

**Endpoint**: `DELETE /api/role/:id`

**Permission**: Administrator

**Path Parameters**:
- `id`: Role ID

**Response Example**:
```json
{
  "code": 200,
  "message": "Delete successful",
  "data": null
}
```

**curl Example**:
```bash
curl -X DELETE http://localhost:8888/api/role/1 \
  -H "Authorization: Bearer <token>"
```

---

### 6. Get All Roles

**Endpoint**: `GET /api/role/all`

**Permission**: Authenticated users

**Response Example**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "Super Administrator",
      "code": "admin"
    },
    {
      "id": 2,
      "name": "Regular User",
      "code": "user"
    }
  ]
}
```

**curl Example**:
```bash
curl -X GET http://localhost:8888/api/role/all \
  -H "Authorization: Bearer <token>"
```

---

## 📋 Menu APIs (Menu)

### 1. Create Menu

**Endpoint**: `POST /api/menu`

**Permission**: Administrator

**Request Parameters**:
```json
{
  "name": "Menu Name",
  "path": "/user",
  "component": "user/index",
  "icon": "user",
  "sort": 1,
  "parent_id": 0,
  "type": 1,
  "status": 1
}
```

**Response Example**:
```json
{
  "code": 200,
  "message": "Create successful",
  "data": null
}
```

**curl Example**:
```bash
curl -X POST http://localhost:8888/api/menu \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "User Management",
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

### 2. Update Menu

**Endpoint**: `PUT /api/menu`

**Permission**: Administrator

**Request Parameters**:
```json
{
  "id": 1,
  "name": "New Menu Name",
  "path": "/new-path",
  "component": "new/component",
  "icon": "new-icon",
  "sort": 2,
  "parent_id": 0,
  "type": 1,
  "status": 1
}
```

**Response Example**:
```json
{
  "code": 200,
  "message": "Update successful",
  "data": null
}
```

**curl Example**:
```bash
curl -X PUT http://localhost:8888/api/menu \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "name": "New Menu Name",
    "sort": 2
  }'
```

---

### 3. Get Menu Details

**Endpoint**: `GET /api/menu/:id`

**Permission**: Authenticated users

**Path Parameters**:
- `id`: Menu ID

**Response Example**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "User Management",
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

**curl Example**:
```bash
curl -X GET http://localhost:8888/api/menu/1 \
  -H "Authorization: Bearer <token>"
```

---

### 4. Get Menu List

**Endpoint**: `GET /api/menu/list`

**Permission**: Authenticated users

**Request Parameters**:
```json
{
  "page": 1,
  "page_size": 10
}
```

**Response Example**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "User Management",
        "path": "/user",
        "sort": 1
      }
    ],
    "total": 10
  }
}
```

**curl Example**:
```bash
curl -X GET "http://localhost:8888/api/menu/list?page=1&page_size=10" \
  -H "Authorization: Bearer <token>"
```

---

### 5. Get Menu Tree

**Endpoint**: `GET /api/menu/tree`

**Permission**: Authenticated users

**Response Example**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "System Management",
      "children": [
        {
          "id": 2,
          "name": "User Management"
        },
        {
          "id": 3,
          "name": "Role Management"
        }
      ]
    }
  ]
}
```

**curl Example**:
```bash
curl -X GET http://localhost:8888/api/menu/tree \
  -H "Authorization: Bearer <token>"
```

---

### 6. Delete Menu

**Endpoint**: `DELETE /api/menu/:id`

**Permission**: Administrator

**Path Parameters**:
- `id`: Menu ID

**Response Example**:
```json
{
  "code": 200,
  "message": "Delete successful",
  "data": null
}
```

**curl Example**:
```bash
curl -X DELETE http://localhost:8888/api/menu/1 \
  -H "Authorization: Bearer <token>"
```

---

### 7. Get Role Menus

**Endpoint**: `GET /api/menu/role/:role_id`

**Permission**: Authenticated users

**Path Parameters**:
- `role_id`: Role ID

**Response Example**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "User Management",
      "path": "/user"
    },
    {
      "id": 2,
      "name": "Role Management",
      "path": "/role"
    }
  ]
}
```

**curl Example**:
```bash
curl -X GET http://localhost:8888/api/menu/role/1 \
  -H "Authorization: Bearer <token>"
```

---

## 📖 Dictionary APIs (Dict)

### 1. Get Dictionary by Code (Public)

**Endpoint**: `GET /api/dict/code/:code`

**Permission**: Public (no authentication required)

**Path Parameters**:
- `code`: Dictionary code

**Response Example**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "Gender Dictionary",
    "code": "gender",
    "items": [
      {
        "label": "Male",
        "value": "1",
        "sort": 1
      },
      {
        "label": "Female",
        "value": "2",
        "sort": 2
      }
    ]
  }
}
```

**curl Example**:
```bash
curl -X GET http://localhost:8888/api/dict/code/gender
```

---

### 2. Get Dictionary Details

**Endpoint**: `GET /api/dict/:id`

**Permission**: Authenticated users

**Path Parameters**:
- `id`: Dictionary ID

**Response Example**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "Gender Dictionary",
    "code": "gender",
    "desc": "Gender Dictionary",
    "status": 1,
    "items": [
      {
        "label": "Male",
        "value": "1",
        "sort": 1
      }
    ]
  }
}
```

**curl Example**:
```bash
curl -X GET http://localhost:8888/api/dict/1 \
  -H "Authorization: Bearer <token>"
```

---

### 3. Get Dictionary List

**Endpoint**: `GET /api/dict/list`

**Permission**: Authenticated users

**Request Parameters**:
```json
{
  "page": 1,
  "page_size": 10,
  "name": "Optional, dictionary name fuzzy search",
  "status": 0
}
```

**Response Example**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "Gender Dictionary",
        "code": "gender",
        "status": 1
      }
    ],
    "total": 5
  }
}
```

**curl Example**:
```bash
curl -X GET "http://localhost:8888/api/dict/list?page=1&page_size=10" \
  -H "Authorization: Bearer <token>"
```

---

### 4. Create Dictionary

**Endpoint**: `POST /api/dict`

**Permission**: Administrator

**Request Parameters**:
```json
{
  "name": "Dictionary Name",
  "code": "dict_code",
  "desc": "Dictionary Description",
  "status": 1,
  "items": [
    {
      "label": "Option 1",
      "value": "1",
      "sort": 1,
      "status": 1
    }
  ]
}
```

**Response Example**:
```json
{
  "code": 200,
  "message": "Create successful",
  "data": null
}
```

**curl Example**:
```bash
curl -X POST http://localhost:8888/api/dict \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Gender Dictionary",
    "code": "gender",
    "desc": "Gender Dictionary",
    "status": 1,
    "items": [
      {
        "label": "Male",
        "value": "1",
        "sort": 1,
        "status": 1
      },
      {
        "label": "Female",
        "value": "2",
        "sort": 2,
        "status": 1
      }
    ]
  }'
```

---

### 5. Update Dictionary

**Endpoint**: `PUT /api/dict`

**Permission**: Administrator

**Request Parameters**:
```json
{
  "id": 1,
  "name": "New Dictionary Name",
  "code": "new_code",
  "desc": "New Description",
  "status": 1,
  "items": [
    {
      "label": "Option 1",
      "value": "1",
      "sort": 1
    }
  ]
}
```

**Response Example**:
```json
{
  "code": 200,
  "message": "Update successful",
  "data": null
}
```

**curl Example**:
```bash
curl -X PUT http://localhost:8888/api/dict \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "name": "New Dictionary Name",
    "items": [
      {
        "label": "Male",
        "value": "1",
        "sort": 1
      }
    ]
  }'
```

---

### 6. Delete Dictionary

**Endpoint**: `DELETE /api/dict/:id`

**Permission**: Administrator

**Path Parameters**:
- `id`: Dictionary ID

**Response Example**:
```json
{
  "code": 200,
  "message": "Delete successful",
  "data": null
}
```

**curl Example**:
```bash
curl -X DELETE http://localhost:8888/api/dict/1 \
  -H "Authorization: Bearer <token>"
```

---

## 📤 Upload APIs (Upload)

### 1. Upload File

**Endpoint**: `POST /api/upload`

**Permission**: Authenticated users

**Request Parameters**: `multipart/form-data`
- `file`: File to upload

**Response Example**:
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

**curl Example**:
```bash
curl -X POST http://localhost:8888/api/upload \
  -H "Authorization: Bearer <token>" \
  -F "file=@/path/to/file.jpg"
```

---

### 2. Delete File

**Endpoint**: `DELETE /api/upload/:file_name`

**Permission**: Authenticated users

**Path Parameters**:
- `file_name`: File name

**Response Example**:
```json
{
  "code": 200,
  "message": "Delete successful",
  "data": null
}
```

**curl Example**:
```bash
curl -X DELETE http://localhost:8888/api/upload/example-123.jpg \
  -H "Authorization: Bearer <token>"
```

---

## 🔒 Permission Description

### Permission Levels

| Permission Level | Description | API Examples |
|---------|------|---------|
| Public | No authentication required | `/api/auth/login`, `/api/auth/register`, `/api/dict/code/:code` |
| Authenticated User | Requires JWT Token | `/api/user/info`, `/api/user/list` |
| Administrator | Requires administrator role | `/api/role`, `/api/menu`, `/api/dict` (write operations) |

### Administrator API List

The following endpoints require administrator permission:
- `POST /api/role` - Create role
- `PUT /api/role` - Update role
- `DELETE /api/role/:id` - Delete role
- `POST /api/menu` - Create menu
- `PUT /api/menu` - Update menu
- `DELETE /api/menu/:id` - Delete menu
- `POST /api/dict` - Create dictionary
- `PUT /api/dict` - Update dictionary
- `DELETE /api/dict/:id` - Delete dictionary

---

## 📊 Error Codes

| Error Code | Description |
|--------|------|
| 200 | Success |
| 400 | Bad request parameters |
| 401 | Unauthorized or token expired |
| 403 | Insufficient permissions |
| 404 | Resource not found |
| 500 | Internal server error |

---

## 💡 Usage Examples

### Complete Workflow Example

```bash
# 1. Login to get token
TOKEN=$(curl -s -X POST http://localhost:8888/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}' \
  | jq -r '.data.token')

# 2. Use token to access protected endpoints
curl -X GET http://localhost:8888/api/user/info \
  -H "Authorization: Bearer $TOKEN"

# 3. Get user list
curl -X GET "http://localhost:8888/api/user/list?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN"

# 4. Get dictionary data (public endpoint, no token required)
curl -X GET http://localhost:8888/api/dict/code/gender
```

---

## 📝 Changelog

- **v1.0.0** - Initial release
  - User authentication and management
  - Role and permission management
  - Menu management
  - Dictionary management
  - File upload

---

## 📞 Technical Support

If you have any questions, please refer to the project documentation or contact the development team.
