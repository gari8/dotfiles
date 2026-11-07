# completion.zsh — コマンド入力補完

# Homebrew でインストールしたツールの補完関数を fpath に追加
# （HOMEBREW_PREFIX は .zprofile の brew shellenv で設定される）
if [[ -n "$HOMEBREW_PREFIX" && -d "$HOMEBREW_PREFIX/share/zsh/site-functions" ]]; then
  fpath=("$HOMEBREW_PREFIX/share/zsh/site-functions" $fpath)
fi

autoload -Uz compinit
compinit

# dots コマンドの補完（compinit の後に読み込む）
if (( $+commands[dots] )); then
  source <(dots completion zsh)
fi

# 補完候補に色をつける
zstyle ':completion:*' list-colors "${LS_COLORS}"
# 補完候補を矢印キーで選択できるようにする
zstyle ':completion:*' menu select
# 大文字小文字を区別せずに補完する
zstyle ':completion:*' matcher-list 'm:{a-z}={A-Z}'
