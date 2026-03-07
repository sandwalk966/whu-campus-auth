# HTTPS 配置快速指南

## 📊 当前状态

### ✅ 已完成的配置

1. **证书挂载** - ✅ 正确
   ```yaml
   volumes:
     - ./ssl:/etc/nginx/ssl:ro  # 证书目录挂载正确
   ```

2. **Nginx HTTPS** - ✅ 已启用
   ```nginx
   server {
       listen 443 ssl http2;
       ssl_certificate /etc/nginx/ssl/fullchain.pem;
       ssl_certificate_key /etc/nginx/ssl/privkey.pem;
   }
   ```

3. **证书路径** - ✅ 正确
   - 证书文件：`/etc/nginx/ssl/fullchain.pem`
   - 私钥文件：`/etc/nginx/ssl/privkey.pem`

### ❌ 不能自动申请证书

当前配置**不支持自动申请 Let's Encrypt 证书**，需要手动生成或申请。

---

## 🚀 使用方案

### 方案一：快速启用（开发/测试环境）⭐ 推荐

使用自签名证书，适合本地开发和测试。

#### 步骤 1：生成自签名证书
```bash
chmod +x scripts/generate-ssl-cert.sh
./scripts/generate-ssl-cert.sh
```

#### 步骤 2：启动服务
```bash
docker-compose up -d
```

#### 步骤 3：验证 HTTPS
```bash
# 测试 HTTP
curl http://localhost/health

# 测试 HTTPS（-k 忽略证书验证）
curl -k https://localhost/health
```

#### 访问方式
- HTTP: `http://localhost`
- HTTPS: `https://localhost`（浏览器会显示证书警告，点击继续即可）

---

### 方案二：生产环境（自动申请 Let's Encrypt 证书）

使用正式的 Let's Encrypt 证书，自动续期。

#### 前提条件
- ✅ 有公网 IP 的服务器
- ✅ 域名已解析到服务器
- ✅ 80 端口可公网访问
- ✅ 防火墙允许 80/443 端口

#### 步骤 1：使用生产环境配置
```bash
# 编辑域名和邮箱
vi scripts/auto-https.sh
# 修改 DOMAIN 和 EMAIL

# 执行自动申请脚本
chmod +x scripts/auto-https.sh
./scripts/auto-https.sh yourdomain.com your@email.com
```

#### 步骤 2：启动生产环境
```bash
docker-compose -f docker-compose.prod.yml up -d
```

#### 步骤 3：验证
```bash
# 测试 HTTPS（正式证书，无警告）
curl -I https://yourdomain.com/health
```

#### 自动续期
Certbot 会自动续期证书（每 12 小时检查一次），无需手动干预。

---

## 📋 配置对比

| 特性 | 方案一（自签名） | 方案二（Let's Encrypt） |
|------|-----------------|----------------------|
| **适用环境** | 开发/测试 | 生产环境 |
| **证书费用** | 免费 | 免费 |
| **浏览器警告** | ⚠️ 有警告 | ✅ 无警告 |
| **自动续期** | ❌ 需手动 | ✅ 自动 |
| **有效期** | 365 天 | 90 天 |
| **配置难度** | 简单 | 中等 |
| **域名要求** | 无 | 需要公网域名 |

---

## 🔍 常见问题

### Q1: 证书文件在哪里？

**自签名证书**：
```
ssl/
├── fullchain.pem  # 证书文件
└── privkey.pem    # 私钥文件
```

**Let's Encrypt 证书**：
```
ssl/
├── live/
│   └── yourdomain.com/
│       ├── fullchain.pem
│       └── privkey.pem
├── archive/
└── renewal/
```

### Q2: 如何查看证书信息？

```bash
# 查看证书详情
openssl x509 -in ssl/fullchain.pem -text -noout

# 查看证书有效期
openssl x509 -in ssl/fullchain.pem -noout -dates

# 查看证书域名
openssl x509 -in ssl/fullchain.pem -noout -subject
```

