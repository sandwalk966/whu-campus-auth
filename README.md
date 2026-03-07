# WHU Campus Auth

武汉大学校园权限管理系统

## 项目结构

```
whu-campus-auth/
├── cmd/                 # 启动入口
│   └── main.go
├── config/              # 配置文件
│   ├── config.go
│   └── config.yaml
├── api/                 # 接口层（controller）
│   ├── user.go
│   ├── role.go
│   ├── menu.go
│   ├── dict.go
│   └── upload.go
├── service/             # 业务逻辑层
│   ├── user_service.go
│   ├── role_service.go
│   ├── menu_service.go
│   ├── dict_service.go
│   └── upload_service.go
├── model/               # 数据库模型 & 请求参数
│   ├── db/              # 数据库表结构体
│   │   ├── user.go
│   │   ├── role.go
│   │   ├── menu.go
│   │   └── dict.go
│   └── req/             # 前端入参结构体
│       ├── user.go
│       ├── role.go
│       ├── menu.go
│       ├── dict.go
│       └── upload.go
├── middleware/          # 中间件
│   ├── jwt.go
│   ├── cors.go
│   └── logger.go
├── utils/               # 工具类
│   ├── jwt.go
│   ├── response.go
│   ├── upload.go
│   └── redis.go
├── router/              # 路由
│   └── router.go
├── initializer/         # 初始化模块
│   ├── database.go      # 数据库初始化
│   ├── migrator.go      # 数据库迁移
│   ├── initializer.go   # 字典数据初始化
│   └── redis.go         # Redis 和日志初始化
├── scripts/             # 运维脚本
│   ├── letsencrypt.sh   # Let's Encrypt 证书管理
│   ├── generate-ssl-cert.sh  # 自签名证书生成
│   ├── monitor-logs.sh  # 日志监控
│   └── monitor-performance.sh  # 性能监控
├── nginx/               # Nginx 配置
│   ├── nginx.conf       # Nginx 配置文件
│   └── healthcheck.sh   # 健康检查脚本
├── ssl/                 # SSL 证书目录（自动生成）
├── uploads/             # 上传文件目录
├── docker-compose.yml   # Docker Compose 配置
├── .env                 # 环境变量
├── .dockerignore        # Docker 构建忽略文件
├── go.mod
├── go.sum
└── .env
```

## 快速开始

### 环境要求

- Go 1.25+
- MySQL 5.7+
- Redis (可选)
- Docker & Docker Compose（推荐）

### 方式一：Docker 部署（推荐）

```bash
# 1. 启动所有服务
docker-compose up -d --build

# 2. 查看日志
docker-compose logs -f

# 3. 访问服务
# HTTP: http://localhost
# HTTPS: https://localhost（自签名证书，浏览器会提示警告）
```

**详细说明**：[DOCKER.md](DOCKER.md)

### 方式二：本地运行

### 安装依赖

```bash
go mod tidy
```

### 配置

编辑 `config.yaml` 文件，配置数据库和 Redis 连接信息。

### 运行

```bash
cd cmd
go run main.go
```

服务器默认启动在 `http://localhost:8888`

## API 接口

### 认证接口
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/register` - 用户注册

### 用户接口
- `GET /api/user/info` - 获取当前用户信息
- `PUT /api/user` - 更新用户信息
- `PUT /api/user/password` - 修改密码
- `GET /api/user/list` - 获取用户列表
- `DELETE /api/user/:id` - 删除用户
- `POST /api/user/assign-roles` - 分配角色

### 角色接口
- `GET /api/role/:id` - 获取角色详情
- `GET /api/role/list` - 获取角色列表
- `GET /api/role/all` - 获取所有角色
- `POST /api/role` - 创建角色
- `PUT /api/role` - 更新角色
- `DELETE /api/role/:id` - 删除角色

### 菜单接口
- `GET /api/menu/:id` - 获取菜单详情
- `GET /api/menu/list` - 获取菜单列表
- `GET /api/menu/tree` - 获取菜单树
- `POST /api/menu` - 创建菜单
- `PUT /api/menu` - 更新菜单
- `DELETE /api/menu/:id` - 删除菜单
- `GET /api/menu/role/:role_id` - 获取角色菜单

### 字典接口
- `GET /api/dict/:id` - 获取字典详情
- `GET /api/dict/list` - 获取字典列表
- `GET /api/dict/code/:code` - 根据编码获取字典
- `POST /api/dict` - 创建字典
- `PUT /api/dict` - 更新字典
- `DELETE /api/dict/:id` - 删除字典

### 上传接口
- `POST /api/upload` - 上传文件
- `DELETE /api/upload/:file_name` - 删除文件

## 数据库表

系统会自动创建以下数据表：
- `sys_user` - 用户表
- `sys_role` - 角色表
- `sys_menu` - 菜单表
- `sys_dict` - 字典表
- `sys_dict_item` - 字典项表
- `user_roles` - 用户角色关联表
- `role_menus` - 角色菜单关联表

## 技术栈

- **框架**: Gin
- **ORM**: GORM
- **认证**: JWT
- **缓存**: Redis (可选)
- **密码加密**: bcrypt
- **反向代理**: Nginx
- **容器化**: Docker & Docker Compose

## 文档

- [Docker 部署说明](DOCKER.md)
- [Let's Encrypt 证书管理](LETS-ENCRYPT.md)
- [Docker 优化指南](DOCKER-OPTIMIZATION.md)
- [运维脚本说明](scripts/README.md)

## License

MIT
