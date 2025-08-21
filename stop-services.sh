#!/bin/bash

# 停止后端服务脚本

echo "🛑 正在停止所有后端服务..."

# 定义服务名称列表
SERVICES=(
    "user-service"
    "auth-service"
    "match-service"
    "chat-service"
    "content-service"
    "location-service"
    "notification-service"
    "admin-service"
    "upload-service"
    "moment-service"
    "news-service"
)

# 停止Go服务进程
echo "📋 查找并停止Go服务进程..."

for service in "${SERVICES[@]}"; do
    echo "正在停止 $service..."
    
    # 查找进程ID
    PIDS=$(ps aux | grep "$service" | grep -v grep | awk '{print $2}')
    
    if [ -n "$PIDS" ]; then
        echo "找到 $service 进程: $PIDS"
        # 优雅停止
        kill $PIDS 2>/dev/null
        sleep 1
        
        # 检查是否还在运行
        REMAINING=$(ps aux | grep "$service" | grep -v grep | awk '{print $2}')
        if [ -n "$REMAINING" ]; then
            echo "强制停止 $service 进程: $REMAINING"
            kill -9 $REMAINING 2>/dev/null
        fi
        echo "✅ $service 已停止"
    else
        echo "ℹ️  $service 未运行"
    fi
done

# 停止所有go run进程
echo ""
echo "🔍 停止所有 'go run' 进程..."
GO_RUN_PIDS=$(ps aux | grep "go run" | grep -v grep | awk '{print $2}')

if [ -n "$GO_RUN_PIDS" ]; then
    echo "找到 go run 进程: $GO_RUN_PIDS"
    kill $GO_RUN_PIDS 2>/dev/null
    sleep 1
    
    # 强制停止仍在运行的进程
    REMAINING_GO=$(ps aux | grep "go run" | grep -v grep | awk '{print $2}')
    if [ -n "$REMAINING_GO" ]; then
        echo "强制停止剩余 go run 进程: $REMAINING_GO"
        kill -9 $REMAINING_GO 2>/dev/null
    fi
    echo "✅ 所有 go run 进程已停止"
else
    echo "ℹ️  没有找到 go run 进程"
fi

# 停止可能的后台shell进程
echo ""
echo "🔍 停止相关的shell进程..."
SHELL_PIDS=$(ps aux | grep -E "cd backend.*go run" | grep -v grep | awk '{print $2}')

if [ -n "$SHELL_PIDS" ]; then
    echo "找到相关shell进程: $SHELL_PIDS"
    kill $SHELL_PIDS 2>/dev/null
    echo "✅ shell进程已停止"
fi

# 检查端口占用情况
echo ""
echo "📊 检查后端服务端口状态:"
PORTS=(8001 8002 8003 8004 8005 8006 8007 8008 8009 8010 8011)

for port in "${PORTS[@]}"; do
    PROCESS=$(lsof -ti :$port 2>/dev/null)
    if [ -n "$PROCESS" ]; then
        echo "⚠️  端口 $port 仍被占用 (PID: $PROCESS)"
        echo "正在释放端口 $port..."
        kill -9 $PROCESS 2>/dev/null
        sleep 0.5
        
        # 再次检查
        STILL_OCCUPIED=$(lsof -ti :$port 2>/dev/null)
        if [ -z "$STILL_OCCUPIED" ]; then
            echo "✅ 端口 $port 已释放"
        else
            echo "❌ 端口 $port 仍被占用"
        fi
    else
        echo "✅ 端口 $port 空闲"
    fi
done

# 最终检查
echo ""
echo "🔍 最终检查..."
REMAINING_SERVICES=$(ps aux | grep -E "(user-service|auth-service|match-service|chat-service|content-service|location-service|notification-service|admin-service|upload-service|moment-service)" | grep -v grep)

if [ -z "$REMAINING_SERVICES" ]; then
    echo "🎉 所有后端服务已成功停止！"
else
    echo "⚠️  仍有服务在运行:"
    echo "$REMAINING_SERVICES"
    echo ""
    echo "如需强制停止，请运行:"
    echo "sudo pkill -f 'service'"
fi

echo ""
echo "📋 服务停止完成！"