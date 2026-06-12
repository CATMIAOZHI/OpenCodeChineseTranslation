package cmd

import (
	"fmt"
	"opencode-cli/internal/core"
	"os"

	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "应用汉化补丁到源码",
	Run: func(cmd *cobra.Command, args []string) {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		silent, _ := cmd.Flags().GetBool("silent")

		i18n, err := core.NewI18n()
		if err != nil {
			fmt.Printf("错误: 初始化失败: %v\n", err)
			os.Exit(1)
		}

		configs, err := i18n.LoadConfig()
		if err != nil {
			fmt.Printf("错误: 加载配置失败: %v\n", err)
			os.Exit(1)
		}

		if !silent {
			if dryRun {
				fmt.Println("模拟应用汉化配置...")
			} else {
				fmt.Println("应用汉化配置...")
			}
			fmt.Printf("找到 %d 个配置文件\n", len(configs))
		}

		stats := struct {
			Files struct {
				Total   int
				Success int
				Skipped int
				Failed  int
			}
			Replacements struct {
				Total   int
				Success int
				Failed  int
			}
		}{}

		for _, config := range configs {
			result := i18n.ApplyConfig(config, dryRun)

			stats.Files.Total++
			if result.Skipped {
				stats.Files.Skipped++
				if !silent {
					fmt.Printf("  ⊘ %s (跳过: %s)\n", config.File, result.SkipReason)
				}
			} else if result.Success {
				stats.Files.Success++
				if !silent {
					fmt.Printf("  ✓ %s (%d/%d 处替换)\n", config.File, result.Replacements.Success, result.Replacements.Total)
				}
			} else {
				stats.Files.Failed++
				if !silent {
					fmt.Printf("  ✗ %s 失败 (%d/%d 处替换)\n", config.File, result.Replacements.Success, result.Replacements.Total)
				}
			}

			stats.Replacements.Total += result.Replacements.Total
			stats.Replacements.Success += result.Replacements.Success
			stats.Replacements.Failed += result.Replacements.Failed
		}

		if !silent {
			fmt.Println("")
			if dryRun {
				fmt.Println("汉化模拟完成:")
			} else {
				fmt.Println("汉化应用完成:")
			}
			fmt.Printf("  📁 文件: %d 成功, %d 跳过, %d 失败\n", stats.Files.Success, stats.Files.Skipped, stats.Files.Failed)
			fmt.Printf("  📝 替换: %d/%d 成功\n", stats.Replacements.Success, stats.Replacements.Total)
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().Bool("dry-run", false, "Simulate the application without modifying files")
	applyCmd.Flags().Bool("silent", false, "Suppress output")
}
