# AGENTS.md

本文档为 cc-switch 项目提供开发指南，供 AI 编码代理参考。

## 项目概述

cc-switch 是一个 Claude Code 模型服务商切换工具，支持智普 GLM 和 MiniMax 等服务商的快速切换。

### 技术栈

- **语言**: Go 1.25+
- **CLI 框架**: github.com/spf13/cobra
- **配置解析**: gopkg.in/yaml.v3
- **TUI**: github.com/charmbracelet/bubbletea + lipgloss

### 项目结构

```
cc-switch/
├── cmd/                    # CLI 命令实现
│   ├── root.go            # 入口和交互式选择
│   ├── use.go             # 切换服务商
│   ├── list.go            # 列出服务商
│   ├── current.go         # 显示当前服务商
│   ├── add.go             # 添加自定义服务商
│   ├── edit.go            # 编辑服务商配置
│   ├── remove.go          # 删除服务商
│   └── backup.go          # 备份管理
├── internal/
│   ├── provider/          # 服务商定义和预设
│   ├── config/            # 工具配置管理
│   ├── claude/            # Claude Code 配置操作
│   ├── backup/            # 备份管理
│   └── tui/               # 交互式终端 UI
├── main.go                # 程序入口
└── go.mod
```

## 构建与测试命令

### 构建

```bash
# 构建二进制
go build -o cc-switch .

# 安装依赖
go mod tidy
```

### 测试

```bash
# 运行所有测试
go test ./...

# 运行测试并显示覆盖率
go test ./... -cover

# 运行测试并显示详细输出
go test ./... -v

# 运行单个测试文件
go test ./internal/config/config_test.go -v

# 运行单个测试函数
go test ./internal/provider -v -run TestGetPresets

# 运行匹配模式的测试
go test ./... -v -run "Test.*Config"

# 生成覆盖率报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Lint 和格式化

```bash
# 格式化代码
go fmt ./...

# 静态检查
go vet ./...

# 格式化 + 检查（推荐在提交前运行）
go fmt ./... && go vet ./...
```

### 运行

```bash
# 直接运行
go run .

# 或使用构建后的二进制
./cc-switch
./cc-switch list
./cc-switch use zhipu
./cc-switch current
```

## 代码风格指南

### 导入顺序

导入按以下顺序分组，组间用空行分隔：

1. 标准库
2. 项目内部包
3. 第三方库

```go
import (
    "fmt"
    "os"
    "path/filepath"

    "cc-switch/internal/provider"

    "github.com/spf13/cobra"
    "gopkg.in/yaml.v3"
)
```

### 命名约定

- **包名**: 小写单词，不使用下划线或驼峰 (如 `config`, `backup`, `claude`)
- **类型名**: 大写驼峰，导出类型首字母大写 (如 `Config`, `Provider`, `ModelConfig`)
- **函数名**: 
  - 导出函数：大写驼峰 (如 `GetConfigPath`, `Load`, `Save`)
  - 内部函数：小写驼峰 (如 `initDefaultConfig`, `switchProvider`)
- **变量名**: 小写驼峰，简短有意义 (如 `cfg`, `p`, `name`)
- **常量/错误**: 大写驼峰或 Err 前缀 (如 `ErrProviderNotFound`)

### 类型定义

```go
// 结构体字段使用 yaml 标签
type Config struct {
    Providers  map[string]provider.Provider `yaml:"providers"`
    Current    string                       `yaml:"current"`
    BackupsDir string                       `yaml:"backups_dir"`
}

// 嵌套结构体单独定义
type ModelConfig struct {
    Opus   string `yaml:"opus"`
    Sonnet string `yaml:"sonnet"`
    Haiku  string `yaml:"haiku"`
}
```

### 错误处理

- 使用 `fmt.Errorf` 创建错误，使用 `%w` 包装底层错误
- 返回错误而非 panic
- 在 CLI 层输出友好错误信息并退出

```go
// 正确：包装错误
if err != nil {
    return fmt.Errorf("加载配置失败: %w", err)
}

// 正确：创建新错误
if p.APIKey == "" {
    return fmt.Errorf("服务商 '%s' 未配置 API Key", name)
}

// CLI 层错误处理
if err := switchProvider(name); err != nil {
    fmt.Fprintf(os.Stderr, "切换失败: %v\n", err)
    os.Exit(1)
}
```

### 注释风格

- 使用中文注释
- 导出函数应添加注释说明
- 注释以函数名开头

```go
// GetConfigPath 返回配置文件路径
// 如果设置了自定义路径则返回自定义路径，否则返回默认路径
func GetConfigPath() string {
    // ...
}
```

### 测试风格

- 测试文件与源文件同目录，命名为 `*_test.go`
- 使用表驱动测试
- 子测试使用 `t.Run()`
- 测试函数命名为 `Test<功能名>`

```go
func TestValidateProvider(t *testing.T) {
    tests := []struct {
        name    string
        provider *provider.Provider
        wantErr bool
    }{
        {
            name: "valid provider",
            provider: &provider.Provider{
                BaseURL: "https://api.test.com",
                APIKey:  "test-key",
                Models: provider.ModelConfig{
                    Opus: "opus", Sonnet: "sonnet", Haiku: "haiku",
                },
            },
            wantErr: false,
        },
        // 更多测试用例...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateProvider(tt.provider)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateProvider() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### 可测试性设计

对于需要测试的包级变量（如路径函数），使用变量函数并在测试时替换：

```go
// 源码中
var getBackupsDir = config.GetBackupsDir

// 测试中
func setupTest(t *testing.T) func() {
    oldGetBackupsDir := getBackupsDir
    getBackupsDir = func() string { return t.TempDir() }
    return func() { getBackupsDir = oldGetBackupsDir }
}
```

## 配置文件位置

- 工具配置: `~/.config/cc-switch/config.yaml`
- Claude Code 配置: `~/.claude/settings.json`
- 备份目录: `~/.config/cc-switch/backups/`

## 预设服务商

| 服务商 | Base URL | 默认模型 |
|--------|----------|----------|
| 智普 GLM | `https://open.bigmodel.cn/api/anthropic` | glm-5, glm-4.7-flash |
| MiniMax | `https://api.minimaxi.com/anthropic` | MiniMax-M2.5 |

## 提交前检查清单

1. `go fmt ./...` - 格式化代码
2. `go vet ./...` - 静态检查
3. `go test ./...` - 所有测试通过
4. 构建成功: `go build -o cc-switch .`
