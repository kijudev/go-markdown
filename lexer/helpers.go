package lexer

import (
	"bufio"
	"io"
)

func peekOneRune(reader *bufio.Reader) (rune, int, error) {
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

func peekRunes(reader *bufio.Reader, count int) ([]rune, int, error) {
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

func lexUntil(reader *bufio.Reader, matcher func(rune) bool, limit int) (string, int, error) {
	var lit string
	span := 0

	for {
		if limit > 0 && span >= limit {
			return lit, span, nil
		}

		r, _, err := reader.ReadRune()

		if err == io.EOF {
			return lit, span, nil
		}

		if err != nil {
			return lit, span, ErrRead
		}

		if !matcher(r) {
			_ = reader.UnreadRune()
			return lit, span, nil
		}

		lit += string(r)
		span++
	}
}

func lexUpTo(reader *bufio.Reader, matcher func(rune) bool, limit int) (string, int, error) {
	var lit string
	span := 0

	for {
		if limit > 0 && span >= limit {
			return lit, span, nil
		}

		r, _, err := reader.ReadRune()

		if err == io.EOF {
			return lit, span, nil
		}

		if err != nil {
			return lit, span, ErrRead
		}

		if matcher(r) {
			_ = reader.UnreadRune()
			return lit, span, nil
		}

		lit += string(r)
		span++
	}
}
