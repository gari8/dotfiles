# aliases.zsh — エイリアス定義
# ssh接続先など個人・マシン固有のものはここではなく ~/.zshrc.local に書く

# ls
alias ls="ls -G"
alias ll="ls -lhG"
alias la="ls -lahG"

# git
alias g="git"
alias gpush="git push origin HEAD"
alias gpull="git pull origin HEAD"

# kubernetes
alias k="kubectl"

# Homebrew の古いバージョンとダウンロードキャッシュを削除する
# （C コンパイラの cc と衝突しないよう cleanup という名前にしている）
alias cleanup="brew cleanup -s && rm -rf ~/Library/Caches/Homebrew"

# シェルを再起動して設定を読み直す
alias relogin='exec $SHELL -l'
