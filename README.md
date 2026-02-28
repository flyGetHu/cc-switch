# cc-switch

Claude Code 模型服务商切换工具

## 安装

```bash
go build -o cc-switch .
```

或下载预编译版本（如有）

## 使用

```bash
# 交互式选择服务商
cc-switch

# 直接切换
cc-switch use zhipu
cc-switch use minimax

# 配置 API Key
cc-switch edit zhipu -k "your-api-key"
cc-switch edit minimax -k "your-api-key"

# 查看当前服务商
cc-switch current

# 列出所有服务商
cc-switch list

# 添加自定义服务商
cc-switch add

# 编辑服务商
cc-switch edit <name>

# 删除服务商
cc-switch remove <name>

# 备份管理
cc-switch backup create
cc-switch backup list
cc-switch backup restore
```

## 预设服务商

| 服务商 | Base URL | Opus | Sonnet | Haiku |
|--------|----------|------|--------|-------|
| 智普 GLM | `https://open.bigmodel.cn/api/anthropic` | glm-5 | glm-5 | glm-4.7-flash |
| MiniMax | `https://api.minimaxi.com/anthropic` | MiniMax-M2.5 | MiniMax-M2.5 | MiniMax-M2.5 |

## 配置文件

- 工具配置: `~/.config/cc-switch/config.yaml`
- Claude Code 配置: `~/.claude/settings.json`
- 备份目录: `~/.config/cc-switch/backups/`

## 开发

```bash
# 运行测试
go test ./...

# 运行测试（带覆盖率）
go test ./... -cover

# 构建
go build -o cc-switch .
```

## License

MIT
