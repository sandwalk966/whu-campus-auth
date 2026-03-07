#!/bin/bash

# Docker Compose CORS 配置脚本
# 使用方法：
#   ./configure-cors.sh production example.com
#   ./configure-cors.sh development

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

show_usage() {
    echo "使用方法："
    echo "  $0 <环境> [域名]"
    echo ""
    echo "环境："
    echo "  production   - 生产环境（需要指定域名）"
    echo "  development  - 开发环境（使用 localhost）"
    echo ""
    echo "示例："
    echo "  $0 production example.com"
    echo "  $0 production yourdomain.com www.yourdomain.com"
    echo "  $0 development"
    exit 1
}

if [ "$#" -lt 1 ]; then
    show_usage
fi

ENVIRONMENT=$1
shift

DOCKER_COMPOSE_FILE="docker-compose.yml"
ENV_FILE=".env"

if [ ! -f "$DOCKER_COMPOSE_FILE" ]; then
    log_error "未找到 docker-compose.yml 文件"
    exit 1
fi

if [ "$ENVIRONMENT" == "production" ]; then
    if [ "$#" -lt 1 ]; then
        log_error "生产环境必须指定域名"
        show_usage
    fi
    
    # 构建域名列表
    DOMAINS=""
    for domain in "$@"; do
        if [ -z "$DOMAINS" ]; then
            DOMAINS="https://$domain"
        else
            DOMAINS="$DOMAINS,https://$domain"
        fi
    done
    
    log_info "配置生产环境 CORS"
    log_info "允许的域名：$DOMAINS"
    
    # 更新 docker-compose.yml
    sed -i.bak "s|- ALLOWED_ORIGINS=.*|- ALLOWED_ORIGINS=$DOMAINS|g" "$DOCKER_COMPOSE_FILE"
    rm -f "$DOCKER_COMPOSE_FILE.bak"
    
    # 更新 .env 文件（如果存在）
    if [ -f "$ENV_FILE" ]; then
        sed -i.bak "s|^ALLOWED_ORIGINS=.*|ALLOWED_ORIGINS=$DOMAINS|g" "$ENV_FILE"
        rm -f "$ENV_FILE.bak"
    fi
    
    log_info "✅ 生产环境配置完成！"
    log_info "重启服务：docker-compose restart app"
    
elif [ "$ENVIRONMENT" == "development" ]; then
    DEV_ORIGINS="http://localhost:3000,http://localhost:5173,http://127.0.0.1:3000,http://127.0.0.1:5173"
    
    log_info "配置开发环境 CORS"
    log_info "允许的域名：$DEV_ORIGINS"
    
    # 更新 docker-compose.yml
    sed -i.bak "s|- ALLOWED_ORIGINS=.*|- ALLOWED_ORIGINS=$DEV_ORIGINS|g" "$DOCKER_COMPOSE_FILE"
    rm -f "$DOCKER_COMPOSE_FILE.bak"
    
    # 更新 .env 文件（如果存在）
    if [ -f "$ENV_FILE" ]; then
        sed -i.bak "s|^ALLOWED_ORIGINS=.*|ALLOWED_ORIGINS=$DEV_ORIGINS|g" "$ENV_FILE"
        rm -f "$ENV_FILE.bak"
    fi
    
    log_info "✅ 开发环境配置完成！"
    log_info "重启服务：docker-compose restart app"
    
else
    log_error "未知的环境：$ENVIRONMENT"
    show_usage
fi

echo ""
echo "当前配置："
grep "ALLOWED_ORIGINS" "$DOCKER_COMPOSE_FILE"
