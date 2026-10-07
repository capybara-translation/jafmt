package formatter

import "unicode"

// delimRole は対になる囲み記号（半角の二重引用符・Markdown の強調）における役割。
// 保護範囲と違い、囲まれた内側も通常どおり整形する。
type delimRole int8

const (
	delimNone  delimRole = iota
	delimOpen            // 開き記号の先頭の文字
	delimClose           // 閉じ記号の末尾の文字
)

// pairedDelimiters は対になっている囲み記号を求め、各文字の役割を返す。
// 誤判定を避けるため、対は同じ行の中でだけ作り、保護範囲（コードなど）の中の記号は数えない。
func pairedDelimiters(rs []rune, spans []int) []delimRole {
	roles := make([]delimRole, len(rs))
	pairQuotes(rs, spans, roles)
	pairEmphasis(rs, spans, roles, '*', func(n int) bool { return n <= 3 }) // *斜体* **太字** ***太字斜体***
	pairEmphasis(rs, spans, roles, '~', func(n int) bool { return n == 2 }) // ~~取り消し~~。単独の ~ は範囲（10~20）に使われるため除く
	return roles
}

// pairQuotes は半角の二重引用符を、行の中で前から順に 2 つずつ対にする。
func pairQuotes(rs []rune, spans []int, roles []delimRole) {
	open := -1
	for i, r := range rs {
		switch {
		case r == '\n':
			open = -1
		case r != '"' || spans[i] != 0:
		case open < 0:
			open = i
		default:
			roles[open], roles[i] = delimOpen, delimClose
			open = -1
		}
	}
}

// pairEmphasis は c の連続（**、~~ など）を、同じ長さの連続と対にする。
// CommonMark と同様に、開きの直後と閉じの直前が空白のものは強調とみなさない（「2 * 3」「* 箇条書き」）。
func pairEmphasis(rs []rune, spans []int, roles []delimRole, c rune, validLen func(int) bool) {
	for i := 0; i < len(rs); {
		n := delimRun(rs, spans, i, c)
		if n == 0 {
			i++
			continue
		}
		if !validLen(n) || i+n >= len(rs) || unicode.IsSpace(rs[i+n]) {
			i += n
			continue
		}
		j := closingRun(rs, spans, i+n, c, n)
		if j < 0 {
			i += n
			continue
		}
		roles[i], roles[j+n-1] = delimOpen, delimClose
		i = j + n
	}
}

// closingRun は rs[from] 以降の同じ行で、閉じとして使える長さ n の c の連続の位置を返す。なければ -1。
func closingRun(rs []rune, spans []int, from int, c rune, n int) int {
	for j := from; j < len(rs) && rs[j] != '\n'; {
		m := delimRun(rs, spans, j, c)
		if m == 0 {
			j++
			continue
		}
		if m == n && !unicode.IsSpace(rs[j-1]) {
			return j
		}
		j += m
	}
	return -1
}

// delimRun は rs[i] から続く c の個数を返す。保護範囲の中の文字は数えない。
func delimRun(rs []rune, spans []int, i int, c rune) int {
	n := 0
	for i+n < len(rs) && rs[i+n] == c && spans[i+n] == 0 {
		n++
	}
	return n
}
