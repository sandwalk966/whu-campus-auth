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

### HTTPS 证书配置（手动）

证书文件需要手动放置到 `ssl/` 目录。

#### 1. 准备证书文件

你需要准备两个证书文件：
- `fullchain.pem` - 证书文件（包含证书链）
- `privkey.pem` - 私钥文件

**证书来源**（任选其一）：
- **Let's Encrypt**（免费，推荐）：使用 Certbot 申请
- **云服务商**：阿里云、腾讯云等提供的免费证书
- **购买商业证书**：DigiCert、GlobalSign 等

#### 2. 放置证书文件

将证书文件放到项目根目录的 `ssl/` 文件夹：

```
whu-campus-auth/
├── ssl/
│   ├── fullchain.pem    ← 你的证书文件
│   └── privkey.pem      ← 你的私钥文件
├── docker-compose.yml
└── nginx/
    └── nginx.conf
```

**注意**：
- ✅ 文件名必须是 `fullchain.pem` 和 `privkey.pem`
- ✅ 确保私钥文件权限安全（不要提交到 Git）
- ✅ `.gitignore` 已配置忽略 `ssl/` 目录

#### 3. 配置域名（可选）

编辑 `nginx/nginx.conf`，修改 `server_name`：

```nginx
# HTTP server
server_name yourdomain.com www.yourdomain.com;

# HTTPS server
server_name yourdomain.com www.yourdomain.com;
```

#### 4. 启用 HTTP→HTTPS 跳转（可选）

编辑 `nginx/nginx.conf`，在 HTTP server 块中取消注释：

```nginx
server {
    listen 80;
    server_name yourdomain.com;
    
    # 取消注释这行
    return 301 https://$server_name$request_uri;
}
```

#### 5. 启动服务

```bash
docker-compose down
docker-compose up -d
```

#### 6. 验证

```bash
# 测试 HTTP
curl -I http://yourdomain.com/health

# 测试 HTTPS
curl -kI https://yourdomain.com/health

# 查看证书信息
openssl s_client -connect yourdomain.com:443 -servername yourdomain.com
```

### 证书续期

证书到期后，只需替换 `ssl/` 目录下的文件，然后重启 Nginx：

```bash
# 替换证书文件
cp /path/to/new/fullchain.pem ./ssl/
cp /path/to/new/privkey.pem ./ssl/

# 重启服务
docker-compose restart nginx
```

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
