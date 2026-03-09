# JWT + Redis 登录方案

## 方案概述

本系统采用 **JWT + Redis** 的双重认证方案，结合了两种技术的优势：

### JWT 优势
- ✅ 无状态认证，减轻服务器压力
- ✅ 自包含 token，无需查询数据库
- ✅ 支持分布式部署
- ✅ 自动过期机制

### Redis 优势
- ✅ 支持强制下线（删除 Redis 中的 token）
- ✅ 支持多点登录控制
- ✅ 实时权限控制
- ✅ 灵活的会话管理

## 实现原理

### 1. 登录流程

```
用户登录
  ↓
验证用户名密码
  ↓
生成 JWT Token
  ↓
存储 Token 到 Redis: key=user:token:{userId}, value=token
  ↓
返回 Token 给前端
```

**关键代码** (`api/user.go`):
```go
func (api *UserAPI) Login(c *gin.Context) {
    // 1. 验证用户名密码
    token, user, err := api.userService.Login(loginReq)
    
    // 2. 将 token 存入 Redis
    tokenKey := fmt.Sprintf("user:token:%d", user.ID)
    api.redisService.client.Set(ctx, tokenKey, token, expiresTime)
    
    // 3. 返回 token
    utils.SuccessWithData(c, gin.H{
        "token": token,
        "expires_in": 86400, // 过期时间（秒）
    })
}
```

### 2. 认证流程

```
请求携带 Token
  ↓
JWT 中间件验证签名和过期时间
  ↓
检查 Token 黑名单
  ↓
从 Redis 获取存储的 Token，验证一致性
  ↓
验证通过，允许访问
```

**关键代码** (`middleware/jwt.go`):
```go
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 解析 JWT Token
        claims, err := j.ParseToken(token)
        
        // 2. 检查黑名单
        isBlacklisted, _ := redisService.IsTokenBlacklisted(token)
        if isBlacklisted {
            Unauthorized("Token 已失效")
            return
        }
        
        // 3. 验证 Redis 中的 Token
        storedToken, _ := redisService.client.Get(ctx, fmt.Sprintf("user:token:%d", claims.ID)).Result()
        if storedToken != token {
            Unauthorized("Token 已过期，请重新登录")
            return
        }
        
        // 4. 验证通过
        c.Set("userId", claims.ID)
        c.Next()
    }
}
```

### 3. 登出流程

```
用户登出
  ↓
获取当前 Token
  ↓
将 Token 加入黑名单（Redis）
  ↓
可选：删除 Redis 中存储的 Token
  ↓
返回成功
```

**关键代码** (`api/auth.go`):
```go
func (api *AuthAPI) Logout(c *gin.Context) {
    // 1. 获取 token
    token := c.GetHeader("x-token")
    
    // 2. 计算剩余有效期
    expiresAt := claims.ExpiresAt.Unix() - time.Now().Unix()
    
    // 3. 加入黑名单
    api.redisService.AddTokenToBlacklist(token, time.Duration(expiresAt) * time.Second)
    
    // 4. 可选：删除存储的 token
    api.redisService.client.Del(ctx, fmt.Sprintf("user:token:%d", userID))
}
```

## Redis 数据结构

### 存储的 Key

| Key 格式 | 说明 | 过期时间 |
|---------|------|---------|
| `user:token:{userId}` | 存储用户当前有效的 token | 与 token 过期时间一致 |
| `blacklist:token:{token}` | 存储已登出的 token（黑名单） | token 剩余有效期 |
| `user:disabled:{userId}` | 用户禁用标记 | 永久或指定时长 |

### 示例

```bash
# 用户登录后
127.0.0.1:6379> GET user:token:1
"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# 用户登出后
127.0.0.1:6379> GET blacklist:token:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
"1"

# 检查用户是否被禁用
127.0.0.1:6379> GET user:disabled:1
(nil)  # 未禁用
```

## 功能特性

### 1. 强制下线

**场景**：管理员禁用用户、用户违规封禁

**实现**：
```go
// 方法 1：删除 Redis 中的 token
redisService.client.Del(ctx, fmt.Sprintf("user:token:%d", userID))

// 方法 2：设置用户禁用标记
redisService.SetUserDisabled(userID, 0) // 0 表示永久
```

**效果**：
- 用户再次请求时，JWT 中间件验证失败
- 返回 401，提示重新登录

### 2. 多点登录控制

**当前策略**：允许一个用户多处登录

**修改为单点登录**：
```go
func (api *UserAPI) Login(c *gin.Context) {
    tokenKey := fmt.Sprintf("user:token:%d", user.ID)
    
    // 登录前检查是否已有 token
    oldToken, _ := api.redisService.client.Get(ctx, tokenKey).Result()
    if oldToken != "" {
        // 可选：将旧 token 加入黑名单
        api.redisService.AddTokenToBlacklist(oldToken, expiresTime)
    }
    
    // 存储新 token（覆盖旧 token）
    api.redisService.client.Set(ctx, tokenKey, newToken, expiresTime)
}
```

