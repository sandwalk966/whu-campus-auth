#!/bin/bash

# HTTPS 证书配置脚本（Linux 服务器）
# 用于申请和配置 Let's Encrypt SSL 证书
# 
# 使用方法：
#   sudo bash setup-https.sh yourdomain.com your@email.com
#   sudo bash setup-https.sh yourdomain.com your@email.com --production

set -e

# 参数检查
if [ "$#" -lt 2 ]; then
    echo "使用方法：$0 <域名> <邮箱> [--production]"
    echo ""
    echo "示例："
    echo "  $0 example.com admin@example.com"
    echo "  $0 example.com admin@example.com --production"
    exit 1
fi

DOMAIN=$1
EMAIL=$2
PRODUCTION=${3:-}

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${CYAN}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

# 检查是否在服务器上运行
log_info "检查运行环境..."
if [ ! -d "/etc/nginx" ]; then
    log_warn "这看起来像是在本地环境"
    log_info "HTTPS 证书需要在服务器上配置"
    log_info "请将此脚本上传到服务器后执行"
    exit 1
fi

# 检查 certbot 是否安装
check_certbot() {
    if command -v certbot &> /dev/null; then
        return 0
    else
        return 1
    fi
}

# 安装 certbot
install_certbot() {
    log_info "正在安装 certbot..."
    
    if [ -f /etc/debian_version ]; then
        # Debian/Ubuntu
        apt-get update
        apt-get install -y certbot python3-certbot-nginx
    elif [ -f /etc/redhat-release ]; then
        # CentOS/RHEL
        yum install -y certbot python3-certbot-nginx
    else
        log_error "不支持的系统，请手动安装 certbot"
        exit 1
    fi
    
    log_success "certbot 安装完成"
}

# 创建验证目录
create_certbot_dir() {
    log_info "创建证书验证目录..."
    mkdir -p /var/www/certbot
    log_success "目录创建完成：/var/www/certbot"
}

# 创建临时 Nginx 配置
create_temp_nginx_config() {
    log_info "创建临时 Nginx 配置用于证书验证..."
    
    cat > /etc/nginx/conf.d/${DOMAIN}.conf << EOF
server {
    listen 80;
    server_name ${DOMAIN} www.${DOMAIN};
    
    # ACME challenge 验证目录
    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }
    
    # 其他请求重定向到维护页面
    location / {
        return 200 "证书申请中，请稍后...\n";
        add_header Content-Type text/plain;
    }
}
EOF
    
    log_success "临时配置创建完成：/etc/nginx/conf.d/${DOMAIN}.conf"
}

# 重新加载 Nginx
reload_nginx() {
    log_info "重新加载 Nginx 配置..."
    nginx -s reload || systemctl reload nginx
    log_success "Nginx 重新加载完成"
}

# 申请证书
request_certificate() {
    log_info "正在申请 SSL 证书..."
    
    CERTBOT_ARGS=(
        "certonly"
        "--webroot"
        "--webroot-path=/var/www/certbot"
        "--email" "$EMAIL"
        "--agree-tos"
        "--no-eff-email"
        "-d" "$DOMAIN"
        "-d" "www.$DOMAIN"
    )
    
    if [ "$PRODUCTION" != "--production" ]; then
        log_info "使用测试环境证书（限流宽松）"
        CERTBOT_ARGS+=("--staging")
    else
        log_info "使用生产环境证书"
    fi
    
    if certbot "${CERTBOT_ARGS[@]}"; then
        log_success "证书申请成功！"
        log_info "证书位置：/etc/letsencrypt/live/$DOMAIN/"
    else
        log_error "证书申请失败"
        exit 1
    fi
}

