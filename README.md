# dotfiles

macOS（Apple Silicon）の開発環境設定。

## 方針

- **GUIアプリ** → Homebrew（`brew/gui.Brewfile`）
- **mise に無いシステム寄りのツール・ライブラリ** → Homebrew（`brew/cli.Brewfile`）
- **言語ランタイム・CLIツール** → mise（`config/mise/config.toml`）
- **シェル設定** → 役割ごとにファイルを分割（`zsh/conf.d/`）
- **マシン固有の設定・秘密情報** → `~/.zshrc.local`（リポジトリには含めない）

## 構成

```
dotfiles/
├── brew/
│   ├── cli.Brewfile         # Homebrew の formula（mise に無いツール・ライブラリ）
│   └── gui.Brewfile         # Homebrew の cask・Mac App Store のアプリ
├── zsh/
│   ├── .zprofile            # Homebrew の shellenv のみ
│   ├── .zshrc               # conf.d/ の読み込みと mise の有効化
│   └── conf.d/              # 役割別のシェル設定
│       ├── aliases.zsh      #   エイリアス
│       ├── completion.zsh   #   補完
│       ├── fzf.zsh          #   fzf のキーバインド・ghq 連携
│       ├── history.zsh      #   履歴
│       ├── options.zsh      #   シェルオプション
│       └── prompt.zsh       #   プロンプト（git ブランチ表示）
├── config/
│   ├── ghostty/config       # Ghostty（Option を Alt として使う・フォント・テーマ）
│   └── mise/config.toml     # mise グローバルツール（node LTS, python, go, rust, gh, kubectl, rg, jq, yq, fd, bat, ghq）
└── dots.toml                # dots の設定（シンボリックリンク・カスタムコマンド）
```

シンボリックリンクの対応（`dots.toml` の `[links]`）:

| リポジトリ内 | リンク先 |
| --- | --- |
| `zsh/.zprofile` | `~/.zprofile` |
| `zsh/.zshrc` | `~/.zshrc` |
| `config/ghostty/` | `~/.config/ghostty` |
| `config/mise/` | `~/.config/mise` |

## 新しい Mac でのセットアップ

```sh
# 0. App Store にサインインしておく（gui.Brewfile の mas 行で Xcode を入れるため）

# 1. Command Line Tools（git が入る）
xcode-select --install

# 2. dots を入れる（~/.local/bin/dots）
curl -fsSL https://raw.githubusercontent.com/gari8/dotfiles-cli/main/install.sh | sh

# 3. このリポジトリを ~/dotfiles に clone して登録
#    まだ SSH 鍵が無いので HTTPS で取る（公開リポジトリなので認証不要）
~/.local/bin/dots init --from https://github.com/gari8/dotfiles.git
#    （SSH 鍵を GitHub に登録済みなら git@github.com:gari8/dotfiles.git でも可）

# 4. セットアップ（まず -n で内容確認を推奨）
~/.local/bin/dots install -n
~/.local/bin/dots install

# 5. シェルを再起動して反映
exec $SHELL -l
```

`dots install` がやること（何度実行しても安全）:

1. `dots.toml` の `[links]` のシンボリックリンクを作成（既存ファイルは `~/.dotfiles_backup/<日時>/` に退避）
2. Homebrew が無ければ確認の上インストールし、`brew bundle` で2つの Brewfile の足りないものを入れる
3. `mise install` でグローバルツールを入れる

push するには GitHub への認証が必要。SSH 鍵を作って GitHub に登録し、`git -C ~/dotfiles remote set-url origin git@github.com:gari8/dotfiles.git` で SSH に切り替える（または `gh auth login`）。

`dots` 自体は mise ではなく `~/.local/bin/dots` に置く。更新するときは手順 2 の `curl ... | sh` をもう一度実行する。

## この設定を取り込む（他の人向け）

