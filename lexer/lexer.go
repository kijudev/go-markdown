package lexer

import (
	"bufio"
	"io"
	"unicode"
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
	case TokenKindRuneHash.Rune():
		return l.lexHash()
	case TokenKindRuneDash.Rune():
		return l.lexDash()
	case TokenKindRuneAsterisk.Rune():
		return l.lexAsterisk()
	case TokenKindRuneUnderscore.Rune():
		return l.lexUnderscore()
	case TokenKindRuneEscape.Rune():
		return l.lexEscape()
	default:
		if unicode.IsSpace(r) {
			return l.lexWhitespace()
		}

		if unicode.IsDigit(r) {
			return l.lexNumber()
		}

		return l.lexString()
	}
}

func (l *Lexer) lexHash() (TokenInfo, error) {
	lit, span, err := lexUntil(l.reader, runeIs(TokenKindRuneHash), 6)

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
	lit, span, err := lexUntil(l.reader, runeIs(TokenKindRuneDash), -1)

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
	lit, span, err := lexUntil(l.reader, runeIs(TokenKindRuneAsterisk), 3)

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
	lit, span, err := lexUntil(l.reader, runeIs(TokenKindRuneUnderscore), 3)

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

func (l *Lexer) lexEscape() (TokenInfo, error) {
	var lit string
	_, _, err := l.reader.ReadRune()

	if err == io.EOF {
		return TokenInfo{}, ErrEOF
	}

	if err != nil {
		return TokenInfo{}, ErrUnread
	}

	r, _, err := l.reader.ReadRune()

	if err == io.EOF {
		return TokenInfo{}, ErrEOF
	}

	if err != nil {
		return TokenInfo{}, ErrUnread
	}

	lit = string(r)

	return TokenInfo{
		Token: Token{TokenKindString, lit},
		Pos:   TokenPos{1, 1, 1},
	}, nil
}

func (l *Lexer) lexWhitespace() (TokenInfo, error) {
	lit, span, err := lexUntil(l.reader, unicode.IsSpace, -1)

	if err != nil {
		return TokenInfo{}, err
	}

	return TokenInfo{
		Token: Token{TokenKindWhitespace, lit},
		Pos:   TokenPos{span, span, span},
	}, nil
}

func (l *Lexer) lexNumber() (TokenInfo, error) {
	lit, span, err := lexUntil(l.reader, unicode.IsDigit, -1)

	if err != nil {
		return TokenInfo{}, err
	}

	r, _, err := peekOneRune(l.reader)

	if err == ErrEOF {
		return TokenInfo{
			Token: Token{TokenKindNumber, lit},
			Pos:   TokenPos{span, span, span},
		}, nil
	}

	if err != nil {
		return TokenInfo{}, err
	}

	if r == '.' {
		_, _, err := l.reader.ReadRune()
		if err != nil {
			return TokenInfo{}, ErrRead
		}

		return TokenInfo{
			Token: Token{TokenKindNumbering, lit},
			Pos:   TokenPos{span, span, span},
		}, nil
	}

	return TokenInfo{
		Token: Token{TokenKindNumber, lit},
		Pos:   TokenPos{span, span, span},
	}, nil
}

func (l *Lexer) lexString() (TokenInfo, error) {
	lit, span, err := lexUntil(l.reader, runeIsNotOneOf(
		TokenKindNewline,
		TokenKindNumber,
		TokenKindRuneHash,
		TokenKindRuneDash,
		TokenKindRuneAsterisk,
		TokenKindRuneUnderscore,
		TokenKindRuneEscape,
	), -1)

	if err != nil {
		return TokenInfo{}, err
	}

	return TokenInfo{
		Token: Token{TokenKindString, lit},
		Pos:   TokenPos{span, span, span},
	}, nil
}
