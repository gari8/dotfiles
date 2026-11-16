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
~/.local/bin/dots init --from git@github.com:gari8/dotfiles.git

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

`dots` 自体は mise ではなく `~/.local/bin/dots` に置く。更新するときは手順 2 の `curl ... | sh` をもう一度実行する。

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
