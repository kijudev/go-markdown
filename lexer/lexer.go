package lexer

import (
	"bufio"
	"errors"
	"io"
)

type Lexer struct {
	reader *bufio.Reader
	posrow int
	poscol int
	pos    int
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		reader: bufio.NewReader(reader),
		posrow: 0,
		poscol: -1,
		pos:    -1,
	}
}

func (l *Lexer) SafeLex() (Token, int, error) {
	r, _, err := l.reader.ReadRune()
	l.poscol++
	l.pos++

	if err != nil && err == io.EOF {
		return Token{EOF, EOF.Literal()}, l.pos, nil
	}

	if err != nil {
		return Token{}, l.pos, errors.Join(ErrRuneRead, err)
	}

	switch r {
	case '_':
		return l.lexUnderscore()
	case '*':
		return l.lexAsterisk()
	case '#':
		return l.lexHash()
	case '\n':
		l.poscol = -1
		l.posrow++
		return Token{NEWLINE, NEWLINE.Literal()}, l.pos, nil
	default:
		return l.lexString()
	}
}

func (l *Lexer) lexUnderscore() (Token, int, error) {
	return l.lexRange('_', UNDERSCORE_1, UNDERSCORE_2)
}

func (l *Lexer) lexAsterisk() (Token, int, error) {
	return l.lexRange('*', ASTERISK_1, ASTERISK_2)
}

func (l *Lexer) lexHash() (Token, int, error) {
	return l.lexRange('#', HASH_1, HASH_2, HASH_3, HASH_4, HASH_5)
}

// func (l *Lexer) lexIndentation() (Token, int, error) {

// }

func (l *Lexer) lexString() (Token, int, error) {
	l.unread()
	var literal string
	span := -1

	for {
		r, _, err := l.reader.ReadRune()
		l.poscol++
		l.pos++
		span++

		if err != nil && err == io.EOF {
			l.unread()
			return Token{STRING, literal}, l.pos - span + 1, nil
		}

		if err != nil {
			return Token{}, l.pos, errors.Join(ErrRuneRead, err)
		}

		if isSyntax(r) {
			l.unread()
			return Token{STRING, literal}, l.pos - span + 1, nil
		}

		literal += string(r)
	}
}

func (l *Lexer) lexRange(target rune, kinds ...TokenKind) (Token, int, error) {
	for i := 1; i <= len(kinds); i++ {
		r, _, err := l.reader.ReadRune()
		l.poscol++
		l.pos++

		if err != nil && err == io.EOF {
			l.unread()
			return Token{kinds[i-1], kinds[i-1].Literal()}, l.pos - i + 1, nil
		}

		if err != nil {
			return Token{}, l.pos, errors.Join(ErrRuneRead, err)
		}

		if r != target {
			l.unread()
			return Token{kinds[i-1], kinds[i-1].Literal()}, l.pos - i + 1, nil
		}
	}

	l.unread()
	tk := kinds[len(kinds)-1]
	return Token{tk, tk.Literal()}, l.pos - len(kinds) + 1, nil
}

func (l *Lexer) unread() {
	err := l.safeUnread()

	if err != nil {
		panic(err)
	}
}

func (l *Lexer) safeUnread() error {
	err := l.reader.UnreadRune()

	if err != nil {
		return ErrRuneUnread
	}

	l.poscol--
	l.pos--

	return nil
}
