package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// defaultDir はビルド時に install.sh が埋め込む dotfiles リポジトリのパス
//
//	go build -ldflags "-X main.defaultDir=/path/to/dotfiles" ./cmd/dots
var defaultDir string

// errCanceled は確認で n を選んだときのエラー（メッセージは表示済みなので出さない）
var errCanceled = errors.New("canceled")

// repoDir は dotfiles リポジトリの場所を返す
func repoDir() (string, error) {
	dir := os.Getenv("DOTFILES_DIR")
	if dir == "" {
		dir = defaultDir
	}
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(home, "dotfiles")
	}
	if _, err := os.Stat(filepath.Join(dir, ".git")); err != nil {
		return "", fmt.Errorf("dotfiles リポジトリが見つかりません: %s（DOTFILES_DIR で指定できます）", dir)
	}
	return dir, nil
}

// run はコマンドを dir で実行する。入出力は端末にそのままつなぐ。
func run(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s %s: %w", name, strings.Join(args, " "), err)
	}
	return nil
}

// output はコマンドを dir で実行し、標準出力を前後の空白を除いて返す
func output(dir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s %s: %w", name, strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(out)), nil
}

// succeeds はコマンドが成功するかだけを調べる（出力は捨てる）
func succeeds(dir, name string, args ...string) bool {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.Run() == nil
}

var stdin = bufio.NewReader(os.Stdin)

// confirm は [y/N] で確認する。--yes なら聞かずに true。
func confirm(prompt string) bool {
	if assumeYes {
		return true
	}
	fmt.Printf("%s [y/N] ", prompt)
	line, _ := stdin.ReadString('\n')
	answer := strings.ToLower(strings.TrimSpace(line))
	return answer == "y" || answer == "yes"
}

// canceled は「〜しませんでした」と表示して errCanceled を返す
func canceled(what string) error {
	fmt.Printf("%sしませんでした\n", what)
	return errCanceled
}
