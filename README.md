# Claude GitHub Bot

一个基于 Claude Code CLI 的 GitHub 自动化修复机器人，支持智能分析 Issue 并自动生成修复 PR。

## 功能特点

- 🔍 **智能分析**：自动分析 GitHub Issue，提供详细的问题诊断
- 🤖 **自动修复**：基于 Claude Code CLI 进行代码修复
- 📋 **PR 管理**：自动创建 Pull Request 并处理反馈
- 💬 **交互式控制**：通过评论与机器人交互
- 🔄 **状态跟踪**：实时跟踪任务进度和状态

## 工作流程

### 1. Issue 创建和分析
```
用户创建 Issue (带 auto-fix 标签) → Bot 自动分析 → 提供修复方案 → 等待用户确认
```

### 2. 自动修复和 PR 创建
```
用户确认 (@bot proceed) → Bot 克隆代码 → Claude 修复 → 创建 PR → 等待审查
```

### 3. 反馈和迭代
```
用户在 PR 中提供反馈 → Bot 根据反馈修改 → 更新 PR → 继续审查
```

## 安装和配置

### 1. 环境要求

- Go 1.21+
- Node.js (用于运行 Claude Code CLI)
- Git
- 已安装 Claude Code CLI: `npm install -g @anthropic-ai/claude-code`

### 2. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，填入相应的 token 和配置
```

### 3. 运行服务

```bash
# 安装依赖
go mod tidy

# 运行服务
go run main.go
```

### 4. GitHub 配置

#### 创建 Personal Access Token
1. 前往 GitHub Settings > Developer settings > Personal access tokens
2. 创建新 token，勾选以下权限：
   - `repo` (完整仓库权限)
   - `issues` (读写 Issues)
   - `pull_requests` (创建和管理 PR)

#### 配置 Webhook
1. 前往目标仓库 Settings > Webhooks
2. 添加新 webhook：
   - Payload URL: `http://your-server.com/webhook`
   - Content type: `application/json`
   - Secret: 设置一个密钥 (与 WEBHOOK_SECRET 相同)
   - 选择触发事件：
     - Issues
     - Issue comments
     - Pull request reviews
     - Pull request review comments

## 使用方法

### 1. 创建 Issue
在 GitHub 仓库中创建 Issue，并添加 `auto-fix` 标签

### 2. Bot 命令
在 Issue 或 PR 评论中使用以下命令：

- `@bot proceed` - 开始修复处理
- `@bot cancel` - 取消当前任务
- `@bot status` - 查看任务状态

### 3. 使用示例

```markdown
# Issue 示例
标题：Login authentication fails with correct credentials
标签：auto-fix

描述：
Users are reporting that they cannot log in even with correct username and password.
The error message shows "Invalid credentials" but the credentials are correct.

Steps to reproduce:
1. Go to login page
2. Enter valid username/password
3. Click login
4. Error appears

Expected: User should be logged in
Actual: Error message "Invalid credentials"
```

Bot 将自动：
1. 分析问题并提供修复方案
2. 等待用户确认 (`@bot proceed`)
3. 自动修复代码并创建 PR
4. 处理用户反馈和进一步修改

## 技术架构

```
GitHub Issue → Webhook → Go Server → Claude Code CLI → Git Operations → GitHub PR
```

### 核心组件

- **Webhook 监听器**：接收 GitHub 事件
- **GitHub API 客户端**：与 GitHub 交互
- **Claude 集成**：调用 Claude Code CLI
- **状态机**：管理任务生命周期
- **Git 操作**：自动化代码管理

### 状态流转

```
ANALYZING → WAITING_APPROVAL → PROCESSING → WAITING_REVIEW → COMPLETED
    ↓              ↓                              ↓
CANCELLED      CANCELLED                    MODIFYING
```

## 安全考虑

- 使用 GitHub Personal Access Token 进行身份验证
- Webhook 使用 Secret 验证请求来源
- 限制 Bot 只能访问特定标签的 Issue
- 所有操作都需要用户明确确认

## 部署

### 1. 本地开发
```bash
export GITHUB_TOKEN=your_token
export WEBHOOK_SECRET=your_secret
export ANTHROPIC_API_KEY=your_anthropic_key
go run main.go
```

### 2. 生产部署
建议使用以下方式部署：
- Docker 容器
- Kubernetes
- 云服务 (AWS Lambda, Google Cloud Functions 等)

### 3. Docker 部署 (可选)
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o bot main.go

FROM node:alpine
RUN npm install -g @anthropic-ai/claude-code
COPY --from=builder /app/bot /usr/local/bin/
EXPOSE 8080
CMD ["bot"]
```

## 故障排除

### 常见问题

1. **Bot 没有响应**
   - 检查 webhook 配置
   - 确认 GITHUB_TOKEN 权限
   - 查看服务器日志

2. **Claude 调用失败**
   - 确认 ANTHROPIC_API_KEY 设置正确
   - 检查 Claude Code CLI 是否正确安装

3. **Git 操作失败**
   - 确认 token 有 repo 权限
   - 检查目标分支是否存在

### 日志查看
```bash
# 查看服务日志
journalctl -u claude-bot -f

# 或者直接运行查看输出
go run main.go
```

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License