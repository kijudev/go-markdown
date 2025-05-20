package main

import (
	"os"

	"github.com/kijudev/go-markdown/lexer"
)

func main() {
	file, _ := os.Open("./README.md")
	defer file.Close()

	l := lexer.NewLexer(file)

	for {
		t, p, err := l.SafeLex()

		if err != nil {
			panic(err)
		}

		println(p, " -> ", t.Literal)

		if t.Kind == lexer.EOF {
			break
		}
	}
}
