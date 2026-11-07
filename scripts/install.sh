#!/usr/bin/env bash
#
# dotfiles セットアップスクリプト
#
# 使い方:
#   scripts/install.sh            # 実行
#   scripts/install.sh --dry-run  # 何が行われるかの確認のみ（変更なし）
#
# 何度実行しても安全（冪等）:
#   - すでに正しいシンボリックリンクが張られていればスキップ
#   - 既存のファイル・ディレクトリは ~/.dotfiles_backup/<日時>/ に退避してからリンク
set -euo pipefail

DOTFILES_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKUP_DIR="$HOME/.dotfiles_backup/$(date +%Y%m%d-%H%M%S)"
DRY_RUN=false

if [[ "${1:-}" == "--dry-run" || "${1:-}" == "-n" ]]; then
  DRY_RUN=true
fi

info() { printf '%s\n' "$*"; }
run() {
  if $DRY_RUN; then
    info "  [dry-run] $*"
  else
    "$@"
  fi
}

# シンボリックリンクを張る。リンク先に既存ファイルがあればバックアップする。
link() {
  local src="$1" dst="$2"

  if [[ ! -e "$src" ]]; then
    info "  SKIP: リンク元がありません: $src"
    return
  fi

  # すでに正しいリンクなら何もしない（冪等）
  if [[ -L "$dst" && "$(readlink "$dst")" == "$src" ]]; then
    # 変数の直後に全角文字が続くと bash 3.2 が変数名の一部と誤読するので ${} で囲む
    info "  OK:   ${dst}（リンク済み）"
    return
  fi

  # 既存のファイル・ディレクトリ・別のリンクはバックアップへ退避
  if [[ -e "$dst" || -L "$dst" ]]; then
    run mkdir -p "$BACKUP_DIR"
    run mv "$dst" "$BACKUP_DIR/"
    info "  BACKUP: $dst -> $BACKUP_DIR/"
  fi

  run mkdir -p "$(dirname "$dst")"
  run ln -s "$src" "$dst"
  info "  LINK: $dst -> $src"
}

if $DRY_RUN; then
  info "== dry-run モード: 実際の変更は行いません =="
fi

info "== シンボリックリンク =="
link "$DOTFILES_DIR/zsh/.zprofile"   "$HOME/.zprofile"
link "$DOTFILES_DIR/zsh/.zshrc"      "$HOME/.zshrc"
link "$DOTFILES_DIR/config/ghostty"  "$HOME/.config/ghostty"
link "$DOTFILES_DIR/config/mise"     "$HOME/.config/mise"

info ""
info "== Homebrew =="
if [[ -x /opt/homebrew/bin/brew ]]; then
  eval "$(/opt/homebrew/bin/brew shellenv)"
  info "  OK: Homebrew インストール済み"
elif $DRY_RUN; then
  info "  [dry-run] Homebrew が未インストールです（本実行時に確認の上インストールします）"
else
  read -r -p "Homebrew が見つかりません。インストールしますか? [y/N] " answer
  if [[ "$answer" == [yY]* ]]; then
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    eval "$(/opt/homebrew/bin/brew shellenv)"
  else
    info "  Homebrew をスキップしました。brew bundle / mise install も行いません。"
    exit 0
  fi
fi

info ""
info "== brew bundle =="
run brew bundle --file="$DOTFILES_DIR/Brewfile"

info ""
info "== mise install =="
if command -v mise >/dev/null 2>&1; then
  run mise install
else
  info "  SKIP: mise が見つかりません（brew bundle 実行後に再実行してください）"
fi

info ""
info "== dots コマンド =="
# cmd/dots をビルドして ~/.local/bin/dots に置く（Go は上の mise install で入る）
# 実行中の dots 自身を上書きしても壊れないよう、一時ファイルにビルドしてから mv で置き換える
DOTS_BIN="$HOME/.local/bin/dots"
if ! command -v mise >/dev/null 2>&1; then
  info "  SKIP: mise が見つかりません"
elif $DRY_RUN; then
  info "  [dry-run] go build ./cmd/dots -> $DOTS_BIN"
else
  mkdir -p "$(dirname "$DOTS_BIN")"
  tmp="$(mktemp "$DOTS_BIN.XXXXXX")"
  if (cd "$DOTFILES_DIR" && mise exec -- go build -ldflags "-X 'main.defaultDir=$DOTFILES_DIR'" -o "$tmp" ./cmd/dots); then
    mv -f "$tmp" "$DOTS_BIN"
    info "  OK:   $DOTS_BIN"
  else
    rm -f "$tmp"
    info "  FAILED: dots のビルドに失敗しました"
    exit 1
  fi
fi

info ""
info "完了。新しいシェルを開くか 'exec \$SHELL -l' で設定を反映してください。"
