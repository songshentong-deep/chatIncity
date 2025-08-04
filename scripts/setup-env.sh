#!/bin/bash

# 设置Go环境变量脚本

echo "🔧 设置Go开发环境变量..."

# 设置Go模块相关环境变量
export GOSUMDB=off
export GOPROXY=https://goproxy.cn,direct
export GO111MODULE=on

# 将环境变量添加到shell配置文件
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
    echo "📝 添加Go环境变量到 $SHELL_CONFIG"
    
    # 检查是否已经存在配置
    if ! grep -q "# Go Environment for Social App" "$SHELL_CONFIG"; then
        cat >> "$SHELL_CONFIG" << 'EOF'

# Go Environment for Social App
export GOSUMDB=off
export GOPROXY=https://goproxy.cn,direct
export GO111MODULE=on
EOF
        echo "✅ 环境变量已添加到 $SHELL_CONFIG"
        echo "请运行 'source $SHELL_CONFIG' 或重新打开终端使配置生效"
    else
        echo "✅ Go环境变量已存在于 $SHELL_CONFIG"
    fi
fi

echo ""
echo "🎉 Go环境配置完成！"
echo ""
echo "当前环境变量:"
echo "GOSUMDB: $GOSUMDB"
echo "GOPROXY: $GOPROXY"
echo "GO111MODULE: $GO111MODULE"