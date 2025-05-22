package lexerold

type TokenKind uint8

const (
	EOF TokenKind = iota
	NEWLINE
	UNDERSCORE_1
	UNDERSCORE_2
	UNDERSCORE_3
	ASTERISK_1
	ASTERISK_2
	ASTERISK_3
	HASH_1
	HASH_2
	HASH_3
	HASH_4
	HASH_5
	HASH_6
	INDENTATION
	STRING
	NUMBERING
	DASH_SINGLE
	DASH_MULTIPLE
	EQ
	ESCAPE
)

type Token struct {
	Kind    TokenKind
	Literal string
}

func isSyntax(r rune) bool {
	if r == '_' || r == '*' || r == '#' || r == '\n' || r == '-' || r == '=' {
		return true
	}

	return false
}

func (tk TokenKind) Literal() string {
	switch tk {
	case EOF:
		return "[EOF]"
	case NEWLINE:
		return "[NEWLINE]"
	case UNDERSCORE_1:
		return "_"
	case UNDERSCORE_2:
		return "__"
	case UNDERSCORE_3:
		return "___"
	case ASTERISK_1:
		return "*"
	case ASTERISK_2:
		return "**"
	case ASTERISK_3:
		return "***"
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
	case HASH_6:
		return "######"
	case INDENTATION:
		return "[INDENTATION]"
	case STRING:
		return "[STRING]"
	case DASH_SINGLE:
		return "-"
	case NUMBERING:
		return "[NUMBERING]"
	case DASH_MULTIPLE:
		return "[DASH_MULTIPLE]"
	case ESCAPE:
		return "[ESCAPE]"
	case EQ:
		return "[EQ]"
	default:
		return "[UNKNOWN]"
	}
}
