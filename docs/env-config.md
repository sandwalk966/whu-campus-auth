# 环境变量配置指南

## 📋 概述

本项目支持三种配置方式，优先级从高到低：

1. **环境变量**（最高优先级）
2. **.env 文件**
3. **config.yaml 文件**（最低优先级）

## 🚀 使用方法

### 1. 复制示例文件

```bash
cp .env.example .env
```

### 2. 编辑 .env 文件

```bash
# 生产环境配置示例
SERVER_PORT=8888
GIN_MODE=release
ALLOWED_ORIGINS=https://yourdomain.com
DB_HOST=localhost
DB_PASSWORD=your-password
JWT_SECRET=your-secret-key
```

### 3. 启动应用

```bash
# 本地运行
go run cmd/main.go

# Docker 运行
docker-compose up -d
```

## 🔧 环境变量详解

### 服务器配置

| 变量名 | 说明 | 默认值 | 示例 |
|--------|------|--------|------|
| `SERVER_PORT` | 服务端口 | 8888 | 8888 |
| `GIN_MODE` | 运行模式 | release | release / debug |

### CORS 配置

| 变量名 | 说明 | 默认值 | 示例 |
|--------|------|--------|------|
| `ALLOWED_ORIGINS` | 允许的域名列表 | 空 | `https://example.com,https://www.example.com` |

**说明**：
- 留空：只允许 HTTPS 请求和 Nginx 反向代理
- 多个域名用逗号分隔
- 生产环境建议填写具体域名

### 数据库配置

| 变量名 | 说明 | 默认值 | 示例 |
|--------|------|--------|------|
| `DB_HOST` | 数据库主机 | mysql | localhost |
| `DB_PORT` | 数据库端口 | 3306 | 3306 |
| `DB_NAME` | 数据库名称 | whu_campus_auth | whu_campus_auth |
| `DB_USER` | 数据库用户 | root | root |
| `DB_PASSWORD` | 数据库密码 | root | your-password |

### Redis 配置

| 变量名 | 说明 | 默认值 | 示例 |
|--------|------|--------|------|
| `REDIS_HOST` | Redis 主机 | redis | localhost |
| `REDIS_PORT` | Redis 端口 | 6379 | 6379 |
| `REDIS_PASSWORD` | Redis 密码 | 空 | your-password |

### JWT 配置

| 变量名 | 说明 | 默认值 | 示例 |
|--------|------|--------|------|
| `JWT_SECRET` | JWT 密钥 | your-secret-key-change-in-production | **必须修改** |
| `JWT_EXPIRE` | JWT 过期时间（秒） | 86400 | 86400（24 小时） |

**⚠️ 安全警告**：
- 生产环境必须修改 `JWT_SECRET`！
- 使用随机字符串：`openssl rand -base64 32`

### 日志配置

| 变量名 | 说明 | 默认值 | 示例 |
|--------|------|--------|------|
| `LOG_LEVEL` | 日志级别 | info | debug / info / warn / error |

### 文件上传配置

| 变量名 | 说明 | 默认值 | 示例 |
|--------|------|--------|------|
| `UPLOAD_PATH` | 上传文件路径 | ./uploads | ./uploads |
| `MAX_UPLOAD_SIZE` | 最大上传大小（字节） | 10485760 | 10485760（10MB） |

## 📊 配置优先级示例

### 示例 1：config.yaml

```yaml
server:
  port: 8888
  mode: release
database:
  host: mysql
  password: root
```

### 示例 2：.env 文件

```env
SERVER_PORT=9000
DB_HOST=localhost
DB_PASSWORD=new-password
```

### 示例 3：环境变量

```bash
export DB_PASSWORD=env-password
```

### 最终配置

```
port: 9000              # 来自 .env（覆盖 config.yaml）
mode: release           # 来自 config.yaml
host: localhost         # 来自 .env（覆盖 config.yaml）
password: env-password  # 来自环境变量（最高优先级）
```

## 🔐 生产环境配置建议

### 1. 使用 .env.production 文件

```bash
# .env.production
SERVER_PORT=8888
GIN_MODE=release
ALLOWED_ORIGINS=https://yourdomain.com
DB_HOST=mysql
DB_PASSWORD=strong-password-123
JWT_SECRET=$(openssl rand -base64 32)
REDIS_HOST=redis
LOG_LEVEL=warn
```

### 2. 在 docker-compose.yml 中引用

```yaml
services:
  app:
    environment:
      - GIN_MODE=release
      - ALLOWED_ORIGINS=https://yourdomain.com
    env_file:
      - .env.production
```

### 3. 不要提交敏感信息

```bash
# .gitignore
.env
.env.*.local
!.env.example
```

## 💡 使用技巧

### 1. 生成随机 JWT 密钥

```bash
# Linux / macOS
openssl rand -base64 32

# Windows PowerShell
[System.Web.Security.Membership]::GeneratePassword(32, 8)
```

### 2. 查看当前环境变量

```bash
# Linux / macOS
env | grep YOUR_PREFIX

# Windows PowerShell
Get-ChildItem Env:
```

### 3. 测试配置是否生效

```bash
# 启动应用
go run cmd/main.go

# 查看日志输出
# 应该显示加载的配置信息
```

### 4. Docker 环境配置

```yaml
# docker-compose.yml
services:
  app:
    build: .
    environment:
      - GIN_MODE=release
      - ALLOWED_ORIGINS=
      - DB_HOST=mysql
      - DB_PASSWORD=${DB_PASSWORD:-root}
      - JWT_SECRET=${JWT_SECRET:-change-me}
```

使用 `${VAR:-default}` 语法提供默认值。

## 🐛 常见问题

### Q1: .env 文件不生效？

**检查**：
1. 文件名是否正确（`.env` 不是 `.env.txt`）
2. 文件路径是否正确（在项目根目录）
3. 格式是否正确（`KEY=VALUE` 无空格）

### Q2: 环境变量优先级混乱？

**记住**：环境变量 > .env > config.yaml

### Q3: Docker 中如何传递环境变量？

```bash
# 方法 1：docker-compose.yml
environment:
  - KEY=value

# 方法 2：--env-file
docker run --env-file .env myapp

# 方法 3：-e 参数
docker run -e KEY=value myapp
```

### Q4: 如何验证配置？

```bash
# 启动时查看日志
go run cmd/main.go

# 日志应该显示：
# - 是否加载 .env 文件
# - 最终使用的配置值
```

## 📚 相关文档

- [config.go](../config/config.go) - 配置加载代码
- [main.go](../cmd/main.go) - 主程序入口
- [.env.example](../.env.example) - 环境变量示例
- [docker-compose.yml](../docker-compose.yml) - Docker 配置

## 🔗 参考资料

- [godotenv](https://github.com/joho/godotenv) - Go 环境文件加载器
- [12-Factor App](https://12factor.net/zh_cn/config) - 配置最佳实践
