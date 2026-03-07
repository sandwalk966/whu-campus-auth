# 为什么不能用 `http://nginx` 作为 CORS 域名

## ❌ 错误的理解

```yaml
# 错误配置！不要这样做！
ALLOWED_ORIGINS=http://nginx
```

## ✅ 正确的理解

### 1. CORS 检查的是什么？

CORS 检查的是 **浏览器的 `Origin` 请求头**，而不是服务器的地址。

```
浏览器访问 https://yourdomain.com
   │
   │ Origin: https://yourdomain.com  ← 这个值来自浏览器
   ▼
Nginx 反向代理
   │
   │ Origin: https://yourdomain.com  ← Nginx 不会修改这个值
   ▼
App (CORS 中间件检查)
   │
   ├─ 检查：Origin 头是 "http://nginx" 吗？
   └─ 否 → 检查白名单 → 通过 ✅
```

### 2. Docker Compose 网络地址

```yaml
services:
  nginx:
    image: nginx:alpine
    # 容器内部地址：http://nginx:80
    # 外部访问地址：http://yourdomain.com
  
  app:
    build: .
    # 容器内部地址：http://app:8888
```

**关键点**：
- **容器间通信**：App 访问 Nginx 用 `http://nginx:80`（Docker 内部 DNS）
- **浏览器访问**：用户访问 `https://yourdomain.com`（公网域名）
- **CORS 检查**：检查的是浏览器的 `Origin: https://yourdomain.com`

### 3. Origin 头的来源

`Origin` 头是**浏览器自动添加的**，值是当前网页的地址：

| 浏览器地址栏 | Origin 头的值 |
|-------------|--------------|
| `https://yourdomain.com` | `https://yourdomain.com` |
| `http://localhost:3000` | `http://localhost:3000` |
| `http://127.0.0.1:5173` | `http://127.0.0.1:5173` |

**永远不会是**：
- ❌ `http://nginx`（这是容器内部地址）
- ❌ `http://app:8888`（这是后端内部地址）
- ❌ `http://mysql:3306`（这是数据库地址）

## 🎯 正确配置

### 生产环境（推荐：留空）

```yaml
ALLOWED_ORIGINS=
```

**效果**：
- ✅ Nginx 转发的请求（无 Origin 头）→ 自动放行
- ❌ 浏览器直接访问后端端口 → 被拒绝
- ✅ 最安全，强制所有流量通过 Nginx

### 生产环境（需要指定域名）

```yaml
ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
```

**效果**：
- ✅ Nginx 转发的请求 → 正常访问
- ✅ `https://yourdomain.com` 的浏览器请求 → 正常访问
- ❌ 其他域名的浏览器请求 → 被拒绝

### 开发环境

```yaml
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

**效果**：
- ✅ 本地前端开发服务器 → 正常访问
- ✅ Nginx 转发的请求 → 正常访问
- ❌ 其他域名的浏览器请求 → 被拒绝

## 📊 实际案例分析

### 案例 1：前端和后端都在 Docker

```yaml
services:
  frontend:
    image: node:alpine
    # 运行 React/Vue 开发服务器
    # 地址：http://frontend:3000
  
  nginx:
    image: nginx:alpine
    # 生产环境反向代理
    # 外部访问：https://yourdomain.com
    # 内部转发：http://app:8888
  
  app:
    build: .
    # 后端服务
    # 内部地址：http://app:8888
```

**CORS 配置**：
- 开发环境：`ALLOWED_ORIGINS=http://localhost:3000`
- 生产环境：`ALLOWED_ORIGINS=` （留空）或 `ALLOWED_ORIGINS=https://yourdomain.com`

**不能用**：`ALLOWED_ORIGINS=http://frontend` ❌

### 案例 2：前端在 Vercel，后端在 Docker

```yaml
# Vercel 部署前端
# 地址：https://your-app.vercel.app

# Docker 部署后端
services:
  nginx:
    # 外部访问：https://api.yourdomain.com
    # 内部转发：http://app:8888
  
  app:
    # 后端服务
```

**CORS 配置**：
```yaml
ALLOWED_ORIGINS=https://your-app.vercel.app
```

**不能用**：`ALLOWED_ORIGINS=http://nginx` ❌

## 🔍 技术原理

### 请求流程

```
┌─────────────────────────────────────────────────────────┐
│  浏览器访问 https://yourdomain.com                        │
│  Origin: https://yourdomain.com                          │
└─────────────────────────────────────────────────────────┘
                        │
                        │ HTTPS + Origin 头
                        ▼
┌─────────────────────────────────────────────────────────┐
│  Nginx (https://yourdomain.com)                          │
│  反向代理到 http://app:8888                              │
│  （不修改 Origin 头）                                     │
└─────────────────────────────────────────────────────────┘
                        │
                        │ HTTP + Origin: https://yourdomain.com
                        ▼
┌─────────────────────────────────────────────────────────┐
│  App (http://app:8888)                                   │
│  CORS 中间件检查 Origin 头                                │
│  if origin == "" → 放行（Nginx 请求）                     │
│  else → 检查白名单 → 放行（https://yourdomain.com）       │
└─────────────────────────────────────────────────────────┘
```

### 为什么 `http://nginx` 不对？

1. **浏览器看不到 `http://nginx`**
   - `nginx` 是 Docker 内部 DNS 名称
   - 浏览器只能访问公网 IP 或域名

2. **Origin 头不会是 `http://nginx`**
   - Origin 头 = 浏览器地址栏的协议 + 域名 + 端口
   - 用户不可能输入 `http://nginx`

3. **Nginx 不会修改 Origin 头**
   - Nginx 反向代理只是转发请求
   - 不会添加或修改 Origin 头

## ✅ 总结

**可以留空 `ALLOWED_ORIGINS`**：
- ✅ 推荐做法
- ✅ Nginx 转发的请求自动放行
- ✅ 更安全

**不能填写 `http://nginx`**：
- ❌ 浏览器的 Origin 头永远不会是 `http://nginx`
- ❌ 这是 Docker 内部地址，外部无法访问
- ❌ 填了也没用，匹配不到任何请求

**正确做法**：
```yaml
# 生产环境（推荐）
ALLOWED_ORIGINS=

# 或者指定实际域名
ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com

# 开发环境
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```
