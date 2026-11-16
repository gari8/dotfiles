package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func newUpgradeCmd() *cobra.Command {
	var dryRun bool
	var exclude []string
	var t target
	cmd := &cobra.Command{
		Use:   "upgrade [--mise|--brew]",
		Short: "更新できるツールを表示し、確認して更新する",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := repoDir()
			if err != nil {
				return err
			}
			return upgrade(dir, t.all() || t.brew, t.all() || t.mise, dryRun, exclude)
		},
	}
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "表示のみ（更新しない）")
	cmd.Flags().StringSliceVarP(&exclude, "exclude", "x", nil, "更新しないツール（カンマ区切り）")
	addTargetFlags(cmd, &t, false, false)
	return cmd
}

// outdatedItem は更新できるツール1つ（名前・今のバージョン・新しいバージョン）
type outdatedItem struct {
	name, current, latest string
}

func upgrade(dir string, doBrew, doMise, dryRun bool, exclude []string) error {
	var brewItems, miseItems []outdatedItem
	var err error
	if doBrew {
		// formula / cask の定義を最新にしてから比べる
		if err := run(dir, "brew", "update", "--quiet"); err != nil {
			return err
		}
		if brewItems, err = brewOutdated(dir); err != nil {
			return err
		}
	}
	if doMise {
		if miseItems, err = miseOutdated(dir); err != nil {
			return err
		}
	}

	brewItems = without(brewItems, exclude)
	miseItems = without(miseItems, exclude)

	if len(brewItems) == 0 && len(miseItems) == 0 {
		fmt.Println("更新はありません")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	printOutdated(w, "brew", brewItems)
	printOutdated(w, "mise", miseItems)
	w.Flush()

	if dryRun {
		return nil
	}
	if !confirm("更新しますか?") {
		return canceled("更新")
	}
	// 表示した一覧と同じものだけを更新するよう、名前を指定する
	if len(brewItems) > 0 {
		if err := run(dir, "brew", append([]string{"upgrade"}, names(brewItems)...)...); err != nil {
			return err
		}
	}
	if len(miseItems) > 0 {
		args := append([]string{miseEnv(dir), "mise", "upgrade", "--global"}, names(miseItems)...)
		if err := run(dir, "env", args...); err != nil {
			return err
		}
	}
	return nil
}

// without は exclude に含まれる名前を取り除く
func without(items []outdatedItem, exclude []string) []outdatedItem {
	skip := map[string]bool{}
	for _, name := range exclude {
		skip[name] = true
	}
	var kept []outdatedItem
	for _, it := range items {
		if !skip[it.name] {
			kept = append(kept, it)
		}
	}
	return kept
}

func names(items []outdatedItem) []string {
	out := make([]string, len(items))
	for i, it := range items {
		out[i] = it.name
	}
	return out
}

func printOutdated(w *tabwriter.Writer, title string, items []outdatedItem) {
	if len(items) == 0 {
		return
	}
	fmt.Fprintln(w, title)
	for _, it := range items {
		if it.current == it.latest {
			// rust の stable などチャンネル指定は名前が変わらない
			fmt.Fprintf(w, "  %s\t%s（新しいビルドあり）\n", it.name, it.current)
		} else {
			fmt.Fprintf(w, "  %s\t%s → %s\n", it.name, it.current, it.latest)
		}
	}
}

func brewOutdated(dir string) ([]outdatedItem, error) {
	out, err := output(dir, "brew", "outdated", "--json=v2")
	if err != nil {
		return nil, err
	}
	type pkg struct {
		Name      string   `json:"name"`
		Installed []string `json:"installed_versions"`
		Current   string   `json:"current_version"`
		Pinned    bool     `json:"pinned"`
	}
	var res struct {
		Formulae []pkg `json:"formulae"`
		Casks    []pkg `json:"casks"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		return nil, fmt.Errorf("brew outdated の出力を読めません: %w", err)
	}
	var items []outdatedItem
	for _, p := range append(res.Formulae, res.Casks...) {
		if p.Pinned { // brew upgrade は pin したものを更新しない
			continue
		}
		items = append(items, outdatedItem{p.Name, strings.Join(p.Installed, ", "), p.Current})
	}
	return items, nil
}

func miseOutdated(dir string) ([]outdatedItem, error) {
	out, err := output(dir, "env", miseEnv(dir), "mise", "outdated", "--global", "--json")
	if err != nil {
		return nil, err
	}
	var res map[string]struct {
		Current string `json:"current"`
		Latest  string `json:"latest"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		return nil, fmt.Errorf("mise outdated の出力を読めません: %w", err)
	}
	var items []outdatedItem
	for name, t := range res {
		items = append(items, outdatedItem{name, t.Current, t.Latest})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].name < items[j].name })
	return items, nil
}

// miseEnv はリポジトリ内の mise 設定をグローバル設定として読ませる環境変数
func miseEnv(dir string) string {
	return "MISE_GLOBAL_CONFIG_FILE=" + filepath.Join(dir, "config", "mise", "config.toml")
}
