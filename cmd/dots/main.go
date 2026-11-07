// dots — dotfiles リポジトリを操作するコマンド
//
// ビルドとインストールは scripts/install.sh が行う（~/.local/bin/dots に置かれる）。
// ソースを変えたら `dots install` でビルドし直す。
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// assumeYes が true なら確認プロンプトをすべて y として扱う（--yes / -y）
var assumeYes bool

func main() {
	root := &cobra.Command{
		Use:   "dots",
		Short: "dotfiles リポジトリを操作する",
		Long: `dots — dotfiles リポジトリを操作するコマンド

どこから実行しても dotfiles リポジトリに対して動く（今いるディレクトリは変わらない）。
コミット・push の前には内容を見せて確認する（-y で確認を省略）。

リポジトリの場所は次の順で決まる:
  1. 環境変数 DOTFILES_DIR
  2. ビルド時に埋め込まれたパス（install.sh を実行したリポジトリ）
  3. ~/dotfiles`,
		Example: `  dots status                 # 未コミットの変更・未 push のコミットを確認
  dots commit "add lazygit"   # 差分を確認してコミット
  dots push                   # 未 push のコミットを確認して push
  dots pull                   # pull して、必要なら install
  dots mise add lazygit       # mise に追加してコミット
  dots brew add --cask arc    # Brewfile に追加してインストールしてコミット
  cd "$(dots path)"           # dotfiles に移動`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.PersistentFlags().BoolVarP(&assumeYes, "yes", "y", false, "確認プロンプトをすべて y として進める")

	root.AddCommand(
		newCommitCmd(),
		newPushCmd(),
		newPullCmd(),
		newStatusCmd(),
		newInstallCmd(),
		newPathCmd(),
		newMiseCmd(),
		newBrewCmd(),
	)

	if err := root.Execute(); err != nil {
		if !errors.Is(err, errCanceled) {
			fmt.Fprintln(os.Stderr, "dots:", err)
		}
		os.Exit(1)
	}
}
