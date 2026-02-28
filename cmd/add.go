package cmd

import (
	"bufio"
	"fmt"
	"os"
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

		fmt.Print("服务商名称 (如 我的提供商): ")
		displayName, _ := reader.ReadString('\n')
		displayName = strings.TrimSpace(displayName)

		fmt.Print("Base URL: ")
		baseURL, _ := reader.ReadString('\n')
		baseURL = strings.TrimSpace(baseURL)

		fmt.Print("Opus 模型: ")
		opusModel, _ := reader.ReadString('\n')
		opusModel = strings.TrimSpace(opusModel)

		fmt.Print("Sonnet 模型: ")
		sonnetModel, _ := reader.ReadString('\n')
		sonnetModel = strings.TrimSpace(sonnetModel)

		fmt.Print("Haiku 模型: ")
		haikuModel, _ := reader.ReadString('\n')
		haikuModel = strings.TrimSpace(haikuModel)

		fmt.Print("API Key: ")
		apiKey, _ := reader.ReadString('\n')
		apiKey = strings.TrimSpace(apiKey)

		p := provider.Provider{
			Name:    displayName,
			BaseURL: baseURL,
			Models: provider.ModelConfig{
				Opus:   opusModel,
				Sonnet: sonnetModel,
				Haiku:  haikuModel,
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
