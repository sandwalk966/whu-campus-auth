# Docker 生产环境优化指南

本文档介绍如何在生产环境中优化 Docker 部署配置。

## 📋 优化内容概览

1. ✅ **HTTPS 配置** - SSL 证书配置
2. ✅ **性能优化** - Nginx 缓存和 Gzip 压缩
3. ✅ **监控告警** - 健康检查和日志监控
4. ✅ **资源限制** - CPU 和内存限制

---

## 1️⃣ HTTPS 配置

### 开发环境 - 生成自签名证书

```bash
# 生成自签名证书（仅用于开发/测试）
chmod +x scripts/generate-ssl-cert.sh
./scripts/generate-ssl-cert.sh
```

**生成文件**：
- `ssl/fullchain.pem` - 证书文件
- `ssl/privkey.pem` - 私钥文件

### 生产环境 - 使用 Let's Encrypt 证书

```bash
# 使用 Certbot 申请免费证书
docker run --rm \
  -v $(pwd)/ssl:/etc/letsencrypt \
  certbot/certbot certonly \
  --standalone \
  -d yourdomain.com \
  --email your@email.com \
  --agree-tos
```

### 启用 HTTPS

**步骤 1**：编辑 `nginx/nginx.conf`，取消注释 HTTPS server 块（第 65-139 行）

**步骤 2**：编辑 `docker-compose.yml`，确保证书挂载已配置

**步骤 3**：重启服务
```bash
docker-compose up -d
```

### 强制跳转 HTTPS（可选）

在 `nginx/nginx.conf` 的 HTTP server 块中添加：
```nginx
return 301 https://$server_name$request_uri;
```

---

## 2️⃣ 性能优化

### Gzip 压缩（已配置）

在 `nginx/nginx.conf` 中已启用 Gzip 压缩：

```nginx
gzip on;
gzip_vary on;
gzip_min_length 1024;
gzip_types text/plain text/css text/xml application/json application/javascript;
```

**效果**：
- 减少 60-80% 的传输数据量
- 加快页面加载速度
- 节省带宽

### 浏览器缓存（已配置）

静态文件缓存 7 天：

```nginx
location /uploads/ {
    add_header Cache-Control "public, max-age=604800";
}
```

### Nginx 连接优化

```nginx
# 在 http 块中添加（可选）
http {
    # 开启高效文件传输
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    
    # 保持连接
    keepalive_timeout 65;
    keepalive_requests 100;
    
    # 代理缓冲
    proxy_buffering on;
    proxy_buffer_size 4k;
    proxy_buffers 8 4k;
}
```

---

## 3️⃣ 监控告警

### 健康检查（已配置）

所有服务都配置了健康检查：

| 服务 | 检查方式 | 间隔 | 超时 |
|------|---------|------|------|
| Nginx | 检查进程和端口 | 30s | 10s |
| App | HTTP 请求 | 30s | 10s |
| MySQL | MySQL ping | 30s | 10s |
| Redis | Redis ping | 30s | 10s |

### 查看健康状态

```bash
# 查看所有容器健康状态
docker inspect --format='{{.Name}}: {{.State.Health.Status}}' $(docker ps -q)

# 查看详细信息
docker inspect whu-campus-auth-app-1 | grep -A 20 Health
```

### 日志监控

**实时监控日志**：
```bash
# 查看应用日志
docker-compose logs -f app

# 查看 Nginx 访问日志
docker-compose exec nginx tail -f /var/log/nginx/access.log

# 查看错误日志
docker-compose exec nginx tail -f /var/log/nginx/error.log
```

**使用监控脚本**：
```bash
chmod +x scripts/monitor-logs.sh
./scripts/monitor-logs.sh whu-campus-auth-app-1
```

**集成告警（示例：钉钉）**：

编辑 `scripts/monitor-logs.sh`，取消注释钉钉告警部分：

```bash
curl -X POST https://oapi.dingtalk.com/robot/send \
  -H 'Content-Type: application/json' \
  -d "{\"msgtype\":\"text\",\"text\":{\"content\":\"⚠️ 应用错误告警：$line\"}}"
```

### 性能监控

```bash
# 查看容器资源使用
chmod +x scripts/monitor-performance.sh
./scripts/monitor-performance.sh

# 实时查看资源使用
docker stats
```

---

## 4️⃣ 资源限制

### 已配置的资源限制

| 服务 | CPU 限制 | 内存限制 | CPU 预留 | 内存预留 |
|------|---------|---------|---------|---------|
| Nginx | 0.5 | 128M | 0.25 | 64M |
| App | 1.0 | 512M | 0.5 | 256M |
| MySQL | 1.0 | 1G | 0.5 | 512M |
| Redis | 0.5 | 256M | 0.25 | 128M |

### 调整资源限制

