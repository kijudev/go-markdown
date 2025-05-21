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

	lexer := lexer.NewLexer(file)
	tokens, err := lexer.Tokenize()

	if err != nil {
		panic(err)
	}

	for _, t := range tokens {
		println(t.Literal)
	}
}
