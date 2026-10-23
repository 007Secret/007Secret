# 007Secret

一个安全的一次性密钥分享平台。用户可以创建并分享敏感信息，生成临时链接和提取码。接收方通过链接和提取码只能查看一次内容，查看后内容将立即销毁。

## 特性

- 一次性查看：内容只能被查看一次，查看后自动销毁
- 安全性：对数据进行加密，使用 Redis 存储，数据在内存中，且设置过期时间
- 数据压缩：所有存储的数据都会被压缩以节省空间
- 简单易用：简洁的用户界面，容易操作
- 附件直传：可附加单个不超过 10 MB 的文件；文件仅通过 WebRTC 在浏览器之间点对点传输，不会由服务器中转或暂存。发送方与接收方必须同时在线并保持页面打开。

## 效果图
<img width="1469" height="1180" alt="image" src="https://github.com/user-attachments/assets/e0fc9370-dffd-4b34-984f-698bbb077d18" />
<img width="1578" height="1048" alt="image" src="https://github.com/user-attachments/assets/2323d23e-971d-4628-990d-b0c43982183f" />
<img width="1475" height="1100" alt="image" src="https://github.com/user-attachments/assets/264df738-7e84-4730-bfa9-f7dc89c49378" />
<img width="1244" height="977" alt="image" src="https://github.com/user-attachments/assets/077d1596-bc1d-4c2a-b643-8d671a65bfe2" />


## 技术栈

- 前端：Vue.js
- 后端：Go
- 数据存储：Redis

## 配置

| 环境变量 | 说明 | 默认值 |
| --- | --- | --- |
| `REDIS_HOST` | Redis 主机名 | `localhost` |
| `REDIS_PORT` | Redis 端口 | `6379` |
| `REDIS_PASSWORD` | Redis 密码 | 空 |
| `REDIS_DB` | Redis 数据库编号（非数字会导致启动失败） | `0` |
| `PORT` | 后端 HTTP 端口 | `8080` |
| `GIN_MODE` | `release` 关闭 Gin 调试日志（需配合 `gin.SetMode`，已在 main 中调用） | 未设置（debug） |
| `LOAD_DOTENV` | 设为 `false` 时不读取 `.env` | `true` |

**两种启动方式，配置来源不同：**

- **Docker**：唯一来源是镜像的 `backend/Dockerfile` 中的 `ENV`（`GIN_MODE=release`、`REDIS_HOST=redis`、`PORT=8080` 等），镜像里同时设置了 `LOAD_DOTENV=false`，因此**容器内不会读取 `.env`**。`docker-compose.yml` 不再重复声明这些变量。
- **手动启动**：`backend/.env` 仅用于本地开发，从 `backend/.env.example` 复制即可，程序启动时自动加载。

> 修改容器端口时需同步三处：`backend/Dockerfile` 的 `PORT` 与 `EXPOSE`、`frontend/Dockerfile` 的 `BACKEND_PORT`。

## 启动方法

### 方式1：Docker Compose（推荐）
配置全部在镜像 Dockerfile 中，**无需创建任何 `.env` 文件**：

1. 在项目根目录执行 `docker compose up -d --build`
2. 访问 `http://localhost:8087`

如需临时覆盖某项配置，用环境变量在运行期覆盖（不改文件）：`docker compose run -e REDIS_PASSWORD=xxx ...`，或在 `docker-compose.yml` 的 `backend` 服务下补一条 `environment`。

### 方式2：手动安装
#### 后端
1. 确保已安装 Go 1.22+ 和 Redis
2. 进入 backend 目录
3. 复制 `backend/.env.example` 为 `backend/.env`，按需修改（默认 `REDIS_HOST=localhost`、`PORT=8080`）
4. 运行 `go run .`

#### 前端
1. 进入 frontend 目录
2. 安装依赖：`npm install`
3. 启动开发服务器：`npm run dev`（默认 `http://localhost:5173`，`/api` 已代理到 `http://localhost:8080`）

## 贡献
欢迎贡献代码和建议。请在 GitHub 上创建 issue 或 pull request。
