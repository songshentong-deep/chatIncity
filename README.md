# 兴趣匹配+地理位置社交App

基于兴趣匹配和地理位置的即时社交App，通过智能算法推荐高契合度用户，主打"快速破冰"和"安全交友"。

## 技术架构

- **移动端**: Flutter + Dart
- **后端**: Go + Gin框架
- **数据库**: PostgreSQL + MongoDB + Redis
- **部署**: Docker + Kubernetes

## 项目结构

```
social-app/
├── backend/                 # Go微服务后端
│   ├── services/           # 各个微服务
│   │   ├── user-service/   # 用户服务
│   │   ├── auth-service/   # 认证服务
│   │   ├── match-service/  # 匹配服务
│   │   ├── chat-service/   # 聊天服务
│   │   ├── content-service/# 内容服务
│   │   ├── location-service/# 位置服务
│   │   ├── notification-service/# 通知服务
│   │   └── admin-service/  # 管理服务
│   ├── shared/             # 共享代码
│   └── docker/             # Docker配置
├── mobile/                 # Flutter移动端
│   ├── lib/                # Dart源码
│   ├── android/            # Android配置
│   └── ios/                # iOS配置
├── docker-compose.yml      # 开发环境配置
└── k8s/                    # Kubernetes部署配置
```

## 快速开始

### 开发环境要求

- Go 1.21+
- Flutter 3.16+
- Docker & Docker Compose
- PostgreSQL 15+
- MongoDB 6.0+
- Redis 7.0+

### 启动开发环境

```bash
# 启动数据库服务
docker-compose up -d postgres mongodb redis

# 启动后端服务
cd backend
go mod tidy
make dev

# 启动移动端
cd mobile
flutter pub get
flutter run
```

## 开发指南

详细的开发文档请参考 `.kiro/specs/interest-location-social-app/` 目录下的规格说明文档。