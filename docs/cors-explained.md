# CORS 工作原理（Docker Compose 环境）

## 📊 请求流程图

```
┌─────────────────────────────────────────────────────────────────┐
│                         客户端请求                               │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │   Nginx         │
                    │   (80/443)      │
                    │   反向代理       │
                    └─────────────────┘
                              │
                ┌─────────────┴─────────────┐
                │                           │
        ┌───────▼────────┐          ┌──────▼──────┐
        │  浏览器访问     │          │ Nginx 代理   │
        │  (有 Origin)    │          │ (无 Origin) │
        └────────────────┘          └─────────────┘
                │                           │
                ▼                           ▼
    ┌───────────────────────┐   ┌─────────────────────┐
    │   CORS 中间件检查      │   │  CORS 中间件自动放行  │
    │                       │   │                     │
    │ 1. 检查 Origin 头      │   │ if origin == "" {   │
    │ 2. 匹配白名单          │   │     c.Next()        │
    │ 3. 通过/拒绝           │   │     return          │
    └───────────────────────┘   │ }                 │
                │               └─────────────────────┘
                ▼                           │
        ┌───────┴───────┐                   │
        │               │                   │
   ┌────▼────┐   ┌─────▼─────┐              │
   │ 白名单内 │   │ 不在名单  │              │
   │  ✅ 放行 │   │  ❌ 拒绝  │              │
   └─────────┘   └───────────┘              │
        │                                   │
        └─────────────┬─────────────────────┘
                      ▼
            ┌─────────────────┐
            │     App         │
            │   (8888)        │
            │   业务逻辑       │
            └─────────────────┘
```

## 🔍 详细场景分析

### 场景 1：浏览器访问（需要 CORS）

```
浏览器 (http://localhost:3000)
   │
   │ Header: Origin: http://localhost:3000
   ▼
Nginx (反向代理，不修改 Origin 头)
   │
   │ Header: Origin: http://localhost:3000
   ▼
App (CORS 中间件检查)
   │
   ├─ 检查：http://localhost:3000 是否在白名单？
   ├─ 是 → 添加 CORS 响应头 → 放行 ✅
   └─ 否 → 返回 403 → 拒绝 ❌
```

**响应头示例**：
```http
Access-Control-Allow-Origin: http://localhost:3000
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Authorization
```

### 场景 2：Nginx 反向代理（无需 CORS）

```
Nginx (内部转发)
   │
   │ (没有 Origin 头，因为是服务器到服务器的请求)
   ▼
App (CORS 中间件)
   │
   ├─ 检查：Origin 头存在吗？
   └─ 不存在 → 直接放行 ✅
```

**为什么 Nginx 没有 Origin 头？**
- `Origin` 是浏览器添加的请求头
- Nginx 作为反向代理，不是浏览器
- Nginx 转发请求时不会自动添加 `Origin` 头

### 场景 3：curl/Postman 测试（无需 CORS）

```bash
# curl 请求（没有 Origin 头）
curl http://localhost:8888/api/auth/login

# App (CORS 中间件)
# ├─ 检查：Origin 头存在吗？
# └─ 不存在 → 直接放行 ✅
```

## 💡 代码实现

### CORS 中间件核心逻辑

```go
func Cors() gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.GetHeader("Origin")
        
        // 关键优化：没有 Origin 头，直接放行
        if origin == "" {
            c.Next()
            return
        }
        
        // 有 Origin 头，检查白名单
        isAllowed := false
        for _, allowed := range allowedOrigins {
            if origin == allowed {
                isAllowed = true
                break
            }
        }
        
        if isAllowed {
            c.Header("Access-Control-Allow-Origin", origin)
        } else {
            c.AbortWithStatus(http.StatusForbidden)
            return
        }
        
        c.Next()
    }
}
```

## 🎯 Docker Compose 网络架构

### 容器间通信

```yaml
services:
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    networks:
      - app-network
  
  app:
    build: .
    expose:
      - "8888"
    networks:
      - app-network

networks:
  app-network:
    driver: bridge
```

### 网络请求路径

```
┌──────────────┐
│   浏览器      │
│ localhost:80 │
└──────┬───────┘
       │ HTTP (带 Origin 头)
       ▼
┌──────────────┐
│   Nginx      │ (容器)
│ 反向代理      │
└──────┬───────┘
       │ 内部网络 (不带 Origin 头)
       │ http://app:8888
       ▼
┌──────────────┐
│    App       │ (容器)
│   :8888      │
└──────────────┘
```

## ✅ 配置建议

### 生产环境

```yaml
# docker-compose.yml
services:
  app:
    environment:
      - GIN_MODE=release
      - ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
```

**说明**：
- ✅ 只配置实际的域名
- ✅ 不需要包含 Nginx（因为 Nginx 请求没有 Origin 头）
- ✅ 不需要包含 `localhost`（生产环境不会有 localhost 访问）

### 开发环境

```yaml
# docker-compose.yml
services:
  app:
    environment:
      - GIN_MODE=debug
      - ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

**说明**：
- ✅ 配置前端开发服务器的端口
- ✅ Nginx 转发仍然不需要 CORS
- ✅ 只有浏览器直接访问后端端口才需要 CORS

## 🔧 常见问题

### Q1: 为什么 Nginx 反向代理不需要 CORS？

**答**：CORS 是浏览器的安全机制，服务器之间的 HTTP 请求（如 Nginx → App）不受 CORS 限制。

### Q2: 如果浏览器直接访问后端端口会怎样？

**答**：会被 CORS 拦截。所以生产环境必须通过 Nginx 统一入口。

### Q3: 如何测试 CORS 配置？

```bash
# 1. 模拟浏览器请求（带 Origin 头）
curl -H "Origin: http://localhost:3000" http://localhost:8888/api/auth/login

# 2. 模拟 Nginx 请求（不带 Origin 头）
curl http://localhost:8888/api/auth/login
```

### Q4: WebSocket 需要 CORS 吗？

**答**：WebSocket 使用 Upgrade 头，不走 CORS 机制，但需要 Nginx 正确配置 Upgrade 头转发。

## 📚 参考资料

- [MDN CORS](https://developer.mozilla.org/zh-CN/docs/Web/HTTP/CORS)
- [Nginx 反向代理](https://nginx.org/en/docs/http/ngx_http_proxy_module.html)
- [Docker 网络](https://docs.docker.com/network/)
