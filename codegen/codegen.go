package codegen

import (
	"CCIG/ast"
	"fmt"
)

func GenerateAsm(ast []ast.Node) {
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

func generateInfixExpr(expr *ast.InfixExpr) {
	generate(expr.Right)
	generate(expr.Left)
	fmt.Println("  pop rdi")
	fmt.Println("  pop rax")

	fmt.Println("  add rax, rdi")
	fmt.Println("  push rax")
}

func generateExprStatement(statement *ast.ExprStatement) {
	generate(statement.Value)
}
