package main

import (
	"os"

	"github.com/kijudev/go-markdown/lexer"
)

func main() {
	file, err := os.Open("./README.md")

	if err != nil {
		panic(err)
	}

	lx := lexer.NewLexer(file)

	tokens, err := lx.Tokenize()
	if err != nil {
		panic(err)
	}

	for _, ti := range tokens {
		t := ti.Token
		pos := ti.Pos

		println(t.Kind.DebugName(), pos.X, pos.Y, t.Lit)
	}
}
