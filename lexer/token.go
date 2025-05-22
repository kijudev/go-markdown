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
	TokenKindHashRune
	TokenKindDashRune
	TokenKindAsteriskRune
	TokenKindUnderscoreRune
	TokenKindDotRune
	TokenKindEscapeRune

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

func (k TokenKind) Rune() rune {
	switch k {
	case TokenKindHashRune:
		return '#'
	case TokenKindDashRune:
		return '-'
	case TokenKindAsteriskRune:
		return '*'
	case TokenKindUnderscoreRune:
		return '_'
	case TokenKindDotRune:
		return '.'
	case TokenKindEscapeRune:
		return '\\'
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
		case TokenKindHashRune:
			if r == '#' {
				return true
			}
		case TokenKindDashRune:
			if r == '-' {
				return true
			}
		case TokenKindAsteriskRune:
			if r == '*' {
				return true
			}
		case TokenKindUnderscoreRune:
			if r == '_' {
				return true
			}
		case TokenKindDotRune:
			if r == '.' {
				return true
			}
		case TokenKindEscapeRune:
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

func runeIsSyntax(r rune) bool {
	return runeIsOneOf(
		TokenKindHashRune,
		TokenKindDashRune,
		TokenKindAsteriskRune,
		TokenKindUnderscoreRune,
		TokenKindDotRune,
		TokenKindEscapeRune,
	)(r)
}

func runeIsNotSyntax(r rune) bool {
	return !runeIsSyntax(r)
}
