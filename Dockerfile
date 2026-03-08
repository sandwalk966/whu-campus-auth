FROM golang:latest

# 安装时区数据（Debian 基础镜像）
RUN apt-get update && apt-get install -y tzdata && rm -rf /var/lib/apt/lists/*

# 设置时区
ENV TZ=Asia/Shanghai

# 设置工作目录
WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# 创建 uploads 目录
RUN mkdir -p /app/uploads

# 暴露端口（内部端口，不直接暴露给宿主机）
EXPOSE 8888

# 启动应用
CMD ["./main"]