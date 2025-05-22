package lexer

import (
	"bufio"
	"io"
)

func peekOneRune(reader bufio.Reader) (rune, int, error) {
	r, _, err := reader.ReadRune()

	if err == io.EOF {
		return -1, 0, ErrEOF
	}

	if err != nil {
		return -1, 0, ErrRead
	}

	_ = reader.UnreadRune()
	return r, 1, nil
}

func peekRunes(reader bufio.Reader, count int) ([]rune, int, error) {
	var runes []rune

	for _ = range count {
		r, _, err := reader.ReadRune()

		if err == io.EOF {
			return runes, len(runes), ErrEOF
		}

		if err != nil {
			return runes, 0, ErrRead
		}

		_ = reader.UnreadRune()
		runes = append(runes, r)
	}

	return runes, len(runes), nil
}

func lexWhen(reader bufio.Reader, matcher func(rune) bool, kind TokenKind) (Token, int, error) {
	token := Token{kind, ""}
	i := -1

	for {
		r, _, err := reader.ReadRune()
		i++

		if err == io.EOF {
			return token, i, nil
		}

		if err != nil {
			return token, i, ErrRead
		}

		if !matcher(r) {
			_ = reader.UnreadRune()
			return token, i, nil
		}

		token.lit += string(r)
	}
}

func lexUntil(reader bufio.Reader, matcher func(rune) bool, kind TokenKind) (Token, int, error) {
	token := Token{kind, ""}
	i := -1

	for {
		r, _, err := reader.ReadRune()
		i++

		if err == io.EOF {
			return token, i, nil
		}

		if err != nil {
			return token, i, ErrRead
		}

		if matcher(r) {
			_ = reader.UnreadRune()
			return token, i, nil
		}

		token.lit += string(r)
	}
}
