# Claude GitHub Bot

一个完全自主的 Claude GitHub Bot，接收所有 GitHub 事件并让 Claude 自主决定如何处理。

## 快速开始

### Docker 部署（推荐）

```bash
# 1. 拉取镜像
docker pull ghcr.io/xhyqaq/auto-coding:latest

# 2. 配置环境变量
cp .env.example .env
# 编辑 .env 填入必要配置

# 3. 启动服务
docker-compose up -d

# 4. 检查状态
curl http://localhost:8888/health
```

### 直接运行

```bash
docker run -d \
  --name claude-github-bot \
  -p 8888:8080 \
  --env-file .env \
  ghcr.io/xhyqaq/auto-coding:latest
```

## 环境配置

### 必需变量
```bash
GITHUB_TOKEN=ghp_your_token_here       # GitHub Personal Access Token (repo权限)
WEBHOOK_SECRET=your_secret_here        # Webhook验证密钥
```

### 可选变量
```bash
ANTHROPIC_API_KEY=your_key_here                        # Claude API Key (如果CLI未登录)
CLAUDE_INSTALL_SOURCE=https://registry.npmmirror.com/  # npm镜像源(国内用户)
BOT_NAME=claude-github-bot                             # Git提交者名称
BOT_EMAIL=bot@example.com                              # Git提交者邮箱
PORT=8080                                              # 服务端口
```

## GitHub 配置

1. **创建 Personal Access Token**
   - 前往 GitHub Settings > Developer settings > Personal access tokens
   - 创建 token，选择 `repo` 权限

2. **配置 Webhook**
   - 前往仓库 Settings > Webhooks
   - 添加 webhook：
     - URL: `http://your-server.com:8888/webhook`
     - Content type: `application/json`
     - Secret: 与 `WEBHOOK_SECRET` 相同
     - 选择 "Send me everything"

## 使用方法

1. 部署 bot 并配置 webhook
2. 在仓库中创建 issue
3. Claude 自动分析并决定如何响应
4. 无需其他操作

## 本地开发

```bash
# 环境要求
go 1.21+

# 运行
go mod tidy
go run main.go
```

## 许可证

MIT License