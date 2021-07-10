package codegen

import (
	"CCIG/ast"
	"CCIG/parser"
	"fmt"
)

type Compiler struct {
    asm string
}

func GenerateAsm(prog parser.Function) string {
    c := Compiler { asm: "" }

	c.emitInstruction("section .text")
	c.emitInstruction("global _main")
    c.emitInstruction("_main:")

	c.emitInstruction("  ; prologue")
	c.emitInstruction("  push rbp")
	c.emitInstruction("  mov rbp, rsp")
	c.emitInstruction(fmt.Sprintf("  sub rsp, %d", prog.StackSize)) // TODO
	c.emitInstruction("")

	for _, node := range prog.Body {
		c.generate(node)

        c.emitInstruction("  pop rax")
	}

	c.emitInstruction("  mov rsp, rbp")
	c.emitInstruction("  pop rbp")
	c.emitInstruction("  ret")
    return c.asm
}

func (c *Compiler) emitInstruction(instruction string) {
	c.asm += instruction + "\n"
}

func (c *Compiler) generate(node ast.Node) {
	switch node.(type) {
	case *ast.VarStatement:
		c.generateVarAssign(node.(*ast.VarStatement))
	case *ast.ExprStatement:
		c.generateExprStatement(node.(*ast.ExprStatement))
	case *ast.NumberExpr:
		c.generateNumber(node.(*ast.NumberExpr))
	case *ast.PrefixExpr:
		c.generatePrefixExpr(node.(*ast.PrefixExpr))
	case *ast.InfixExpr:
		c.generateInfixExpr(node.(*ast.InfixExpr))
	case *ast.IdentifierExpr:
		c.generateIdentifierExpr(node.(*ast.IdentifierExpr))
	default:
		panic("TODO")
	}
	c.emitInstruction("  ;-----")
}

func (c *Compiler) generateVarAssign(stmt *ast.VarStatement) {
	c.generateOffset(stmt)

	c.generate(stmt.Initializer)

	c.emitInstruction("  pop rdi")
	c.emitInstruction("  pop rax")
	c.emitInstruction("  mov [rax], rdi")
	c.emitInstruction("  push rdi")
}

func (c *Compiler) generateExprStatement(statement *ast.ExprStatement) {
	c.generate(statement.Value)
}

func (c *Compiler) generateNumber(number *ast.NumberExpr) {
	c.emitInstruction(fmt.Sprintf("  push %d\n", int(number.Value))) // TODO Float
}

func (c *Compiler) generatePrefixExpr(expr *ast.PrefixExpr) {
    c.emitInstruction("  push 0")
	c.generate(expr.Right)

    c.emitInstruction("  pop rdi")
    c.emitInstruction("  pop rax")

	c.generateBinaryOperator(expr.Operator)

    c.emitInstruction("  push rax")
}

func (c *Compiler) generateInfixExpr(expr *ast.InfixExpr) {
	c.generate(expr.Left)
	c.generate(expr.Right)

    c.emitInstruction("  pop rdi")
    c.emitInstruction("  pop rax")

	c.generateBinaryOperator(expr.Operator)

    c.emitInstruction("  push rax")
}

func (c *Compiler) generateIdentifierExpr(expr *ast.IdentifierExpr) {
	c.generateOffset(expr)
	c.emitInstruction("  pop rax")
	c.emitInstruction("  mov rax, [rax]")
	c.emitInstruction("  push rax")
}

func (c *Compiler) generateBinaryOperator(operator ast.BinaryOperator) {
	switch operator {
	case ast.Add:
        c.emitInstruction("  add rax, rdi")
	case ast.Subtract:
        c.emitInstruction("  sub rax, rdi")
	case ast.Multiply:
        c.emitInstruction("  imul rax, rdi")
	case ast.Divide:
        c.emitInstruction("  cqo")
        c.emitInstruction("  idiv rdi")
        c.emitInstruction("  mov rax, rdx") // TODO Works?
	default:
		panic("TODO")
	}
}

func (c *Compiler) generateOffset(node ast.Node) {
	// TODO FIX
	if _, ok := node.(*ast.VarStatement); ok {
		offset := node.(*ast.VarStatement).Offset

		c.emitInstruction("  mov rax, rbp")
		c.emitInstruction(fmt.Sprintf("  sub rax, %d", offset))
		c.emitInstruction("  push rax")
		return
	}
	if _, ok := node.(*ast.IdentifierExpr); ok {
		offset := node.(*ast.IdentifierExpr).Offset

		c.emitInstruction("  mov rax, rbp")
		c.emitInstruction(fmt.Sprintf("  sub rax, %d", offset))
		c.emitInstruction("  push rax")
		return
	}
	panic("TODO")
}
