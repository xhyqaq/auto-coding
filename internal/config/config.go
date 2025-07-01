package config

import (
	"os"

	"github.com/xhy/auto-coding/internal/types"
)

// Load 加载配置
func Load() *types.BotConfig {
	return &types.BotConfig{
		GitHubToken:   os.Getenv("GITHUB_TOKEN"),
		WebhookSecret: os.Getenv("WEBHOOK_SECRET"),
		Port:          getEnvWithDefault("PORT", "8080"),
		AnthropicKey:  os.Getenv("ANTHROPIC_API_KEY"),
		BotName:       getEnvWithDefault("BOT_NAME", "claude-github-bot"),
		BotEmail:      getEnvWithDefault("BOT_EMAIL", "bot@example.com"),
		
		// Claude CLI 配置
		ClaudeCommand:       getEnvWithDefault("CLAUDE_COMMAND", "npx @anthropic-ai/claude-code"),
		ClaudeInstallSource: getEnvWithDefault("CLAUDE_INSTALL_SOURCE", "https://registry.npmjs.org/"),
		
		// Docker 配置
		DockerRegistry:  getEnvWithDefault("DOCKER_REGISTRY", "ghcr.io"),
		DockerImageName: getEnvWithDefault("DOCKER_IMAGE_NAME", "claude-github-bot"),
		
		// GitHub App 配置
		GitHubAppID:         os.Getenv("GITHUB_APP_ID"),
		GitHubAppPrivateKey: os.Getenv("GITHUB_APP_PRIVATE_KEY"),
		GitHubAppKeyFile:    os.Getenv("GITHUB_APP_KEY_FILE"),
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}