编辑 `docker-compose.yml`，修改 `deploy.resources` 部分：

```yaml
services:
  app:
    deploy:
      resources:
        limits:
          cpus: '2.0'      # 最大 2 核 CPU
          memory: 1G       # 最大 1G 内存
        reservations:
          cpus: '1.0'      # 预留 1 核 CPU
          memory: 512M     # 预留 512M 内存
```

### 查看资源使用

```bash
# 实时查看
docker stats

# 查看历史（需要额外工具）
docker exec whu-campus-auth-app-1 top
```

---

## 🚀 部署流程

### 1. 准备环境

```bash
# 生成 SSL 证书（开发环境）
./scripts/generate-ssl-cert.sh

# 或者使用生产证书
# 将证书文件放在 ssl/ 目录
```

### 2. 配置优化

- ✅ 编辑 `nginx/nginx.conf`，启用 HTTPS 配置
- ✅ 编辑 `docker-compose.yml`，调整资源限制（如需要）

### 3. 构建并启动

```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f
```

### 4. 验证部署

```bash
# 检查所有服务状态
docker-compose ps

# 检查健康状态
docker inspect --format='{{.Name}}: {{.State.Health.Status}}' $(docker ps -q)

# 测试 HTTPS
curl -k https://localhost/health

# 测试 Gzip
curl -H "Accept-Encoding: gzip" -I http://localhost/api/auth/login
```

### 5. 监控和维护

```bash
# 定期查看性能
./scripts/monitor-performance.sh

# 监控日志
./scripts/monitor-logs.sh

# 查看资源使用
docker stats
```

---

## 📊 性能测试

### 使用 Apache Bench 测试

```bash
# 安装 Apache Bench
# Windows: 下载 httpd 工具包
# Linux: sudo apt-get install apache2-utils

# 测试
ab -n 1000 -c 10 http://localhost/api/auth/login

# 测试 HTTPS
ab -n 1000 -c 10 -Z TLSv1.2 https://localhost/api/auth/login
```

### 查看性能指标

```bash
# Nginx 访问日志（包含响应时间）
docker-compose exec nginx tail -f /var/log/nginx/access.log

# 应用性能日志
docker-compose logs -f app | grep "response time"
```

---

## 🔧 故障排查

### 1. 健康检查失败

```bash
# 查看健康检查日志
docker inspect whu-campus-auth-app-1 | grep -A 30 Health

# 手动执行健康检查
docker exec whu-campus-auth-app-1 wget -q --spider http://localhost:8888/
echo $?  # 应该返回 0
```

### 2. 资源不足

```bash
# 查看资源使用
docker stats

# 增加资源限制
# 编辑 docker-compose.yml，增加 limits 和 reservations
```

### 3. SSL 证书问题

```bash
# 检查证书文件
ls -la ssl/

# 测试 Nginx 配置
docker-compose exec nginx nginx -t

# 查看证书信息
openssl x509 -in ssl/fullchain.pem -text -noout
```

### 4. 日志过大

```bash
# 查看日志大小
docker inspect --format='{{.LogPath}}' $(docker ps -q) | xargs du -h

# 清理日志
docker-compose down -v
docker-compose up -d

# 或者配置日志轮转（推荐）
```

---

## 📈 进一步优化建议

### 1. 日志轮转

创建 `nginx/logrotate.conf`：
```
/var/log/nginx/*.log {
    daily
    missingok
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 www-data adm
    sharedscripts
    postrotate
        [ -f /var/run/nginx.pid ] && kill -USR1 `cat /var/run/nginx.pid`
    endscript
}
```

### 2. 自动扩缩容

使用 Docker Swarm 或 Kubernetes 实现自动扩缩容：

```yaml
# Docker Swarm 示例
services:
  app:
    deploy:
      replicas: 3
      update_config:
        parallelism: 1
        delay: 10s
      restart_policy:
        condition: on-failure
```

### 3. 集中式日志

使用 ELK Stack（Elasticsearch, Logstash, Kibana）或 Loki + Grafana

### 4. 应用性能监控（APM）

集成 Prometheus + Grafana 或 New Relic、DataDog 等监控工具

---

## ✅ 检查清单

部署前请确认：

- [ ] SSL 证书已生成或配置
- [ ] Nginx HTTPS 配置已启用
- [ ] 健康检查已配置
- [ ] 资源限制已调整
- [ ] 日志监控已配置
- [ ] 性能优化已启用（Gzip、缓存）
- [ ] 备份策略已制定
- [ ] 告警通知已配置

---

## 📚 参考资源

- [Nginx 官方文档](https://nginx.org/en/docs/)
- [Docker 最佳实践](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)
- [Let's Encrypt](https://letsencrypt.org/)
- [Prometheus 监控](https://prometheus.io/)
