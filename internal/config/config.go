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
		BotName:       getEnvWithDefault("BOT_NAME", "agent-auto-coding"),
		BotEmail:      getEnvWithDefault("BOT_EMAIL", "xhyovo@qq.com"),
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}