package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// target は --mise / --brew / --cask のどれを指定したか
type target struct {
	mise, brew, cask bool
}

// addTargetFlags は --mise / --brew / --cask を登録する（required なら1つ必須）
func addTargetFlags(cmd *cobra.Command, t *target, withCask, required bool) {
	cmd.Flags().BoolVar(&t.mise, "mise", false, "mise のツール")
	cmd.Flags().BoolVar(&t.brew, "brew", false, "Homebrew（list / upgrade では cask と mas も含む）")
	opts := "--mise / --brew"
	if withCask {
		cmd.Flags().BoolVar(&t.cask, "cask", false, "Homebrew の cask")
		opts += " / --cask"
	}
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		n := 0
		for _, on := range []bool{t.mise, t.brew, t.cask} {
			if on {
				n++
			}
		}
		if n > 1 {
			return fmt.Errorf("%s はどれか1つだけ指定してください", opts)
		}
		if required && n == 0 {
			return fmt.Errorf("%s のどれかを指定してください", opts)
		}
		return nil
	}
}

// all は何も指定されていない（list / upgrade で両方を対象にする）
func (t target) all() bool { return !t.mise && !t.brew && !t.cask }

func newAddCmd() *cobra.Command {
	var t target
	var comment string
	cmd := &cobra.Command{
		Use:   "add (--mise|--brew|--cask) <名前>...",
		Short: "インストールして設定に追加し、コミットする",
		Example: `  dots add --mise lazygit terraform@1.9
  dots add --cask arc -m "ブラウザ"`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if t.mise {
				if comment != "" {
					return fmt.Errorf("-m は --brew / --cask のときだけ使えます")
				}
				return miseEdit("use", "add", args)
			}
			return brewAdd(kindOf(t.cask), args, comment)
		},
	}
	addTargetFlags(cmd, &t, true, true)
	cmd.Flags().StringVarP(&comment, "comment", "m", "", "Brewfile の行末コメント")
	return cmd
}

func newRmCmd() *cobra.Command {
	var t target
	cmd := &cobra.Command{
		Use:     "rm (--mise|--brew|--cask) <名前>...",
		Aliases: []string{"remove"},
		Short:   "設定から削除してコミットする",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if t.mise {
				return miseEdit("unuse", "remove", args)
			}
			return brewRemove(kindOf(t.cask), args)
		},
	}
	addTargetFlags(cmd, &t, true, true)
	return cmd
}
