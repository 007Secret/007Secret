# 使用官方 Golang 镜像作为构建环境
FROM golang:1.22.1-alpine3.19 AS builder

# 设置工作目录
WORKDIR /app

# 复制 go mod 和 sum 文件
COPY go.* ./

# 下载依赖
RUN go mod download

# 复制源代码和配置文件
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 使用特定版本的 alpine 作为运行环境
FROM alpine:3.19

WORKDIR /app

# 从构建阶段复制二进制文件和配置文件
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"] 