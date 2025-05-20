package lexer

type TokenKind uint8

const (
	EOF TokenKind = iota
	UNDERSCORE_1
	UNDERSCORE_2
	ASTERISK_1
	ASTERISK_2
	HASH_1
	HASH_2
	HASH_3
	HASH_4
	HASH_5
	TAB
	STRING
)

type Token struct {
	Kind    TokenKind
	Literal string
}

func isSyntax(r rune) bool {
	if r == '_' || r == '*' || r == '#' {
		return true
	}

	return false
}

func (tk TokenKind) Literal() string {
	switch tk {
	case EOF:
		return "[EOF]"
	case UNDERSCORE_1:
		return "_"
	case UNDERSCORE_2:
		return "__"
	case ASTERISK_1:
		return "*"
	case ASTERISK_2:
		return "**"
	case HASH_1:
		return "#"
	case HASH_2:
		return "##"
	case HASH_3:
		return "###"
	case HASH_4:
		return "####"
	case HASH_5:
		return "#####"
	case TAB:
		return "[TAB]"
	case STRING:
		return "[STRING]"
	default:
		return "[UNKNOWN]"
	}
}
