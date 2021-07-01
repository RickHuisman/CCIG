package main

import (
    "os"
	"CCIG/codegen"
	"CCIG/parser"
	"CCIG/tokenizer"
)

func main() {
    source := `10 * 2;`
	//source := `var a = 10;`
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

    asm := codegen.GenerateAsm(ast)
    writeAsm(asm)
}

func writeAsm(asm string) {
    f, err := os.Create("temp.asm")
    if err != nil {
        panic(err)
    }
    f.WriteString(asm)
}
