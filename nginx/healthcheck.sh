#!/bin/sh

# Nginx 健康检查脚本
# 用于 Docker healthcheck

# 检查 Nginx 进程
if ! pgrep -x "nginx" > /dev/null; then
    echo "Nginx process not running"
    exit 1
fi

# 检查 Nginx 配置
if ! nginx -t > /dev/null 2>&1; then
    echo "Nginx configuration test failed"
    exit 1
fi

# 检查端口监听
if ! netstat -tln | grep -q ":80"; then
    echo "Port 80 not listening"
    exit 1
fi

echo "Nginx is healthy"
exit 0
