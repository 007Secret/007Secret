# 构建阶段
FROM node:20.11.1-alpine3.19 as build-stage
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# 生产阶段
FROM nginx:1.25.4-alpine
COPY --from=build-stage /app/dist /usr/share/nginx/html

# 后端地址配置：镜像内 ENV 即唯一来源（compose 不再重复声明）
ENV BACKEND_HOST=backend \
    BACKEND_PORT=8080

# 使用官方镜像的 envsubst 模板机制：容器启动时 20-envsubst-on-templates.sh
# 会读取 .template 文件，把 ${BACKEND_HOST} / ${BACKEND_PORT} 替换为环境变量值，
# 输出到 /etc/nginx/conf.d/default.conf（此目录是官方镜像的默认输出目录）。
COPY nginx.conf /etc/nginx/conf.d/default.conf.template
# 官方镜像自带的 default.conf 会与渲染产物同名，必须删掉
RUN rm -f /etc/nginx/conf.d/default.conf

# 只替换这两个变量，避免误伤 nginx 自身的 $uri / $host 等变量
ENV NGINX_ENVSUBST_FILTER='^(BACKEND_HOST|BACKEND_PORT)$'

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
