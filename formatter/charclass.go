package formatter

import (
	"unicode"
	"unicode/utf8"
)

// isJapanese は漢字・ひらがな・カタカナ・長音記号「ー」かを返す。
// 全角記号（、。「」など）や半角カタカナは含めない。
func isJapanese(r rune) bool {
	switch {
	case r == 'ー': // Unicode の用字は Common なので個別に扱う
		return true
	case r >= 0xFF61 && r <= 0xFF9F: // 半角カタカナ
		return false
	}
	return unicode.In(r, unicode.Han, unicode.Hiragana, unicode.Katakana)
}

// isAlnum は半角英数字（アクセント付きラテン文字を含む）かを返す。
// 全角英数字は含めない。
func isAlnum(r rune) bool {
	if r < utf8.RuneSelf {
		return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || '0' <= r && r <= '9'
	}
	if r >= 0xFF21 && r <= 0xFF5A { // 全角英字
		return false
	}
	return unicode.Is(unicode.Latin, r) && unicode.IsLetter(r)
}

// isLeadingAttached は語頭で英数字に付く記号（#123、@user、$100）かを返す。
func isLeadingAttached(r rune) bool {
	switch r {
	case '#', '@', '$':
		return true
	}
	return false
}

// isTrailingAttached は語末で英数字に付く記号（50%、45°、C#、C++）かを返す。
func isTrailingAttached(r rune) bool {
	switch r {
	case '%', '°', '#', '+':
		return true
	}
	return false
}

// isOpeningBracket は直後をトークンの先頭とみなす開き括弧・引用符かを返す。
func isOpeningBracket(r rune) bool {
	switch r {
	case '(', '[', '{', '<', '"', '\'',
		'「', '『', '（', '【', '〈', '《', '［', '｛', '“', '‘':
		return true
	}
	return false
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

// isFullwidthSymbol は日本語・英数字以外の非 ASCII 文字（、。」など）かを返す。
func isFullwidthSymbol(r rune) bool {
	return r >= utf8.RuneSelf && !isJapanese(r) && !unicode.IsLetter(r) && !unicode.IsDigit(r)
}

// isOpenBracket と isCloseBracket は、外側の日本語との間にスペースを入れる半角括弧かを返す。
func isOpenBracket(r rune) bool  { return r == '(' || r == '[' }
func isCloseBracket(r rune) bool { return r == ')' || r == ']' }

// isSpanTerminator はファイルパスの終わりとみなす文字（空白と全角記号）かを返す。
func isSpanTerminator(r rune) bool {
	return unicode.IsSpace(r) || isFullwidthSymbol(r)
}
