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
		pos:    TokenPos{0, 0, 0},
	}
}

func (l *Lexer) TokenizeDebug() ([]TokenInfoDebug, error) {
	tokens, err := l.Tokenize()
	if err != nil {
		return []TokenInfoDebug{}, err
	}

	var infos []TokenInfoDebug
	for _, t := range tokens {
		infos = append(infos, NewTokenInfoDebug(t))
	}

	return infos, nil
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
	l.pos = TokenPos{0, 0, 0}
}

func (l *Lexer) LexDebug(TokenInfoDebug, error) (TokenInfoDebug, error) {
	info, err := l.Lex()
	return NewTokenInfoDebug(info), err
}

func (l *Lexer) Lex() (TokenInfo, error) {
	r, _, err := peekOneRune(l.reader)

	if err != nil {
		return TokenInfo{}, err
	}

	switch r {
	case RuneKindHash.Rune():
		return l.lexHash()
	case RuneKindDash.Rune():
		return l.lexDash()
	case RuneKindAsterisk.Rune():
		return l.lexAsterisk()
	case RuneKindUnderscore.Rune():
		return l.lexUnderscore()
	case RuneKindEscape.Rune():
		return l.lexEscape()
	default:
		if runeIs(RuneKindNewline)(r) {
			return l.lexNewline()
		}

		if runeIs(RuneKindWhitespace)(r) {
			return l.lexWhitespace()
		}

		if runeIs(RuneKindNumber)(r) {
			return l.lexNumber()
		}

		return l.lexString()
	}
}

func (l *Lexer) lexHash() (TokenInfo, error) {
	lit, span, err := lexUntil(l.reader, runeIs(RuneKindHash), 6)
	pos := l.pos

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

	l.pos.X += span
	l.pos.Abs += span

	return TokenInfo{
		Token: Token{kind, lit},
		Pos:   pos,
		Span:  span,
		Raw:   lit,
	}, nil
}

func (l *Lexer) lexDash() (TokenInfo, error) {
	lit, span, err := lexUntil(l.reader, runeIs(RuneKindDash), -1)
	pos := l.pos

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

	l.pos.X += span
	l.pos.Abs += span

	return TokenInfo{
		Token: Token{kind, lit},
		Pos:   pos,
		Span:  span,
		Raw:   lit,
	}, nil
}

func (l *Lexer) lexAsterisk() (TokenInfo, error) {
	lit, span, err := lexUntil(l.reader, runeIs(RuneKindAsterisk), 3)
	pos := l.pos

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

	l.pos.X += span
	l.pos.Abs += span

	return TokenInfo{
		Token: Token{kind, lit},
		Pos:   pos,
		Span:  span,
		Raw:   lit,
	}, nil
}

func (l *Lexer) lexUnderscore() (TokenInfo, error) {
	lit, span, err := lexUntil(l.reader, runeIs(RuneKindUnderscore), 3)
	pos := l.pos

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

	l.pos.X += span
	l.pos.Abs += span

	return TokenInfo{
		Token: Token{kind, lit},
		Pos:   pos,
		Span:  span,
		Raw:   lit,
	}, nil
}

func (l *Lexer) lexEscape() (TokenInfo, error) {
	var lit string
	_, _, err := l.reader.ReadRune()
	pos := l.pos

	if err == io.EOF {
		return TokenInfo{}, ErrEOF
	}

	if err != nil {
		return TokenInfo{}, ErrUnread
	}

	r, _, err := l.reader.ReadRune()

	if err == io.EOF {
		l.pos.X++
		l.pos.Y++

		return TokenInfo{
			Token: Token{TokenKindString, "\\"},
			Pos:   pos,
			Span:  1,
			Raw:   "\\",
		}, nil
	}

	if err != nil {
		return TokenInfo{}, ErrUnread
	}

	lit = string(r)

	if runeIsOneOf(
		RuneKindHash,
		RuneKindDash,
		RuneKindAsterisk,
		RuneKindUnderscore,
	)(r) {
		pos.X++
		pos.Abs++
		l.pos.X += 2
		l.pos.Abs += 2

		return TokenInfo{
			Token: Token{TokenKindString, lit},
			Pos:   pos,
			Span:  1,
			Raw:   "\\" + lit,
		}, nil
	}

	l.pos.X += 2
	l.pos.Abs += 2

	return TokenInfo{
		Token: Token{TokenKindString, "\\" + lit},
		Pos:   pos,
		Span:  2,
		Raw:   "\\" + lit,
	}, nil
}

func (l *Lexer) lexNewline() (TokenInfo, error) {
	_, _, err := l.reader.ReadRune()
	pos := l.pos

	if err == io.EOF {
		return TokenInfo{}, ErrEOF
	}

	if err != nil {
		return TokenInfo{}, ErrRead
	}

	l.pos.X = 0
	l.pos.Y++
	l.pos.Abs++

	return TokenInfo{
		Token: Token{TokenKindNewline, "\n"},
		Pos:   pos,
		Span:  1,
		Raw:   "\n",
	}, nil
}

func (l *Lexer) lexWhitespace() (TokenInfo, error) {
	lit, span, err := lexUntil(l.reader, unicode.IsSpace, -1)
	pos := l.pos

	if err != nil {
		return TokenInfo{}, err
	}

	l.pos.X += span
	l.pos.Abs += span

	return TokenInfo{
		Token: Token{TokenKindWhitespace, lit},
		Pos:   pos,
		Span:  span,
		Raw:   lit,
	}, nil
}

func (l *Lexer) lexNumber() (TokenInfo, error) {
	lit, span, err := lexUntil(l.reader, unicode.IsDigit, -1)
	pos := l.pos

	if err != nil {
		return TokenInfo{}, err
	}

	r, _, err := peekOneRune(l.reader)

	if err == ErrEOF {
		return TokenInfo{
			Token: Token{TokenKindNumber, lit},
			Pos:   pos,
			Span:  span,
			Raw:   lit,
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

		l.pos.X += span + 1
		l.pos.Abs += span + 1

		return TokenInfo{
			Token: Token{TokenKindNumbering, lit},
			Pos:   pos,
			Span:  span,
			Raw:   lit + ".",
		}, nil
	}

	l.pos.X += span
	l.pos.Abs += span

	return TokenInfo{
		Token: Token{TokenKindNumber, lit},
		Pos:   pos,
		Span:  span,
		Raw:   lit,
	}, nil
}

func (l *Lexer) lexString() (TokenInfo, error) {
	lit, span, err := lexUntil(l.reader, runeIsNotOneOf(
		RuneKindNewline,
		RuneKindNumber,
		RuneKindHash,
		RuneKindDash,
		RuneKindAsterisk,
		RuneKindUnderscore,
		RuneKindEscape,
	), -1)
	pos := l.pos

	if err != nil {
		return TokenInfo{}, err
	}

	l.pos.X += span
	l.pos.Abs += span

	return TokenInfo{
		Token: Token{TokenKindString, lit},
		Pos:   pos,
		Span:  span,
		Raw:   lit,
	}, nil
}
