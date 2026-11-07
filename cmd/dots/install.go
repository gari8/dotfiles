package main

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newInstallCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "install",
		Short: "セットアップスクリプト（scripts/install.sh）を実行する",
		Long: `セットアップスクリプト scripts/install.sh を実行する。

シンボリックリンクの作成、brew bundle、mise install、dots 自身のビルドを行う。
何度実行しても安全（冪等）。新しい Mac では dots がまだ無いので scripts/install.sh を直接実行する。`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := repoDir()
			if err != nil {
				return err
			}
			return install(dir, dryRun)
		},
	}
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "何が行われるかの確認のみ（変更なし）")
	return cmd
}

func install(dir string, dryRun bool) error {
	args := []string{}
	if dryRun {
		args = append(args, "--dry-run")
	}
	return run(dir, filepath.Join(dir, "scripts", "install.sh"), args...)
}

func newPathCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "path",
		Short:   "dotfiles リポジトリのパスを表示する",
		Example: `  cd "$(dots path)"`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := repoDir()
			if err != nil {
				return err
			}
			fmt.Println(dir)
			return nil
		},
	}
}
