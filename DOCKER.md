# Docker 部署说明

## 架构说明

```
┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│   Nginx     │ ───→ │    App      │ ───→ │   MySQL     │
│   (80/443)  │      │   (8888)    │      │   (3306)    │
└─────────────┘      └─────────────┘      └─────────────┘
                            │
                            ↓
                     ┌─────────────┐
                     │    Redis    │
                     │   (6379)    │
                     └─────────────┘
```

## 服务端口

| 服务 | 容器端口 | 宿主机端口 | 说明 |
|------|---------|-----------|------|
| Nginx | 80 | 80 | HTTP 入口 |
| Nginx | 443 | 443 | HTTPS 入口（需配置证书） |
| App | 8888 | 不直接暴露 | Go 应用（仅内网访问） |
| MySQL | 3306 | 3306 | 数据库 |
| Redis | 6379 | 6379 | 缓存 |

## 快速开始

### 1. 构建并启动所有服务

```bash
docker-compose up -d --build
```

### 2. 查看日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f app
docker-compose logs -f nginx
```

### 3. 停止服务

```bash
docker-compose down
```

### 4. 重启服务

```bash
docker-compose restart
```

## 访问方式

### API 接口
```bash
# 通过 Nginx 访问（推荐）
http://localhost/api/auth/login
http://localhost/api/user/info

# 直接访问应用（不推荐，仅调试用）
http://localhost:8888/api/auth/login
```

### 静态文件
```bash
# 上传的文件通过 Nginx 访问
http://localhost/uploads/avatar.jpg
```

### 数据库
```bash
# MySQL
mysql -h 127.0.0.1 -P 3306 -u root -p

# Redis
redis-cli -h 127.0.0.1 -p 6379
```

## Nginx 配置说明

### 反向代理
- 所有 `/api/` 请求代理到 Go 应用
- 所有 `/uploads/` 请求代理到 Go 应用的静态文件目录
- 自动透传客户端真实 IP（X-Forwarded-For）

### 健康检查
```bash
curl http://localhost/health
# 返回：healthy
```

## 数据持久化

以下数据通过 Docker Volume 持久化：

- `mysql-data`: MySQL 数据库文件
- `redis-data`: Redis 数据文件
- `./uploads`: 上传的文件（本地目录挂载）

## 环境变量

### App 服务
```yaml
environment:
  - GIN_MODE=release          # Gin 运行模式
  - ALLOWED_ORIGINS=...       # CORS 允许的域名
```

### MySQL 服务
```yaml
environment:
  - MYSQL_ROOT_PASSWORD=root  # root 密码
  - MYSQL_DATABASE=whu_campus_auth  # 初始数据库
```

## 常见问题

### 1. 修改配置后重新加载 Nginx
```bash
docker-compose exec nginx nginx -s reload
```

### 2. 查看容器 IP
```bash
docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' <container_name>
```

### 3. 进入容器调试
```bash
# 进入 App 容器
docker-compose exec app sh

# 进入 Nginx 容器
docker-compose exec nginx sh
```

### 4. 重置数据库
```bash
docker-compose down -v  # 删除所有 volume（谨慎使用！）
docker-compose up -d
```

## 生产环境建议

### 1. HTTPS 配置
在 `nginx/nginx.conf` 中添加 SSL 证书配置：

```nginx
server {
    listen 443 ssl http2;
    server_name yourdomain.com;
    
    ssl_certificate /etc/nginx/ssl/fullchain.pem;
    ssl_certificate_key /etc/nginx/ssl/privkey.pem;
    
    # ... 其他配置
}
```

然后挂载证书：
```yaml
volumes:
  - ./ssl:/etc/nginx/ssl:ro
```

### 2. 日志轮转
配置 Nginx 日志轮转，避免日志文件过大。

### 3. 资源限制
在 `docker-compose.yml` 中添加资源限制：
```yaml
services:
  app:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
```

### 4. 健康检查
添加健康检查配置：
```yaml
services:
  app:
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8888/"]
      interval: 30s
      timeout: 10s
      retries: 3
```

## 性能优化

### 1. Nginx 缓存
对于静态文件，可以在 Nginx 中配置浏览器缓存：

```nginx
location /uploads/ {
    # ... 代理配置
    
    # 浏览器缓存 7 天
    add_header Cache-Control "public, max-age=604800";
}
```

### 2. Gzip 压缩
在 Nginx 中启用 Gzip：

```nginx
http {
    gzip on;
    gzip_types text/plain application/json application/javascript text/css;
    gzip_min_length 1000;
}
```

## 监控

### 1. 查看容器资源使用
```bash
docker stats
```

### 2. 查看 Nginx 访问日志
```bash
docker-compose exec nginx tail -f /var/log/nginx/access.log
```

### 3. 查看应用错误日志
```bash
docker-compose logs -f app | grep ERROR
```