どの方法も、先に [dots](https://github.com/gari8/dotfiles-cli) を入れておく。

```sh
curl -fsSL https://raw.githubusercontent.com/gari8/dotfiles-cli/main/install.sh | sh
```

### 丸ごと使う

GitHub でこのリポジトリを fork し、**自分の fork** を取り込む（このリポジトリを直接 `--from` に渡すと、push 先が私のリポジトリになってしまう）。

```sh
dots init --from https://github.com/<あなた>/dotfiles.git
dots install -n     # 何が起きるか確認してから
dots install
```

- `dots init --from` は、リンクを張る前に `~/.zshrc` などどのファイルを置き換えるかと、カスタムコマンドの一覧を表示して確認する（置き換えたファイルは `~/.dotfiles_backup/<日時>/` に退避される）
- `dots install` は Brewfile のアプリと mise のツールを**全部**入れる。要らないものは先に `brew/*.Brewfile` と `config/mise/config.toml` から消しておく

### カスタムコマンドだけ欲しい

`dots cmd status` / `commit` / `push` / `pull` / `sync` など（`dots.toml` に定義）を、自分の dotfiles に取り込む。中身を表示して、1つずつ確認してから追加される。

```sh
dots pick https://github.com/gari8/dotfiles.git                 # 一覧
dots pick https://github.com/gari8/dotfiles.git commit push sync
```

### アプリや CLI ツールの一覧だけ欲しい

Brewfile や mise の設定を、ファイルごと自分の dotfiles にコピーする（自分のファイルがすでにあれば、置き換えるか確認される）。

```sh
dots init --brewfile-gui https://github.com/gari8/dotfiles.git//brew/gui.Brewfile           # GUI アプリ
dots init --brewfile-cli https://github.com/gari8/dotfiles.git//brew/cli.Brewfile           # Homebrew の CLI・ライブラリ
dots init --mise-config https://github.com/gari8/dotfiles.git//config/mise/config.toml      # CLI ツール・言語ランタイム
```

ファイルごと置き換わるので、いくつかだけ欲しいときは中身を見て `dots add --cask <名前>` / `dots add --mise <名前>` で1つずつ足す方がよい。

### シェルの設定だけ欲しい

`zsh/conf.d/` のファイル（エイリアス・補完・履歴・プロンプト・fzf）は独立しているので、必要なものだけ自分の設定にコピーする。

## dots コマンド

[gari8/dotfiles-cli](https://github.com/gari8/dotfiles-cli) の `dots` で、dotfiles をどこからでも操作する。確認は `-y` で省略できる。

| コマンド | やること |
| --- | --- |
| `dots install` | リンク作成・`brew bundle`・`mise install`（`-n` で確認のみ） |
| `dots init [dir]` | dotfiles の場所を登録し、リンクを作成 |
| `dots list [--mise\|--brew]` | 管理しているツールと導入状況（✓ / ✗）を表示 |
| `dots upgrade [--mise\|--brew]` | 更新できるものを表示し、y で更新（`-n` で表示のみ、`-x a,b` で除外） |
| `dots add --mise\|--brew\|--cask <名前>` | インストールして設定に追加 |
| `dots rm --mise\|--brew\|--cask <名前>` | 設定から削除（アンインストールは確認） |
| `dots path` | dotfiles のパスを表示（`cd "$(dots path)"`） |
| `dots cmd <名前>` | `dots.toml` のカスタムコマンドを実行（`dots cmd list` で一覧） |

git の操作は `dots cmd status` / `commit` / `push` / `pull` / `sync` など（`dots.toml` に定義。書き方は同ファイル冒頭）。
詳しくは `dots --help`。

## キーバインド（fzf）

| キー | やること |
| --- | --- |
| `Ctrl-G` | ghq のリポジトリを選んで `cd`（clone は `ghq get <owner>/<repo>`） |
| `Ctrl-R` | コマンド履歴を検索 |
| `Ctrl-T` | ファイルを選んで挿入 |
| `Alt-C` | サブディレクトリを選んで `cd` |

## ツールを追加するときのルール

| 追加したいもの | 書く場所 | 反映コマンド |
| --- | --- | --- |
| GUIアプリ（cask） | `brew/gui.Brewfile`（`dots add --cask`） | `dots install` |
| mise に無いシステム寄りのCLI | `brew/cli.Brewfile`（`dots add --brew`） | 同上 |
| 言語ランタイム・開発用CLI | `config/mise/config.toml` | `dots install` |
| エイリアス・シェル設定 | `zsh/conf.d/*.zsh`（役割に合うファイル、無ければ新規作成） | `exec $SHELL -l` |
| マシン固有・秘密情報 | `~/.zshrc.local`（コミットしない） | 同上 |

追加は `dots add` を使うと、インストールと設定への追記をまとめて行える。
コミットは `dots cmd commit <メッセージ>` → `dots cmd push`。

- 名前は推測せず `brew search <名前>` / `mise registry | grep <名前>` で確認してから書く
- プロジェクト固有のランタイムは、ここではなく各リポジトリの `mise.toml` に書く

## 秘密情報について

APIキー・トークン・ssh接続先などは **絶対にこのリポジトリに入れない**。
`~/.zshrc.local` に書けば `.zshrc` の最後で自動的に読み込まれる（`.gitignore` 済み）。
