package bot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/v58/github"

	"github.com/xhy/auto-coding/internal/claude"
	ghclient "github.com/xhy/auto-coding/internal/github"
	"github.com/xhy/auto-coding/internal/types"
)

// Bot 实例
type Bot struct {
	config       *types.BotConfig
	githubClient *ghclient.Client
	claudeClient *claude.Client
}

// NewBot 创建新的 Bot 实例
func NewBot(config *types.BotConfig) *Bot {
	return &Bot{
		config:       config,
		githubClient: ghclient.NewClient(config),
		claudeClient: claude.NewClient(config),
	}
}

// createWorkspace 创建工作空间
func (b *Bot) createWorkspace() (string, error) {
	workspace, err := ioutil.TempDir("", "claude-workspace-*")
	if err != nil {
		return "", fmt.Errorf("failed to create workspace: %v", err)
	}
	log.Printf("Created workspace: %s", workspace)
	return workspace, nil
}

// cleanupWorkspace 清理工作空间
func (b *Bot) cleanupWorkspace(workspace string) {
	if err := os.RemoveAll(workspace); err != nil {
		log.Printf("Failed to cleanup workspace %s: %v", workspace, err)
	} else {
		log.Printf("Cleaned up workspace: %s", workspace)
	}
}

// getBotCapabilities 获取 Bot 能力
func (b *Bot) getBotCapabilities() types.BotCapabilities {
	return types.BotCapabilities{
		CanCreatePR:       true,
		CanModifyIssues:   true,
		CanManageLabels:   true,
		CanTriggerActions: true,
		HasFullRepoAccess: true,
		CanCloneRepo:      true,
		CanPushChanges:    true,
	}
}

// setupClaudeContext 设置 Claude 上下文
func (b *Bot) setupClaudeContext(workspace string, context *types.GitHubContext) error {
	// 创建上下文文件，供 Claude 参考
	contextFile := filepath.Join(workspace, ".claude-context.json")
	contextData, err := json.MarshalIndent(context, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal context: %v", err)
	}

	if err := ioutil.WriteFile(contextFile, contextData, 0644); err != nil {
		return fmt.Errorf("failed to write context file: %v", err)
	}

	log.Printf("Created context file: %s", contextFile)
	return nil
}

// eventNeedsRepository 判断事件是否需要克隆仓库
func (b *Bot) eventNeedsRepository(eventType string) bool {
	// 需要仓库访问的事件类型
	repoEvents := map[string]bool{
		"issues":                      true,  // issue 创建、更新
		"issue_comment":               true,  // issue 评论
		"pull_request":                true,  // PR 相关
		"pull_request_review":         true,  // PR 审查
		"pull_request_review_comment": true,  // PR 评论
		"push":                        false, // push 通常不需要额外操作
		"create":                      false, // 创建分支等
		"delete":                      false, // 删除操作
		"star":                        false, // 星标操作
		"watch":                       false, // 关注操作
		"fork":                        false, // fork 操作
	}

	// 默认为需要仓库（保守策略）
	if needed, exists := repoEvents[eventType]; exists {
		return needed
	}

	log.Printf("Unknown event type %s, assuming repository needed", eventType)
	return true
}

// isSelfGeneratedEvent 检查是否是机器人自己产生的事件（避免自回复循环）
func (b *Bot) isSelfGeneratedEvent(eventType string, payload map[string]interface{}) bool {
	// 检查不同类型事件的作者
	switch eventType {
	case "issue_comment":
		if comment, ok := payload["comment"].(map[string]interface{}); ok {
			if user, ok := comment["user"].(map[string]interface{}); ok {
				if login, ok := user["login"].(string); ok {
					// 检查是否是机器人用户
					if b.isBotUser(login) {
						log.Printf("Ignoring self-generated issue comment from: %s", login)
						return true
					}
				}
			}
		}
	case "pull_request_review_comment":
		if comment, ok := payload["comment"].(map[string]interface{}); ok {
			if user, ok := comment["user"].(map[string]interface{}); ok {
				if login, ok := user["login"].(string); ok {
					if b.isBotUser(login) {
						log.Printf("Ignoring self-generated PR comment from: %s", login)
						return true
					}
				}
			}
		}
	}
	return false
}

// isBotUser 判断用户是否是机器人
func (b *Bot) isBotUser(username string) bool {
	botKeywords := []string{"bot", "claude", "github-actions", "[bot]"}
	username = strings.ToLower(username)

	for _, keyword := range botKeywords {
		if strings.Contains(username, keyword) {
			return true
		}
	}
	return false
}

// detectWebhookSource 检测webhook来源（GitHub或Gitee）
func (b *Bot) detectWebhookSource(r *http.Request) string {
	// 检查GitHub特有的header
	if r.Header.Get("X-GitHub-Event") != "" {
		return "github"
	}
	
	// 检查Gitee特有的header
	if r.Header.Get("X-Gitee-Event") != "" || r.Header.Get("X-Gitee-Token") != "" {
		return "gitee"
	}
	
	// 默认按GitHub处理
	return "github"
}

