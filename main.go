package main

import (
	"bufio"
	"strings"
)

func main() {
	rs := strings.NewReader("# lubię placki 123")
	r := bufio.NewReader(rs)

	for {
		rune, size, err := r.ReadRune()

		if err != nil {
			panic(err)
		}

		println(size, string(rune))
	}
}
