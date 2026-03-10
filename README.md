# WHU Campus Auth

Wuhan University Campus Permission Management System

## Project Structure

```
whu-campus-auth/
├── cmd/                 # Entry point
│   └── main.go
├── config/              # Configuration
│   ├── config.go
│   └── config.yaml
├── dao/                 # Data Access Layer
│   ├── base_dao.go
│   ├── user_dao.go
│   ├── role_dao.go
│   ├── menu_dao.go
│   └── dict_dao.go
├── api/                 # Controller Layer
│   ├── user.go
│   ├── role.go
│   ├── menu.go
│   ├── dict.go
│   └── upload.go
├── service/             # Business Logic Layer
│   ├── user_service.go
│   ├── role_service.go
│   ├── menu_service.go
│   ├── dict_service.go
│   └── upload_service.go
├── model/               # Models & Request DTOs
│   ├── db/              # Database models
│   │   ├── user.go
│   │   ├── role.go
│   │   ├── menu.go
│   │   └── dict.go
│   └── req/             # Request DTOs
│       ├── user.go
│       ├── role.go
│       ├── menu.go
│       ├── dict.go
│       └── upload.go
├── middleware/          # Middleware
│   ├── jwt.go           # JWT Authentication
│   ├── casbin_rbac.go   # RBAC Authorization
│   ├── db.go            # Database Connection
│   └── logger.go        # Logging
├── utils/               # Utilities
│   ├── jwt.go
│   ├── response.go
│   ├── upload.go
│   └── redis.go
├── router/              # Router Configuration
│   ├── auth_routes.go   # Authentication routes
│   ├── user_routes.go   # User routes
│   ├── role_routes.go   # Role routes
│   ├── menu_routes.go   # Menu routes
│   ├── dict_routes.go   # Dictionary routes
│   └── static_routes.go # Static file routes
├── initializer/         # Initialization Modules
│   ├── initializer.go   # Unified initialization
│   ├── database.go      # Database initialization
│   ├── redis.go         # Redis initialization
│   ├── migrator.go      # Database migration
│   ├── admin_initializer.go  # Admin account initialization
│   ├── deps.go          # Dependency injection
│   ├── api.go           # API layer initialization
│   ├── dao.go           # DAO layer initialization
│   ├── service.go       # Service layer initialization
│   └── router.go        # Router initialization
├── scripts/             # Operations Scripts
│   ├── letsencrypt.sh   # Let's Encrypt certificate management
│   ├── generate-ssl-cert.sh  # Self-signed certificate generation
│   ├── monitor-logs.sh  # Log monitoring
│   └── monitor-performance.sh  # Performance monitoring
├── nginx/               # Nginx Configuration
│   ├── nginx.conf       # Nginx configuration file
│   └── healthcheck.sh   # Health check script
├── ssl/                 # SSL certificates (auto-generated)
├── uploads/             # Uploaded files
├── frontend/            # Frontend Project
│   ├── src/
│   │   ├── api/         # API calls
│   │   ├── layouts/     # Layout components
│   │   ├── router/      # Router configuration
│   │   ├── stores/      # State management (Pinia)
│   │   ├── utils/       # Utility functions
│   │   └── views/       # Page components
│   ├── index.html
│   ├── package.json
│   └── vite.config.js
├── docker-compose.yml   # Docker Compose configuration
├── .env                 # Environment variables
├── .dockerignore        # Docker ignore file
├── go.mod
├── go.sum
└── .env
```

## Requirements

- Go 1.25+
- MySQL 5.7+
- Redis (optional)
- Node.js 16+ (frontend)
- Docker & Docker Compose (recommended)

## Default Administrator Account

The system automatically creates a default administrator account on first startup:

| Field | Value |
|-------|-------|
| **Username** | `admin` |
| **Password** | `your password(set in the .env file)` |
| **Role** | Super Administrator |


For more details, see [Default Administrator Account Documentation](docs/DEFAULT_ADMIN.md).

## Deployment

### Method 1: Docker Deployment (Recommended)

#### Prerequisites

- Docker installed
- Docker Compose installed

#### Step 1: Clone the Repository

```bash
git clone <repository-url>
cd whu-campus-auth
```

#### Step 2: Configure Environment Variables

Copy the `.env.example` file to `.env` and modify the values:

```bash
cp .env.example .env
```

Edit `.env` file with your preferred editor:

```env
# Database Configuration
MYSQL_ROOT_PASSWORD=your_root_password
MYSQL_DATABASE=whu_campus_auth
MYSQL_USER=whu_user
MYSQL_PASSWORD=your_password

# Redis Configuration (optional)
REDIS_PASSWORD=your_redis_password

# Application Configuration
GIN_MODE=release
SERVER_PORT=8888
```

#### Step 3: Start All Services

```bash
# Build and start all containers
docker-compose up -d --build

# View logs
docker-compose logs -f

# Check service status
docker-compose ps
```

#### Step 4: Access the Application

- **HTTP**: http://localhost
- **HTTPS**: https://localhost (self-signed certificate, browser will show warning)
- **Frontend**: http://localhost:3000 (development mode)

#### Step 5: Stop Services

```bash
# Stop all containers
docker-compose down

# Stop and remove volumes (use with caution)
docker-compose down -v
```

**Detailed Docker Documentation**: [DOCKER.md](DOCKER.md)

---

### Method 2: Local Backend Deployment

#### Step 1: Install Dependencies

```bash
go mod tidy
```

#### Step 2: Configure Database

Make sure MySQL is running, then edit `config/config.yaml`:

```yaml
database:
  host: localhost
  port: 3306
  user: root
  password: your_password
  dbname: whu_campus_auth
  
redis:
  enabled: true
  host: localhost
  port: 6379
  password: ""
```

