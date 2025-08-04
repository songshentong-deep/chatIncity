# 社交App常用命令文档

## 📋 目录
- [环境检查](#环境检查)
- [服务启动](#服务启动)
- [编译构建](#编译构建)
- [开发调试](#开发调试)
- [服务管理](#服务管理)
- [数据库操作](#数据库操作)
- [代码质量](#代码质量)
- [故障排查](#故障排查)

## 🔍 环境检查

### 检查开发环境
```bash
# 运行环境检查脚本
./scripts/check-env.sh

# 手动检查各组件版本
go version                    # 检查Go版本 (需要1.21+)
flutter --version            # 检查Flutter版本 (需要3.16+)
docker --version             # 检查Docker版本
docker-compose --version     # 检查Docker Compose版本
```

### 检查端口占用
```bash
# 检查关键端口是否被占用
lsof -i :5432    # PostgreSQL
lsof -i :27017   # MongoDB
lsof -i :6379    # Redis
lsof -i :8001    # 用户服务
lsof -i :8002    # 认证服务
lsof -i :8003    # 匹配服务
lsof -i :8004    # 聊天服务
```

## 🚀 服务启动

### 快速启动开发环境
```bash
# 一键启动完整开发环境
make dev

# 分步启动
make docker-up      # 启动数据库服务
make backend-dev    # 启动后端服务
make mobile-dev     # 启动移动端应用
```

### 启动数据库服务
```bash
# 启动所有数据库服务
docker-compose up -d postgres mongodb redis

# 启动单个数据库服务
docker-compose up -d postgres    # 仅启动PostgreSQL
docker-compose up -d mongodb     # 仅启动MongoDB
docker-compose up -d redis       # 仅启动Redis
```

### 启动后端微服务
```bash
# 启动所有后端服务 (后台运行)
make backend-dev

# 手动启动单个服务
cd backend
go run ./services/user-service         # 用户服务 (8001)
go run ./services/auth-service         # 认证服务 (8002)
go run ./services/match-service        # 匹配服务 (8003)
go run ./services/chat-service         # 聊天服务 (8004)
go run ./services/content-service      # 内容服务 (8005)
go run ./services/location-service     # 位置服务 (8006)
go run ./services/notification-service # 通知服务 (8007)
go run ./services/admin-service        # 管理服务 (8008)
```

### 启动移动端应用
```bash
# 启动Flutter开发模式
cd mobile
flutter pub get     # 安装依赖
flutter run         # 启动应用

# 指定设备启动
flutter devices     # 查看可用设备
flutter run -d chrome              # 在Chrome中运行
flutter run -d android             # 在Android设备中运行
flutter run -d ios                 # 在iOS设备中运行
```

## 🔨 编译构建

### 构建所有服务
```bash
# 一键构建所有服务
make build

# 构建后端服务
make build-backend

# 构建移动端应用
make build-mobile
```

### 构建后端服务
```bash
cd backend

# 构建所有微服务
go mod tidy
go build -o bin/user-service ./services/user-service
go build -o bin/auth-service ./services/auth-service
go build -o bin/match-service ./services/match-service
go build -o bin/chat-service ./services/chat-service
go build -o bin/content-service ./services/content-service
go build -o bin/location-service ./services/location-service
go build -o bin/notification-service ./services/notification-service
go build -o bin/admin-service ./services/admin-service

# 构建单个服务
go build -o bin/user-service ./services/user-service
```

### 构建移动端应用
```bash
cd mobile

# 安装依赖
flutter pub get

# 构建APK (Android)
flutter build apk                    # 构建release APK
flutter build apk --debug           # 构建debug APK
flutter build apk --split-per-abi   # 按架构分别构建

# 构建iOS应用
flutter build ios                   # 构建iOS应用
flutter build ipa                   # 构建IPA文件

# 构建Web应用
flutter build web                   # 构建Web版本
```

## 🛠 开发调试

### 热重载开发
```bash
# Flutter热重载 (开发过程中按 'r' 热重载，按 'R' 热重启)
cd mobile && flutter run

# Go服务热重载 (需要安装air)
go install github.com/cosmtrek/air@latest
cd backend && air
```

### 代码生成
```bash
# 生成所有代码
make generate

# 生成Go代码
cd backend && go generate ./...

# 生成Flutter代码
cd mobile && flutter packages pub run build_runner build
cd mobile && flutter packages pub run build_runner build --delete-conflicting-outputs
```

### 运行测试
```bash
# 运行所有测试
make test

# 运行后端测试
cd backend && go test -v ./...
cd backend && go test -v ./services/user-service/...

# 运行移动端测试
cd mobile && flutter test
cd mobile && flutter test test/unit/
cd mobile && flutter test integration_test/
```

## 📊 服务管理

### 查看运行中的服务
```bash
# 查看Docker容器状态
docker-compose ps
docker ps

# 查看Go进程
ps aux | grep "go run"
pgrep -f "services/"

# 查看端口占用
netstat -tulpn | grep LISTEN
lsof -i -P -n | grep LISTEN
```

### 停止服务
```bash
# 停止所有Docker服务
make docker-down
docker-compose down

# 停止特定Docker服务
docker-compose stop postgres
docker-compose stop mongodb
docker-compose stop redis

# 停止Go后端服务
pkill -f "services/"
pkill -f "go run"

# 停止Flutter应用
# 在运行Flutter的终端按 Ctrl+C
```

### 重启服务
```bash
# 重启Docker服务
docker-compose restart postgres
docker-compose restart mongodb
docker-compose restart redis

# 重启所有服务
make docker-down && make dev
```

## 🗄 数据库操作

### 数据库连接
```bash
# 连接PostgreSQL
docker exec -it social-app-postgres psql -U postgres -d social_app

# 连接MongoDB
docker exec -it social-app-mongodb mongosh -u admin -p password

# 连接Redis
docker exec -it social-app-redis redis-cli -a password
```

### 数据库迁移
```bash
# 运行数据库迁移
make migrate
cd backend && go run ./cmd/migrate

# 查看迁移状态
cd backend && go run ./cmd/migrate -status

# 回滚迁移
cd backend && go run ./cmd/migrate -down
```

### 数据库备份与恢复
```bash
# 备份PostgreSQL
docker exec social-app-postgres pg_dump -U postgres social_app > backup.sql

# 恢复PostgreSQL
docker exec -i social-app-postgres psql -U postgres social_app < backup.sql

# 备份MongoDB
docker exec social-app-mongodb mongodump --uri="mongodb://admin:password@localhost:27017/social_app"

# 恢复MongoDB
docker exec social-app-mongodb mongorestore --uri="mongodb://admin:password@localhost:27017/social_app" dump/
```

## 🔍 代码质量

### 代码格式化
```bash
# 格式化所有代码
make format

# 格式化Go代码
cd backend && go fmt ./...
cd backend && goimports -w .

# 格式化Flutter代码
cd mobile && flutter format .
cd mobile && dart format .
```

### 代码检查
```bash
# 检查所有代码
make lint

# 检查Go代码
cd backend && golangci-lint run
cd backend && go vet ./...

# 检查Flutter代码
cd mobile && flutter analyze
cd mobile && dart analyze
```

### 依赖管理
```bash
# 安装所有依赖
make install

# 更新Go依赖
cd backend && go mod tidy
cd backend && go mod download
cd backend && go get -u ./...

# 更新Flutter依赖
cd mobile && flutter pub get
cd mobile && flutter pub upgrade
cd mobile && flutter pub outdated
```

## 🚨 故障排查

### 查看日志
```bash
# 查看Docker容器日志
docker-compose logs postgres
docker-compose logs mongodb
docker-compose logs redis
docker-compose logs -f --tail=100 postgres  # 实时查看最近100行

# 查看应用日志
cd backend && go run ./services/user-service 2>&1 | tee logs/user-service.log
cd mobile && flutter run --verbose
```

### 清理缓存
```bash
# 清理构建缓存
make clean

# 清理Go缓存
cd backend && go clean -cache -modcache -testcache
rm -rf backend/bin/

# 清理Flutter缓存
cd mobile && flutter clean
cd mobile && flutter pub cache repair
rm -rf mobile/build/
```

### 重置开发环境
```bash
# 完全重置环境
make docker-down
docker system prune -a
make clean
make install
make dev
```

### 常见问题解决
```bash
# 端口被占用
sudo lsof -ti:8001 | xargs kill -9  # 强制关闭占用8001端口的进程

# Docker容器无法启动
docker-compose down
docker system prune
docker-compose up -d

# Go模块问题
cd backend && rm go.sum && go mod tidy

# Flutter依赖问题
cd mobile && flutter clean && flutter pub get
```

## 📱 移动端特定命令

### 设备管理
```bash
# 查看连接的设备
flutter devices

# 启动Android模拟器
flutter emulators                    # 查看可用模拟器
flutter emulators --launch Pixel_4  # 启动指定模拟器

# iOS模拟器
open -a Simulator                    # 打开iOS模拟器
xcrun simctl list devices           # 查看iOS设备列表
```

### 性能分析
```bash
# 性能分析
flutter run --profile               # 性能模式运行
flutter run --release              # 发布模式运行

# 构建分析
flutter build apk --analyze-size   # 分析APK大小
flutter build appbundle --analyze-size  # 分析AAB大小
```

## 🔧 实用工具命令

### 快速命令别名 (可添加到 ~/.zshrc)
```bash
# 添加到 ~/.zshrc 文件中
alias social-dev="cd /path/to/social-app && make dev"
alias social-backend="cd /path/to/social-app && make backend-dev"
alias social-mobile="cd /path/to/social-app/mobile && flutter run"
alias social-logs="cd /path/to/social-app && docker-compose logs -f"
alias social-clean="cd /path/to/social-app && make clean"
```

### 监控脚本
```bash
# 创建服务监控脚本
cat > monitor.sh << 'EOF'
#!/bin/bash
while true; do
    echo "=== $(date) ==="
    echo "Docker服务状态:"
    docker-compose ps
    echo "端口占用情况:"
    lsof -i :8001,8002,8003,8004 | grep LISTEN
    echo "内存使用:"
    docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}"
    sleep 30
done
EOF
chmod +x monitor.sh
```

---

## 📞 获取帮助

如果遇到问题，可以：
1. 运行 `make help` 查看可用命令
2. 运行 `./scripts/check-env.sh` 检查环境配置
3. 查看项目 README.md 文档
4. 检查 `.kiro/specs/` 目录下的详细规格文档

---

*最后更新: 2025年1月*