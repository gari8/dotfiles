# history.zsh — コマンド履歴

HISTFILE="$HOME/.zsh_history"
HISTSIZE=100000
SAVEHIST=100000

# 複数のシェル間で履歴を共有する
setopt share_history
# 直前と同じコマンドは履歴に残さない
setopt hist_ignore_dups
# 先頭にスペースを付けたコマンドは履歴に残さない（秘密情報を含むコマンド用）
setopt hist_ignore_space
# 履歴に実行時刻も記録する
setopt extended_history
