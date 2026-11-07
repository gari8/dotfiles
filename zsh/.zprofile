# ~/.zprofile — ログインシェルで一度だけ読まれる。Homebrew の shellenv のみ。
if [[ -x /opt/homebrew/bin/brew ]]; then
  eval "$(/opt/homebrew/bin/brew shellenv)"
fi
