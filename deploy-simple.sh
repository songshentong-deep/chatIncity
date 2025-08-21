#!/bin/bash

set -e

echo "🚀 开始简化版生产环境部署..."

# 配置变量
APP_DIR="/opt/social-app"
BACKUP_DIR="/opt/backups/social-app"
DOCKER_COMPOSE_FILE="docker-compose.simple.yml"
LOG_FILE="/opt/social-app/logs/deploy.log"

# 创建日志文件
mkdir -p "$(dirname "$LOG_FILE")"
exec 1> >(tee -a "$LOG_FILE")
exec 2>&1

echo "$(date): 开始部署流程"

# 检查必要的命令
command -v docker >/dev/null 2>&1 || { echo "❌ Docker未安装"; exit 1; }
command -v docker-compose >/dev/null 2>&1 || { echo "❌ Docker Compose未安装"; exit 1; }

# 创建备份目录
mkdir -p $BACKUP_DIR

# 备份当前版本
if [ -d "$APP_DIR" ] && [ "$(ls -A $APP_DIR)" ]; then
    echo "📦 备份当前版本..."
    BACKUP_FILE="$BACKUP_DIR/backup-$(date +%Y%m%d-%H%M%S).tar.gz"
    tar -czf "$BACKUP_FILE" -C $APP_DIR . --exclude='logs' --exclude='data' --exclude='.git'
    echo "✅ 备份完成: $BACKUP_FILE"
fi

# 进入应用目录
cd $APP_DIR

# 拉取最新代码
echo "📥 拉取最新代码..."
git fetch origin
git reset --hard origin/main
echo "✅ 代码更新完成"

# 检查配置文件
if [ ! -f ".env" ]; then
    echo "❌ 缺少 .env 配置文件"
    exit 1
fi

# 停止旧服务
echo "⏹️ 停止旧服务..."
docker-compose -f $DOCKER_COMPOSE_FILE down --remove-orphans || true

# 清理未使用的镜像
echo "🧹 清理旧镜像..."
docker image prune -f

# 构建新镜像
echo "🔨 构建应用镜像..."
docker-compose -f $DOCKER_COMPOSE_FILE build --no-cache

# 启动新服务
echo "▶️ 启动新服务..."
docker-compose -f $DOCKER_COMPOSE_FILE up -d

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 60

# 健康检查
echo "🔍 执行健康检查..."
MAX_RETRIES=10
RETRY_COUNT=0

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if curl -f -s http://localhost/health > /dev/null; then
        echo "✅ 健康检查通过"
        break
    else
        RETRY_COUNT=$((RETRY_COUNT + 1))
        echo "⏳ 健康检查失败，重试 $RETRY_COUNT/$MAX_RETRIES..."
        sleep 10
    fi
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    echo "❌ 健康检查失败，回滚到上一版本"
    
    # 回滚
    docker-compose -f $DOCKER_COMPOSE_FILE down
    
    # 恢复备份
    if [ -f "$BACKUP_FILE" ]; then
        echo "🔄 恢复备份..."
        tar -xzf "$BACKUP_FILE" -C $APP_DIR
        docker-compose -f $DOCKER_COMPOSE_FILE up -d
    fi
    
    exit 1
fi

# 清理旧备份（保留最近5个）
echo "🧹 清理旧备份..."
ls -t $BACKUP_DIR/backup-*.tar.gz | tail -n +6 | xargs -r rm

echo "✅ 部署完成！"
echo "$(date): 部署流程完成"

# 显示服务状态
echo ""
echo "📊 服务状态:"
docker-compose -f $DOCKER_COMPOSE_FILE ps

echo ""
echo "🌐 服务访问地址:"
echo "- API网关: http://localhost"
echo "- 用户服务: http://localhost/api/v1/user/"
echo "- 聊天服务: http://localhost/api/v1/chat/"
echo "- 动态服务: http://localhost/api/v1/moment/"
echo "- 上传服务: http://localhost/api/v1/upload/"
echo "- 新闻服务: http://localhost/api/v1/news/"

echo ""
echo "📋 管理命令:"
echo "- 查看日志: docker-compose -f $DOCKER_COMPOSE_FILE logs -f"
echo "- 重启服务: docker-compose -f $DOCKER_COMPOSE_FILE restart"
echo "- 进入容器: docker-compose -f $DOCKER_COMPOSE_FILE exec social-app sh"
echo "- 查看进程: docker-compose -f $DOCKER_COMPOSE_FILE exec social-app supervisorctl status"