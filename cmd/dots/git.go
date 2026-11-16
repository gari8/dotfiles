package main

import "fmt"

// git の操作はカスタムコマンド（dots.toml の dots cmd）で行う。
// ここには mise / brew の追加・削除が内部で使うコミットだけを置く。

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
