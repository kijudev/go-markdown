package lexer

import (
	"bufio"
	"io"
)

type Lexer struct {
	reader *bufio.Reader
	posrow uint
	poscol uint
	pos    uint
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		reader: bufio.NewReader(reader),
		posrow: 0,
		poscol: 0,
		pos:    0,
	}
}

func (l *Lexer) SafeLex() (Token, uint, error) {
	r, _, err := l.reader.ReadRune()

	if err != nil && err == io.EOF {
		return Token{EOF, ""}, l.pos, nil
	}
}
