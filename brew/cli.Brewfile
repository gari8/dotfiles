# Homebrew の formula — mise に無いシステム寄りのツールとライブラリ（dots add --brew で追記される）
# 反映: dots install

brew "git"
brew "mise"
brew "tree"       # ディレクトリ構成をツリー表示（mise の registry に無いので brew で入れる）
brew "mas"        # Mac App Store の CLI（gui.Brewfile の mas 行に必要。App Store へのサインインが前提）

# OpenCascade（OCCT）— CAD 用 C++ ライブラリ。8.0.1 を使う
# homebrew/core はまだ 7.9.3 なので、自分の tap（github.com/gari8/homebrew-tap）の formula で 8.0.1 を入れる。
# 版を変えるときは tap の Formula/opencascade.rb を更新する（bottle が無いのでソースからビルドされる）
tap "gari8/tap"
brew "gari8/tap/opencascade"
