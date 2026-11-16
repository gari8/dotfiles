package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	var t target
	cmd := &cobra.Command{
		Use:     "list [--mise|--brew]",
		Aliases: []string{"ls"},
		Short:   "管理しているツールと導入状況を表示する",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := repoDir()
			if err != nil {
				return err
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			defer w.Flush()
			if t.all() || t.mise {
				if err := listMise(w, dir); err != nil {
					return err
				}
			}
			if t.all() || t.brew {
				if err := listBrew(w, dir); err != nil {
					return err
				}
			}
			return nil
		},
	}
	addTargetFlags(cmd, &t, false, false)
	return cmd
}

// listItem は一覧の1行（名前・入っているか・補足）
type listItem struct {
	name, note string
	installed  bool
}

// printGroup は見出しと各行を indent 分だけ下げて表示する
func printGroup(w *tabwriter.Writer, indent, title string, items []listItem) {
	if len(items) == 0 {
		return
	}
	fmt.Fprintf(w, "%s%s\n", indent, title)
	for _, it := range items {
		mark := "✓"
		if !it.installed {
			mark = "✗"
		}
		fmt.Fprintf(w, "%s  %s %s\t%s\n", indent, mark, it.name, it.note)
	}
}

func listMise(w *tabwriter.Writer, dir string) error {
	config := filepath.Join(dir, "config", "mise", "config.toml")
	// mise ls には --path が無いので、グローバル設定としてリポジトリ内のファイルを読ませる
	out, err := output(dir, "env", "MISE_GLOBAL_CONFIG_FILE="+config, "mise", "ls", "--global", "--json")
	if err != nil {
		return err
	}
	var tools map[string][]struct {
		Version   string `json:"version"`
		Installed bool   `json:"installed"`
	}
	if err := json.Unmarshal([]byte(out), &tools); err != nil {
		return fmt.Errorf("mise ls の出力を読めません: %w", err)
	}

	names := make([]string, 0, len(tools))
	for name := range tools {
		names = append(names, name)
	}
	sort.Strings(names)

	var items []listItem
	for _, name := range names {
		for _, v := range tools[name] {
			items = append(items, listItem{name: name, note: v.Version, installed: v.Installed})
		}
	}
	printGroup(w, "", "mise", items)
	fmt.Fprintln(w)
	return nil
}

// Brewfile の行: brew "name"  # コメント / mas "Name", id: 123  # コメント
var (
	brewLine = regexp.MustCompile(`^\s*(brew|cask)\s+"([^"]+)"[^#]*(?:#\s*(.*))?$`)
	masLine  = regexp.MustCompile(`^\s*mas\s+"([^"]+)",\s*id:\s*(\d+)[^#]*(?:#\s*(.*))?$`)
)

func listBrew(w *tabwriter.Writer, dir string) error {
	lines, err := readLines(filepath.Join(dir, "Brewfile"))
	if err != nil {
		return err
	}
	formulae := installedSet(dir, 0, "brew", "list", "--formula", "-1")
	casks := installedSet(dir, 0, "brew", "list", "--cask", "-1")
	apps := installedSet(dir, 0, "mas", "list") // 1列目が App Store の id

	groups := map[string][]listItem{}
	for _, l := range lines {
		if m := brewLine.FindStringSubmatch(l); m != nil {
			set := formulae
			if m[1] == "cask" {
				set = casks
			}
			// tap 付きの名前（owner/tap/name）は brew list では name だけになる
			groups[m[1]] = append(groups[m[1]], listItem{name: m[2], note: m[3], installed: set[path.Base(m[2])]})
		} else if m := masLine.FindStringSubmatch(l); m != nil {
			groups["mas"] = append(groups["mas"], listItem{name: m[1], note: m[3], installed: apps[m[2]]})
		}
	}
	// formula / cask / mas はどれも Brewfile から brew bundle で入るので brew の下にまとめる
	fmt.Fprintln(w, "brew")
	printGroup(w, "  ", "formula", groups["brew"])
	printGroup(w, "  ", "cask", groups["cask"])
	printGroup(w, "  ", "mas", groups["mas"])
	return nil
}

// installedSet はコマンド出力の各行の col 列目を集める（コマンドが失敗したら空）
func installedSet(dir string, col int, name string, args ...string) map[string]bool {
	set := map[string]bool{}
	out, err := output(dir, name, args...)
	if err != nil {
		return set
	}
	for _, l := range strings.Split(out, "\n") {
		if f := strings.Fields(l); len(f) > col {
			set[f[col]] = true
		}
	}
	return set
}
