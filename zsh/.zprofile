# The following lines were added by Docker Desktop to add commands to your PATH.
export PATH="$PATH:/Users/hayato.hagari/.docker/bin"
# End of Docker Desktop section.

# ~/.zprofile — ログインシェルで一度だけ読まれる。Homebrew の shellenv と PATH の設定。
if [[ -x /opt/homebrew/bin/brew ]]; then
  eval "$(/opt/homebrew/bin/brew shellenv)"
fi

# 自分で入れたコマンド（install.sh がビルドする dots など）
export PATH="$HOME/.local/bin:$PATH"
