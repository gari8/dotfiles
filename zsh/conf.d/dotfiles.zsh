# dotfiles.zsh — dotfiles リポジトリを操作する関数
# 引数を受け取るのでエイリアスではなく関数にしている

# dotfiles リポジトリの場所（DOTFILES_ZSH_DIR は .zshrc で設定される）
DOTFILES_DIR="${DOTFILES_ZSH_DIR:h}"

# セットアップスクリプトを実行する（dotinstall --dry-run で確認のみ）
# シングルクォートなので、実行するたびにその時点の DOTFILES_DIR を使う
alias dotinstall='"$DOTFILES_DIR/scripts/install.sh"'

# dotfiles の変更を確認してコミットする（cd しないので今いる場所はそのまま）
#   dotsave "add lazygit"
dotsave() {
  if [[ -z "$1" ]]; then
    echo "usage: dotsave <コミットメッセージ>" >&2
    return 1
  fi

  if [[ -z "$(git -C "$DOTFILES_DIR" status --porcelain)" ]]; then
    echo "dotfiles に変更はありません"
    return 0
  fi

  # 変更内容を見せてから確認する（新規ファイルは status に出る）
  git -C "$DOTFILES_DIR" status --short
  git -C "$DOTFILES_DIR" diff

  if read -q "?この内容でコミットしますか? [y/N] "; then
    echo
    git -C "$DOTFILES_DIR" add -A && git -C "$DOTFILES_DIR" commit -m "$*"
  else
    echo
    echo "コミットしませんでした"
    return 1
  fi
}

# mise でツールをグローバルに追加し、そのまま dotfiles にコミットする
#   mise-add lazygit
#   mise-add cargo:cargo-watch
mise-add() {
  if [[ $# -eq 0 ]]; then
    echo "usage: mise-add <ツール名> [<ツール名>...]" >&2
    return 1
  fi
  mise use -g "$@" && dotsave "mise: add $*"
}

# dotfiles のまだ push していないコミットを確認して push する（cd しないので今いる場所はそのまま）
#   dotpush
dotpush() {
  # 未コミットの変更があれば先に dotsave するよう促す（push 自体は止めない）
  if [[ -n "$(git -C "$DOTFILES_DIR" status --porcelain)" ]]; then
    echo "未コミットの変更があります（push されません。先に dotsave でコミットできます）"
    git -C "$DOTFILES_DIR" status --short
  fi

  # upstream が無い（初回 push）ときは origin に同名ブランチで push する
  if ! git -C "$DOTFILES_DIR" rev-parse --abbrev-ref --symbolic-full-name '@{u}' >/dev/null 2>&1; then
    if read -q "?upstream がありません。origin に push しますか? [y/N] "; then
      echo
      git -C "$DOTFILES_DIR" push -u origin HEAD
    else
      echo
      echo "push しませんでした"
      return 1
    fi
    return
  fi

  local commits
  commits="$(git -C "$DOTFILES_DIR" log --oneline '@{u}..HEAD')"
  if [[ -z "$commits" ]]; then
    echo "push するコミットはありません"
    return 0
  fi

  echo "$commits"
  if read -q "?これらのコミットを push しますか? [y/N] "; then
    echo
    git -C "$DOTFILES_DIR" push
  else
    echo
    echo "push しませんでした"
    return 1
  fi
}
