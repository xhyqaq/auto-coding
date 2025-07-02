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

# 检查 Claude CLI 是否可用
echo "🔍 检查 Claude CLI 可用性..."

if command -v claude >/dev/null 2>&1; then
    echo "✅ Claude CLI 已就绪"
    claude --version 2>/dev/null || echo "Claude CLI 可用"
else
    echo "⚠️  警告: Claude CLI 未找到"
    echo "   请确保已安装 Claude CLI 并挂载了 ~/.claude 目录"
    echo "   继续启动服务..."
fi

echo "🎯 启动 Claude GitHub Bot 服务..."

# 启动应用程序
exec "$@"