# The following lines were added by Docker Desktop to add commands to your PATH.
export PATH="$PATH:/Users/hayato.hagari/.docker/bin"
# End of Docker Desktop section.

# ~/.zprofile — ログインシェルで一度だけ読まれる。Homebrew の shellenv のみ。
if [[ -x /opt/homebrew/bin/brew ]]; then
  eval "$(/opt/homebrew/bin/brew shellenv)"
fi
