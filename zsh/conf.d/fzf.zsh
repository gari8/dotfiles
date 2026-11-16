# fzf.zsh — fzf（あいまい検索）のキーバインドと ghq 連携
# completion.zsh の compinit より後に読み込まれる（conf.d はファイル名順）

if (( $+commands[fzf] )); then
  # fzf 標準のキーバインドと補完
  #   Ctrl-R: コマンド履歴を絞り込んで選ぶ
  #   Ctrl-T: ファイルを選んでコマンドラインに挿入する
  #   Alt-C:  サブディレクトリを選んで cd する
  source <(fzf --zsh)

  # 絞り込み画面は画面の下 40% に、候補は上から順に出す
  export FZF_DEFAULT_OPTS="--height 40% --reverse --border"
fi

# Ctrl-G: ghq で管理しているリポジトリを選んで cd する
if (( $+commands[fzf] && $+commands[ghq] )); then
  ghq-fzf-cd() {
    local root repo
    root="$(ghq root)"
    # 右側に選択中のリポジトリの README（無ければファイル一覧）を表示する
    repo="$(ghq list | fzf --prompt 'repo> ' \
      --preview "bat --color=always --style=plain $root/{}/README.md 2>/dev/null || ls -A $root/{}")"

    if [[ -n "$repo" ]]; then
      # cd をコマンドとして実行し、履歴に残してプロンプト（git ブランチ表示）も更新する
      BUFFER="cd ${(q)root}/${(q)repo}"
      zle accept-line
    elif [[ -z "$(ghq list)" ]]; then
      zle -M "ghq のリポジトリがありません（ghq get <owner>/<repo> で clone できます）"
    fi
    zle reset-prompt
  }
  zle -N ghq-fzf-cd
  bindkey '^G' ghq-fzf-cd
fi
