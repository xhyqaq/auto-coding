# 多阶段构建：构建阶段
FROM golang:1.21-alpine AS builder

# 安装必要的构建工具
RUN apk add --no-cache git ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o claude-github-bot .

# 运行阶段
FROM node:18-alpine

# 安装系统依赖
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    bash \
    curl \
    && rm -rf /var/cache/apk/*

# 设置时区
ENV TZ=Asia/Shanghai

# 创建应用用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup

# 设置工作目录
WORKDIR /app

# 安装 Claude Code CLI (支持自定义安装源)
ARG CLAUDE_INSTALL_SOURCE=""
ENV CLAUDE_INSTALL_SOURCE=${CLAUDE_INSTALL_SOURCE}

# 安装 Claude CLI
RUN if [ -n "$CLAUDE_INSTALL_SOURCE" ]; then \
        npm config set registry $CLAUDE_INSTALL_SOURCE; \
    fi && \
    npm install -g @anthropic-ai/claude-code && \
    npm cache clean --force

# 从构建阶段复制二进制文件
COPY --from=builder /app/claude-github-bot .

# 创建必要的目录
RUN mkdir -p /app/workspaces && \
    chown -R appuser:appgroup /app

# 切换到应用用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# 启动应用
CMD ["./claude-github-bot"]