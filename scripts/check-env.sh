#!/bin/bash

# 社交App开发环境检查脚本

echo "🔍 检查开发环境..."

# 检查Go环境
echo "📦 检查Go环境..."
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | cut -d' ' -f3)
    echo "✅ Go已安装: $GO_VERSION"
else
    echo "❌ Go未安装，请安装Go 1.21+"
    exit 1
fi

# 检查Flutter环境
echo "📱 检查Flutter环境..."
if command -v flutter &> /dev/null; then
    FLUTTER_VERSION=$(flutter --version | head -n1 | cut -d' ' -f2)
    echo "✅ Flutter已安装: $FLUTTER_VERSION"
    
    # 检查Flutter doctor
    echo "🔧 运行Flutter doctor..."
    flutter doctor --android-licenses > /dev/null 2>&1
    flutter doctor
else
    echo "❌ Flutter未安装，请安装Flutter 3.16+"
    exit 1
fi

# 检查Docker环境
echo "🐳 检查Docker环境..."
if command -v docker &> /dev/null; then
    DOCKER_VERSION=$(docker --version | cut -d' ' -f3 | cut -d',' -f1)
    echo "✅ Docker已安装: $DOCKER_VERSION"
else
    echo "❌ Docker未安装，请安装Docker"
    exit 1
fi

if command -v docker-compose &> /dev/null; then
    COMPOSE_VERSION=$(docker-compose --version | cut -d' ' -f4)
    echo "✅ Docker Compose已安装: $COMPOSE_VERSION"
elif docker compose version &> /dev/null; then
    COMPOSE_VERSION=$(docker compose version --short)
    echo "✅ Docker Compose已安装: $COMPOSE_VERSION"
else
    echo "❌ Docker Compose未安装，请安装Docker Compose"
    exit 1
fi

# 检查端口占用
echo "🔌 检查端口占用..."
PORTS=(5432 27017 6379 8080 8001 8002 8003 8004 8005)
for port in "${PORTS[@]}"; do
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null; then
        echo "⚠️  端口 $port 已被占用"
    else
        echo "✅ 端口 $port 可用"
    fi
done

# 检查必要的目录
echo "📁 检查项目结构..."
DIRS=(
    "backend/services"
    "backend/shared/models"
    "backend/shared/config"
    "mobile/lib"
    "mobile/lib/core"
)

for dir in "${DIRS[@]}"; do
    if [ -d "$dir" ]; then
        echo "✅ 目录存在: $dir"
    else
        echo "❌ 目录不存在: $dir"
        mkdir -p "$dir"
        echo "📁 已创建目录: $dir"
    fi
done

# 检查配置文件
echo "⚙️  检查配置文件..."
if [ -f "backend/.env" ]; then
    echo "✅ 后端环境配置文件存在"
else
    echo "⚠️  后端环境配置文件不存在，请复制 .env.example 到 .env"
    if [ -f "backend/.env.example" ]; then
        cp backend/.env.example backend/.env
        echo "📋 已复制示例配置文件"
    fi
fi

echo ""
echo "🎉 环境检查完成！"
echo ""
echo "📝 下一步操作："
echo "1. 复制并配置 backend/.env 文件"
echo "2. 运行 'make docker-up' 启动数据库服务"
echo "3. 运行 'make backend-dev' 启动后端服务"
echo "4. 运行 'make mobile-dev' 启动移动端应用"