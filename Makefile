# 社交App开发工具

.PHONY: help dev build test clean docker-up docker-down

# 默认目标
help:
	@echo "可用命令:"
	@echo "  dev          - 启动开发环境"
	@echo "  build        - 构建所有服务"
	@echo "  test         - 运行测试"
	@echo "  clean        - 清理构建文件"
	@echo "  docker-up    - 启动Docker服务"
	@echo "  docker-down  - 停止Docker服务"
	@echo "  mobile-dev   - 启动Flutter开发"
	@echo "  backend-dev  - 启动Go后端开发"

# 启动开发环境
dev: docker-up
	@echo "开发环境已启动"
	@echo "PostgreSQL: localhost:5432"
	@echo "MongoDB: localhost:27017"
	@echo "Redis: localhost:6379"

# 构建所有服务
build:
	@echo "构建Go后端服务..."
	cd backend && go mod tidy
	cd backend && go build -o bin/user-service ./services/user-service
	cd backend && go build -o bin/auth-service ./services/auth-service
	cd backend && go build -o bin/match-service ./services/match-service
	cd backend && go build -o bin/chat-service ./services/chat-service
	@echo "构建Flutter应用..."
	cd mobile && flutter pub get
	cd mobile && flutter build apk

# 运行测试
test:
	@echo "运行Go后端测试..."
	cd backend && go test -v ./...
	@echo "运行Flutter测试..."
	cd mobile && flutter test

# 清理构建文件
clean:
	@echo "清理Go构建文件..."
	cd backend && rm -rf bin/
	@echo "清理Flutter构建文件..."
	cd mobile && flutter clean

# 启动Docker服务
docker-up:
	@echo "启动Docker服务..."
	docker-compose up -d postgres mongodb redis
	@echo "等待数据库启动..."
	sleep 2

# 停止Docker服务
docker-down:
	@echo "停止Docker服务..."
	docker-compose down

## 后端

# 启动Go后端开发
backend-dev:
	@echo "启动Go后端开发模式..."
	cd backend && go mod tidy
	@echo "启动用户服务 (端口8001)..."
	cd backend && go run ./services/user-service &
	@echo "启动认证服务 (端口8002)..."
	cd backend && go run ./services/auth-service &
	@echo "启动匹配服务 (端口8003)..."
	cd backend && go run ./services/match-service &
	@echo "启动聊天服务 (端口8004)..."
	cd backend && go run ./services/chat-service &
	@echo "启动内容服务 (端口8005)..."
	cd backend && go run ./services/content-service &
	@echo "启动位置服务 (端口8006)..."
	cd backend && go run ./services/location-service &
	@echo "启动通知服务 (端口8007)..."
	cd backend && go run ./services/notification-service &
	@echo "启动管理服务 (端口8008)..."
	cd backend && go run ./services/admin-service &
	@echo "所有后端服务已启动！"

# 停止后端服务
stop-backend:
	@echo "停止所有后端服务..."
	./scripts/stop-services.sh

# 重启后端服务
restart-backend: stop-backend backend-dev

## 前端
# Flutter开发服务
flutter-start:
	@echo "启动Flutter开发服务..."
	./scripts/flutter-dev.sh start

flutter-stop:
	@echo "停止Flutter开发服务..."
	./scripts/flutter-dev.sh stop

flutter-restart:
	@echo "重启Flutter开发服务..."
	./scripts/flutter-dev.sh restart

flutter-build:
	@echo "构建Flutter应用..."
	./scripts/flutter-dev.sh build web

# Flutter调试命令
flutter-debug:
	@echo "Flutter调试工具..."
	./scripts/flutter-dev.sh debug info

flutter-diagnose:
	@echo "Flutter环境诊断..."
	./scripts/flutter-dev.sh diagnose

flutter-clean:
	@echo "深度清理Flutter..."
	./scripts/flutter-dev.sh debug clean

flutter-fix:
	@echo "修复Flutter常见问题..."
	./scripts/flutter-dev.sh fix common

# 启动Flutter开发
mobile-dev:
	@echo "启动Flutter开发模式..."
	cd mobile && flutter pub get
	cd mobile && flutter run

# 数据库迁移
migrate:
	@echo "运行数据库迁移..."
	cd backend && go run ./cmd/migrate

# 格式化代码
format:
	@echo "格式化Go代码..."
	cd backend && go fmt ./...
	@echo "格式化Flutter代码..."
	cd mobile && flutter format .

# 代码检查
lint:
	@echo "检查Go代码..."
	cd backend && golangci-lint run
	@echo "检查Flutter代码..."
	cd mobile && flutter analyze

# 安装依赖
install:
	@echo "安装Go依赖..."
	cd backend && go mod download
	@echo "安装Flutter依赖..."
	cd mobile && flutter pub get