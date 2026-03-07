# Let's Encrypt HTTPS 证书管理指南

## 📊 架构说明

### 优化后的设计

**之前的问题**：
- ❌ 运行一个常驻的 Certbot 容器（浪费资源）
- ❌ 两个 docker-compose 文件（维护困难）
- ❌ Certbot 容器一直运行，每 12 小时检查一次

**现在的方案**：
- ✅ **按需运行** Certbot（不占用资源）
- ✅ **单一 docker-compose 文件**（简单易维护）
- ✅ **脚本化管理**（清晰透明）

---

## 🚀 快速开始

### 方案 A：开发环境（自签名证书）

适合本地开发和测试。

```bash
# 1. 生成自签名证书
chmod +x scripts/generate-ssl-cert.sh
./scripts/generate-ssl-cert.sh

# 2. 启动服务
docker-compose up -d

# 3. 访问
curl -k https://localhost/health
```

---

### 方案 B：生产环境（Let's Encrypt 证书）⭐ 推荐

适合有公网域名的生产环境。

#### 前提条件
- ✅ 服务器有公网 IP
- ✅ 域名已解析到服务器（A 记录）
- ✅ 80 和 443 端口开放
- ✅ 防火墙允许 HTTP/HTTPS 流量

#### 步骤 1：申请证书

```bash
# 使用自动脚本（推荐）
chmod +x scripts/letsencrypt.sh
./scripts/letsencrypt.sh yourdomain.com your@email.com apply

# 示例：
./scripts/letsencrypt.sh example.com admin@example.com apply
```

**脚本会执行**：
1. 临时停止 Nginx（释放 80 端口）
2. 运行 Certbot 容器申请证书
3. 重启 Nginx
4. 生成生产环境 Nginx 配置
5. 自动应用配置

#### 步骤 2：重启服务

```bash
docker-compose up -d
```

#### 步骤 3：验证

```bash
# 测试 HTTPS
curl -I https://yourdomain.com/health

# 查看证书信息
echo | openssl s_client -connect yourdomain.com:443 2>/dev/null | openssl x509 -noout -dates
```

---

## 📋 证书管理

### 1. 申请新证书

```bash
./scripts/letsencrypt.sh <域名> <邮箱> apply
```

**示例**：
```bash
./scripts/letsencrypt.sh example.com admin@example.com apply
```

### 2. 续期证书

Let's Encrypt 证书有效期为 **90 天**，需要定期续期。

```bash
# 手动续期
./scripts/letsencrypt.sh <域名> <邮箱> renew
```

**示例**：
```bash
./scripts/letsencrypt.sh example.com admin@example.com renew
```

### 3. 生成 Nginx 配置

```bash
./scripts/letsencrypt.sh <域名> <邮箱> config
```

**示例**：
```bash
./scripts/letsencrypt.sh example.com admin@example.com config
```

---

## ⏰ 自动续期配置

### 方法 1：Cron 定时任务（推荐）

```bash
# 编辑 crontab
crontab -e

# 添加以下行（每天凌晨 2 点检查续期）
0 2 * * * cd /path/to/whu-campus-auth && ./scripts/letsencrypt.sh yourdomain.com your@email.com renew >> /var/log/certbot-renew.log 2>&1
```

### 方法 2：Systemd Timer

创建 `/etc/systemd/system/certbot-timer.service`：

```ini
[Unit]
Description=Let's Encrypt Certificate Renewal
After=network.target

[Service]
Type=oneshot
User=root
WorkingDirectory=/path/to/whu-campus-auth
ExecStart=/path/to/whu-campus-auth/scripts/letsencrypt.sh yourdomain.com your@email.com renew
```

创建 `/etc/systemd/system/certbot-timer.timer`：

```ini
[Unit]
Description=Run Let's Encrypt Renewal Daily

[Timer]
OnCalendar=*-*-* 02:00:00
Persistent=true

[Install]
WantedBy=timers.target
```

启用定时器：
```bash
sudo systemctl daemon-reload
sudo systemctl enable certbot-timer.timer
sudo systemctl start certbot-timer.timer
```

---

## 📁 文件结构

```
whu-campus-auth/
├── ssl/                          # SSL 证书目录
│   ├── live/
│   │   └── yourdomain.com/
│   │       ├── fullchain.pem     # 证书文件
│   │       └── privkey.pem       # 私钥文件
│   ├── archive/                  # 历史证书
│   └── renewal/                  # 续期配置
│
├── certbot-www/                  # ACME 验证目录
│   └── .well-known/acme-challenge/
│
├── certbot-logs/                 # Certbot 日志
│   └── letsencrypt.log
│
├── scripts/
│   ├── letsencrypt.sh           # Let's Encrypt 管理脚本 ⭐
│   └── generate-ssl-cert.sh     # 自签名证书脚本
│
├── nginx/
│   ├── nginx.conf               # 当前 Nginx 配置
│   └── nginx.prod.conf          # 生产环境配置（自动生成）
│
└── docker-compose.yml           # Docker 配置
```

