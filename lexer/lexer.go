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
	case '#':
		return l.lexHash()
	default:
		return l.lexString()
	}
}

func (l *Lexer) lexHash() (TokenInfo, error) {
	lit, span, err := lexWhen(l.reader, func(r rune) bool { return r == '#' }, 6)

	if err != nil {
		return TokenInfo{}, err
	}

	return TokenInfo{
		Token: Token{TokenKindHash, lit},
		Pos:   TokenPos{span, span, span},
	}, nil
}

func (l *Lexer) lexString() (TokenInfo, error) {
	lit, span, err := lexWhen(l.reader, func(r rune) bool { return true }, -1)

	if err != nil {
		return TokenInfo{}, err
	}

	return TokenInfo{
		Token: Token{TokenKindString, lit},
		Pos:   TokenPos{span, span, span},
	}, nil
}
