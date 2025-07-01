package claude

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/xhy/auto-coding/internal/types"
)

// Client Claude 集成客户端
type Client struct {
	config *types.BotConfig
}

// NewClient 创建新的 Claude 客户端
func NewClient(config *types.BotConfig) *Client {
	return &Client{
		config: config,
	}
}

// Invoke 调用 Claude Code CLI - 完全自主模式
func (c *Client) Invoke(workspace string) error {
	log.Printf("Invoking Claude in autonomous mode for workspace: %s", workspace)

	// 检查是否有仓库目录
	repoDir := workspace + "/repository"
	workDir := workspace
	contextPath := ".claude-context.json"

	if _, err := os.Stat(repoDir); err == nil {
		// 仓库存在，在仓库目录内运行
		workDir = repoDir
		contextPath = "../.claude-context.json"
		log.Printf("Repository found, working in: %s", workDir)
	} else {
		// 仓库不存在，在工作空间根目录运行
		log.Printf("No repository found, working in: %s", workDir)
	}

	// 🔥 根据事件类型的具体工作流程 prompt
	prompt := fmt.Sprintf(`你是一个 GitHub 机器人，收到了一个事件。

检查 %s 了解事件详情和类型。

🚨 强制性规则 🚨
❌ 绝不推送到 main/master 分支
❌ 绝不直接关闭 issue
❌ 绝不使用 "Claude Bot" 作为提交作者
❌ 绝不重复评论！每个 issue 只评论一次通知 PR
❌ 绝不在 PR 下添加无意义的评论
❌ 绝不回复表扬、感谢、客套话
❌ 绝不回复自己之前的评论（避免循环）
❌ 绝不在创建 PR 后再次评论解释
❌ Issues 处理完毕后绝不再追加任何评论

🔄 根据事件类型的处理方式：

【Issues 事件 - 新需求】
1. 配置 Bot 专用身份:
   git config user.name "%s"
   git config user.email "%s"
2. 创建新分支: git checkout -b fix-issue-N
3. 实现功能/修复
4. 提交: git add . && git commit -m "描述"
5. 推送分支: git push origin fix-issue-N
6. 创建 PR: gh pr create --title "标题" --body "Fixes #N"
7. 在原始 issue 评论: gh issue comment N --body "已创建 PR #X，请审查"
8. 🚨 完成！绝不在 PR 下评论！绝不追加说明！

【PR Review/Comment 事件 - 极度谨慎处理】
🚨🚨🚨 CRITICAL: 检查评论作者！
- 如果评论者是你自己（机器人），立即停止！直接忽略！
- 如果评论者用户名包含 "bot"、"claude"、"github-actions"，立即忽略！

仅对人类用户的评论做以下判断：

A) 明确的代码修改请求（如"请添加错误处理"、"修复这个bug"）：
1. 配置 Bot 身份: git config user.name "%s" && git config user.email "%s"
2. 找到 PR 对应分支名称
3. 切换到该分支: git checkout existing-branch-name  
4. 根据反馈修改代码
5. 提交: git add . && git commit -m "根据反馈修改: 具体描述"
6. 推送: git push origin existing-branch-name
7. 🚨 用代码更改回应，绝不添加任何文字回复！

B) 其他所有情况（技术询问、表扬、感谢、讨论、新需求等）：
1. 🚨 完全忽略！不要回应！
2. 🚨 这包括："谢谢"、"很好"、"这代码什么意思"、"能否添加新功能"等
3. 🚨 默认策略：什么都不做！

【Issue Comment 事件】
- 如果是对现有 issue 的澄清，简单回复即可
- 如果需要修改现有 PR，按照 PR Comment 流程处理

关键：不同事件用不同处理方式，不要混淆！`,
		contextPath,
		c.config.BotName, c.config.BotEmail, // Issues 事件的身份配置
		c.config.BotName, c.config.BotEmail) // PR 修改事件的身份配置

	// 创建 Claude Code CLI 命令，跳过所有权限检查
	cmd := exec.Command("npx", "@anthropic-ai/claude-code", "--dangerously-skip-permissions")
	cmd.Dir = workDir

	// 设置环境变量，给 Claude 完整的上下文
	cmd.Env = append(os.Environ(),
		"GITHUB_TOKEN="+c.config.GitHubToken,
		"ANTHROPIC_API_KEY="+c.config.AnthropicKey,
		"CLAUDE_WORKSPACE="+workspace,
		"CLAUDE_MODE=autonomous",
	)

	// 通过 stdin 提供 prompt
	cmd.Stdin = strings.NewReader(prompt)

	// 让 Claude 完全自主运行
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Claude execution failed: %v, output: %s", err, string(output))
		return err
	}

	log.Printf("Claude completed successfully: %s", string(output))
	return nil
}
