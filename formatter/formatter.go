// Package formatter は日本語と半角英数字の間に半角スペースを入れる整形を行う。
//
// ルールは Microsoft Japanese Localization Style Guide 4.1.11 "Symbols & spaces" に基づく。
package formatter

import (
	"strings"
	"unicode/utf8"
)

// text は整形対象のルーン列と、整形してはいけない範囲の情報を持つ。
type text struct {
	rs    []rune
	spans []int       // spans[i] は rs[i] が属する保護範囲の番号（0 は範囲外）
	kinds []spanKind  // kinds[id] は保護範囲の種類
	roles []delimRole // roles[i] は rs[i] の囲み記号としての役割
}

// spaceRule は t.rs[i-1] と t.rs[i] の間に半角スペースを入れるべきかを判定する。
type spaceRule func(t *text, i int) bool

// spaceRules のいずれかが true を返す位置にスペースを入れる。
// numberThenUnit が needsSpace を参照して初期化が循環するため、init で設定する。
var spaceRules []spaceRule

func init() {
	spaceRules = []spaceRule{
		japaneseThenWord,
		wordThenJapanese,
		japaneseThenBracket,
		bracketThenJapanese,
		questionThenWord,
		japaneseThenCode,
		codeThenJapanese,
		japaneseThenDelim,
		delimThenJapanese,
		numberThenUnit,
	}
}

// Format は s を整形した結果を返す。
// 不正な UTF-8 を含む入力は壊さないよう、そのまま返す。
func Format(s string) string {
	if !utf8.ValidString(s) {
		return s
	}
	t := &text{rs: []rune(s)}
	t.spans, t.kinds = protectedSpans(t.rs)
	t.roles = pairedDelimiters(t.rs, t.spans)

	var b strings.Builder
	b.Grow(len(s) + len(s)/4)
	for i, r := range t.rs {
		if i > 0 && !t.insideSameSpan(i) && t.needsSpace(i) {
			b.WriteByte(' ')
		}
		b.WriteRune(r)
	}
	return b.String()
}

func (t *text) insideSameSpan(i int) bool {
	return t.spans[i] != 0 && t.spans[i-1] == t.spans[i]
}

func (t *text) needsSpace(i int) bool {
	for _, rule := range spaceRules {
		if rule(t, i) {
			return true
		}
	}
	return false
}

// japaneseThenWord: 日本語の直後に英数字の語が始まる（「日本語abc」「は#123」）。
func japaneseThenWord(t *text, i int) bool {
	return isJapanese(t.rs[i-1]) && startsWord(t.rs, i)
}

// wordThenJapanese: 英数字の語の直後に日本語が続く（「abc日本語」「50%の」）。
func wordThenJapanese(t *text, i int) bool {
	return endsWord(t.rs, i-1) && isJapanese(t.rs[i])
}

// japaneseThenBracket: 日本語の直後に半角の開き括弧が続く（「日本語(Japanese)」「の[詳細]」）。
// 括弧の内側（「(タイトル)」「[新規]」）には入れない。
func japaneseThenBracket(t *text, i int) bool {
	return isJapanese(t.rs[i-1]) && isOpenBracket(t.rs[i])
}

// bracketThenJapanese: 半角の閉じ括弧の直後に日本語が続く（「(S)を」「[新規]を」）。
func bracketThenJapanese(t *text, i int) bool {
	return isCloseBracket(t.rs[i-1]) && isJapanese(t.rs[i])
}

// questionThenWord: 日本語の文末の「?」「!」の直後に英数字の語が始まる（「しますか?Excel」）。
func questionThenWord(t *text, i int) bool {
	j := i - 1
	for j >= 0 && (t.rs[j] == '?' || t.rs[j] == '!') {
		j--
	}
	return j < i-1 && j >= 0 && isJapanese(t.rs[j]) && startsWord(t.rs, i)
}

// japaneseThenCode: 日本語の直後にインラインコードが始まる（「これは`code`」）。
func japaneseThenCode(t *text, i int) bool {
	return isJapanese(t.rs[i-1]) && t.spanStartsAt(i, spanCode)
}

// codeThenJapanese: インラインコードの直後に日本語が続く（「`code`です」）。
func codeThenJapanese(t *text, i int) bool {
	return t.spanEndsAt(i-1, spanCode) && isJapanese(t.rs[i])
}

// japaneseThenDelim: 日本語の直後に引用符や強調が始まる（「これは"OK"」「これは**bold**」）。
func japaneseThenDelim(t *text, i int) bool {
	return isJapanese(t.rs[i-1]) && t.roles[i] == delimOpen
}

// delimThenJapanese: 引用符や強調の直後に日本語が続く（「"test"と」「**bold**です」）。
func delimThenJapanese(t *text, i int) bool {
	return t.roles[i-1] == delimClose && isJapanese(t.rs[i])
}

// spanStartsAt は t.rs[i] が kind の保護範囲の先頭かを返す。
func (t *text) spanStartsAt(i int, kind spanKind) bool {
	id := t.spans[i]
	return id != 0 && t.kinds[id] == kind && (i == 0 || t.spans[i-1] != id)
}

// spanEndsAt は t.rs[i] が kind の保護範囲の末尾かを返す。
func (t *text) spanEndsAt(i int, kind spanKind) bool {
	id := t.spans[i]
	return id != 0 && t.kinds[id] == kind && (i == len(t.rs)-1 || t.spans[i+1] != id)
}

// startsWord は rs[i] から英数字の語が始まるかを返す。語頭に付く記号（#、@、$）は飛ばして判定する。
func startsWord(rs []rune, i int) bool {
	for i < len(rs) && isLeadingAttached(rs[i]) {
		i++
	}
	return i < len(rs) && isAlnum(rs[i])
}

// endsWord は rs[i] で英数字の語が終わるかを返す。語末に付く記号（%、°、#、+）は飛ばして判定する。
func endsWord(rs []rune, i int) bool {
	for i >= 0 && isTrailingAttached(rs[i]) {
		i--
	}
	return i >= 0 && isAlnum(rs[i])
}
