# Docker Compose 环境配置指南

## 📋 概述

本项目使用 Docker Compose 进行容器化部署，包含以下服务：
- **Nginx**：反向代理和静态文件服务器（端口 80/443）
- **App**：Go 后端应用（内部端口 8888）
- **MySQL**：数据库（端口 3306）
- **Redis**：缓存（端口 6379）

## 🔐 CORS 配置说明

### 为什么需要配置 CORS？

浏览器会检查跨域请求，但 **Nginx 反向代理不会触发 CORS**！

**请求流程**：
```
浏览器 → Nginx (80/443) → App (8888)
         ↑                    ↑
      有 Origin 头          无 Origin 头
      需要 CORS 检查       自动放行
```

- **浏览器直接访问**：带有 `Origin` 头，需要 CORS 白名单
- **Nginx 反向代理**：没有 `Origin` 头，自动放行

### 如何配置 CORS？

在 `docker-compose.yml` 中修改 `ALLOWED_ORIGINS` 环境变量：

```yaml
services:
  app:
    environment:
      - ALLOWED_ORIGINS=
```

**重要：`ALLOWED_ORIGINS` 可以为空！**

- **为空（推荐）**：只有通过 Nginx 反向代理的请求才能访问
  - ✅ Nginx → App：无 Origin 头，自动放行
  - ❌ 浏览器直接访问后端端口：有 Origin 头，被拒绝
  - ✅ 更安全，强制所有请求通过 Nginx

- **填写域名**：允许指定的浏览器访问
  - ✅ Nginx → App：正常访问
  - ✅ 指定域名的浏览器：正常访问
  - ❌ 其他域名的浏览器：被拒绝

### 不同环境的配置

#### 1️⃣ 生产环境（已部署到服务器）

```yaml
environment:
  - GIN_MODE=release
  - ALLOWED_ORIGINS=https://api.example.com,https://www.example.com
```

**重要**：
- ✅ 必须使用 HTTPS 协议
- ✅ 必须明确指定允许的域名
- ✅ 不要使用通配符 `*`（不安全）
- ✅ 如果前端和后端在同一域名下，只需配置该域名

#### 2️⃣ 开发环境（本地调试）

```yaml
environment:
  - GIN_MODE=debug
  - ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173,http://127.0.0.1:3000,http://127.0.0.1:5173
```

**说明**：
- 本地开发常用的前端端口：
  - React/Vue 开发服务器：`3000` 或 `5173`（Vite）
  - 同时支持 `localhost` 和 `127.0.0.1`

#### 3️⃣ 混合环境（本地开发 + 远程 API）

如果前端在本地，后端在远程服务器：

```yaml
environment:
  - ALLOWED_ORIGINS=http://localhost:3000,https://dev.example.com
```

## 🚀 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd whu-campus-auth
```

### 2. 配置环境变量

复制示例配置文件：

```bash
cp .env.example .env
```

编辑 `.env` 文件，修改配置：

```bash
# 修改为你的实际域名
ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com

# 修改 JWT 密钥（生产环境必须修改！）
JWT_SECRET=your-very-secret-key-here
```

### 3. 配置 Nginx

编辑 `nginx/nginx.conf`，修改域名：

```nginx
server_name yourdomain.com www.yourdomain.com;
```

### 4. 启动服务

```bash
docker-compose up -d
```

### 5. 查看日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f app
docker-compose logs -f nginx
```

### 6. 停止服务

```bash
docker-compose down
```

## 🔧 常见问题

### Q1: 前端访问后端报 CORS 错误

**错误信息**：
```
Access to XMLHttpRequest at 'https://api.example.com/api/auth/login' from origin 'http://localhost:3000' has been blocked by CORS policy
```

**解决方案**：
1. 检查 `docker-compose.yml` 中的 `ALLOWED_ORIGINS` 是否包含前端域名
2. 重启应用容器：`docker-compose restart app`
3. 清除浏览器缓存

### Q2: 生产环境是否需要包含 localhost？

**不需要**！生产环境的前端和后端都在服务器上，通过 Nginx 反向代理，不存在跨域问题。

**正确配置**：
```yaml
ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
```

### Q3: 如何支持多个子域名？

```yaml
ALLOWED_ORIGINS=https://example.com,https://www.example.com,https://api.example.com,https://admin.example.com
```

### Q4: 前端部署在 CDN，后端在服务器？

```yaml
ALLOWED_ORIGINS=https://cdn.example.com,https://api.example.com
```

## 📊 Docker Compose 网络架构

```
┌─────────────────────────────────────────────────────┐
│                   Docker Network                     │
│                                                      │
│  ┌─────────────┐     ┌─────────────┐                │
│  │   Nginx     │────▶│     App     │                │
│  │  (80/443)   │     │   (8888)    │                │
│  └─────┬───────┘     └─────┬───────┘                │
│        │                   │                        │
│        │           ┌───────┴───────┐                │
│        │           │               │                │
│        │      ┌────▼────┐   ┌─────▼────┐           │
│        │      │  MySQL  │   │  Redis   │           │
│        │      │ (3306)  │   │ (6379)   │           │
│        │      └─────────┘   └──────────┘           │
│        │                                           │
└────────┼───────────────────────────────────────────┘
         │
         ▼
    互联网/客户端
```

### 网络通信流程

1. **客户端请求** → Nginx（80/443 端口）
2. **Nginx 反向代理** → App 容器（8888 端口）
3. **App 访问数据库** → MySQL 容器（内部网络）
4. **App 访问缓存** → Redis 容器（内部网络）

**注意**：
- App 的 8888 端口只在 Docker 内部网络暴露
- 外部只能通过 Nginx 访问（80/443 端口）
- MySQL 和 Redis 不暴露给外部（生产环境建议）

## 🔒 安全建议

### 1. 生产环境不要暴露数据库端口

修改 `docker-compose.yml`：

```yaml
services:
  mysql:
    # 注释掉 ports 配置
    # ports:
    #   - "3306:3306"
    volumes:
      - mysql-data:/var/lib/mysql

  redis:
    # 注释掉 ports 配置
    # ports:
    #   - "6379:6379"
    volumes:
      - redis-data:/data
```

### 2. 使用环境变量文件

创建 `.env.production` 文件（不要提交到 Git）：

```bash
# .env.production
ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
JWT_SECRET=your-very-secret-key
DB_PASSWORD=strong-password
REDIS_PASSWORD=strong-password
```

在 `docker-compose.yml` 中引用：

```yaml
services:
  app:
    env_file:
      - .env.production
```

### 3. 启用 HTTPS

参考 `scripts/` 目录下的 HTTPS 配置脚本：

```bash
# 申请 Let's Encrypt 证书
sudo bash scripts/setup-https.sh yourdomain.com admin@yourdomain.com
```

## 📝 检查清单

部署前请确认：

- [ ] 修改了 `ALLOWED_ORIGINS` 为实际域名
- [ ] 修改了 `JWT_SECRET` 为随机密钥
- [ ] 修改了数据库密码
- [ ] 配置了 Nginx 的 `server_name`
- [ ] 申请了 HTTPS 证书
- [ ] 注释掉了数据库的外部端口映射
- [ ] 设置了合理的资源限制

## 🔗 相关文档

- [CORS 中间件代码](../middleware/cors.go)
- [HTTPS 配置脚本](../scripts/README.md)
- [Nginx 配置模板](../scripts/nginx-https-template.conf)
