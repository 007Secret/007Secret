# 007Secret

一个一次性文本与附件分享平台。用户可以分享密码、密钥、配置说明等内容，生成临时链接和提取码。文本成功提取后从服务器删除；附件仅通过 WebRTC 在浏览器之间点对点传输，接收完成后由用户手动下载。

## 特性

- 一次性查看：文本只能成功提取一次，提取后自动从服务器删除
- 加密存储：文本与附件元数据加密后存入 Redis，默认 24 小时过期
- 数据压缩：存储内容经过压缩处理
- 简单易用：简洁的用户界面，容易操作
- 附件直传：可附加单个不超过 10 MB 的文件；文件仅通过 WebRTC 在浏览器之间点对点传输，不会由服务器中转或暂存。发送方与接收方必须同时在线并保持页面打开。
- 文件大小提示：选择附件后显示大小，超过 10 MB 时立即提示并阻止生成链接；秘密内容为空时自动填入文件名，不覆盖已输入的文本
- 独立提取文本：WebRTC 连接等待或失败时仍可查看文本，附件区域单独显示连接状态并支持重试
- 手动下载附件：附件接收完成后显示“下载文件”按钮，用户点击后保存，不会在提取文本时自动触发下载

## 使用方法

1. 输入秘密内容，可选择一个附件。文本为空时，选择附件会自动将文件名填入秘密内容。
2. 确认附件名称与大小后生成分享链接，将链接与提取码发送给接收者。附件上限为 `10 × 1024 × 1024` 字节，界面显示为 10 MB。
3. 接收者输入提取码并点击“解密并查看”，文本立即显示，同时从服务器删除；有附件时，页面会另外建立 WebRTC 连接并接收文件。
4. 附件接收完成后，点击“下载文件”保存。若显示“WebRTC 未连接”，仍可读取或复制文本，并通过“重新尝试连接”单独重试附件，无需再次提取文本。

传输附件时，发送方必须保留创建分享的原始页面，接收方也需保持当前页面打开。提取后的文本、附件下载链接和重试状态保留在当前页面中；刷新或关闭页面后无法再次提取文本，请先复制需要的内容并下载附件。文件接收完成后，可在当前页面内再次点击下载。

## 效果图

<img width="1469" height="1180" alt="image" src="https://github.com/user-attachments/assets/e0fc9370-dffd-4b34-984f-698bbb077d18" />
<img width="1578" height="1048" alt="image" src="https://github.com/user-attachments/assets/2323d23e-971d-4628-990d-b0c43982183f" />
<img width="1475" height="1100" alt="image" src="https://github.com/user-attachments/assets/264df738-7e84-4730-bfa9-f7dc89c49378" />
<img width="1244" height="977" alt="image" src="https://github.com/user-attachments/assets/077d1596-bc1d-4c2a-b643-8d671a65bfe2" />


## 技术栈

- 前端：Vue 3、TypeScript、Vite、Tailwind CSS
- 后端：Go 1.22+、Gin、WebSocket 信令服务
- 数据存储：Redis
- 附件传输：WebRTC DataChannel

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

当前 Compose 使用固定镜像标签 `redis:8.2-alpine`，并设置 `pull_policy: missing`：本地已有该标签的镜像时直接使用，仅在缺少时拉取。镜像标签和拉取策略见 [docker-compose.yml](docker-compose.yml)。

如需临时覆盖某项配置，用环境变量在运行期覆盖（不改文件）：`docker compose run -e REDIS_PASSWORD=xxx ...`，或在 `docker-compose.yml` 的 `backend` 服务下补一条 `environment`。

### 方式2：手动安装

#### 后端

1. 确保已安装 Go 1.22+ 和 Redis
2. 进入 backend 目录
3. 复制 `.env.example` 为 `.env`，按需修改（默认 `REDIS_HOST=localhost`、`PORT=8080`）
4. 运行 `go run .`

#### 前端

1. 进入 frontend 目录
2. 安装依赖：`npm install`
3. 启动开发服务器：`npm run dev`（默认 `http://localhost:5173`，`/api` 已代理到 `http://localhost:8080`）

## 开发检查

前端在 `frontend` 目录执行：

```sh
npm run test:unit -- --run
npm run type-check
npm run build
```

后端在 `backend` 目录执行：

```sh
go test ./...
```

## 贡献

欢迎贡献代码和建议。请在 GitHub 上创建 issue 或 pull request。
