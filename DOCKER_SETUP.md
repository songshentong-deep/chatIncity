# Docker 环境配置指南

## 当前状态
✅ Docker CLI 已安装 (版本: 28.3.2)  
❌ Docker 守护进程未运行  
✅ Docker 环境变量已配置

## 环境变量配置
已添加到 `~/.zshrc`:
```bash
# Docker Environment
export DOCKER_HOST=unix:///var/run/docker.sock
export DOCKER_BUILDKIT=1
export COMPOSE_DOCKER_CLI_BUILD=1

# Docker 别名
alias dps="docker ps"
alias dimg="docker images"
alias dlog="docker logs"
alias dexec="docker exec -it"
alias dstop="docker stop \$(docker ps -q)"
alias dclean="docker system prune -f"

# Docker Compose 别名
alias dcup="docker-compose up -d"
alias dcdown="docker-compose down"
alias dcps="docker-compose ps"
alias dclogs="docker-compose logs -f"
```

## 启动 Docker 的方法

### 方法1: 安装 Docker Desktop (推荐)
```bash
# 下载并安装 Docker Desktop
# 访问: https://www.docker.com/products/docker-desktop
# 或使用 Homebrew (需要网络连接):
brew install --cask docker

# 启动 Docker Desktop
open -a Docker
```

### 方法2: 使用 Colima (轻量级替代方案)
```bash
# 已安装 colima，启动时需要网络连接下载镜像
colima start

# 如果网络有问题，可以尝试:
colima start --network-address
```

### 方法3: 手动启动 Docker 守护进程 (高级用户)
```bash
# 在某些情况下可以尝试
sudo dockerd
```

## 验证 Docker 是否正常工作
```bash
# 检查 Docker 状态
docker info

# 运行测试容器
docker run hello-world

# 检查 Docker Compose
docker-compose --version
```

## 启动项目服务

一旦 Docker 正常运行，就可以启动项目服务：

```bash
# 启动数据库服务
docker-compose up -d postgres mongodb redis

# 或使用 Makefile
make docker-up

# 检查服务状态
docker-compose ps
```

## 故障排查

### 问题1: Docker 守护进程未运行
```bash
# 检查 Docker 进程
ps aux | grep docker

# 检查 Docker 服务状态
brew services list | grep docker
```

### 问题2: 端口被占用
```bash
# 检查端口占用
lsof -i :5432  # PostgreSQL
lsof -i :27017 # MongoDB
lsof -i :6379  # Redis
```

### 问题3: 网络连接问题
```bash
# 设置代理 (如果需要)
export HTTP_PROXY=http://proxy:port
export HTTPS_PROXY=http://proxy:port

# 或使用国内镜像
# 编辑 ~/.docker/daemon.json
{
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn",
    "https://hub-mirror.c.163.com"
  ]
}
```

## 下一步操作

1. **解决网络连接问题** - 确保可以访问 Docker Hub 或使用镜像
2. **选择并安装 Docker 运行时** - Docker Desktop 或 Colima
3. **启动 Docker 服务**
4. **运行项目服务**: `make docker-up`
5. **启动后端服务**: `make backend-dev`

## 临时解决方案

如果 Docker 暂时无法启动，可以：

1. **使用本地数据库**:
   ```bash
   # 安装本地 PostgreSQL
   brew install postgresql
   brew services start postgresql
   
   # 安装本地 MongoDB
   brew install mongodb-community
   brew services start mongodb-community
   
   # 安装本地 Redis
   brew install redis
   brew services start redis
   ```

2. **修改后端配置** 连接到本地数据库而不是 Docker 容器

---

*配置完成后，请运行 `source ~/.zshrc` 使环境变量生效*