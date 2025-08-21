#!/bin/bash

# 重启指定端口的服务脚本
# 用法: ./restart-service.sh <端口号> [服务名]

if [ $# -eq 0 ]; then
    echo "❌ 请指定端口号"
    echo "用法: $0 <端口号> [服务名]"
    echo "示例: $0 8009 upload-service"
    exit 1
fi

PORT=$1
SERVICE_NAME=${2:-"service"}

echo "🔄 重启端口 $PORT 上的服务..."

# 1. 停止占用指定端口的进程
echo "🛑 停止端口 $PORT 上的进程..."
PIDS=$(lsof -ti :$PORT 2>/dev/null)

if [ -n "$PIDS" ]; then
    echo "找到占用端口 $PORT 的进程: $PIDS"
    
    # 优雅停止
    kill $PIDS 2>/dev/null
    sleep 2
    
    # 检查是否还在运行
    REMAINING=$(lsof -ti :$PORT 2>/dev/null)
    if [ -n "$REMAINING" ]; then
        echo "强制停止进程: $REMAINING"
        kill -9 $REMAINING 2>/dev/null
        sleep 1
    fi
    
    # 最终检查
    FINAL_CHECK=$(lsof -ti :$PORT 2>/dev/null)
    if [ -z "$FINAL_CHECK" ]; then
        echo "✅ 端口 $PORT 已释放"
    else
        echo "❌ 端口 $PORT 仍被占用，无法重启"
        exit 1
    fi
else
    echo "ℹ️  端口 $PORT 未被占用"
fi

# 2. 根据端口号启动对应服务
echo "🚀 启动服务..."

# 获取脚本所在目录的父目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BACKEND_DIR="$PROJECT_ROOT/backend"

case $PORT in
    8001)
        echo "启动用户服务..."
        cd "$BACKEND_DIR" && go run ./services/user-service &
        ;;
    8002)
        echo "启动认证服务..."
        cd "$BACKEND_DIR" && go run ./services/auth-service &
        ;;
    8003)
        echo "启动匹配服务..."
        cd "$BACKEND_DIR" && go run ./services/match-service &
        ;;
    8004)
        echo "启动聊天服务..."
        cd "$BACKEND_DIR" && go run ./services/chat-service &
        ;;
    8005)
        echo "启动内容服务..."
        cd "$BACKEND_DIR" && go run ./services/content-service &
        ;;
    8006)
        echo "启动位置服务..."
        cd "$BACKEND_DIR" && go run ./services/location-service &
        ;;
    8007)
        echo "启动通知服务..."
        cd "$BACKEND_DIR" && go run ./services/notification-service &
        ;;
    8008)
        echo "启动管理服务..."
        cd "$BACKEND_DIR" && go run ./services/admin-service &
        ;;
    8009)
        echo "启动文件上传服务..."
        cd "$BACKEND_DIR" && go run ./services/upload-service &
        ;;
    8010)
        echo "启动动态服务..."
        cd "$BACKEND_DIR" && go run ./services/moment-service &
        ;;
    8011)
        echo "启动动态服务..."
        ./scripts/start-news-service.sh &
        ;;
    *)
        echo "❌ 未知端口 $PORT，无法自动启动服务"
        echo "支持的端口: 8001-8011"
        exit 1
        ;;
esac

# 3. 等待服务启动
echo "⏳ 等待服务启动..."
sleep 3

# 4. 检查服务是否成功启动
NEW_PID=$(lsof -ti :$PORT 2>/dev/null)
if [ -n "$NEW_PID" ]; then
    echo "✅ 服务已成功启动在端口 $PORT (PID: $NEW_PID)"
else
    echo "❌ 服务启动失败，请检查日志"
    exit 1
fi

echo "🎉 端口 $PORT 服务重启完成！"