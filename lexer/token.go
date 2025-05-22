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
	TokenKindWhitespace TokenKind = iota
	TokenKindNewline
	TokenKindString
	TokenKindNumber

	TokenKindHash
	TokenKindDash
	TokenKindAsterisk
	TokenKindUnderscore
	TokenKindDot
	TokenKindEscape
	TokenKindEQ
)

func IsOneOf(r rune, kinds TokenKind) bool {
	for k := range kinds {
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
		case TokenKindHash:
			if r == '#' {
				return true
			}
		case TokenKindDash:
			if r == '-' {
				return true
			}
		case TokenKindAsterisk:
			if r == '*' {
				return true
			}
		case TokenKindUnderscore:
			if r == '_' {
				return true
			}
		case TokenKindDot:
			if r == '.' {
				return true
			}
		case TokenKindEscape:
			if r == '\\' {
				return true
			}
		case TokenKindEQ:
			if r == '=' {
				return true
			}
		}
	}

	return false
}
