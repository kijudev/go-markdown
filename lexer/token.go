package lexer

type TokenKind uint8
type runeKind uint8

type Token struct {
	Kind TokenKind
	Lit  string
}

type TokenPos struct {
	X   int
	Y   int
	Abs int
}

type TokenInfo struct {
	Token Token
	Pos   TokenPos
	Span  int
	Raw   string
}

const (
	// Base
	TokenKindWhitespace TokenKind = iota
	TokenKindNewline
	TokenKindString
	TokenKindNumber

	TokenKindHeading1
	TokenKindHeading2
	TokenKindHeading3
	TokenKindHeading4
	TokenKindHeading5
	TokenKindHeading6
	TokenKindDash
	TokenKindSpacer
	TokenKindBold
	TokenKindItalic
	TokenKindBoldItalic
	TokenKindNumbering
)

var tokenKindDebugNames = []string{
	"WHITESPACE",
	"NEWLINE",
	"STRING",
	"NUMBER",

	"HEADING_1",
	"HEADING_2",
	"HEADING_3",
	"HEADING_4",
	"HEADING_5",
	"HEADING_6",
	"DASH",
	"SPACER",
	"BOLD",
	"ITALIC",
	"BOLD_ITALIC",
	"NUMBERING",
}

func (k TokenKind) DebugName() string {
	if int(k) > len(tokenKindDebugNames) {
		panic("Invalid TokenKind.")
	}

	return tokenKindDebugNames[k]
}