// validateGitHubWebhook 验证GitHub webhook
func (b *Bot) validateGitHubWebhook(r *http.Request) ([]byte, error) {
	return github.ValidatePayload(r, []byte(b.config.WebhookSecret))
}

// validateGiteeWebhook 验证Gitee webhook
func (b *Bot) validateGiteeWebhook(r *http.Request) ([]byte, error) {
	// 读取请求体
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %v", err)
	}

	// Gitee使用X-Gitee-Token header进行验证
	giteeToken := r.Header.Get("X-Gitee-Token")
	if giteeToken != "" {
		// 验证token是否匹配
		if giteeToken != b.config.WebhookSecret {
			return nil, fmt.Errorf("gitee token mismatch")
		}
	} else {
		// 如果没有X-Gitee-Token，使用HMAC验证（一些Gitee配置）
		signature := r.Header.Get("X-Gitee-Signature")
		if signature == "" {
			return nil, fmt.Errorf("no gitee signature found")
		}
		
		// 验证HMAC签名
		mac := hmac.New(sha256.New, []byte(b.config.WebhookSecret))
		mac.Write(body)
		expectedSignature := hex.EncodeToString(mac.Sum(nil))
		
		if signature != expectedSignature {
			return nil, fmt.Errorf("gitee signature mismatch")
		}
	}
	
	return body, nil
}

// getEventType 获取事件类型
func (b *Bot) getEventType(r *http.Request, source string) string {
	switch source {
	case "github":
		return github.WebHookType(r)
	case "gitee":
		return r.Header.Get("X-Gitee-Event")
	default:
		return github.WebHookType(r)
	}
}

// HandleGitHubEvent 处理 GitHub 事件 - 完全委托给 Claude
func (b *Bot) HandleGitHubEvent(eventType string, payload interface{}) error {
	// 检查是否是机器人自己产生的事件（避免自回复循环）
	payloadMap := payload.(map[string]interface{})
	if b.isSelfGeneratedEvent(eventType, payloadMap) {
		log.Printf("Skipping self-generated event: %s", eventType)
		return nil
	}

	// 创建工作空间
	workspace, err := b.createWorkspace()
	if err != nil {
		return fmt.Errorf("failed to create workspace: %v", err)
	}
	defer b.cleanupWorkspace(workspace)

	// 创建上下文
	context := &types.GitHubContext{
		EventType:    eventType,
		Payload:      payload.(map[string]interface{}),
		Workspace:    workspace,
		Capabilities: b.getBotCapabilities(),
		Timestamp:    time.Now(),
	}

	// 对于大部分事件都需要克隆仓库（除了一些纯通知事件）
	needsRepo := b.eventNeedsRepository(eventType)
	if needsRepo {
		if repo, ok := context.Payload["repository"].(map[string]interface{}); ok {
			if cloneURL, ok := repo["clone_url"].(string); ok {
				if err := b.githubClient.CloneRepository(cloneURL, workspace); err != nil {
					log.Printf("Failed to clone repository: %v", err)
					// 对于需要仓库的事件，如果克隆失败就不继续
					return fmt.Errorf("repository required but clone failed: %v", err)
				}
			}
		}
	}

	// 设置 Claude 上下文
	if err := b.setupClaudeContext(workspace, context); err != nil {
		return fmt.Errorf("failed to setup claude context: %v", err)
	}

	// 让 Claude 完全自主处理
	return b.claudeClient.Invoke(workspace)
}

// WebhookHandler Webhook 处理器 - 支持GitHub和Gitee
func (b *Bot) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	// 检测webhook来源
	source := b.detectWebhookSource(r)
	log.Printf("Detected webhook source: %s", source)

	// 根据来源验证payload
	var payload []byte
	var err error
	
	switch source {
	case "github":
		payload, err = b.validateGitHubWebhook(r)
	case "gitee":
		payload, err = b.validateGiteeWebhook(r)
	default:
		payload, err = b.validateGitHubWebhook(r)
	}
	
	if err != nil {
		log.Printf("Invalid %s payload: %v", source, err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// 获取事件类型
	eventType := b.getEventType(r, source)
	log.Printf("Received %s event: %s", source, eventType)

	// 解析 payload
	parsedPayload, err := b.githubClient.ParseWebhookPayload(payload)
	if err != nil {
		log.Printf("Failed to parse webhook payload: %v", err)
		http.Error(w, "Failed to parse webhook", http.StatusBadRequest)
		return
	}

	// 在payload中添加来源信息
	parsedPayload["webhook_source"] = source

	// 异步处理所有事件，不做任何过滤
	go func() {
		if err := b.HandleGitHubEvent(eventType, parsedPayload); err != nil {
			log.Printf("Failed to handle %s event %s: %v", source, eventType, err)
		}
	}()

	w.WriteHeader(http.StatusOK)
}
