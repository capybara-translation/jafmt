package formatter

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// units は数字との間にスペースを入れる単位。大文字・小文字を区別する。
// 「1990s」「2B」「4A」のような単位以外の用法と紛れやすい s・h・t・A は含めない。
// % と °（角度）はガイドで数字に付けると定められているため含めない。
var units = map[string]bool{
	// 長さ
	"nm": true, "mm": true, "cm": true, "m": true, "km": true,
	// 重さ
	"mg": true, "g": true, "kg": true,
	// 体積
	"mL": true, "ml": true, "L": true, "cc": true,
	// データ量・通信速度
	"B": true, "KB": true, "kB": true, "MB": true, "GB": true, "TB": true, "PB": true,
	"KiB": true, "MiB": true, "GiB": true, "TiB": true,
	"bit": true, "bps": true, "kbps": true, "Kbps": true, "Mbps": true, "Gbps": true,
	// 周波数
	"Hz": true, "kHz": true, "MHz": true, "GHz": true,
	// 時間
	"ns": true, "ms": true, "sec": true, "min": true,
	// 電気
	"V": true, "mV": true, "mA": true, "mAh": true, "W": true, "kW": true, "Wh": true, "kWh": true, "dB": true,
	// 画面・印刷
	"px": true, "pt": true, "dpi": true, "ppi": true, "fps": true,
	// マイクロ（μ はギリシャ文字の U+03BC。マイクロ記号 U+00B5 の µ も同じ扱いにする）
	"μm": true, "μg": true, "μL": true, "μs": true, "μA": true, "μV": true,
}

// numberThenUnit: 数値の直後に単位が続く（「50kg」→「50 kg」）。
// 「margin:16px」「image_100px.png」のような識別子の一部は対象外にするため、
// 数値と単位がそれぞれ独立したトークンになっている場合だけスペースを入れる。
func numberThenUnit(t *text, i int) bool {
	rs := t.rs
	if !isDigit(rs[i-1]) {
		return false
	}
	j := i
	if j < len(rs) && isMicro(rs[j]) {
		j++
	}
	for j < len(rs) && isASCIILetter(rs[j]) {
		j++
	}
	unit := strings.ReplaceAll(string(rs[i:j]), "µ", "μ")
	return units[unit] && isUnitEnd(rs, j) && isNumberStart(t, i-1)
}

func isMicro(r rune) bool {
	return r == 'μ' || r == 'µ'
}

// isNumberStart は t.rs[i] で終わる数値（「1,000」「1.5」）が独立したトークンとして始まっているかを返す。
func isNumberStart(t *text, i int) bool {
	rs := t.rs
	for i >= 0 && (isDigit(rs[i]) || rs[i] == '.' || rs[i] == ',') {
		i--
	}
	if !isDigit(rs[i+1]) {
		return false
	}
	// 別のルールで数値の直前にスペースが入る場合（「か?5GB」）も、冪等性を保つため区切りとみなす。
	if i < 0 || isNumberBoundary(rs[i]) || t.needsSpace(i+1) {
		return true
	}
	// 符号や範囲（「-5」「10-20」「10~20」）は、その前が区切りか数字なら数値の始まりとみなす。
	switch rs[i] {
	case '-', '+', '~':
		return i == 0 || isNumberBoundary(rs[i-1]) || isDigit(rs[i-1])
	}
	return false
}

// isNumberBoundary は数値の直前にあってよい文字かを返す。
func isNumberBoundary(r rune) bool {
	return unicode.IsSpace(r) || isJapanese(r) || isFullwidthSymbol(r) || isOpenBracket(r)
}

// isUnitEnd は rs[j] が単位の直後にあってよい文字か（単位がトークンとして終わっているか）を返す。
// 「.」は文末のピリオドとみなせる場合だけ許し、「100px.png」のような拡張子は除く。
func isUnitEnd(rs []rune, j int) bool {
	if j == len(rs) {
		return true
	}
	r := rs[j]
	switch r {
	case ')', ']', '/', '!', '?', ',', ';', ':': // 「(50kg)」「100km/h」「50kg,」など
		return true
	case '.':
		return j+1 == len(rs) || unicode.IsSpace(rs[j+1]) || rs[j+1] >= utf8.RuneSelf
	}
	return unicode.IsSpace(r) || isJapanese(r) || isFullwidthSymbol(r)
}
