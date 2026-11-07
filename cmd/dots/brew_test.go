package main

import (
	"strings"
	"testing"
)

const testBrewfile = `# Brewfile

# --- CLI / システム寄りのツール ---
brew "git"
brew "tree"       # ツリー表示

# --- GUIアプリ ---
cask "ghostty"    # ターミナル

# --- Mac App Store ---
mas "Xcode", id: 497799835`

func TestInsertIntoSection(t *testing.T) {
	lines := strings.Split(testBrewfile, "\n")
	lines = insertIntoSection(lines, formula.section(), `brew "hello"`)
	lines = insertIntoSection(lines, cask.section(), `cask "arc"  # ブラウザ`)
	got := strings.Join(lines, "\n")

	want := `# Brewfile

# --- CLI / システム寄りのツール ---
brew "git"
brew "tree"       # ツリー表示
brew "hello"

# --- GUIアプリ ---
cask "ghostty"    # ターミナル
cask "arc"  # ブラウザ

# --- Mac App Store ---
mas "Xcode", id: 497799835`
	if got != want {
		t.Errorf("got:\n%s\n\nwant:\n%s", got, want)
	}
}

func TestInsertIntoSectionMissing(t *testing.T) {
	lines := insertIntoSection([]string{`brew "git"`}, "GUI", `cask "arc"`)
	if got := strings.Join(lines, "\n"); got != "brew \"git\"\ncask \"arc\"" {
		t.Errorf("セクションが無いときは末尾に足す: got %q", got)
	}
}

func TestContainsEntry(t *testing.T) {
	lines := strings.Split(testBrewfile, "\n")
	cases := []struct {
		k    brewKind
		name string
		want bool
	}{
		{formula, "tree", true},
		{formula, "tre", false}, // 前方一致では一致しない
		{cask, "tree", false},   // 種類が違えば一致しない
		{cask, "ghostty", true},
		{formula, "ghostty", false},
	}
	for _, c := range cases {
		if got := containsEntry(lines, c.k, c.name); got != c.want {
			t.Errorf("containsEntry(%s, %q) = %v, want %v", c.k, c.name, got, c.want)
		}
	}
}
