package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

// customCmd は dots.toml の [cmd.<名前>] 1つ分
type customCmd struct {
	Run     string `toml:"run"`     // 実行するコマンド（sh -c で実行。引数の渡し方は script を参照）
	Desc    string `toml:"desc"`    // 一覧に出す説明
	Dir     string `toml:"dir"`     // 実行場所: root（dotfiles。既定）/ cwd（今いる場所）
	Confirm bool   `toml:"confirm"` // 実行前に確認する
}

// exitCodeError は子プロセスの終了コードをそのまま dots の終了コードにする
type exitCodeError int

func (e exitCodeError) Error() string { return fmt.Sprintf("exit status %d", int(e)) }

func newCmdCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cmd [<名前> [引数]...]",
		Short: "dots.toml のカスタムコマンドを実行する（dots cmd list で一覧）",
		// 引数はフラグも含めてカスタムコマンドにそのまま渡す
		DisableFlagParsing: true,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
			if len(args) > 0 {
				return nil, cobra.ShellCompDirectiveDefault
			}
			dir, err := repoDir()
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			cmds, _ := loadCustom(dir)
			names := []string{"list\t一覧を表示"}
			for name, c := range cmds {
				names = append(names, name+"\t"+c.describe())
			}
			return names, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			args = stripYes(args)
			dir, err := repoDir()
			if err != nil {
				return err
			}
			cmds, err := loadCustom(dir)
			if err != nil {
				return err
			}
			if len(args) == 0 || args[0] == "-h" || args[0] == "--help" || args[0] == "list" || args[0] == "ls" {
				listCustom(cmds)
				return nil
			}
			if reservedNames[args[0]] {
				return fmt.Errorf("dots cmd %s はまだありません（予約済み）", args[0])
			}
			c, ok := cmds[args[0]]
			if !ok {
				return fmt.Errorf("カスタムコマンド %q はありません（dots cmd list で一覧）", args[0])
			}
			return runCustom(dir, args[0], c, args[1:])
		},
	}
}

// stripYes は先頭の -y / --yes を取り除いて assumeYes に反映する
// （DisableFlagParsing なので dots cmd -y push の -y は自分で読む）
func stripYes(args []string) []string {
	for len(args) > 0 && (args[0] == "-y" || args[0] == "--yes") {
		assumeYes = true
		args = args[1:]
	}
	return args
}

// reservedNames は dots cmd 自身のサブコマンド用に予約した名前（カスタムコマンド名には使えない）
var reservedNames = map[string]bool{
	"list": true, "ls": true,
	"add": true, "rm": true, "remove": true,
}

// loadCustom は dots.toml を読む。ファイルが無ければ空。
func loadCustom(dir string) (map[string]customCmd, error) {
	var cfg struct {
		Cmd map[string]customCmd `toml:"cmd"`
	}
	path := filepath.Join(dir, "dots.toml")
	md, err := toml.DecodeFile(path, &cfg)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("dots.toml を読めません: %w", err)
	}
	if keys := md.Undecoded(); len(keys) > 0 {
		return nil, fmt.Errorf("dots.toml に不明な項目があります: %v", keys)
	}
	for name, c := range cfg.Cmd {
		if reservedNames[name] {
			return nil, fmt.Errorf("dots.toml: cmd.%s は予約名なので使えません", name)
		}
		if c.Run == "" {
			return nil, fmt.Errorf("dots.toml: cmd.%s に run がありません", name)
		}
		if c.Dir != "" && c.Dir != "root" && c.Dir != "cwd" {
			return nil, fmt.Errorf("dots.toml: cmd.%s の dir は root か cwd です", name)
		}
	}
	return cfg.Cmd, nil
}

func (c customCmd) describe() string {
	if c.Desc != "" {
		return c.Desc
	}
	return c.Run
}

func listCustom(cmds map[string]customCmd) {
	if len(cmds) == 0 {
		fmt.Println("カスタムコマンドはありません（dots.toml に [cmd.<名前>] を書く）")
		return
	}
	names := make([]string, 0, len(cmds))
	for name := range cmds {
		names = append(names, name)
	}
	sort.Strings(names)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()
	for _, name := range names {
		fmt.Fprintf(w, "  %s\t%s\n", name, cmds[name].describe())
	}
}

// argRef は run の中で引数を参照している箇所（$1 / ${1} / $@ / $* / $#）
var argRef = regexp.MustCompile(`\$\{?[1-9@*#]`)

// script は1行のコマンドで引数を参照していなければ、末尾に "$@" を付けて引数を渡す
// （複数行のスクリプトや $1 などを使うものは、引数の位置を自分で決める）
func script(run string) string {
	run = strings.TrimSpace(run)
	if strings.Contains(run, "\n") || argRef.MatchString(run) {
		return run
	}
	return run + ` "$@"`
}

func runCustom(dir, name string, c customCmd, args []string) error {
	if c.Confirm {
		line := strings.TrimSpace(c.Run + " " + strings.Join(args, " "))
		if !confirm(line + " を実行しますか?") {
			return canceled("実行")
		}
	}
	// 引数は $1.. / "$@" で参照できる（$0 にはコマンド名が入る）
	ex := exec.Command("sh", "-c", script(c.Run), name)
	ex.Args = append(ex.Args, args...)
	if c.Dir != "cwd" {
		ex.Dir = dir
	}
	pwd, _ := os.Getwd()
	ex.Env = append(os.Environ(), "DOTS_ROOT="+dir, "DOTS_PWD="+pwd)
	ex.Stdin, ex.Stdout, ex.Stderr = os.Stdin, os.Stdout, os.Stderr

	err := ex.Run()
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return exitCodeError(ee.ExitCode())
	}
	return err
}
