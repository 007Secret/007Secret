# 使用官方 Golang 镜像作为构建环境
FROM golang:1.22.1-alpine3.19 AS builder

# 设置工作目录
WORKDIR /app

# 复制 go mod 和 sum 文件
COPY go.* ./

# 下载依赖
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download

# 复制源代码和配置文件
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 使用特定版本的 alpine 作为运行环境
FROM alpine:3.19

WORKDIR /app

# 只复制二进制文件：镜像中不包含 .env
# 容器内配置的唯一来源就是下面的 ENV（docker-compose.yml 不再重复声明）
COPY --from=builder /app/main .

ENV GIN_MODE=release \
    REDIS_HOST=redis \
    REDIS_PORT=6379 \
    REDIS_PASSWORD= \
    REDIS_DB=0 \
    PORT=8080

# 容器内不读取 .env（config.Init 依据此开关判断），配置只能来自环境变量
ENV LOAD_DOTENV=false

# 暴露端口（如需改端口，改上面的 PORT 与这里、以及前端 BACKEND_PORT 保持一致）
EXPOSE 8080

# 运行应用
CMD ["./main"]
