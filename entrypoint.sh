#!/bin/bash
set -e

echo "🚀 启动 Claude GitHub Bot..."

# 检查必要的环境变量
if [ -z "$GITHUB_TOKEN" ]; then
    echo "❌ 错误: GITHUB_TOKEN 环境变量是必需的"
    exit 1
fi

if [ -z "$WEBHOOK_SECRET" ]; then
    echo "❌ 错误: WEBHOOK_SECRET 环境变量是必需的"
    exit 1
fi

# 动态安装 Claude CLI
echo "🔧 正在安装 Claude CLI..."

# 默认安装命令
DEFAULT_CLAUDE_INSTALL_CMD="npm install -g @anthropic-ai/claude-code"

# 使用用户提供的命令或默认命令
CLAUDE_INSTALL_CMD="${CLAUDE_INSTALL_CMD:-$DEFAULT_CLAUDE_INSTALL_CMD}"

echo "📦 安装命令: $CLAUDE_INSTALL_CMD"

# 设置 npm registry（如果指定）
if [ -n "$CLAUDE_INSTALL_SOURCE" ]; then
    echo "🌏 设置 npm registry: $CLAUDE_INSTALL_SOURCE"
    npm config set registry "$CLAUDE_INSTALL_SOURCE"
fi

# 执行安装命令
if eval "$CLAUDE_INSTALL_CMD"; then
    echo "✅ Claude CLI 安装完成"
else
    echo "❌ Claude CLI 安装失败"
    exit 1
fi

# 验证安装
if command -v claude >/dev/null 2>&1; then
    echo "✅ Claude CLI 验证成功"
    claude --version 2>/dev/null || echo "Claude CLI 已安装"
else
    echo "❌ Claude CLI 未找到，请检查安装命令"
    exit 1
fi

# 清理 npm 缓存
npm cache clean --force 2>/dev/null || true

echo "🎯 启动 Claude GitHub Bot 服务..."

# 启动应用程序
exec "$@"