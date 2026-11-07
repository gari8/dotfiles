# dotfiles

macOS（Apple Silicon）の開発環境設定。

## 方針

- **GUIアプリ・システム寄りのツール** → Homebrew（`Brewfile`）
- **言語ランタイム・CLIツール** → mise（`config/mise/config.toml`）
- **シェル設定** → 役割ごとにファイルを分割（`zsh/conf.d/`）
- **マシン固有の設定・秘密情報** → `~/.zshrc.local`（リポジトリには含めない）

## 構成

```
dotfiles/
├── Brewfile                 # Homebrew で入れるもの（GUIアプリ・システムツール）
├── zsh/
│   ├── .zprofile            # Homebrew の shellenv のみ
│   ├── .zshrc               # conf.d/ の読み込みと mise の有効化
│   └── conf.d/              # 役割別のシェル設定
│       ├── aliases.zsh      #   エイリアス
│       ├── completion.zsh   #   補完
│       ├── history.zsh      #   履歴
│       ├── options.zsh      #   シェルオプション
│       └── prompt.zsh       #   プロンプト（git ブランチ表示）
├── config/
│   ├── ghostty/config       # Ghostty（Option を Alt として使う・フォント・テーマ）
│   └── mise/config.toml     # mise グローバルツール（node LTS, python, go, rust, gh, kubectl, rg, jq, yq, fd, bat, ghq）
├── cmd/dots/                # dots コマンド（Go）
└── scripts/
    └── install.sh           # セットアップスクリプト（冪等・バックアップ付き）
```

シンボリックリンクの対応:

| リポジトリ内 | リンク先 |
| --- | --- |
| `zsh/.zprofile` | `~/.zprofile` |
| `zsh/.zshrc` | `~/.zshrc` |
| `config/ghostty/` | `~/.config/ghostty` |
| `config/mise/` | `~/.config/mise` |

## 新しい Mac でのセットアップ

```sh
# 0. App Store にサインインしておく（Brewfile の mas 行で Xcode を入れるため）

# 1. Command Line Tools（git が入る。clone 時に自動でプロンプトが出る場合もある）
xcode-select --install

# 2. リポジトリを取得
git clone git@github.com:gari8/dotfiles.git ~/dotfiles

# 3. セットアップ（まず dry-run で内容確認を推奨）
~/dotfiles/scripts/install.sh --dry-run
~/dotfiles/scripts/install.sh

# 4. シェルを再起動して反映
exec $SHELL -l
```

`install.sh` がやること:

1. 設定ファイルのシンボリックリンクを作成（既存ファイルは `~/.dotfiles_backup/<日時>/` に退避）
2. Homebrew が無ければ確認の上インストール
3. `brew bundle` で Brewfile の内容をインストール
4. `mise install` でグローバルツールをインストール
5. `dots` コマンドをビルドして `~/.local/bin/dots` に置く

何度実行しても安全です。セットアップ後は `dots install` でどこからでも再実行できます。

## dots コマンド

dotfiles をどこからでも操作できるコマンド。コミット・push の前に内容を見せて確認する（`-y` で省略）。

| コマンド | やること |
| --- | --- |
| `dots status` | 未コミットの変更と未 push のコミットを表示 |
| `dots commit "<メッセージ>"` | 差分を確認してコミット |
| `dots push` | 未 push のコミットを確認して push |
| `dots pull` | リモートの変更を取り込み、`dots install` するか確認 |
| `dots install` | `install.sh` を実行（`--dry-run` で確認のみ） |
| `dots mise add <名前>` / `rm` | mise のツールを追加・削除してコミット |
| `dots brew add <名前>` / `rm` | Brewfile に追加・削除してコミット（cask は `--cask`） |
| `dots path` | dotfiles のパスを表示（`cd "$(dots path)"`） |

詳しくは `dots --help` / `dots <コマンド> --help`。

## ツールを追加するときのルール

| 追加したいもの | 書く場所 | 反映コマンド |
| --- | --- | --- |
| GUIアプリ（cask） | `Brewfile` | `brew bundle --file=~/dotfiles/Brewfile` |
| システム寄りのCLI | `Brewfile` | 同上 |
| 言語ランタイム・開発用CLI | `config/mise/config.toml` | `mise install` |
| エイリアス・シェル設定 | `zsh/conf.d/*.zsh`（役割に合うファイル、無ければ新規作成） | `exec $SHELL -l` |
| マシン固有・秘密情報 | `~/.zshrc.local`（コミットしない） | 同上 |

追加は `dots mise add` / `dots brew add` を使うと、インストールからコミットまでまとめて行える。
手で書き換えたときは `dots commit` でコミットし、`dots push` で push する。

- 名前は推測せず `brew search <名前>` / `mise registry | grep <名前>` で確認してから書く
- プロジェクト固有のランタイムは、ここではなく各リポジトリの `mise.toml` に書く

## 秘密情報について

APIキー・トークン・ssh接続先などは **絶対にこのリポジトリに入れない**。
`~/.zshrc.local` に書けば `.zshrc` の最後で自動的に読み込まれる（`.gitignore` 済み）。
