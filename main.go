package main

import (
	"encoding/json"
	"os"

	"github.com/kijudev/go-markdown/lexer"
)

func main() {
	readmeFile, err := os.Open("./README.md")
	if err != nil {
		panic(err)
	}
	defer readmeFile.Close()

	lx := lexer.NewLexer(readmeFile)

	tokens, err := lx.Tokenize()
	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(tokens)
	if err != nil {
		panic(err)
	}

	jsonFile, err := os.Create("./readme.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(b)
}