### 3. Token 自动续期

**当前实现**：
- JWT 中间件检测到 token 即将过期（缓冲期内）
- 自动生成新 token 并通过响应头返回
- 前端更新存储的 token

**响应头**：
```
new-token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
new-expires-at: 2026-03-09T12:00:00+08:00
```

**前端处理**：
```javascript
// utils/request.js
request.interceptors.response.use(response => {
    const newToken = response.headers['new-token']
    if (newToken) {
        localStorage.setItem('token', newToken)
    }
    return response
})
```

### 4. 禁用用户

**场景**：用户违规、离职员工

**实现**：
```go
// 设置用户禁用标记
redisService.SetUserDisabled(userID, 0) // 永久禁用

// 或者设置禁用时长
redisService.SetUserDisabled(userID, 24*time.Hour) // 禁用 24 小时
```

**效果**：
- JWT 中间件检查到禁用标记
- 拒绝所有请求，返回 401
- 即使用户有有效 token 也无法访问

## 安全性分析

### 优势

1. **双重验证**：
   - JWT 验证签名和过期时间
   - Redis 验证 token 有效性

2. **实时控制**：
   - 可以随时使 token 失效
   - 支持强制下线

3. **灵活权限**：
   - 支持动态权限调整
   - 支持临时禁用

### 注意事项

1. **Redis 故障**：
   - Redis 不可用时，JWT 仍然可以验证
   - 但无法进行强制下线等高级功能
   - 建议配置 Redis 哨兵或集群

2. **Token 存储**：
   - 前端应安全存储 token（localStorage/ sessionStorage）
   - 建议使用 HTTPS 传输
   - 设置合理的过期时间

3. **性能考虑**：
   - 每次请求都需要查询 Redis
   - 增加了一次网络开销
   - 但 Redis 性能很高，影响可忽略

## 使用示例

### 1. 登录

```bash
curl -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'

# 响应
{
  "code": 200,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 86400
  }
}
```

### 2. 访问受保护接口

```bash
curl -X GET http://localhost/api/user/info \
  -H "x-token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 3. 登出

```bash
curl -X POST http://localhost/api/auth/logout \
  -H "x-token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 4. 强制用户下线（管理员）

```bash
# 删除 Redis 中的 token
docker compose exec redis redis-cli DEL user:token:1

# 或设置用户禁用
docker compose exec redis redis-cli SET user:disabled:1 "1"
```

## 监控和管理

### 查看在线用户

```bash
# 查看所有 user:token:* key
docker compose exec redis redis-cli KEYS "user:token:*"

# 查看特定用户的 token
docker compose exec redis redis-cli GET "user:token:1"
```

### 清理过期数据

Redis 会自动清理过期数据，无需手动操作。

### 查看黑名单

```bash
# 查看黑名单中的 token 数量
docker compose exec redis redis-cli KEYS "blacklist:token:*" | wc -l
```

## 配置说明

### JWT 配置 (`utils/jwt.go`)

```go
type JWT struct {
    SigningKey  []byte        // 签名密钥
    ExpiresTime time.Duration // token 有效期（默认 24 小时）
    BufferTime  time.Duration // 缓冲时间（默认 1 小时）
}
```

### Redis 配置 (`config.yaml`)

```yaml
redis:
  addr: "redis:6379"
  password: ""
  db: 0
```

## 故障排查

### 问题 1：登录后立即提示 token 失效

**原因**：Redis 中的 token 未正确存储

**解决**：
1. 检查 Redis 是否正常运行
2. 查看应用日志：`docker compose logs app | grep "存储 token"`
3. 验证 Redis 连接配置

### 问题 2：无法强制用户下线

**原因**：JWT 中间件未正确验证 Redis 中的 token

**解决**：
1. 检查 JWT 中间件代码
2. 确认 Redis key 格式正确
3. 查看 Redis 中是否有存储 token

### 问题 3：Token 续期不生效

**原因**：前端未正确处理响应头

**解决**：
1. 检查前端 axios 拦截器
2. 确认响应头名称正确
3. 查看浏览器 Network 面板

## 总结

✅ **已实现功能**：
- JWT + Redis 双重认证
- 强制下线
- Token 自动续期
- 用户禁用
- Token 黑名单

✅ **优势**：
- 安全性高（双重验证）
- 灵活性强（实时控制）
- 性能好（Redis 缓存）
- 易扩展（支持分布式）

✅ **适用场景**：
- 需要强制下线的系统
- 对安全性要求较高的系统
- 多端登录控制的系统
- 实时权限管理的系统
