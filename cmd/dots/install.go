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
		Short: "scripts/install.sh を実行する",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := repoDir()
			if err != nil {
				return err
			}
			return install(dir, dryRun)
		},
	}
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "確認のみ（変更しない）")
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
		Short:   "dotfiles のパスを表示する",
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
