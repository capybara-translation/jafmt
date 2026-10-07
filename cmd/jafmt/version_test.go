package main

import (
	"runtime/debug"
	"testing"
)

func TestResolveVersion(t *testing.T) {
	tests := []struct {
		name      string
		ldVersion string
		info      *debug.BuildInfo
		want      string
	}{
		{"ldflags の値を優先", "v1.2.3", buildInfo("v9.9.9"), "v1.2.3"},
		{"ビルド情報なし", "dev", nil, "dev"},
		{"タグ付きのモジュールバージョン（go install @vX.Y.Z）", "dev", buildInfo("v0.1.0"), "v0.1.0"},
		{"(devel) は隠す", "dev", buildInfo("(devel)"), "dev"},
		{"空は隠す", "dev", buildInfo(""), "dev"},
		// Go Modules Reference にある 3 種類の疑似バージョン。
		// Go 1.24 以降、タグのないコミットを go build すると埋め込まれる。
		{"疑似バージョン/タグなし", "dev", buildInfo("v0.0.0-20260928154706-9667b859f76d"), "dev"},
		{"疑似バージョン/タグより後のコミット", "dev", buildInfo("v0.1.1-0.20260928154707-bf85eeda22a3"), "dev"},
		{"疑似バージョン/プレリリースより後のコミット", "dev", buildInfo("v0.2.0-rc.1.0.20260928154707-bf85eeda22a3"), "dev"},
		{"未コミットの変更がある疑似バージョン", "dev", buildInfo("v0.1.1-0.20260928154707-bf85eeda22a3+dirty"), "dev"},
		{"未コミットの変更があるタグ", "dev", buildInfo("v0.1.0+dirty"), "dev"},
		{"プレリリースのタグはリリース扱い", "dev", buildInfo("v0.2.0-rc.1"), "v0.2.0-rc.1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveVersion(tt.ldVersion, tt.info); got != tt.want {
				t.Errorf("resolveVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}

func buildInfo(version string) *debug.BuildInfo {
	return &debug.BuildInfo{Main: debug.Module{Version: version}}
}
