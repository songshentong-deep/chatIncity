# 服务管理指南

## 🚀 快速启动

### 一键启动所有服务
```bash
# 启动数据库 + 后端 + 前端
make dev                    # 启动数据库
make start-backend         # 启动后端服务
make flutter-start         # 启动前端服务
```

### 一键停止所有服务
```bash
make stop-backend          # 停止后端服务
make flutter-stop          # 停止前端服务
make docker-down           # 停止数据库服务
```

## 🛠️ 后端服务管理

### 使用脚本管理
```bash
# 启动所有后端服务
./scripts/start-services.sh

# 停止所有后端服务
./scripts/stop-services.sh

# 查看服务状态
ps aux | grep -E "(user-service|auth-service)" | grep -v grep
```

### 使用Makefile管理
```bash
# 启动后端服务
make start-backend

# 停止后端服务
make stop-backend

# 重启后端服务
make restart-backend
```

### 单独管理服务
```bash
# 启动单个服务
cd backend && go run ./services/user-service

# 查看服务日志
tail -f backend/logs/user-service.log

# 检查端口占用
lsof -i :8001,8002,8003,8004,8005,8006,8007,8008
```

## 📱 前端服务管理

### 使用脚本管理
```bash
# 启动Flutter开发服务
./scripts/flutter-dev.sh start chrome

# 停止Flutter服务
./scripts/flutter-dev.sh stop

# 重启Flutter服务
./scripts/flutter-dev.sh restart

# 构建应用
./scripts/flutter-dev.sh build web
```

### 使用Makefile管理
```bash
# 启动Flutter服务
make flutter-start

# 停止Flutter服务
make flutter-stop

# 重启Flutter服务
make flutter-restart

# 构建Flutter应用
make flutter-build
```

### Flutter开发命令
```bash
# 查看可用设备
./scripts/flutter-dev.sh devices

# 运行测试
./scripts/flutter-dev.sh test

# 代码分析
./scripts/flutter-dev.sh analyze

# 代码格式化
./scripts/flutter-dev.sh format
```

## 🗄️ 数据库服务管理

### Docker数据库服务
```bash
# 启动数据库服务
make docker-up
docker-compose up -d postgres mongodb redis

# 停止数据库服务
make docker-down
docker-compose down

# 查看数据库状态
docker-compose ps

# 查看数据库日志
docker-compose logs postgres
docker-compose logs mongodb
docker-compose logs redis
```

### 数据库连接
```bash
# 连接PostgreSQL
docker exec -it social-app-postgres psql -U chat -d social_app

# 连接MongoDB
docker exec -it social-app-mongodb mongosh -u admin -p password

# 连接Redis
docker exec -it social-app-redis redis-cli -a password
```

## 📊 服务状态检查

### 后端服务状态
```bash
# 检查所有后端服务端口
echo "后端服务状态:" && for port in 8001 8002 8003 8004 8005 8006 8007 8008; do echo -n "端口 $port: "; lsof -i :$port 2>/dev/null | grep LISTEN | awk '{print $1}' || echo "未运行"; done

# 检查服务进程
ps aux | grep -E "(user-service|auth-service|match-service|chat-service|content-service|location-service|notification-service|admin-service)" | grep -v grep

# 检查Go进程
ps aux | grep "go run" | grep -v grep
```

### 前端服务状态
```bash
# 检查Flutter进程
ps aux | grep flutter | grep -v grep

# 检查Web服务端口
lsof -i :3000  # Flutter Web默认端口
```

### 数据库服务状态
```bash
# 检查Docker容器
docker ps

# 检查数据库端口
lsof -i :5432,27017,6379
```

## 🔧 故障排查

### 后端服务问题

#### 端口被占用
```bash
# 查找占用端口的进程
lsof -ti :8001

# 强制释放端口
kill -9 $(lsof -ti :8001)
```

#### 服务启动失败
```bash
# 查看服务日志
tail -f backend/logs/user-service.log

# 手动启动服务查看错误
cd backend && go run ./services/user-service
```

#### 数据库连接失败
```bash
# 检查数据库是否运行
docker ps | grep postgres

# 测试数据库连接
cd backend && go run test_db.go
```

### 前端服务问题

#### Flutter启动失败
```bash
# 清理Flutter缓存
cd mobile && flutter clean

# 重新获取依赖
cd mobile && flutter pub get

# 检查Flutter环境
flutter doctor
```

#### 设备连接问题
```bash
# 查看可用设备
flutter devices

# 启动Android模拟器
flutter emulators --launch Pixel_4

# 启动iOS模拟器
open -a Simulator
```

### 数据库服务问题

#### Docker启动失败
```bash
# 检查Docker状态
docker info

# 重启Docker服务
# macOS: 重启Docker Desktop
# Linux: sudo systemctl restart docker
```

#### 数据库初始化失败
```bash
# 删除数据卷重新初始化
docker-compose down -v
docker-compose up -d postgres

# 手动运行数据库迁移
cd backend && go run ./cmd/migrate
```

## 📋 常用命令速查

### 开发环境启动
```bash
# 完整启动流程
make docker-up          # 1. 启动数据库
make start-backend      # 2. 启动后端
make flutter-start      # 3. 启动前端
```

### 开发环境停止
```bash
# 完整停止流程
make flutter-stop       # 1. 停止前端
make stop-backend       # 2. 停止后端
make docker-down        # 3. 停止数据库
```

### 服务重启
```bash
# 重启后端服务
make restart-backend

# 重启前端服务
make flutter-restart

# 重启数据库服务
make docker-down && make docker-up
```

### 日志查看
```bash
# 后端服务日志
tail -f backend/logs/user-service.log

# 数据库日志
docker-compose logs -f postgres

# Flutter日志
# 在Flutter运行的终端中查看
```

## 🎯 开发工作流

### 日常开发
1. 启动数据库: `make docker-up`
2. 启动后端: `make start-backend`
3. 启动前端: `make flutter-start`
4. 开始开发...
5. 停止服务: `make stop-backend && make flutter-stop`

### 代码更新后
1. 停止相关服务: `make stop-backend`
2. 重新启动: `make start-backend`
3. 如果前端有更新，热重载会自动生效

### 数据库更新后
1. 停止后端: `make stop-backend`
2. 运行迁移: `cd backend && go run ./cmd/migrate`
3. 重启后端: `make start-backend`

---

## 📞 获取帮助

如果遇到问题：
1. 查看相应的日志文件
2. 运行 `./scripts/flutter-dev.sh help` 查看Flutter命令帮助
3. 检查 `COMMANDS.md` 文档
4. 查看 `DOCKER_SETUP.md` 数据库配置文档