package cmd

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"cc-switch/internal/config"
	"cc-switch/internal/provider"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "添加自定义服务商",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		var name string
		if len(args) > 0 {
			name = args[0]
		} else {
			fmt.Print("服务商标识 (如 my-provider): ")
			name, _ = reader.ReadString('\n')
			name = strings.TrimSpace(name)
		}

		if name == "" {
			fmt.Println("错误: 服务商标识不能为空")
			return
		}

		// 验证服务商标识格式：只允许字母、数字、下划线、连字符
		if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(name) {
			fmt.Println("错误: 服务商标识只能包含字母、数字、下划线和连字符")
			return
		}

		fmt.Print("服务商名称 (如 我的提供商): ")
		displayName, _ := reader.ReadString('\n')
		displayName = strings.TrimSpace(displayName)

		fmt.Print("Base URL: ")
		baseURL, _ := reader.ReadString('\n')
		baseURL = strings.TrimSpace(baseURL)

		// 验证 BaseURL 格式
		if _, err := url.ParseRequestURI(baseURL); err != nil {
			fmt.Printf("错误: Base URL 格式无效: %v\n", err)
			return
		}

		fmt.Print("Opus 模型: ")
		opusModel, _ := reader.ReadString('\n')
		opusModel = strings.TrimSpace(opusModel)

		fmt.Print("Sonnet 模型: ")
		sonnetModel, _ := reader.ReadString('\n')
		sonnetModel = strings.TrimSpace(sonnetModel)

		fmt.Print("Haiku 模型: ")
		haikuModel, _ := reader.ReadString('\n')
		haikuModel = strings.TrimSpace(haikuModel)

		fmt.Print("SmallFast 模型 (可选，回车使用 default_haiku): ")
		smallFastModel, _ := reader.ReadString('\n')
		smallFastModel = strings.TrimSpace(smallFastModel)
		if smallFastModel == "" {
			smallFastModel = haikuModel
		}

		fmt.Print("Default 模型 (可选，回车使用 default_sonnet): ")
		defaultModel, _ := reader.ReadString('\n')
		defaultModel = strings.TrimSpace(defaultModel)
		if defaultModel == "" {
			defaultModel = sonnetModel
		}

		fmt.Print("API Key: ")
		apiKey, _ := reader.ReadString('\n')
		apiKey = strings.TrimSpace(apiKey)

		p := provider.Provider{
			Name:    displayName,
			BaseURL: baseURL,
			Models: provider.ModelConfig{
				DefaultOpus:   opusModel,
				DefaultSonnet: sonnetModel,
				DefaultHaiku:  haikuModel,
				SmallFast:     smallFastModel,
				DefaultModel:  defaultModel,
			},
			APIKey: apiKey,
		}

		if err := config.AddProvider(name, p); err != nil {
			fmt.Fprintf(os.Stderr, "添加失败: %v\n", err)
			return
		}

		fmt.Printf("已添加服务商: %s\n", name)
	},
}