### Q3: 证书过期了怎么办？

**自签名证书**：
```bash
# 重新生成
./scripts/generate-ssl-cert.sh

# 重启 Nginx
docker-compose restart nginx
```

**Let's Encrypt 证书**：
```bash
# 手动续期（通常会自动续期）
docker-compose -f docker-compose.prod.yml run --rm certbot renew

# 重载 Nginx
docker-compose exec nginx nginx -s reload
```

### Q4: 如何验证 HTTPS 配置是否正确？

```bash
# 1. 检查 Nginx 配置
docker-compose exec nginx nginx -t

# 2. 测试 HTTP 跳转
curl -I http://localhost/health

# 3. 测试 HTTPS
curl -k https://localhost/health

# 4. 查看 SSL 信息
echo | openssl s_client -connect localhost:443 2>/dev/null | openssl x509 -noout -dates
```

### Q5: 浏览器显示"不安全"怎么办？

这是正常的，因为使用的是自签名证书。

**解决方案**：
1. 开发环境：点击"高级" → "继续访问"
2. 生产环境：使用 Let's Encrypt 正式证书

---

## 📝 Nginx 配置说明

### 证书路径配置
```nginx
# 证书文件路径（已配置正确）
ssl_certificate /etc/nginx/ssl/fullchain.pem;
ssl_certificate_key /etc/nginx/ssl/privkey.pem;
```

### SSL 优化配置
```nginx
# 协议版本（已启用 TLS 1.2/1.3）
ssl_protocols TLSv1.2 TLSv1.3;

# 加密套件（已配置现代加密套件）
ssl_ciphers 'ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:...';

# 会话缓存
ssl_session_cache shared:SSL:10m;
ssl_session_timeout 1d;
```

### HTTP 强制跳转（可选）
```nginx
# 在 HTTP server 块中添加（生产环境建议启用）
return 301 https://$server_name$request_uri;
```

---

## 🔧 故障排查

### 1. Nginx 启动失败

```bash
# 查看错误日志
docker-compose logs nginx

# 检查证书文件是否存在
ls -la ssl/

# 测试 Nginx 配置
docker-compose exec nginx nginx -t
```

### 2. HTTPS 无法访问

```bash
# 检查 443 端口是否监听
docker-compose exec nginx netstat -tln | grep 443

# 检查防火墙
# Windows: 防火墙 → 允许应用通过防火墙 → Docker
# Linux: sudo ufw allow 443/tcp
```

### 3. 证书路径错误

检查 `docker-compose.yml` 中的挂载：
```yaml
volumes:
  - ./ssl:/etc/nginx/ssl:ro  # 确保路径正确
```

---

## ✅ 检查清单

使用前请确认：

- [ ] 证书文件已生成（`ls -la ssl/`）
- [ ] Nginx 配置已启用 HTTPS
- [ ] 证书路径配置正确
- [ ] 443 端口已开放
- [ ] 防火墙允许 HTTPS 流量
- [ ] （生产环境）域名已解析到服务器

---

## 📚 相关文档

- [DOCKER-OPTIMIZATION.md](DOCKER-OPTIMIZATION.md) - 完整优化指南
- [DOCKER.md](DOCKER.md) - Docker 部署说明
- [nginx/nginx.conf](nginx/nginx.conf) - Nginx 配置文件

---

## 🎯 快速总结

### 开发环境
```bash
# 1. 生成证书
./scripts/generate-ssl-cert.sh

# 2. 启动服务
docker-compose up -d

# 3. 访问
https://localhost
```

### 生产环境
```bash
# 1. 自动申请证书
./scripts/auto-https.sh yourdomain.com your@email.com

# 2. 启动生产环境
docker-compose -f docker-compose.prod.yml up -d

# 3. 访问
https://yourdomain.com
```

现在你的 Nginx 已经配置好 HTTPS，证书路径也正确！🎉
