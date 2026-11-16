package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// brewKind は Brewfile の行の種類（brew = formula / cask）
type brewKind string

const (
	formula brewKind = "brew"
	cask    brewKind = "cask"
)

func kindOf(isCask bool) brewKind {
	if isCask {
		return cask
	}
	return formula
}

// brewFlags は brew install / info / uninstall に渡す種類のフラグ
func (k brewKind) brewFlags() []string {
	if k == cask {
		return []string{"--cask"}
	}
	return []string{"--formula"}
}

// section は追加先セクションの見出し（"# --- " で始まる行）に含まれる文字列
func (k brewKind) section() string {
	if k == cask {
		return "GUI"
	}
	return "CLI"
}

// label はコミットメッセージ用の種類名（brew: add cask arc のように使う）
func (k brewKind) label() string {
	if k == cask {
		return "cask "
	}
	return ""
}

func entryPattern(k brewKind, name string) *regexp.Regexp {
	return regexp.MustCompile(`^\s*` + string(k) + `\s+"` + regexp.QuoteMeta(name) + `"`)
}

func brewAdd(k brewKind, names []string, comment string) error {
	dir, err := repoDir()
	if err != nil {
		return err
	}
	brewfile := filepath.Join(dir, "Brewfile")
	lines, err := readLines(brewfile)
	if err != nil {
		return err
	}

	// すでに Brewfile にあるものは飛ばし、残りは Homebrew にあるか確認する
	var targets []string
	for _, name := range names {
		if containsEntry(lines, k, name) {
			fmt.Printf("%s %q はすでに Brewfile にあります\n", k, name)
			continue
		}
		if !succeeds(dir, "brew", append([]string{"info"}, append(k.brewFlags(), name)...)...) {
			return fmt.Errorf("%s %q が見つかりません", k, name)
		}
		targets = append(targets, name)
	}
	if len(targets) == 0 {
		return nil
	}

	if err := run(dir, "brew", append([]string{"install"}, append(k.brewFlags(), targets...)...)...); err != nil {
		return err
	}

	for _, name := range targets {
		entry := fmt.Sprintf("%s %q", k, name)
		if comment != "" {
			entry += "  # " + comment
		}
		lines = insertIntoSection(lines, k.section(), entry)
	}
	if err := writeLines(brewfile, lines); err != nil {
		return err
	}
	return commit(dir, "brew: add "+k.label()+strings.Join(targets, " "))
}

func brewRemove(k brewKind, names []string) error {
	dir, err := repoDir()
	if err != nil {
		return err
	}
	brewfile := filepath.Join(dir, "Brewfile")
	lines, err := readLines(brewfile)
	if err != nil {
		return err
	}

	for _, name := range names {
		if !containsEntry(lines, k, name) {
			return fmt.Errorf("%s %q は Brewfile にありません", k, name)
		}
		re := entryPattern(k, name)
		kept := lines[:0:0]
		for _, l := range lines {
			if !re.MatchString(l) {
				kept = append(kept, l)
			}
		}
		lines = kept
	}
	if err := writeLines(brewfile, lines); err != nil {
		return err
	}
	if err := commit(dir, "brew: remove "+k.label()+strings.Join(names, " ")); err != nil {
		return err
	}

	// Brewfile から消してもアプリ本体は残るので、消すかどうかは別に確認する
	if !confirm(strings.Join(names, " ") + " をアンインストールしますか?") {
		fmt.Println("アンインストールはしませんでした")
		return nil
	}
	return run(dir, "brew", append([]string{"uninstall"}, append(k.brewFlags(), names...)...)...)
}

func containsEntry(lines []string, k brewKind, name string) bool {
	re := entryPattern(k, name)
	for _, l := range lines {
		if re.MatchString(l) {
			return true
		}
	}
	return false
}

// insertIntoSection は見出しに key を含むセクションの最後の行の後に entry を入れる。
// セクションが見つからなければファイルの末尾に足す。
func insertIntoSection(lines []string, key, entry string) []string {
	start := -1
	for i, l := range lines {
		if strings.HasPrefix(l, "# --- ") && strings.Contains(l, key) {
			start = i
			break
		}
	}
	if start < 0 {
		return append(lines, entry)
	}

	// 次のセクション見出し（またはファイル末尾）の手前にある空行を飛ばした位置に入れる
	end := len(lines)
	for i := start + 1; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "# --- ") {
			end = i
			break
		}
	}
	at := end
	for at > start+1 && strings.TrimSpace(lines[at-1]) == "" {
		at--
	}

	out := make([]string, 0, len(lines)+1)
	out = append(out, lines[:at]...)
	out = append(out, entry)
	return append(out, lines[at:]...)
}

func readLines(path string) ([]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimRight(string(b), "\n"), "\n"), nil
}

func writeLines(path string, lines []string) error {
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644)
}
