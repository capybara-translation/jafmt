# jafmt

日本語と半角英数字の間に半角スペースを入れる整形ツール。stdin から読み込み、整形結果を stdout に出力する。

```console
$ echo 'MacでGoを使う。第3章、50%の確率' | jafmt
Mac で Go を使う。第 3 章、50% の確率
```

macOS の Automator クイックアクションと組み合わせると、任意のアプリで選択したテキストをその場で整形できる。

## ルール

根拠は [Microsoft Japanese Localization Style Guide](https://download.microsoft.com/download/a/8/2/a822a118-18d4-4429-b857-1b65ab388315/jpn-jpn-StyleGuide.pdf) の 4.1.11 "Symbols & spaces"。
各ルールの具体例は [formatter/formatter_test.go](formatter/formatter_test.go) のテストケースが正となる。

### スペースを入れる

- 日本語（漢字・ひらがな・カタカナ・長音記号「ー」）と半角英数字の間に半角スペースを入れる（`Wordを使う` → `Word を使う`）
- 数字と日本語の間にも入れる（`第3章` → `第 3 章`、`10個` → `10 個`）
- 英数字に付く次の半角記号は語の一部として扱い、日本語との間にスペースを入れる
  - 語末: `%` `°` `#` `+`（`50%の` → `50% の`、`C#を` → `C# を`、`C++で` → `C++ で`）
  - 語頭: `#` `@` `$`（`は#123` → `は #123`、`は@user` → `は @user`）
- 語中の記号はそのまま語の一部になる（`node.jsと` → `node.js と`、`Wi-Fiを` → `Wi-Fi を`）
- アクセント付きのラテン文字（`é` など）も英字として扱う
- 結合文字（濁点が分かれた NFD の `が` など）や異体字セレクタ（`葛󠄀` など）が付いた文字は、それが付いている元の文字で判定する
- 日本語と半角括弧 `()` `[]` の外側の間に入れる（`日本語(Japanese)です` → `日本語 (Japanese) です`、`[新規]をクリック` → `[新規] をクリック`）
- 日本語の後の `?` `!` と、続く英数字の間に入れる（`保存しますか?Excelを` → `保存しますか? Excel を`）
- 日本語とインラインコードの外側の間に入れる（``これは`code`です`` → ``これは `code` です``）
- 日本語と半角の二重引用符の外側の間に入れる（`"test"と入力` → `"test" と入力`）
- 日本語と Markdown の強調 `*斜体*` `**太字**` `***太字斜体***` `~~取り消し~~` の外側の間に入れる（`これは**bold**です` → `これは **bold** です`）
  - 引用符と強調は、同じ行の中で開きと閉じが揃っている場合だけ対象にする
  - 開きの直後や閉じの直前が空白の記号は強調とみなさない（`2 * 3`、箇条書きの `* 項目`）
  - 単独の `~` は範囲（`10~20`）に使われるため、`~~` だけを対象にする。`_斜体_` は識別子（`foo_bar`）と紛れるため対象外
- 数字と単位の間に入れる（`50kg` → `50 kg`、`1.5GB` → `1.5 GB`、`100km/h` → `100 km/h`）
  - 対象は次の単位だけ（大文字・小文字を区別する）
    - 長さ: `nm` `mm` `cm` `m` `km`
    - 重さ: `mg` `g` `kg`
    - 体積: `mL` `ml` `L` `cc`
    - データ量・通信速度: `B` `KB` `kB` `MB` `GB` `TB` `PB` `KiB` `MiB` `GiB` `TiB` `bit` `bps` `kbps` `Kbps` `Mbps` `Gbps`
    - 周波数: `Hz` `kHz` `MHz` `GHz`
    - 時間: `ns` `ms` `sec` `min`
    - 電気: `V` `mV` `mA` `mAh` `W` `kW` `Wh` `kWh` `dB`
    - 画面・印刷: `px` `pt` `dpi` `ppi` `fps`
    - マイクロ: `μm` `μg` `μL` `μs` `μA` `μV`（マイクロ記号 `µ` も同じ扱い）
  - 数値と単位がそれぞれ独立した語の場合だけ入れる。`margin:16px` や `image_100px.png` のような識別子・ファイル名の一部は変えない
  - 単位以外の意味と紛れやすい `s`（`1990s`）・`h`・`t`・`A` は含めない。写真・投影の文脈の `mm` もガイドでは例外だが、文脈を判別できないため常にスペースを入れる

### スペースを入れない

- 全角記号（、。「」『』（）【】！？・など）の前後
- 半角スラッシュの両側（`日本/US`、`3/14`）
- 日本語の直後の `?` `!` `:` `...`（`更新しますか?`、`フォント:`）
- `?` `!` の後に日本語が続く箇所（`本当?はい`）
- 半角括弧・二重引用符・強調の内側（`(タイトル)`、`[新規]`、`"日本語"`、`**太字**`）
- 数字と `%`・`°`（`50%`、`45°`）、一覧にない単位（`4K`、`5G`）
- 上記以外の半角記号だけが日本語と接している箇所（`割合を%で` など）
- 全角英数字（`Ｗｏｒｄ`）と半角カタカナ（`ｶﾀｶﾅ`）は対象外

### 壊さない

- 既存のスペース・全角スペース・タブ・改行・行頭のインデントはそのまま残す（スペースを追加するだけで、削除・置換はしない）
- 2 回適用しても結果は変わらない（冪等）
- 次の内部には手を入れない
  - インラインコード（`` `code` ``）とコードフェンス（` ``` `）
  - URL（`https://...`）。URL に使える半角文字が続く範囲を URL とみなす。日本語は URL に含めないため、日本語を含む URL（`wiki/日本語abc`）は途中にスペースが入る
  - ファイルパス（`/usr/...`、`~/...`、`./...`、`../...`、`C:\...`）。空白・全角記号（、。」など）・二重引用符までをパスとみなす。`/` で始まるものは行頭・空白・開き括弧の直後にある場合だけパスとみなす
- 不正な UTF-8 を含む入力はそのまま出力する

### 未対応

- 数字と単位以外の半角文字同士の規則（`Ctrl+Alt` → `Ctrl + Alt` など）は扱わない
- 半角の一重引用符（`'test'`）。アポストロフィ（`don't`）と区別できないため

## インストール

### Homebrew（macOS / Linux）

```sh
brew install --cask capybara-translation/tap/jafmt
```

### go install

Go 1.25 以降が必要。

```sh
go install github.com/capybara-translation/jafmt/cmd/jafmt@latest
```

### ビルド済みバイナリ

[Releases](https://github.com/capybara-translation/jafmt/releases) から使っている環境のアーカイブをダウンロードする。同じリリースの `checksums.txt` で検証する。

### インストール先の確認

Automator から呼ぶときに絶対パスが必要なので、確認しておく。

```sh
command -v jafmt
```

| インストール方法 | 通常のインストール先 |
|---|---|
| Homebrew（Apple シリコン） | `/opt/homebrew/bin/jafmt` |
| Homebrew（Intel） | `/usr/local/bin/jafmt` |
| go install | `~/go/bin/jafmt`（`$(go env GOPATH)/bin`） |

`jafmt version` でインストールしたバージョンを確認できる。

## Automator でクイックアクションを作成する

1. Automator を開き、「新規書類」→「クイックアクション」を選ぶ
2. 上部の設定を次のようにする
   - 「ワークフローが受け取る現在の項目」: **テキスト**
   - 「検索対象」: **すべてのアプリケーション**
   - 「出力で選択されたテキストを置き換える」: **チェックを入れる**
3. 左のライブラリから「シェルスクリプトを実行」を右側へドラッグする
4. 「シェルスクリプトを実行」を次のように設定する
   - 「シェル」: `/bin/zsh`
   - 「入力の引き渡し方法」: **stdin へ**
   - スクリプト本文: [インストール先](#インストール先の確認)の絶対パスを書く（Automator のシェルは `.zshrc` を読まず、`/opt/homebrew/bin` や `~/go/bin` が `PATH` に含まれないため）

     ```sh
     /opt/homebrew/bin/jafmt
     ```

5. 「選択したテキストを整形する」などの名前で保存する（`~/Library/Services/` に保存される）

テキストを選択して右クリックし、「サービス」（アプリによってはメニュー直下）から保存した名前を選ぶと、選択範囲が整形済みテキストに置き換わる。アプリのメニューバーの「アプリ名 → サービス」からも実行できる。

## キーボードショートカットを割り当てる

1. 「システム設定」→「キーボード」を開く
2. 「キーボードショートカット...」ボタンを押す
3. 左の一覧から「サービス」を選ぶ
4. 「テキスト」の下にある作成したクイックアクションにチェックが入っていることを確認し、右端をダブルクリックしてショートカットキーを押す

アプリ側のショートカットと重なると、アプリ側が優先されて動かないことがある。`⌃⌥⌘` を組み合わせるなど、他と重なりにくいキーを選ぶ。

## 開発

テストとビルドは次のとおり。外部依存はない。

```sh
go test ./...
go build ./cmd/jafmt
```

- 整形ロジックは [formatter](formatter/) パッケージ、CLI は [cmd/jafmt](cmd/jafmt/) で入出力だけを担当する
- 整形ロジックは正規表現を使わず、ルーン単位で前後の文字種を判定する。文字種の判定は [charclass.go](formatter/charclass.go)、スペースを入れる規則は [formatter.go](formatter/formatter.go) の `spaceRules`、単位の一覧は [units.go](formatter/units.go)、引用符と強調の対応づけは [delimiters.go](formatter/delimiters.go)、整形しない範囲は [protect.go](formatter/protect.go) にある
- ルールを追加・変更するときは、先に [formatter_test.go](formatter/formatter_test.go) のテーブルへケースを追加してから実装する
- 冪等性と「スペースの挿入以外はしない」ことはファズテストでも確認できる

```sh
go test -run '^$' -fuzz=FuzzFormat -fuzztime=60s ./formatter
```

### CI とリリース

- push と pull request のたびに、GitHub Actions で gofmt・`go vet`・テスト（macOS / Linux / Windows）・60 秒のファズテスト・govulncheck・GoReleaser の試しビルドを実行する
- `v*` のタグを push すると、テストと govulncheck が通った後に GoReleaser が GitHub Release を作成し、[capybara-translation/homebrew-tap](https://github.com/capybara-translation/homebrew-tap) の cask を更新する
- tap へ push するため、リポジトリのシークレット `HOMEBREW_TAP_TOKEN` に tap へ書き込めるトークンを設定しておく

```sh
git tag v0.1.0
git push origin v0.1.0
```

## ライセンス

MIT。[LICENSE](LICENSE) を参照。
