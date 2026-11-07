# ~/.zshrc — 対話シェルの設定。本体は zsh/conf.d/ に役割ごとに分割してある。

# path 配列の重複エントリを自動で除去する
typeset -U path PATH

# このファイル（シンボリックリンク解決後）の場所から dotfiles の zsh/ を特定する
DOTFILES_ZSH_DIR="${${(%):-%N}:A:h}"

# mise — 言語ランタイム・CLIツールの有効化
if command -v mise >/dev/null 2>&1; then
  eval "$(mise activate zsh)"
fi

# 役割別の設定を読み込む（(N): マッチなしでもエラーにしない）
for _conf in "$DOTFILES_ZSH_DIR"/conf.d/*.zsh(N); do
  source "$_conf"
done
unset _conf

# マシン固有の設定・秘密情報はリポジトリ外の ~/.zshrc.local に書く
[[ -f "$HOME/.zshrc.local" ]] && source "$HOME/.zshrc.local"

# Unity CLI
. "/Users/hayato.hagari/.unity/env"
