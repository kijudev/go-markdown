package lexer

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"unicode"
)

type Lexer struct {
	readerSource io.Reader
	reader       *bufio.Reader
	posrow       int
	poscol       int
	pos          int
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		readerSource: reader,
		reader:       bufio.NewReader(reader),
		posrow:       0,
		poscol:       -1,
		pos:          -1,
	}
}

func (l *Lexer) Tokenize() ([]Token, error) {
	l.Reset()
	var tokens []Token

	for {
		t, _, err := l.Lex()
		i := len(tokens)

		if err != nil {
			return []Token{}, err
		}

		if t.Kind == EOF {
			break
		}

		if t.Kind == STRING {
			t.Literal = strings.Trim(t.Literal, " ")
		}

		if i == 0 {
			tokens = append(tokens, t)
			continue
		}

		tokens = append(tokens, t)
	}

	return tokens, nil
}

func (l *Lexer) Reset() {
	l.reader.Reset(l.readerSource)
	l.posrow = 0
	l.poscol = -1
	l.pos = -1
}

func (l *Lexer) Lex() (Token, int, error) {
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
	case '\\':
		return l.lexEscape()
	case '_':
		return l.lexUnderscore()
	case '*':
		return l.lexAsterisk()
	case '#':
		return l.lexHash()
	case '\n':
		return l.lexNewline()
	case '-':
		return l.lexDash()
	case '=':
		return l.lexEq()
	default:
		if l.poscol == 0 && unicode.IsSpace(r) {
			return l.lexIndentation()
		}

		if unicode.IsDigit(r) {
			return l.lexNumbering()
		}

		return l.lexString()
	}
}

func (l *Lexer) lexEscape() (Token, int, error) {
	return Token{ESCAPE, ESCAPE.Literal()}, l.pos, nil
}

func (l *Lexer) lexNumbering() (Token, int, error) {
	t, p, err := l.lexConsecutive(unicode.IsDigit, STRING)

	if err != nil {
		return Token{}, p, err
	}

	r, _, err := l.reader.ReadRune()
	l.poscol++
	l.pos++

	if r == '.' {
		return Token{NUMBERING, t.Literal}, p, nil
	}

	l.mustUnread()
	return t, p, nil
}

func (l *Lexer) lexDash() (Token, int, error) {
	t, p, err := l.lexConsecutive(func(r rune) bool { return r == '-' }, DASH_MULTIPLE)

	if err != nil {
		return Token{}, p, err
	}

	if len(t.Literal) == 1 {
		return Token{DASH_SINGLE, DASH_SINGLE.Literal()}, p, nil
	}

	return t, p, nil
}

func (l *Lexer) lexEq() (Token, int, error) {
	return l.lexConsecutive(func(r rune) bool { return r == '=' }, EQ)
}

func (l *Lexer) lexNewline() (Token, int, error) {

	l.poscol = -1
	l.posrow++
	return Token{NEWLINE, NEWLINE.Literal()}, l.pos, nil
}

func (l *Lexer) lexUnderscore() (Token, int, error) {
	return l.lexRange('_', UNDERSCORE_1, UNDERSCORE_2, UNDERSCORE_3)
}

func (l *Lexer) lexAsterisk() (Token, int, error) {
	return l.lexRange('*', ASTERISK_1, ASTERISK_2, ASTERISK_3)
}

func (l *Lexer) lexHash() (Token, int, error) {
	return l.lexRange('#', HASH_1, HASH_2, HASH_3, HASH_4, HASH_5, HASH_6)
}

func (l *Lexer) lexIndentation() (Token, int, error) {
	return l.lexConsecutive(unicode.IsSpace, INDENTATION)
}

func (l *Lexer) lexString() (Token, int, error) {
	l.mustUnread()
	var literal string
	span := -1

	for {
		r, _, err := l.reader.ReadRune()
		l.poscol++
		l.pos++
		span++

		if err != nil && err == io.EOF {
			l.mustUnread()
			return Token{STRING, literal}, l.pos - span + 1, nil
		}

		if err != nil {
			return Token{}, l.pos, errors.Join(ErrRuneRead, err)
		}

		if isSyntax(r) {
			l.mustUnread()
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
			l.mustUnread()
			return Token{kinds[i-1], kinds[i-1].Literal()}, l.pos - i + 1, nil
		}

		if err != nil {
			return Token{}, l.pos, errors.Join(ErrRuneRead, err)
		}

		if r != target {
			l.mustUnread()
			return Token{kinds[i-1], kinds[i-1].Literal()}, l.pos - i + 1, nil
		}
	}

	l.mustUnread()
	tk := kinds[len(kinds)-1]
	return Token{tk, tk.Literal()}, l.pos - len(kinds) + 1, nil
}

func (l *Lexer) lexConsecutive(checkFn func(rune) bool, kind TokenKind) (Token, int, error) {
	l.mustUnread()
	var literal string
	span := -1

	for {
		r, _, err := l.reader.ReadRune()
		l.poscol++
		l.pos++

		if err != nil && err == io.EOF {
			l.mustUnread()
			return Token{kind, literal}, l.pos - span + 1, nil
		}

		if err != nil {
			return Token{}, l.pos, errors.Join(ErrRuneRead, err)
		}

		if !checkFn(r) {
			l.mustUnread()
			return Token{kind, literal}, l.pos - span + 1, nil
		}

		literal += string(r)
	}
}

func (l *Lexer) mustUnread() {
	err := l.unread()

	if err != nil {
		panic(err)
	}
}

func (l *Lexer) unread() error {
	err := l.reader.UnreadRune()

	if err != nil {
		return ErrRuneUnread
	}

	l.poscol--
	l.pos--

	return nil
}
