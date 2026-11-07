# Brewfile — GUIアプリとシステム寄りのツールはここで管理する
# 反映: brew bundle --file=Brewfile（install.sh からも実行される）

# --- CLI / システム寄りのツール ---
brew "git"
brew "mise"
brew "tree"       # ディレクトリ構成をツリー表示（mise の registry に無いので brew で入れる）
brew "mas"        # Mac App Store の CLI（下の mas 行に必要。App Store へのサインインが前提）

# --- ライブラリ ---
# OpenCascade（OCCT）— CAD 用 C++ ライブラリ。追加時点の版は 7.9.3
# Brewfile では版を固定できないので、7.9.x に留めたい場合は `brew pin opencascade` を使う
brew "opencascade"

# --- GUIアプリ ---
cask "ghostty"    # ターミナルエミュレータ
cask "clipy"      # クリップボード履歴（初回起動時にアクセシビリティの許可が必要）
cask "rectangle"  # キーボードショートカットでウィンドウを整列（同上）
cask "unity-hub"  # Unity Hub — エディタ本体はプロジェクトに合わせて Hub から入れる
cask "unity-cli"  # Unity 公式 CLI（unity コマンド。beta 版）
cask "blender"         # 3DCG
cask "bambu-studio"    # Bambu Lab の 3D プリンタ用スライサー
cask "jlcone"          # JLC ONE — JLCPCB の見積もり・注文・注文追跡
cask "docker-desktop"  # Docker Desktop（docker / docker compose コマンドも入る）
cask "discord"         # チャット

# --- Mac App Store ---
mas "Xcode", id: 497799835  # Unity の iOS / Mac ビルドに必要。容量が大きく時間がかかる
