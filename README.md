# 007Secret

一个安全的一次性密钥分享平台。用户可以创建并分享敏感信息，生成临时链接和提取码。接收方通过链接和提取码只能查看一次内容，查看后内容将立即销毁。

## 特性

- 一次性查看：内容只能被查看一次，查看后自动销毁
- 数据压缩：所有存储的数据都会被压缩以节省空间
- 安全性：使用 Redis 存储，数据在内存中，且设置过期时间（意味着24小时未查看也会自动销毁）
- 简单易用：简洁的用户界面，容易操作

## 技术栈

- 前端：Vue.js
- 后端：Go
- 数据存储：Redis

## 启动方法

### 方式1：Docker Compose（推荐）
1. 进入 backend 目录
2. 复制 .env.example 为 .env 并配置环境变量，并将`REDIS_HOST=localhost`改为`REDIS_HOST=redis`
3. 回到项目根目录，执行`docker-compose up -d`

### 方式2：下载源码
`git clone https://github.com/007Secret/007Secret.git`

### 后端
1. 确保已安装 Go 1.16+ 和 Redis
2. 进入 backend 目录
3. 复制 `.env.example` 为 `.env` 并配置环境变量
4. 运行 `go run main.go`

### 前端
1. 进入 frontend 目录
2. 安装依赖：`npm install`
3. 启动开发服务器：`npm run dev` 

## 贡献
欢迎贡献代码和建议。请在 GitHub 上创建 issue 或 pull request。