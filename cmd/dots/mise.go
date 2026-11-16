package main

import (
	"path/filepath"
	"strings"
)

// miseEdit は mise use / unuse でリポジトリ内の設定ファイルを直接書き換えてコミットする
// （~/.config/mise はリポジトリへのリンクだが、DOTS_ROOT を変えたときも正しく動くようパスを渡す）
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
