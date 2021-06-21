package main

import (
	"CCIG/codegen"
	"CCIG/parser"
	"CCIG/tokenizer"
)

func main() {
	source := `5 + 10;`
	run(source)
}

func run(source string) {
	t := tokenizer.NewTokenizer(source)
	tokens := t.Tokenize()

	p := parser.NewParser(tokens)
	ast := p.Parse()

	codegen.GenerateAsm(ast)
}
