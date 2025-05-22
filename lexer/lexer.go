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
	case TokenKindHash.Rune():
		return l.lexHash()
	case TokenKindDash.Rune():
		return l.lexDash()
	case TokenKindAsterisk.Rune():
		return l.lexAsterisk()
	case TokenKindUnderscore.Rune():
		return l.lexUnderscore()
	case TokenKindDot.Rune():
		return l.lexDot()
	case TokenKindEscape.Rune():
		return l.lexEscape()
	case TokenKindEQ.Rune():
		return l.lexEQ()
	default:
		return l.lexString()
	}
}

func (l *Lexer) lexHash() (TokenInfo, error) {
	lit, span, err := lexUntil(l.reader, runeIs(TokenKindHash), 6)

	if err != nil {
		return TokenInfo{}, err
	}

	return TokenInfo{
		Token: Token{TokenKindHash, lit},
		Pos:   TokenPos{span, span, span},
	}, nil
}

func (l *Lexer) lexDash() (TokenInfo, error) {
	lit, span, err := lexUntil(l.reader, runeIs(TokenKindDash), -1)

	if err != nil {
		return TokenInfo{}, err
	}

	return TokenInfo{
		Token: Token{TokenKindDash, lit},
		Pos:   TokenPos{span, span, span},
	}, nil
}

func (l *Lexer) lexAsterisk() (TokenInfo, error) {
	lit, span, err := lexUntil(l.reader, runeIs(TokenKindAsterisk), 3)

	if err != nil {
		return TokenInfo{}, err
	}

	return TokenInfo{
		Token: Token{TokenKindAsterisk, lit},
		Pos:   TokenPos{span, span, span},
	}, nil
}

func (l *Lexer) lexUnderscore() (TokenInfo, error) {
	lit, span, err := lexUntil(l.reader, runeIs(TokenKindUnderscore), 3)

	if err != nil {
		return TokenInfo{}, err
	}

	return TokenInfo{
		Token: Token{TokenKindUnderscore, lit},
		Pos:   TokenPos{span, span, span},
	}, nil
}

func (l *Lexer) lexDot() (TokenInfo, error) {
	return TokenInfo{
		Token: Token{TokenKindDot, "."},
		Pos:   TokenPos{1, 1, 1},
	}, nil
}

func (l *Lexer) lexEscape() (TokenInfo, error) {
	return TokenInfo{
		Token: Token{TokenKindEscape, "\\"},
		Pos:   TokenPos{1, 1, 1},
	}, nil
}

func (l *Lexer) lexEQ() (TokenInfo, error) {
	lit, span, err := lexUntil(l.reader, runeIs(TokenKindEQ), -1)

	if err != nil {
		return TokenInfo{}, err
	}

	return TokenInfo{
		Token: Token{TokenKindEQ, lit},
		Pos:   TokenPos{span, span, span},
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
