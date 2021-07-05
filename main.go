package main

import (
	"CCIG/codegen"
	"CCIG/parser"
	"CCIG/tokenizer"
	"os"
	"os/exec"
)

func main() {
	source := `var b = 5; var c = b + 10; c + b;`
	run(source)
}

func run(source string) {
	t := tokenizer.NewTokenizer(source)
	tokens := t.Tokenize()

	p := parser.NewParser(tokens)
	ast := p.Parse()

    asm := codegen.GenerateAsm(ast)
    writeAsm(asm)
    buildAndLink()
}

func writeAsm(asm string) {
    f, err := os.Create("temp.asm")
    if err != nil {
        panic(err)
    }
	_, err = f.WriteString(asm)
	if err != nil {
		panic(err)
	}
}

func buildAndLink() {
	cmd := exec.Command("nasm", "-f", "macho64", "temp.asm")
	err := cmd.Run()

	if err != nil {
		panic(err)
	}

	cmd2 := exec.Command("gcc", "-arch", "x86_64", "-o", "temp", "temp.o")
	err2 := cmd2.Run()

	if err2 != nil {
		panic(err2)
	}
}
