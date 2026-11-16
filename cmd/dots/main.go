// dots — dotfiles リポジトリを操作するコマンド
//
// ビルドとインストールは scripts/install.sh が行う（~/.local/bin/dots に置かれる）。
// ソースを変えたら `dots install` でビルドし直す。
package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// assumeYes が true なら確認プロンプトをすべて y として扱う（--yes / -y）
var assumeYes bool

func main() {
	root := &cobra.Command{
		Use:           "dots",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.PersistentFlags().BoolVarP(&assumeYes, "yes", "y", false, "確認を省略する")
	// 補完スクリプトは completion.zsh が読み込むだけなので、help の一覧には出さない
	root.CompletionOptions.HiddenDefaultCmd = true
	// help は --help / -h で足りるので、サブコマンドの一覧には出さない
	// （cobra の一覧テンプレートは help を隠していても表示するので、その条件を外す）
	root.InitDefaultHelpCmd()
	for _, c := range root.Commands() {
		if c.Name() == "help" {
			c.Hidden = true
		}
	}
	root.SetUsageTemplate(strings.ReplaceAll(root.UsageTemplate(),
		`(or .IsAvailableCommand (eq .Name "help"))`, `.IsAvailableCommand`))

	root.AddCommand(
		newUpgradeCmd(),
		newListCmd(),
		newInstallCmd(),
		newPathCmd(),
		newAddCmd(),
		newRmCmd(),
		newCmdCmd(),
	)

	if err := root.Execute(); err != nil {
		// カスタムコマンドの終了コードはそのまま返す（メッセージは子が出している）
		var code exitCodeError
		if errors.As(err, &code) {
			os.Exit(int(code))
		}
		if !errors.Is(err, errCanceled) {
			fmt.Fprintln(os.Stderr, "dots:", err)
		}
		os.Exit(1)
	}
}
