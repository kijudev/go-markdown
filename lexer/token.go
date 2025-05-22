package lexer

import "unicode"

type TokenKind uint8

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
}

const (
	// Base
	TokenKindWhitespace TokenKind = iota
	TokenKindNewline
	TokenKindString
	TokenKindNumber

	// Syntax (rune)
	TokenKindRuneHash
	TokenKindRuneDash
	TokenKindRuneAsterisk
	TokenKindRuneUnderscore
	TokenKindRuneEscape

	// Semantic
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

	"RUNE_HASH",
	"RUNE_DASH",
	"RUNE_ASTERISK",
	"RUNE_UNDERSCORE",
	"RUNE_ESCAPE",

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
	return tokenKindDebugNames[k]
}

func (k TokenKind) Rune() rune {
	switch k {
	case TokenKindRuneHash:
		return '#'
	case TokenKindRuneDash:
		return '-'
	case TokenKindRuneAsterisk:
		return '*'
	case TokenKindRuneUnderscore:
		return '_'
	case TokenKindRuneEscape:
		return '\\'
	default:
		panic("Cannot access the rune value from non-rune tokens.")
	}
}

func runeIs(k TokenKind) func(rune) bool {
	return func(r rune) bool {
		switch k {
		case TokenKindWhitespace:
			if unicode.IsSpace(r) {
				return true
			}
		case TokenKindNewline:
			if r == '\n' {
				return true
			}
		case TokenKindString:
			if unicode.IsLetter(r) {
				return true
			}
		case TokenKindNumber:
			if unicode.IsNumber(r) {
				return true
			}
		case TokenKindRuneHash:
			if r == '#' {
				return true
			}
		case TokenKindRuneDash:
			if r == '-' {
				return true
			}
		case TokenKindRuneAsterisk:
			if r == '*' {
				return true
			}
		case TokenKindRuneUnderscore:
			if r == '_' {
				return true
			}
		case TokenKindRuneEscape:
			if r == '\\' {
				return true
			}
		}

		return false
	}
}

func runeIsOneOf(kinds ...TokenKind) func(rune) bool {
	return func(r rune) bool {
		for _, k := range kinds {
			if runeIs(k)(r) {
				return true
			}
		}

		return false
	}
}

func runeIsNotOneOf(kinds ...TokenKind) func(rune) bool {
	return func(r rune) bool {
		for _, k := range kinds {
			if runeIs(k)(r) {
				return false
			}
		}

		return true
	}
}
