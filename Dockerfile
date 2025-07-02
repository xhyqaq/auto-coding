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
    wget \
    jq \
    && rm -rf /var/cache/apk/*

# 安装 GitHub CLI (Alpine方式)
RUN ARCH=$(uname -m) && \
    if [ "$ARCH" = "x86_64" ]; then ARCH="amd64"; fi && \
    if [ "$ARCH" = "aarch64" ]; then ARCH="arm64"; fi && \
    GH_VERSION=$(curl -s https://api.github.com/repos/cli/cli/releases/latest | jq -r '.tag_name' | cut -c2-) && \
    wget "https://github.com/cli/cli/releases/download/v${GH_VERSION}/gh_${GH_VERSION}_linux_${ARCH}.tar.gz" -O gh.tar.gz && \
    tar -xzf gh.tar.gz && \
    mv "gh_${GH_VERSION}_linux_${ARCH}/bin/gh" /usr/local/bin/ && \
    chmod +x /usr/local/bin/gh && \
    rm -rf gh*

# 设置时区
ENV TZ=Asia/Shanghai

# 创建应用用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup

# 设置工作目录
WORKDIR /app

# Claude CLI需要用户手动安装并挂载配置目录或提供API Key

# 从构建阶段复制二进制文件
COPY --from=builder /app/claude-github-bot .

# 复制启动脚本
COPY entrypoint.sh .

# 设置权限和目录
RUN chmod +x entrypoint.sh && \
    mkdir -p /app/workspaces && \
    chown -R appuser:appgroup /app

# 切换到应用用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# 设置启动脚本
ENTRYPOINT ["./entrypoint.sh"]

# 启动应用
CMD ["./claude-github-bot"]