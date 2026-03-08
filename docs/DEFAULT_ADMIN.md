# 默认管理员账户说明

## 🎉 功能概述

系统现在会在**首次启动时自动创建默认管理员账户**，无需手动操作。

---

## 👤 默认账户信息

### 管理员账户

| 字段 | 值 |
|------|-----|
| **用户名** | `admin` |
| **密码** | `admin123` |
| **昵称** | `系统管理员` |
| **角色** | `超级管理员` |
| **角色编码** | `admin` |
| **邮箱** | `admin@system.local` |

### 测试账户（开发环境）

如果你调用了 `CreateDefaultUser(db)`，还会创建测试账户：

| 字段 | 值 |
|------|-----|
| **用户名** | `test` |
| **密码** | `123456` |
| **昵称** | `测试用户` |
| **角色** | `普通用户` |
| **角色编码** | `user` |

---

## 🚀 启动流程

### 第一次启动

```
1. 启动应用
   ↓
2. 初始化数据库
   ↓
3. 执行数据库迁移（创建表）
   ↓
4. 调用 InitAdminUser(db)
   ├─ 检查管理员角色是否存在
   │  └─ 不存在 → 创建 `admin` 角色
   ├─ 检查管理员账户是否存在
   │  └─ 不存在 → 创建 `admin` 账户
   │      └─ 分配 `admin` 角色
   └─ 初始化字典数据
   ↓
5. 启动完成
```

**日志输出**：
```
[INFO] 管理员角色创建成功 id=1
[INFO] 默认管理员账户创建成功 id=1 username=admin role_id=1
[INFO] 服务器启动成功，监听地址：:8888
```

### 后续启动

```
1. 启动应用
   ↓
2. 初始化数据库
   ↓
3. 调用 InitAdminUser(db)
   ├─ 检查管理员角色 → ✅ 已存在，跳过
   ├─ 检查管理员账户 → ✅ 已存在，跳过
   └─ 检查角色分配 → ✅ 已分配，跳过
   ↓
4. 启动完成
```

**日志输出**：
```
[INFO] 管理员角色已存在 id=1
[INFO] 管理员账户已存在 id=1
[INFO] 服务器启动成功，监听地址：:8888
```

---

## 🔐 安全建议

### ✅ 生产环境必须修改密码

**第一次启动后**，立即修改管理员密码：

1. **使用修改密码接口**：
   ```bash
   curl -X PUT http://localhost:8888/api/user/password \
     -H "Authorization: Bearer <admin_token>" \
     -H "Content-Type: application/json" \
     -d '{
       "old_password": "admin123",
       "new_password": "your-strong-password"
     }'
   ```

2. **或者在数据库中直接修改**：
   ```sql
   UPDATE users SET password = '<bcrypt_hash>' WHERE username = 'admin';
   ```

### ✅ 删除测试账户（如果创建了）

```sql
DELETE FROM users WHERE username = 'test';
```

---

## 📝 代码实现

### 初始化函数

文件：[`initializer/admin_initializer.go`](../initializer/admin_initializer.go)

```go
func InitAdminUser(db *gorm.DB) {
    // 1. 检查并创建管理员角色
    // 2. 检查并创建管理员账户
    // 3. 分配管理员角色
    // 4. 初始化字典数据
}
```

### 调用位置

文件：[`cmd/main.go`](../cmd/main.go)

```go
func main() {
    // ... 初始化配置 ...
    
    // 初始化数据库
    db := initializer.InitDatabase(cfg)
    initializer.AutoMigrate(db)
    
    // 初始化默认管理员账户（先于字典初始化）
    initializer.InitAdminUser(db)
    
    // ... 其他初始化 ...
}
```

---

## 🔄 修改默认配置

### 修改默认用户名和密码

编辑 [`initializer/admin_initializer.go`](../initializer/admin_initializer.go)：

```go
// 第 54-55 行
adminUser = dbModel.User{
    Username: "admin",        // ← 修改为你的用户名
    Password: string(hashedPassword), // ← 修改密码
    // ...
}

// 第 47 行
hashedPassword, err := bcrypt.GenerateFromPassword(
    []byte("admin123"),  // ← 修改为你的密码
    bcrypt.DefaultCost
)
```

### 添加更多默认用户

在 [`initializer/admin_initializer.go`](../initializer/admin_initializer.go) 的 `CreateDefaultUser` 函数中添加：

```go
func CreateDefaultUsers(db *gorm.DB) {
    // 创建测试用户
    CreateDefaultUser(db)
    
    // 创建其他默认用户
    // ...
}
```

---

## 🛠️ 故障排查

### 问题 1：管理员账户未创建

**检查日志**：
```bash
docker-compose logs app | grep "管理员"
```

**可能原因**：
- 数据库迁移失败
- 密码哈希生成失败
- 角色分配失败

**解决方案**：
1. 查看完整日志
2. 检查数据库连接
3. 手动执行 SQL 创建

### 问题 2：无法登录

**可能原因**：
- 密码错误
- 角色未正确分配
- JWT 配置问题

**解决方案**：
```sql
-- 检查管理员账户
SELECT * FROM users WHERE username = 'admin';

-- 检查角色
SELECT * FROM roles WHERE code = 'admin';

-- 检查角色分配
SELECT * FROM user_roles WHERE user_id = 1;
```

### 问题 3：忘记密码

**重置密码**：
```go
// 在数据库中直接修改
UPDATE users 
SET password = '$2a$10$...'  -- 使用 bcrypt 生成的新哈希
WHERE username = 'admin';
```

或者使用密码重置工具生成新密码。

---

## 📚 相关文档

- [API 文档](./API.md) - 登录、修改密码等接口
- [安全修复报告](./SECURITY_FIX.md) - 权限控制说明
- [Docker Compose 配置](./docker-compose-setup.md) - 部署指南

---

## ✅ 总结

| 项目 | 状态 |
|------|------|
| 自动创建管理员角色 | ✅ |
| 自动创建管理员账户 | ✅ |
| 自动分配角色 | ✅ |
| 密码加密存储 | ✅ |
| 重复启动安全 | ✅ |
| 日志记录完整 | ✅ |

**系统现在可以立即使用了！** 🎉

默认管理员：`admin` / `admin123`
