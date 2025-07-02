package github

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/google/go-github/v58/github"
	"golang.org/x/oauth2"

	"github.com/xhy/auto-coding/internal/types"
)

// Client GitHub 客户端包装器
type Client struct {
	client *github.Client
	config *types.BotConfig
}

// NewClient 创建新的 GitHub 客户端
func NewClient(config *types.BotConfig) *Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GitHubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &Client{
		client: github.NewClient(tc),
		config: config,
	}
}

// GetClient 获取原始 GitHub 客户端
func (c *Client) GetClient() *github.Client {
	return c.client
}

// CloneRepository 克隆仓库到工作空间
func (c *Client) CloneRepository(repoURL, workspace string) error {
	repoDir := filepath.Join(workspace, "repository")
	
	// 使用GitHub token进行认证的clone URL
	// 将 https://github.com/owner/repo 转换为 https://token@github.com/owner/repo
	authenticatedURL := c.addTokenToURL(repoURL)
	
	cmd := exec.Command("git", "clone", authenticatedURL, repoDir)
	cmd.Env = append(cmd.Env, "GIT_TERMINAL_PROMPT=0")

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to clone repository: %v, output: %s", err, string(output))
	}

	log.Printf("Successfully cloned repository to %s", repoDir)
	return nil
}

// addTokenToURL 在GitHub URL中添加token认证
func (c *Client) addTokenToURL(repoURL string) string {
	// 如果是GitHub URL，添加token认证
	if len(repoURL) > 19 && repoURL[:19] == "https://github.com/" {
		return fmt.Sprintf("https://%s@github.com/%s", c.config.GitHubToken, repoURL[19:])
	}
	// 如果已经包含认证信息或不是GitHub URL，直接返回
	return repoURL
}

// ParseWebhookPayload 解析 webhook payload
func (c *Client) ParseWebhookPayload(payload []byte) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, fmt.Errorf("failed to parse payload: %v", err)
	}
	return data, nil
}
