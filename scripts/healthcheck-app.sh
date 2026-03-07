#!/bin/bash

# 应用健康检查脚本
# 检查 Go 应用的健康状态

set -e

# 检查应用端口
if ! curl -f -s http://localhost:8888/ > /dev/null; then
    echo "Application not responding on port 8888"
    exit 1
fi

echo "Application is healthy"
exit 0
