package main

import (
	"fmt"

	"zero-backend/cmd/cli/runner"

	"github.com/241x/zero-kit/logger"
	"github.com/spf13/cobra"
)

func syncApiCmd(r *runner.SyncApiRunner) *cobra.Command {
	var (
		filePath   string
		deleteFlag bool
	)

	cmd := &cobra.Command{
		Use:   "sync-api",
		Short: "同步 OpenAPI 接口到数据库",
		Long:  `解析 docs/admin/index.json 文件，将接口信息同步到 rbac_api 表`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			l := logger.Ctx(ctx)
			l.Info("开始同步 API 接口...")
			if deleteFlag {
				l.Info("已启用删除模式，将删除文档中不存在的接口")
			}

			result, err := r.Run(ctx, filePath, deleteFlag)
			if err != nil {
				return err
			}

			cmd.Println("\n===== 同步结果 =====")
			cmd.Printf("文档中的API数: %d\n", result.TotalInDoc)
			cmd.Printf("新增: %d\n", result.Added)
			cmd.Printf("更新: %d\n", result.Updated)
			cmd.Printf("删除: %d\n", result.Deleted)
			cmd.Printf("跳过: %d\n", result.Skipped)

			if len(result.Errors) > 0 {
				cmd.Println("\n===== 错误列表 =====")
				for i, e := range result.Errors {
					cmd.Printf("%d. %s\n", i+1, e)
				}
				return fmt.Errorf("同步过程中有 %d 个错误", len(result.Errors))
			}

			cmd.Println("\n同步完成!")
			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "OpenAPI 文档路径 (默认 docs/admin/index.json)")
	cmd.Flags().BoolVarP(&deleteFlag, "delete", "d", false, "删除数据库中存在但文档中不存在的接口")

	return cmd
}
