package codegen

import (
	"CCIG/ast"
	"fmt"
)

func GenerateAsm(ast []ast.Node) {
	fmt.Println("section .text")
	fmt.Println("global _main")
	fmt.Println("_main:")

	fmt.Println("  push rbp");
	fmt.Println("  mov rbp, rsp");
	fmt.Println("  sub rsp, 208");

	for _, node := range ast {
		generate(node)

		fmt.Println("  pop rax")
	}

	fmt.Println("  mov rsp, rbp")
	fmt.Println("  pop rbp")
	fmt.Println("  ret")
}

func generate(node ast.Node) {
	switch node.(type) {
	case *ast.VarStatement:
		generateVarAssign(node.(*ast.VarStatement))
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

func generateVarAssign(stmt *ast.VarStatement) {
	//generateAddr(stmt)

	fmt.Println("  mov rax, rbp")
	fmt.Printf("  sub rax, %d\n", getOffset(rune(stmt.Name[0])))
	fmt.Println("  push rax")

	fmt.Println("  pop rax")
	fmt.Println("  mov rax, [rax]")
	fmt.Println("  push rax")

	fmt.Println("")

	generate(stmt.Initializer)

	fmt.Println("  pop rdi")
	fmt.Println("  pop rax")
	fmt.Println("  push rax") // TODO?

	fmt.Println("")
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

func generateAddr(stmt *ast.VarStatement) {
	r := []rune(stmt.Name)
	offset := getOffset(r[0])
	fmt.Printf("  lea rbp[%d], rax\n", -offset)
}

func getOffset(r rune) int {
	return int((r - 'a' + 1) * 8)
}
