# 动态服务 Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git ca-certificates tzdata

# 复制go mod文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o moment-service ./services/moment-service

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata curl
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/moment-service .
COPY --from=builder /app/shared ./shared

# 创建非root用户
RUN addgroup -g 1001 appgroup && adduser -u 1001 -G appgroup -s /bin/sh -D appuser
RUN chown -R appuser:appgroup /app
USER appuser

EXPOSE 8010

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8010/health || exit 1

CMD ["./moment-service"]