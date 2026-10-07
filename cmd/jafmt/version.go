package main

import (
	"regexp"
	"runtime/debug"
	"strings"
)

// version はリリース時に -ldflags "-X main.version=vX.Y.Z" で設定される（GoReleaser が行う）。
// 設定されていなければ、currentVersion は Go がバイナリに埋め込むモジュールのバージョンを使う。
// `go install ...@vX.Y.Z` でインストールした場合はそのタグになる。
var version = "dev"

// pseudoVersionSuffix は Go の疑似バージョンに共通する、末尾のタイムスタンプとコミットに一致する。
//
//	vX.0.0-yyyymmddhhmmss-abcdefabcdef          （タグがまだない）
//	vX.Y.(Z+1)-0.yyyymmddhhmmss-abcdefabcdef    （vX.Y.Z より後のコミット）
//	vX.Y.Z-pre.0.yyyymmddhhmmss-abcdefabcdef    （プレリリースより後のコミット）
var pseudoVersionSuffix = regexp.MustCompile(`[-.]\d{14}-[0-9a-f]{12}$`)

// currentVersion はこのバイナリのバージョンを返す。
func currentVersion() string {
	info, _ := debug.ReadBuildInfo()
	return resolveVersion(version, info)
}

// resolveVersion は ldflags の値を優先し、なければビルド情報からバージョンを決める。
// Go 1.24 以降、タグのないコミットを go build すると疑似バージョンが埋め込まれ、
// 未コミットの変更があると "+dirty" が付く。これらはリリースと見分けがつかず、
// 不具合報告でどのビルドか曖昧になるため "dev" として扱う。
func resolveVersion(ldVersion string, info *debug.BuildInfo) string {
	if ldVersion != "dev" {
		return ldVersion
	}
	if info == nil {
		return "dev"
	}
	v := info.Main.Version
	if v == "" || v == "(devel)" || strings.Contains(v, "+dirty") || pseudoVersionSuffix.MatchString(v) {
		return "dev"
	}
	return v
}
