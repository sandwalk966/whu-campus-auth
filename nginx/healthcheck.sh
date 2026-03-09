#!/bin/sh

# Nginx 健康检查脚本
# 用于 Docker healthcheck

# 检查 Nginx 进程
if ! pgrep -x "nginx" > /dev/null; then
    echo "Nginx process not running"
    exit 1
fi

# 检查端口监听（使用 ss 替代 netstat，如果不可用则使用 lsof）
if command -v ss >/dev/null 2>&1; then
    if ! ss -tln | grep -q ":80"; then
        echo "Port 80 not listening (using ss)"
        exit 1
    fi
elif command -v lsof >/dev/null 2>&1; then
    if ! lsof -i:80 | grep -q LISTEN; then
        echo "Port 80 not listening (using lsof)"
        exit 1
    fi
else
    # 如果没有可用的网络工具，尝试简单访问
    if ! curl -sf http://localhost/health > /dev/null 2>&1; then
        echo "Cannot access health endpoint"
        exit 1
    fi
fi

echo "Nginx is healthy"
exit 0
