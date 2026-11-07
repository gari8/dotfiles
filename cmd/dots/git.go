package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newCommitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "commit <メッセージ>",
		Short: "差分を確認してコミットする（git add -A してコミット）",
		Example: `  dots commit "add lazygit"
  dots commit ghostty: フォントを変更   # クォートしなくてもよい`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := repoDir()
			if err != nil {
				return err
			}
			return commit(dir, strings.Join(args, " "))
		},
	}
}

// commit は変更内容を見せてから確認し、すべての変更をコミットする
func commit(dir, message string) error {
	st, err := output(dir, "git", "status", "--porcelain")
	if err != nil {
		return err
	}
	if st == "" {
		fmt.Println("dotfiles に変更はありません")
		return nil
	}

	// 新規ファイルは diff に出ないので status も見せる
	if err := run(dir, "git", "status", "--short"); err != nil {
		return err
	}
	if err := run(dir, "git", "diff"); err != nil {
		return err
	}

	if !confirm("この内容でコミットしますか?") {
		return canceled("コミット")
	}
	if err := run(dir, "git", "add", "-A"); err != nil {
		return err
	}
	return run(dir, "git", "commit", "-m", message)
}

func newPushCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "push",
		Short: "まだ push していないコミットを確認して push する",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := repoDir()
			if err != nil {
				return err
			}
			return push(dir)
		},
	}
}

func push(dir string) error {
	// 未コミットの変更は push されないことを知らせる（push 自体は止めない）
	if st, _ := output(dir, "git", "status", "--short"); st != "" {
		fmt.Println("未コミットの変更があります（push されません。先に dots commit でコミットできます）")
		fmt.Println(st)
	}

	// upstream が無い（初回 push）ときは origin に同名ブランチで push する
	if !hasUpstream(dir) {
		if !confirm("upstream がありません。origin に push しますか?") {
			return canceled("push")
		}
		return run(dir, "git", "push", "-u", "origin", "HEAD")
	}

	commits, err := output(dir, "git", "log", "--oneline", "@{u}..HEAD")
	if err != nil {
		return err
	}
	if commits == "" {
		fmt.Println("push するコミットはありません")
		return nil
	}
	fmt.Println(commits)
	if !confirm("これらのコミットを push しますか?") {
		return canceled("push")
	}
	return run(dir, "git", "push")
}

func newPullCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pull",
		Short: "リモートの変更を取り込む（fast-forward のみ。反映はしない）",
		Long: `リモートの変更を取り込む（git pull --ff-only）。

fast-forward できない（手元とリモートの両方にコミットがある）ときは何もせず止まる。
Brewfile や mise の設定は反映しないので、反映もするなら dots sync を使う。`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := repoDir()
			if err != nil {
				return err
			}
			changed, err := pull(dir)
			if err != nil {
				return err
			}
			if changed {
				fmt.Println("\n取り込んだ設定を反映するには dots install を実行してください")
			}
			return nil
		},
	}
}

func newSyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "リモートの変更を取り込み、dots install で反映する（pull + install）",
		Long: `dots pull でリモートの変更を取り込んでから、dots install で反映する。

install は冪等なので、取り込む変更が無くても実行する（手元の設定を反映し直すのにも使える）。
pull が失敗したときは install しない。`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := repoDir()
			if err != nil {
				return err
			}
			if _, err := pull(dir); err != nil {
				return err
			}
			fmt.Println()
			return install(dir, false)
		},
	}
}

// pull は fast-forward でリモートの変更を取り込み、HEAD が進んだかを返す
func pull(dir string) (bool, error) {
	if st, _ := output(dir, "git", "status", "--short"); st != "" {
		fmt.Println("未コミットの変更があります（取り込む変更と重なると pull は失敗します）")
		fmt.Println(st)
	}

	before, err := output(dir, "git", "rev-parse", "HEAD")
	if err != nil {
		return false, err
	}
	// 取り込んだファイルの一覧は git pull が表示する
	if err := run(dir, "git", "pull", "--ff-only"); err != nil {
		return false, err
	}
	after, err := output(dir, "git", "rev-parse", "HEAD")
	if err != nil {
		return false, err
	}
	return before != after, nil
}

func newStatusCmd() *cobra.Command {
	var fetch bool
	cmd := &cobra.Command{
		Use:   "status",
		Short: "未コミットの変更と、まだ push していないコミットを表示する",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := repoDir()
			if err != nil {
				return err
			}
			if fetch {
				if err := run(dir, "git", "fetch", "--quiet"); err != nil {
					return err
				}
			}
			// 1行目にブランチと ahead / behind が出る
			if err := run(dir, "git", "status", "--short", "--branch"); err != nil {
				return err
			}
			if !hasUpstream(dir) {
				fmt.Println("（upstream がありません。dots push で origin に push できます）")
				return nil
			}
			commits, err := output(dir, "git", "log", "--oneline", "@{u}..HEAD")
			if err != nil {
				return err
			}
			if commits != "" {
				fmt.Println("\nまだ push していないコミット:")
				fmt.Println(commits)
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&fetch, "fetch", "f", false, "先に git fetch してリモートの最新と比べる")
	return cmd
}

func hasUpstream(dir string) bool {
	return succeeds(dir, "git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
}
