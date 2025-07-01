package types

import (
	"time"

	"github.com/google/go-github/v58/github"
)

// BotConfig 配置
type BotConfig struct {
	GitHubToken   string
	WebhookSecret string
	Port          string
	AnthropicKey  string
	BotName       string // Bot 提交时使用的名称
	BotEmail      string // Bot 提交时使用的邮箱

	// Claude CLI 配置
	ClaudeCommand       string // Claude CLI 命令
	ClaudeInstallSource string // Claude CLI 安装源

	// Docker 配置
	DockerRegistry  string // Docker 镜像仓库
	DockerImageName string // Docker 镜像名称

	// GitHub App 配置
	GitHubAppID         string // GitHub App ID
	GitHubAppPrivateKey string // GitHub App Private Key 内容
	GitHubAppKeyFile    string // GitHub App Private Key 文件路径
}

// GitHubContext GitHub 事件上下文
type GitHubContext struct {
	EventType    string                 `json:"event_type"`
	Payload      map[string]interface{} `json:"payload"`
	Repository   *github.Repository     `json:"repository,omitempty"`
	Workspace    string                 `json:"workspace"`
	Capabilities BotCapabilities        `json:"capabilities"`
	Timestamp    time.Time              `json:"timestamp"`
}

// BotCapabilities Bot 能力
type BotCapabilities struct {
	CanCreatePR       bool `json:"can_create_pr"`
	CanModifyIssues   bool `json:"can_modify_issues"`
	CanManageLabels   bool `json:"can_manage_labels"`
	CanTriggerActions bool `json:"can_trigger_actions"`
	HasFullRepoAccess bool `json:"has_full_repo_access"`
	CanCloneRepo      bool `json:"can_clone_repo"`
	CanPushChanges    bool `json:"can_push_changes"`
}
