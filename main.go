package main

import (
	"CCIG/codegen"
	"CCIG/parser"
	"CCIG/tokenizer"
)

func main() {
	source := `var a = 10;`
	//source := `5 + 10;`
	//source := `
	//var x = 10;
	//var y = x;`
	run(source)
}

func run(source string) {
	t := tokenizer.NewTokenizer(source)
	tokens := t.Tokenize()

	p := parser.NewParser(tokens)
	ast := p.Parse()

	codegen.GenerateAsm(ast)
}
