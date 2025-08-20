# 统一微服务 Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装构建依赖
RUN apk add --no-cache git ca-certificates tzdata

# 复制go mod文件
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# 复制源代码
COPY backend/ ./

# 构建所有微服务
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/user-service ./services/user-service && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/chat-service ./services/chat-service && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/moment-service ./services/moment-service && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/upload-service ./services/upload-service && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/news-service ./services/news-service

# 运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata curl supervisor nginx

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/bin/ ./bin/
COPY --from=builder /app/shared ./shared/

# 复制配置文件
COPY docker/supervisord.conf /etc/supervisord.conf
COPY docker/nginx.conf /etc/nginx/nginx.conf

# 创建必要的目录
RUN mkdir -p /app/logs /app/uploads /var/log/supervisor /run/nginx

# 创建非root用户
RUN addgroup -g 1001 appgroup && \
    adduser -u 1001 -G appgroup -s /bin/sh -D appuser && \
    chown -R appuser:appgroup /app /var/log/supervisor

# 暴露端口
EXPOSE 80 8001 8004 8009 8010 8011

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
  CMD curl -f http://localhost/health || exit 1

# 启动supervisor管理所有服务
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]