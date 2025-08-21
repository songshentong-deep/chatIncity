#!/bin/bash

# 简化版健康检查脚本

echo "🔍 执行简化版健康检查..."
echo "时间: $(date)"
echo "================================"

# 检查容器状态
echo "🐳 检查容器状态..."
CONTAINER_NAME="social-app-all-services"

if docker ps --format '{{.Names}}' | grep -q "^$CONTAINER_NAME$"; then
    echo "✅ 容器运行中"
    
    # 检查容器健康状态
    HEALTH_STATUS=$(docker inspect --format='{{.State.Health.Status}}' $CONTAINER_NAME 2>/dev/null || echo "no-healthcheck")
    echo "容器健康状态: $HEALTH_STATUS"
    
    if [ "$HEALTH_STATUS" = "healthy" ] || [ "$HEALTH_STATUS" = "no-healthcheck" ]; then
        echo "✅ 容器健康"
    else
        echo "⚠️ 容器不健康"
        CONTAINER_FAILED=true
    fi
else
    echo "❌ 容器未运行"
    CONTAINER_FAILED=true
fi

echo ""

# 检查API网关
echo "🌐 检查API网关..."
if curl -f -s --max-time 10 "http://localhost/health" > /dev/null; then
    echo "✅ API网关正常"
else
    echo "❌ API网关失败"
    GATEWAY_FAILED=true
fi

echo ""

# 检查各个微服务
echo "🔧 检查微服务状态..."
SERVICES=(
    "user:http://localhost/api/v1/user/health"
    "chat:http://localhost/api/v1/chat/health"
    "moment:http://localhost/api/v1/moment/health"
    "upload:http://localhost/api/v1/upload/health"
    "news:http://localhost/api/v1/news/health"
)

for service in "${SERVICES[@]}"; do
    name=$(echo $service | cut -d: -f1)
    url=$(echo $service | cut -d: -f2-)
    
    echo -n "检查 $name 服务... "
    
    if curl -f -s --max-time 10 "$url" > /dev/null; then
        echo "✅ 正常"
    else
        echo "❌ 失败"
        SERVICE_FAILED=true
    fi
done

echo ""

# 检查Supervisor进程状态
echo "👥 检查Supervisor进程状态..."
if docker exec $CONTAINER_NAME supervisorctl status 2>/dev/null; then
    echo "✅ Supervisor状态正常"
else
    echo "❌ 无法获取Supervisor状态"
    SUPERVISOR_FAILED=true
fi

echo ""

# 检查系统资源
echo "📊 检查系统资源..."

# 检查磁盘空间
DISK_USAGE=$(df / | awk 'NR==2 {print $5}' | sed 's/%//')
echo -n "磁盘使用率: $DISK_USAGE%... "
if [ $DISK_USAGE -lt 90 ]; then
    echo "✅ 正常"
else
    echo "⚠️ 磁盘空间不足"
    RESOURCE_FAILED=true
fi

# 检查内存使用
MEMORY_USAGE=$(free | awk 'NR==2{printf "%.0f", $3*100/$2}')
echo -n "内存使用率: $MEMORY_USAGE%... "
if [ $MEMORY_USAGE -lt 90 ]; then
    echo "✅ 正常"
else
    echo "⚠️ 内存使用率过高"
    RESOURCE_FAILED=true
fi

# 检查容器资源使用
echo ""
echo "🐳 容器资源使用情况:"
docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}" $CONTAINER_NAME

echo ""

# 检查日志中的错误
echo "📝 检查最近的错误日志..."
ERROR_COUNT=$(docker logs $CONTAINER_NAME --since="5m" 2>&1 | grep -i "error\|fatal\|panic" | wc -l)
echo -n "最近5分钟错误日志数量: $ERROR_COUNT... "
if [ $ERROR_COUNT -lt 10 ]; then
    echo "✅ 正常"
else
    echo "⚠️ 错误日志较多"
    LOG_FAILED=true
fi

echo ""
echo "================================"

# 汇总结果
if [ "$CONTAINER_FAILED" = true ] || [ "$GATEWAY_FAILED" = true ] || [ "$SERVICE_FAILED" = true ] || [ "$SUPERVISOR_FAILED" = true ] || [ "$RESOURCE_FAILED" = true ] || [ "$LOG_FAILED" = true ]; then
    echo "❌ 健康检查发现问题"
    
    echo ""
    echo "🔧 故障排查建议:"
    
    if [ "$CONTAINER_FAILED" = true ]; then
        echo "- 检查容器状态: docker ps -a"
        echo "- 查看容器日志: docker logs $CONTAINER_NAME"
    fi
    
    if [ "$GATEWAY_FAILED" = true ] || [ "$SERVICE_FAILED" = true ]; then
        echo "- 检查服务进程: docker exec $CONTAINER_NAME supervisorctl status"
        echo "- 重启服务: docker exec $CONTAINER_NAME supervisorctl restart all"
    fi
    
    if [ "$RESOURCE_FAILED" = true ]; then
        echo "- 清理磁盘空间: docker system prune -f"
        echo "- 检查内存使用: docker stats"
    fi
    
    if [ "$LOG_FAILED" = true ]; then
        echo "- 查看详细日志: docker logs $CONTAINER_NAME --tail 100"
    fi
    
    exit 1
else
    echo "✅ 所有检查项目正常"
    echo ""
    echo "🎉 系统运行良好！"
    exit 0
fi