---

## 🔍 常见问题

### Q1: 证书在哪里？

**Let's Encrypt 证书**：
```bash
ls -la ssl/live/yourdomain.com/
# fullchain.pem  - 证书
# privkey.pem    - 私钥
```

**自签名证书**：
```bash
ls -la ssl/
# fullchain.pem  - 证书
# privkey.pem    - 私钥
```

### Q2: 如何查看证书有效期？

```bash
# 查看证书信息
openssl x509 -in ssl/live/yourdomain.com/fullchain.pem -noout -dates

# 输出示例：
# notBefore=Mar  7 00:00:00 2026 GMT
# notAfter=Jun  5 00:00:00 2026 GMT  ← 90 天后过期
```

### Q3: 如何检查续期是否成功？

```bash
# 测试续期（dry-run）
./scripts/letsencrypt.sh yourdomain.com your@email.com renew

# 查看续期日志
cat certbot-logs/letsencrypt.log
```

### Q4: 续期后需要重启服务吗？

脚本会**自动重载 Nginx**，无需手动重启：
```bash
docker-compose exec nginx nginx -s reload
```

### Q5: 可以申请多个域名吗？

可以！支持多域名和通配符：

```bash
# 多域名
./scripts/letsencrypt.sh example.com admin@example.com apply
# 然后在脚本中添加 -d www.example.com -d api.example.com

# 通配符域名（需要 DNS 验证）
# 目前脚本使用 HTTP 验证，不支持通配符
```

### Q6: 80 端口被占用怎么办？

脚本会**临时停止 Nginx**，申请完成后自动重启：

```bash
# 停止 Nginx（释放 80 端口）
docker-compose stop nginx

# 申请证书
docker run --rm -v $(pwd)/ssl:/etc/letsencrypt -p 80:80 certbot/certbot ...

# 重启 Nginx
docker-compose start nginx
```

---

## 📊 两种方案对比

| 特性 | 自签名证书 | Let's Encrypt |
|------|-----------|--------------|
| **费用** | 免费 | 免费 |
| **有效期** | 365 天 | 90 天 |
| **浏览器警告** | ⚠️ 有警告 | ✅ 无警告 |
| **自动续期** | ❌ 手动 | ✅ 自动 |
| **适用环境** | 开发/测试 | 生产环境 |
| **配置难度** | 简单 | 中等 |
| **域名要求** | 无 | 公网域名 |

---

## 🎯 最佳实践

### 1. 开发环境
```bash
# 使用自签名证书
./scripts/generate-ssl-cert.sh
docker-compose up -d
```

### 2. 生产环境
```bash
# 申请 Let's Encrypt 证书
./scripts/letsencrypt.sh example.com admin@example.com apply

# 设置自动续期（crontab）
0 2 * * * cd /path && ./scripts/letsencrypt.sh example.com admin@example.com renew
```

### 3. 监控证书过期
```bash
# 检查证书是否即将过期（30 天内）
openssl x509 -in ssl/live/example.com/fullchain.pem -noout -checkend 2592000
if [ $? -ne 0 ]; then
    echo "⚠️  证书即将过期，请立即续期！"
    # 发送告警通知
fi
```

---

## 🔧 故障排查

### 证书申请失败

```bash
# 1. 检查域名解析
ping example.com
nslookup example.com

# 2. 检查 80 端口
telnet example.com 80

# 3. 查看 Certbot 日志
cat certbot-logs/letsencrypt.log

# 4. 手动测试
docker run --rm -it certbot/certbot --version
```

### Nginx 无法启动

```bash
# 检查配置
docker-compose exec nginx nginx -t

# 查看日志
docker-compose logs nginx

# 检查证书路径
ls -la ssl/live/example.com/
```

---

## ✅ 检查清单

申请证书前请确认：

- [ ] 域名已解析到服务器
- [ ] 80 端口可公网访问
- [ ] 防火墙已开放 80/443
- [ ] 邮箱地址正确
- [ ] 服务器时间准确

---

## 📚 相关资源

- [Let's Encrypt 官方文档](https://letsencrypt.org/docs/)
- [Certbot 使用指南](https://certbot.eff.org/)
- [SSL Labs 测试](https://www.ssllabs.com/ssltest/) - 测试 HTTPS 配置

---

## 🎉 总结

### 现在的优势

✅ **简单**：一个脚本搞定所有操作  
✅ **清晰**：没有常驻容器，资源占用为零  
✅ **灵活**：支持多种场景（开发/生产）  
✅ **易维护**：单一 docker-compose 文件  

### 快速命令

```bash
# 开发环境
./scripts/generate-ssl-cert.sh && docker-compose up -d

# 生产环境
./scripts/letsencrypt.sh example.com admin@example.com apply
```

现在你的 HTTPS 配置既简单又强大！🎉
