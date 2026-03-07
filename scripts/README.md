# 运维脚本说明

本目录包含项目运维相关的脚本文件。

## 📁 脚本文件

### 1. **configure-cors.sh** - CORS 配置脚本

快速配置 CORS（跨域资源共享）环境变量。

#### 使用方法

```bash
# 生产环境配置（指定域名）
./configure-cors.sh production example.com www.example.com

# 开发环境配置（使用 localhost）
./configure-cors.sh development
```

#### 功能

- ✅ 自动更新 `docker-compose.yml` 中的 `ALLOWED_ORIGINS`
- ✅ 如果存在 `.env` 文件，同步更新
- ✅ 支持多个域名配置

---

### 2. **healthcheck-app.sh** - 应用健康检查

检查 Go 应用的健康状态。

#### 使用方法

```bash
./healthcheck-app.sh
```

#### 功能

- ✅ 检查应用端口（8888）是否响应
- ✅ 可用于手动健康检查或集成到监控系统

---

### 3. **letsencrypt.sh** - Let's Encrypt 证书管理 ⭐

自动申请和续期 Let's Encrypt HTTPS 证书。

#### 使用方法

```bash
# 申请新证书
./letsencrypt.sh yourdomain.com your@email.com apply

# 续期证书
./letsencrypt.sh yourdomain.com your@email.com renew

# 生成 Nginx 配置
./letsencrypt.sh yourdomain.com your@email.com config
```

#### 功能

- ✅ 自动申请 Let's Encrypt 证书
- ✅ 自动续期证书
- ✅ 自动生成 Nginx 配置
- ✅ 按需运行，不占用资源

详细说明：[LETS-ENCRYPT.md](../LETS-ENCRYPT.md)

---

### 4. **generate-ssl-cert.sh** - 自签名证书生成

生成自签名 SSL 证书（开发/测试环境使用）。

#### 使用方法

```bash
./generate-ssl-cert.sh
```

#### 功能

- ✅ 生成自签名证书
- ✅ 适用于开发/测试环境
- ✅ 证书有效期 365 天

---

### 5. **monitor-logs.sh** - 日志监控

监控 Docker 容器日志，检测错误和异常。

#### 使用方法

```bash
./monitor-logs.sh [container_name]
```

#### 功能

- ✅ 实时检测 ERROR/FATAL/PANIC 日志
- ✅ 检测性能问题（timeout/slow/memory）
- ✅ 支持集成告警通知

---

### 6. **monitor-performance.sh** - 性能监控

监控 Docker 容器的资源使用情况。

#### 使用方法

```bash
./monitor-performance.sh
```

#### 功能

- ✅ 显示容器资源使用（CPU、内存、网络）
- ✅ 检查服务状态
- ✅ 检查健康状态

---

## 🔧 Nginx 配置

### HTTP + HTTPS 配置

Nginx 同时配置了 HTTP（80 端口）和 HTTPS（443 端口）。

**配置文件位置**：`nginx/nginx.conf`

**当前配置**：
- ✅ HTTP 服务器（80 端口）
- ✅ HTTPS 服务器（443 端口）
- ✅ 反向代理到后端应用
- ✅ 静态文件服务
- ✅ 健康检查端点
- ✅ Gzip 压缩
- ✅ 浏览器缓存

### HTTPS 证书配置

#### 开发环境：自签名证书

```bash
./generate-ssl-cert.sh
docker-compose up -d
```

#### 生产环境：Let's Encrypt 证书

```bash
./letsencrypt.sh yourdomain.com your@email.com apply
docker-compose up -d
```

详细说明：[LETS-ENCRYPT.md](../LETS-ENCRYPT.md)

---

## 📝 其他运维命令

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务
docker-compose logs -f nginx
docker-compose logs -f app
```

### 性能监控

```bash
# 查看容器资源使用
docker stats

# 查看容器状态
docker ps
```

### 健康检查

```bash
# 检查应用
curl http://localhost/health

# 检查 Nginx
docker-compose exec nginx nginx -t
```

---

## 🔗 相关文档

- [Docker Compose 配置](../docs/docker-compose-setup.md)
- [环境变量配置](../docs/env-config.md)
- [CORS 配置指南](../docs/cors-explained.md)
