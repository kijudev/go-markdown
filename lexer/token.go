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

	// Syntax
	TokenKindHash
	TokenKindDash
	TokenKindAsterisk
	TokenKindUnderscore
	TokenKindDot
	TokenKindEscape
	TokenKindEQ
)

func (k TokenKind) Rune() rune {
	switch k {
	case TokenKindHash:
		return '#'
	case TokenKindDash:
		return '-'
	case TokenKindAsterisk:
		return '*'
	case TokenKindUnderscore:
		return '_'
	case TokenKindDot:
		return '.'
	case TokenKindEscape:
		return '\\'
	case TokenKindEQ:
		return '='
	default:
		return -1
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

func runeIsSyntax(r rune) bool {
	return runeIsOneOf(
		TokenKindHash,
		TokenKindDash,
		TokenKindAsterisk,
		TokenKindUnderscore,
		TokenKindDot,
		TokenKindEscape,
		TokenKindEQ,
	)(r)
}

func runeIsNotSyntax(r rune) bool {
	return !runeIsSyntax(r)
}
