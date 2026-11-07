package main

import (
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func newMiseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mise",
		Short: "mise のグローバルツール（config/mise/config.toml）を追加・削除する",
	}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "add <ツール>...",
			Short: "ツールを追加してインストールし、そのままコミットする",
			Long: `ツールを config/mise/config.toml に追加してインストールし、差分を確認してコミットする。
名前は先に mise registry | grep <名前> で確認する。バージョンは ツール@版 で指定できる。`,
			Example: `  dots mise add lazygit
  dots mise add cargo:cargo-watch
  dots mise add terraform@1.9`,
			Args: cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return miseEdit("use", "add", args)
			},
		},
		&cobra.Command{
			Use:     "rm <ツール>...",
			Aliases: []string{"remove"},
			Short:   "ツールを削除し、そのままコミットする",
			Long:    `ツールを config/mise/config.toml から削除し（使われなくなった版はアンインストール）、差分を確認してコミットする。`,
			Example: `  dots mise rm lazygit`,
			Args:    cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return miseEdit("unuse", "remove", args)
			},
		},
	)
	return cmd
}

// miseEdit は mise use / unuse でリポジトリ内の設定ファイルを直接書き換えてコミットする
// （~/.config/mise はリポジトリへのリンクだが、DOTFILES_DIR を変えたときも正しく動くようパスを渡す）
func miseEdit(subcmd, verb string, tools []string) error {
	dir, err := repoDir()
	if err != nil {
		return err
	}
	config := filepath.Join(dir, "config", "mise", "config.toml")
	args := append([]string{subcmd, "--path", config}, tools...)
	if err := run(dir, "mise", args...); err != nil {
		return err
	}
	return commit(dir, "mise: "+verb+" "+strings.Join(tools, " "))
}
