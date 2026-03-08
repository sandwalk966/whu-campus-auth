# 多阶段构建：先构建前端，再部署到 Nginx
FROM node:iron-trixie-slim AS builder

WORKDIR /app

# 复制 package.json
COPY ../frontend/package*.json ./

# 安装依赖
RUN npm install --registry=https://registry.npmmirror.com

# 复制前端源代码
COPY ./frontend .

# 构建生产版本
RUN npm run build

# 最终阶段：Nginx
FROM nginx:latest

# 复制构建产物到 Nginx
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制 Nginx 配置
COPY nginx/nginx.conf /etc/nginx/conf.d/default.conf

# 复制健康检查脚本
COPY nginx/healthcheck.sh /etc/nginx/healthcheck.sh

EXPOSE 80 443

CMD ["nginx", "-g", "daemon off;"]
