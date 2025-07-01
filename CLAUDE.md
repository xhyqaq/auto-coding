# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **fully autonomous Claude GitHub Bot** that receives ALL GitHub events and lets Claude decide how to handle them. There are NO restrictions, NO predefined workflows, and NO limitations on Claude's decision-making. Claude has complete autonomy to analyze issues, implement fixes, and create PRs as it sees fit.

## Development Commands

### Building and Running
```bash
# Install dependencies
go mod tidy

# Run the application
go run main.go

# Build binary
go build -o claude-github-bot main.go
```

### Environment Setup
```bash
# Copy environment template
cp .env.example .env
# Edit .env with your tokens and configuration
```

### Required Environment Variables
- `GITHUB_TOKEN`: GitHub Personal Access Token with full repo permissions
- `WEBHOOK_SECRET`: Secret for webhook validation
- `PORT`: Server port (default: 8080)

### Prerequisites
- Go 1.21+
- Node.js (for Claude Code CLI)
- Claude Code CLI: `npm install -g @anthropic-ai/claude-code`
- 必须需要用中文回答我

## Autonomous Architecture

### Core Philosophy
**Zero Constraints, Maximum Autonomy**: Claude receives ALL GitHub events and decides independently how to respond. The Go service is merely a delivery mechanism.

### Key Components
- **Event Forwarder** (`main.go:webhookHandler`): Receives ALL GitHub events without filtering
- **Workspace Manager** (`main.go:createWorkspace`): Creates clean working environments for Claude
- **Context Provider** (`main.go:setupClaudeContext`): Provides full event context to Claude
- **Claude Invoker** (`main.go:invokeClaude`): Launches Claude in complete autonomous mode

### Event Handling
1. **All Events Accepted**: No filtering by event type, labels, or content
2. **Full Context Provided**: Complete event payload and repository access
3. **Repository Cloning**: For issue events, automatically clones repository for Claude
4. **Autonomous Execution**: Claude runs with zero constraints or guidance

### Claude's Full Capabilities
When Claude receives an event, it has access to:
- **Full Repository Access**: Can read, modify, and understand entire codebase
- **GitHub CLI (gh)**: Direct access to GitHub operations via `gh pr create`, `gh issue comment`, etc.
- **Git Operations**: Can create branches, commit changes, push to GitHub
- **File System Access**: Can modify any files in the workspace
- **No Approval Required**: Can take any action it deems appropriate

### Available Tools
- **GitHub CLI (gh)**: Authenticated and ready to use for all GitHub operations
- **Git**: Standard git commands for version control
- **File System**: Full read/write access to repository files

## Claude's Autonomous Workflow for Issues

When an issue is created, Claude should:

1. **Analyze the Issue**: Understand the problem, requirements, or request
2. **Examine the Codebase**: Review relevant files and understand the architecture  
3. **Make Independent Decisions**: Decide if/how to address the issue
4. **Configure Git Identity**: Set git user as repository owner using `gh api user`
5. **Create Feature Branch**: Never work on main/master directly
6. **Implement Solutions**: Write code, fix bugs, add features as needed
7. **Commit Changes**: Create meaningful commits with proper author info
8. **Create Pull Request**: Use `gh pr create` linking to issue
9. **Notify on Issue**: Comment on original issue about the PR

## Development Guidelines

This file contains guidance for developing the Claude GitHub Bot itself.

### Code Structure
- Keep the bot logic simple and focused
- Separate concerns: webhook handling, event processing, Claude integration
- Use clear, descriptive function and variable names

### Bot Configuration
- Environment variables should be clearly documented
- Support multiple deployment environments
- Ensure graceful error handling for missing configurations

### Testing Strategy  
- Test webhook payload parsing
- Mock Claude API calls for unit tests
- Test workspace creation and cleanup

### Deployment Considerations
- Webhook URL configuration
- Token security and rotation
- Resource cleanup (temporary workspaces)
- Logging and monitoring

## Bot Workflow Notes

**Note**: The actual GitHub event handling workflows are defined in the bot's runtime prompts, not in this file. This file is for developing the bot itself.

## Architecture Notes

The bot creates temporary workspaces for each GitHub event and provides Claude with:
- Event context via `.claude-context.json`
- Cloned repository (when needed)
- GitHub API access via environment tokens

The goal is minimal constraints and maximum Claude autonomy for handling GitHub events.