package formatter

type spanKind int

const (
	spanNone spanKind = iota
	spanCode
	spanURL
	spanPath
)

// protectedSpans は整形してはいけない範囲（インラインコード・URL・ファイルパス）を求める。
// 戻り値の ids[i] は rs[i] が属する範囲の番号で、どの範囲にも属さなければ 0。
// kinds[id] はその範囲の種類。
func protectedSpans(rs []rune) (ids []int, kinds []spanKind) {
	ids = make([]int, len(rs))
	kinds = []spanKind{spanNone}
	for i := 0; i < len(rs); {
		kind, end := spanCode, codeSpanEnd(rs, i)
		if end < 0 {
			kind, end = spanURL, urlEnd(rs, i)
		}
		if end < 0 {
			kind, end = spanPath, pathEnd(rs, i)
		}
		if end < 0 {
			i++
			continue
		}
		kinds = append(kinds, kind)
		for j := i; j < end; j++ {
			ids[j] = len(kinds) - 1
		}
		i = end
	}
	return ids, kinds
}

// codeSpanEnd は rs[i] から始まるバッククォートの範囲の終端（排他）を返す。
// 開きと同じ長さのバッククォート列で閉じる（Markdown のコードスパンとフェンス）。
// 閉じがなければ -1 を返し、バッククォートは通常の文字として扱う。
func codeSpanEnd(rs []rune, i int) int {
	if rs[i] != '`' {
		return -1
	}
	n := backtickRun(rs, i)
	for j := i + n; j < len(rs); {
		if rs[j] != '`' {
			j++
			continue
		}
		m := backtickRun(rs, j)
		if m == n {
			return j + m
		}
		j += m
	}
	return -1
}

func backtickRun(rs []rune, i int) int {
	n := 0
	for i+n < len(rs) && rs[i+n] == '`' {
		n++
	}
	return n
}

// urlEnd は rs[i] から始まる URL（scheme://...）の終端（排他）を返す。URL でなければ -1。
// URL に使える半角文字が続く範囲を URL とみなし、日本語は含めない。
func urlEnd(rs []rune, i int) int {
	if !isASCIILetter(rs[i]) || i > 0 && isSchemeChar(rs[i-1]) {
		return -1
	}
	j := i
	for j < len(rs) && isSchemeChar(rs[j]) {
		j++
	}
	if !hasPrefix(rs[j:], "://") {
		return -1
	}
	j += 3
	for j < len(rs) && isURLChar(rs[j]) {
		j++
	}
	return j
}

// pathEnd は rs[i] から始まるファイルパスの終端（排他）を返す。パスでなければ -1。
// 「/」で始まるものは「日本/US」のような区切りと区別するため、トークンの先頭にある場合だけパスとみなす。
// 日本語を含むパスも多いため、空白・全角記号・二重引用符までをパスとみなす。
func pathEnd(rs []rune, i int) int {
	tokenStart := i == 0 || isPathTerminator(rs[i-1]) || isOpeningBracket(rs[i-1])
	var prefixLen int
	switch {
	case hasPrefix(rs[i:], "~/"), hasPrefix(rs[i:], "./"):
		prefixLen = 2
	case hasPrefix(rs[i:], "../"):
		prefixLen = 3
	case isASCIILetter(rs[i]) && hasPrefix(rs[i+1:], `:\`) && (i == 0 || !isAlnum(rs[i-1])):
		prefixLen = 3
	case rs[i] == '/' && tokenStart:
		prefixLen = 1
	default:
		return -1
	}
	j := i + prefixLen
	if j >= len(rs) || isPathTerminator(rs[j]) {
		return -1
	}
	for j < len(rs) && !isPathTerminator(rs[j]) {
		j++
	}
	return j
}

func hasPrefix(rs []rune, prefix string) bool {
	i := 0
	for _, p := range prefix {
		if i >= len(rs) || rs[i] != p {
			return false
		}
		i++
	}
	return true
}

func isASCIILetter(r rune) bool {
	return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z'
}

func isSchemeChar(r rune) bool {
	return isASCIILetter(r) || '0' <= r && r <= '9' || r == '+' || r == '-' || r == '.'
}

// isURLChar は RFC 3986 で URL に使える文字（英数字・予約文字・非予約文字・%）かを返す。
func isURLChar(r rune) bool {
	if isASCIILetter(r) || '0' <= r && r <= '9' {
		return true
	}
	switch r {
	case '-', '.', '_', '~', ':', '/', '?', '#', '[', ']', '@',
		'!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=', '%':
		return true
	}
	return false
}
