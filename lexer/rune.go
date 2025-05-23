package lexer

import "unicode"

type RuneKind uint8

const (
	RuneKindWhitespace RuneKind = iota
	RuneKindNewline
	RuneKindNumber
	RuneKindString

	RuneKindHash
	RuneKindDash
	RuneKindAsterisk
	RuneKindUnderscore
	RuneKindEscape
)

var runeKindDebugNames = []string{
	"WHITESPACE",
	"NEWLINE",
	"NUMBER",
	"STRING",

	"HASH",
	"DASH",
	"ASTERISK",
	"UNERSCORE",
	"ESCAPE",
}

func (k RuneKind) DebugName() string {
	if int(k) >= len(runeKindDebugNames) {
		panic("Invalid RuneKind.")
	}

	return runeKindDebugNames[k]
}

func (k RuneKind) Rune() rune {
	switch k {
	case RuneKindHash:
		return '#'
	case RuneKindDash:
		return '-'
	case RuneKindAsterisk:
		return '*'
	case RuneKindUnderscore:
		return '_'
	case RuneKindEscape:
		return '\\'
	default:
		panic("Invalid RuneKind.")
	}
}

func runeIs(k RuneKind) func(rune) bool {
	return func(r rune) bool {
		switch k {
		case RuneKindWhitespace:
			if unicode.IsSpace(r) {
				return true
			}
		case RuneKindNewline:
			if r == '\n' {
				return true
			}
		case RuneKindString:
			if unicode.IsLetter(r) {
				return true
			}
		case RuneKindNumber:
			if unicode.IsNumber(r) {
				return true
			}
		case RuneKindHash:
			if r == '#' {
				return true
			}
		case RuneKindDash:
			if r == '-' {
				return true
			}
		case RuneKindAsterisk:
			if r == '*' {
				return true
			}
		case RuneKindUnderscore:
			if r == '_' {
				return true
			}
		case RuneKindEscape:
			if r == '\\' {
				return true
			}
		}

		return false
	}
}

func runeIsOneOf(kinds ...RuneKind) func(rune) bool {
	return func(r rune) bool {
		for _, k := range kinds {
			if runeIs(k)(r) {
				return true
			}
		}

		return false
	}
}

func runeIsNotOneOf(kinds ...RuneKind) func(rune) bool {
	return func(r rune) bool {
		for _, k := range kinds {
			if runeIs(k)(r) {
				return false
			}
		}

		return true
	}
}
