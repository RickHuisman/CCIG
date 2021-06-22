package codegen

import (
	"CCIG/ast"
	"fmt"
)

func GenerateAsm(ast []ast.Node) {
	fmt.Println("section .text")
	fmt.Println("global _main")
	fmt.Println("_main")

	for _, node := range ast {
		generate(node)
	}

	fmt.Println("  pop rax")
	fmt.Println("  ret")
}

func generate(node ast.Node) {
	switch node.(type) {
	case *ast.NumberExpr:
		generateNumber(node.(*ast.NumberExpr))
	case *ast.PrefixExpr:
		generatePrefixExpr(node.(*ast.PrefixExpr))
	case *ast.InfixExpr:
		generateInfixExpr(node.(*ast.InfixExpr))
	case *ast.ExprStatement:
		generateExprStatement(node.(*ast.ExprStatement))
	default:
		panic("TODO")
	}
}

func generateNumber(number *ast.NumberExpr) {
	fmt.Printf("  push %d\n", int(number.Value)) // TODO Float
}

func generatePrefixExpr(expr *ast.PrefixExpr) {
	fmt.Println("  push 0")
	generate(expr.Right)

	fmt.Println("  pop rdi")
	fmt.Println("  pop rax")

	generateBinaryOperator(expr.Operator)

	fmt.Println("  push rax")
}

func generateInfixExpr(expr *ast.InfixExpr) {
	generate(expr.Right)
	generate(expr.Left)
	fmt.Println("  pop rdi")
	fmt.Println("  pop rax")

	generateBinaryOperator(expr.Operator)

	fmt.Println("  push rax")
}

func generateExprStatement(statement *ast.ExprStatement) {
	generate(statement.Value)
}

func generateBinaryOperator(operator ast.BinaryOperator) {
	switch operator {
	case ast.Add:
		fmt.Println("  add rax, rdi")
	case ast.Subtract:
		fmt.Println("  sub rax, rdi")
	case ast.Multiply:
		fmt.Println("  imul rax, rdi")
	case ast.Divide:
		fmt.Println("  cqo")
		fmt.Println("  idiv rdi")
	default:
		panic("TODO")
	}
}