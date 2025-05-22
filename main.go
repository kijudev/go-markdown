package main

import (
	"strings"

	"github.com/kijudev/go-markdown/lexer"
)

func main() {
	lx := lexer.NewLexer(strings.NewReader("# === a"))

	tokens, err := lx.Tokenize()
	if err != nil {
		panic(err)
	}

	for _, ti := range tokens {
		t := ti.Token
		pos := ti.Pos

		println(pos.Abs, t.Lit)
	}
}
