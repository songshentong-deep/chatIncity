#!/bin/bash

# 安装Air工具
if ! command -v air &> /dev/null; then
    echo "安装Air工具..."
    curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
fi

# 启动热更新服务
echo "启动热更新服务..."
air -c .air.toml