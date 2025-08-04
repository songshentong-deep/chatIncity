#!/bin/bash

# Docker 环境配置脚本

echo "🐳 配置 Docker 环境..."

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装"
    echo "请安装 Docker Desktop: https://www.docker.com/products/docker-desktop"
    exit 1
fi

echo "✅ Docker CLI 已安装: $(docker --version)"

# 检查 Docker 守护进程是否运行
if ! docker info &> /dev/null; then
    echo "❌ Docker 守护进程未运行"
    echo ""
    echo "请按以下步骤操作："
    echo "1. 安装 Docker Desktop (如果未安装):"
    echo "   brew install --cask docker"
    echo ""
    echo "2. 启动 Docker Desktop:"
    echo "   open -a Docker"
    echo ""
    echo "3. 或者手动从应用程序中启动 Docker Desktop"
    echo ""
    echo "4. 等待 Docker Desktop 完全启动后再运行此脚本"
    exit 1
fi

echo "✅ Docker 守护进程正在运行"

# 设置 Docker 环境变量
export DOCKER_HOST=unix:///var/run/docker.sock
export DOCKER_BUILDKIT=1
export COMPOSE_DOCKER_CLI_BUILD=1

# 将环境变量添加到 shell 配置文件
SHELL_CONFIG=""
if [ -f ~/.zshrc ]; then
    SHELL_CONFIG=~/.zshrc
elif [ -f ~/.bashrc ]; then
    SHELL_CONFIG=~/.bashrc
elif [ -f ~/.bash_profile ]; then
    SHELL_CONFIG=~/.bash_profile
fi

if [ -n "$SHELL_CONFIG" ]; then
    echo ""
    echo "📝 添加 Docker 环境变量到 $SHELL_CONFIG"
    
    # 检查是否已经存在配置
    if ! grep -q "# Docker Environment" "$SHELL_CONFIG"; then
        cat >> "$SHELL_CONFIG" << 'EOF'

# Docker Environment
export DOCKER_HOST=unix:///var/run/docker.sock
export DOCKER_BUILDKIT=1
export COMPOSE_DOCKER_CLI_BUILD=1

# Docker 别名
alias dps='docker ps'
alias dimg='docker images'
alias dlog='docker logs'
alias dexec='docker exec -it'
alias dstop='docker stop $(docker ps -q)'
alias dclean='docker system prune -f'

# Docker Compose 别名
alias dcup='docker-compose up -d'
alias dcdown='docker-compose down'
alias dcps='docker-compose ps'
alias dclogs='docker-compose logs -f'
EOF
        echo "✅ 环境变量已添加到 $SHELL_CONFIG"
        echo "请运行 'source $SHELL_CONFIG' 或重新打开终端使配置生效"
    else
        echo "✅ Docker 环境变量已存在于 $SHELL_CONFIG"
    fi
fi

# 验证 Docker 配置
echo ""
echo "🔍 验证 Docker 配置:"
echo "Docker 版本: $(docker --version)"
echo "Docker Compose 版本: $(docker-compose --version 2>/dev/null || docker compose version 2>/dev/null || echo '未安装')"
echo "Docker 上下文: $(docker context show)"
echo "Docker 信息:"
docker info --format "{{.ServerVersion}}" 2>/dev/null && echo "✅ Docker 服务正常" || echo "❌ Docker 服务异常"

echo ""
echo "🎉 Docker 环境配置完成！"