# 更新 Nginx 配置
update_nginx_config() {
    log_info "更新 Nginx 配置..."
    
    # 检查项目中的 nginx 配置文件
    NGINX_CONFIG=""
    if [ -f "nginx/nginx.conf" ]; then
        NGINX_CONFIG="nginx/nginx.conf"
    elif [ -f "../nginx/nginx.conf" ]; then
        NGINX_CONFIG="../nginx/nginx.conf"
    fi
    
    if [ -n "$NGINX_CONFIG" ]; then
        log_info "找到 Nginx 配置文件：$NGINX_CONFIG"
        log_warn "请手动修改 $NGINX_CONFIG 中的 server_name 为：$DOMAIN"
    fi
    
    # 创建配置示例
    cat > nginx-https-example.conf << EOF
server {
    listen 443 ssl http2;
    server_name ${DOMAIN} www.${DOMAIN};

    # SSL 证书配置
    ssl_certificate /etc/letsencrypt/live/${DOMAIN}/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/${DOMAIN}/privkey.pem;

    # SSL 优化
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers 'ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384';
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 1d;
    ssl_session_tickets off;

    # OCSP Stapling
    ssl_stapling on;
    ssl_stapling_verify on;

    # 其他配置...
}
EOF
    
    log_success "HTTPS 配置示例已保存到：nginx-https-example.conf"
}

# 配置自动续期
setup_auto_renewal() {
    log_info "配置证书自动续期..."
    
    # 创建续期脚本
    cat > /usr/local/bin/renew-cert-${DOMAIN}.sh << 'EOF'
#!/bin/bash
certbot renew --webroot --webroot-path=/var/www/certbot
docker-compose restart nginx
EOF
    
    chmod +x /usr/local/bin/renew-cert-${DOMAIN}.sh
    
    # 添加到 crontab
    CRON_JOB="0 2 * * * /usr/local/bin/renew-cert-${DOMAIN}.sh"
    
    # 移除旧的相同域名的 cron 任务，添加新的
    (crontab -l 2>/dev/null | grep -v "$DOMAIN" || true; echo "$CRON_JOB") | crontab -
    
    log_success "自动续期配置完成"
    log_info "续期脚本：/usr/local/bin/renew-cert-${DOMAIN}.sh"
    log_info "每天凌晨 2 点自动检查续期"
}

# 主流程
main() {
    echo -e "${CYAN}"
    echo "========================================"
    echo "HTTPS 证书配置脚本"
    echo "========================================"
    echo -e "${NC}"
    echo "域名：$DOMAIN"
    echo "邮箱：$EMAIL"
    if [ "$PRODUCTION" == "--production" ]; then
        echo "环境：生产环境"
    else
        echo "环境：测试环境"
    fi
    echo "========================================"
    echo ""
    
    # 1. 检查并安装 certbot
    if ! check_certbot; then
        log_warn "certbot 未安装"
        install_certbot
    else
        log_success "certbot 已安装"
    fi
    
    # 2. 创建验证目录
    create_certbot_dir
    
    # 3. 创建临时 Nginx 配置
    create_temp_nginx_config
    
    # 4. 重新加载 Nginx
    reload_nginx
    
    # 5. 申请证书
    request_certificate
    
    # 6. 更新 Nginx 配置
    update_nginx_config
    
    # 7. 配置自动续期
    setup_auto_renewal
    
    echo ""
    echo -e "${GREEN}========================================"
    echo "HTTPS 证书配置完成！"
    echo "========================================${NC}"
    echo ""
    echo "下一步操作："
    echo "1. 修改 nginx/nginx.conf 中的 server_name 为：${DOMAIN}"
    echo "2. 确保 docker-compose.yml 中挂载了证书目录："
    echo "   volumes:"
    echo "     - ./ssl:/etc/letsencrypt:ro"
    echo "3. 重启 Nginx："
    echo "   docker-compose restart nginx"
    echo ""
    echo "证书文件位置："
    echo "  证书：/etc/letsencrypt/live/${DOMAIN}/fullchain.pem"
    echo "  私钥：/etc/letsencrypt/live/${DOMAIN}/privkey.pem"
    echo ""
    echo "测试 HTTPS："
    echo "  curl -k https://${DOMAIN}/health"
    echo ""
    echo "========================================"
}

# 执行
main
