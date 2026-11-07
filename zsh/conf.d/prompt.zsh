# prompt.zsh — プロンプト表示（右側に git ブランチ状態）

autoload -Uz vcs_info
autoload -Uz colors && colors
setopt prompt_subst

# trueで作業ブランチの状態（未add/未commit）に応じて表示を変える
zstyle ':vcs_info:*' check-for-changes false
# addしていない変更がある場合の表示
zstyle ':vcs_info:*' unstagedstr "%F{red}%B＋%b%f"
# commitしていない変更がある場合の表示
zstyle ':vcs_info:*' stagedstr "%F{yellow}★ %f"
# デフォルトの表示
zstyle ':vcs_info:*' formats "%u%c%F{green}【 %b 】%f"
# コンフリクトなど特別な状態のときの表示
zstyle ':vcs_info:*' actionformats '【%b | %a】'

precmd() { vcs_info }

PROMPT='%F{cyan}%~%f %# '
RPROMPT='${vcs_info_msg_0_}'
