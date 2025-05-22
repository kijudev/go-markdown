package lexer

import (
	"bufio"
	"io"
)

type Lexer struct {
	source io.Reader
	reader *bufio.Reader
	pos    TokenPos
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		source: reader,
		reader: bufio.NewReader(reader),
	}
}

func (l *Lexer) Tokenize() ([]TokenInfo, error) {
	l.Reset()
	var tokens []TokenInfo

	for {
		ti, err := l.Lex()

		if err == ErrEOF {
			break
		}

		if err != nil {
			return tokens, err
		}

		tokens = append(tokens, ti)
	}

	return tokens, nil
}

func (l *Lexer) Reset() {
	l.reader.Reset(l.source)
	l.pos = TokenPos{
		X:   0,
		Y:   0,
		Abs: 0,
	}
}

func (l *Lexer) Lex() (TokenInfo, error) {
	r, _, err := peekOneRune(l.reader)

	if err != nil {
		return TokenInfo{}, err
	}

	switch r {
	case TokenKindHashRune.Rune():
		return l.lexHash()
	case TokenKindDashRune.Rune():
		return l.lexDash()
	case TokenKindAsteriskRune.Rune():
		return l.lexAsterisk()
	case TokenKindUnderscoreRune.Rune():
		return l.lexUnderscore()
	case TokenKindDotRune.Rune():
		return l.lexDot()
	case TokenKindEscapeRune.Rune():
		return l.lexEscape()
	default:
		return l.lexString()
	}
}

func (l *Lexer) lexHash() (TokenInfo, error) {
	lit, span, err := lexUntil(l.reader, runeIs(TokenKindHashRune), 6)

	if err != nil {
		return TokenInfo{}, err
	}

	var kind TokenKind
	switch span {
	case 1:
		kind = TokenKindHeading1
	case 2:
		kind = TokenKindHeading2
	case 3:
		kind = TokenKindHeading3
	case 4:
		kind = TokenKindHeading4
	case 5:
		kind = TokenKindHeading5
	case 6:
		kind = TokenKindHeading6
	default:
		kind = TokenKindString
	}

	return TokenInfo{
		Token: Token{kind, lit},
		Pos:   TokenPos{span, span, span},
	}, nil
}

func (l *Lexer) lexDash() (TokenInfo, error) {
	lit, span, err := lexUntil(l.reader, runeIs(TokenKindDashRune), -1)

	if err != nil {
		return TokenInfo{}, err
	}

	var kind TokenKind
	switch span {
	case 1:
		kind = TokenKindDash
	case 2:
		kind = TokenKindString
	default:
		kind = TokenKindSpacer
	}

	return TokenInfo{
		Token: Token{kind, lit},
		Pos:   TokenPos{span, span, span},
	}, nil
}

func (l *Lexer) lexAsterisk() (TokenInfo, error) {
	lit, span, err := lexUntil(l.reader, runeIs(TokenKindAsteriskRune), 3)

	if err != nil {
		return TokenInfo{}, err
	}

	var kind TokenKind
	switch span {
	case 1:
		kind = TokenKindItalic
	case 2:
		kind = TokenKindBold
	case 3:
		kind = TokenKindBoldItalic
	default:
		kind = TokenKindString
	}

	return TokenInfo{
		Token: Token{kind, lit},
		Pos:   TokenPos{span, span, span},
	}, nil
}

func (l *Lexer) lexUnderscore() (TokenInfo, error) {
	lit, span, err := lexUntil(l.reader, runeIs(TokenKindUnderscoreRune), 3)

	if err != nil {
		return TokenInfo{}, err
	}

	return TokenInfo{
		Token: Token{TokenKindUnderscoreRune, lit},
		Pos:   TokenPos{span, span, span},
	}, nil
}

func (l *Lexer) lexDot() (TokenInfo, error) {
	return TokenInfo{
		Token: Token{TokenKindDotRune, "."},
		Pos:   TokenPos{1, 1, 1},
	}, nil
}

func (l *Lexer) lexEscape() (TokenInfo, error) {
	return TokenInfo{
		Token: Token{TokenKindEscapeRune, "\\"},
		Pos:   TokenPos{1, 1, 1},
	}, nil
}

func (l *Lexer) lexString() (TokenInfo, error) {
	lit, span, err := lexUntil(l.reader, runeIsNotSyntax, -1)

	if err != nil {
		return TokenInfo{}, err
	}

	return TokenInfo{
		Token: Token{TokenKindString, lit},
		Pos:   TokenPos{span, span, span},
	}, nil
}