#### Step 3: Run the Application

```bash
cd cmd
go run main.go
```

The server will start at `http://localhost:8888` by default.

#### Step 4: Production Build

```bash
cd cmd
go build -o whu-campus-auth.exe main.go
./whu-campus-auth.exe
```

---

### Method 3: Local Frontend Deployment

#### Step 1: Install Dependencies

```bash
cd frontend
npm install
```

#### Step 2: Configure API Endpoint

Edit `frontend/.env` or `frontend/vite.config.js` to set the backend API URL:

```javascript
// vite.config.js
export default {
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8888',
        changeOrigin: true
      }
    }
  }
}
```

#### Step 3: Run Development Server

```bash
npm run dev
```

The frontend will start at `http://localhost:3000` by default.

#### Step 4: Production Build

```bash
npm run build
```

The built files will be in the `frontend/dist` directory.

---

### Method 4: Production Deployment with Nginx

#### Step 1: Build Frontend

```bash
cd frontend
npm install
npm run build
```

#### Step 2: Configure Nginx

Use the provided Nginx configuration in `nginx/nginx.conf`:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # Frontend static files
    location / {
        root /path/to/whu-campus-auth/frontend/dist;
        try_files $uri $uri/ /index.html;
    }

    # Backend API proxy
    location /api/ {
        proxy_pass http://localhost:8888;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

#### Step 3: Start Backend Service

```bash
cd cmd
go build -o whu-campus-auth main.go
nohup ./whu-campus-auth &
```

#### Step 4: Start Nginx

```bash
sudo nginx -c /path/to/nginx/nginx.conf
```

---

### Method 5: HTTPS with Let's Encrypt

#### Prerequisites

- Domain name pointing to your server
- Ports 80 and 443 open

#### Step 1: Install Certificate

```bash
cd scripts
chmod +x letsencrypt.sh
sudo ./letsencrypt.sh install your-domain.com
```

#### Step 2: Configure Nginx for HTTPS

Update `nginx/nginx.conf` with SSL configuration (see `nginx/nginx.conf` for details).

#### Step 3: Auto-renew Certificate

```bash
# Add to crontab
sudo crontab -e

# Add this line for monthly renewal
0 0 1 * * /path/to/scripts/letsencrypt.sh renew
```

**Detailed HTTPS Documentation**: [LETS-ENCRYPT.md](LETS-ENCRYPT.md)

---

## API Endpoints

### Authentication
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration

### User Management
- `GET /api/user/info` - Get current user info
- `POST /api/user` - Create user
- `PUT /api/user` - Update user info
- `PUT /api/user/password` - Change password
- `GET /api/user/list` - Get user list
- `DELETE /api/user/:id` - Delete user
- `POST /api/user/assign-roles` - Assign roles
- `POST /api/user/avatar` - Upload avatar

### Role Management
- `GET /api/role/:id` - Get role details
- `GET /api/role/list` - Get role list
- `GET /api/role/all` - Get all roles
- `POST /api/role` - Create role
- `PUT /api/role` - Update role
- `DELETE /api/role/:id` - Delete role

### Menu Management
- `GET /api/menu/:id` - Get menu details
- `GET /api/menu/list` - Get menu list
- `GET /api/menu/tree` - Get menu tree
- `POST /api/menu` - Create menu
- `PUT /api/menu` - Update menu
- `DELETE /api/menu/:id` - Delete menu
- `GET /api/menu/role/:role_id` - Get role menus

### Dictionary Management
- `GET /api/dict/:id` - Get dictionary details
- `GET /api/dict/list` - Get dictionary list
- `GET /api/dict/code/:code` - Get dictionary by code
- `POST /api/dict` - Create dictionary
- `PUT /api/dict` - Update dictionary
- `DELETE /api/dict/:id` - Delete dictionary

### File Upload
- `POST /api/upload` - Upload file
- `DELETE /api/upload/:file_name` - Delete file

## Database Tables

The system automatically creates the following tables:
- `sys_user` - User table
- `sys_role` - Role table
- `sys_menu` - Menu table
- `sys_dict` - Dictionary table
- `sys_dict_item` - Dictionary item table
- `user_roles` - User-Role relationship table
- `role_menus` - Role-Menu relationship table

## Tech Stack

### Backend
- **Framework**: Gin
- **ORM**: GORM
- **Authentication**: JWT
- **Cache**: Redis (optional)
- **Password Encryption**: bcrypt
- **Logging**: Zap
- **Authorization**: Casbin RBAC

### Frontend
- **Framework**: Vue 3
- **Build Tool**: Vite
- **UI Library**: Element Plus
- **State Management**: Pinia
- **Router**: Vue Router
- **HTTP Client**: Axios

### DevOps
- **Reverse Proxy**: Nginx
- **Containerization**: Docker & Docker Compose
- **HTTPS**: Let's Encrypt

## Documentation

- [Docker Deployment Guide](DOCKER.md)
- [Let's Encrypt Certificate Management](LETS-ENCRYPT.md)
- [Docker Optimization Guide](DOCKER-OPTIMIZATION.md)
- [Operations Scripts](scripts/README.md)

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Check MySQL service is running
   - Verify database credentials in `config.yaml`
   - Ensure database exists

2. **Port Already in Use**
   - Change `SERVER_PORT` in `.env` or `config.yaml`
   - Check if another service is using the port

3. **Frontend Cannot Connect to Backend**
   - Verify backend API URL in frontend configuration
   - Check CORS settings
   - Ensure backend service is running

4. **Docker Build Fails**
   - Clean Docker cache: `docker system prune`
   - Rebuild without cache: `docker-compose build --no-cache`

## License

MIT
