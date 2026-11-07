# options.zsh — シェルの挙動に関するオプション

export CLICOLOR=1

# ディレクトリ名だけで cd できる
setopt auto_cd
# cd したディレクトリを自動で pushd の履歴に積む（cd - + Tab で履歴から選べる）
setopt auto_pushd
# pushd の履歴に重複を積まない
setopt pushd_ignore_dups
# Ctrl+D でシェルを終了しない
setopt ignore_eof
