# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

cc-switch 是一个 Claude Code 模型服务商切换工具，通过修改 `~/.claude/settings.json` 中的环境变量来实现不同服务商（智普 GLM、MiniMax 等）的快速切换。

## 常用命令

```bash
# 构建
go build -o cc-switch .

# 运行测试
go test ./...

# 运行单个测试
go test ./internal/config -v -run TestLoad
go test ./internal/provider -v -run TestGetPresets

# 格式化 + 静态检查（提交前运行）
go fmt ./... && go vet ./...

# 直接运行
go run .
```

## 架构

```
cmd/                    # Cobra CLI 命令（root, use, list, current, add, edit, remove, backup）
internal/
├── provider/          # Provider 和 ModelConfig 类型定义，预设服务商
├── config/            # 工具配置管理（~/.config/cc-switch/config.yaml）
├── claude/            # Claude Code settings.json 读写，ApplyProvider() 写入环境变量
├── backup/            # 备份管理
└── tui/               # Bubbletea 交互式选择器
```

**核心流程**：`use` 命令调用 `config.Load()` 获取服务商配置，然后调用 `claude.ApplyProvider()` 将 BaseURL、APIKey 和模型名称写入 `~/.claude/settings.json` 的 env 字段。

## 配置文件位置

- 工具配置: `~/.config/cc-switch/config.yaml`
- Claude Code 配置: `~/.claude/settings.json`
- 备份目录: `~/.config/cc-switch/backups/`

## 代码风格

详见 [AGENTS.md](AGENTS.md)，包括导入顺序、命名约定、错误处理和测试风格。
