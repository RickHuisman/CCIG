package codegen

import (
	"CCIG/ast"
	"CCIG/parser"
	"fmt"
)

type Compiler struct {
	ifElseCount int
	asm         string
}

func GenerateAsm(prog parser.Function) string {
	c := Compiler{
		ifElseCount: 0,
		asm:         "",
	}

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

	c.emitInstruction(".L.return:")
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
	case *ast.ReturnStatement:
		c.generateReturnStatement(node.(*ast.ReturnStatement))
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
	case *ast.BlockStatement:
		c.generateBlockStatement(node.(*ast.BlockStatement))
	case *ast.IfElseStatement:
		c.generateIfElseStatement(node.(*ast.IfElseStatement))
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

func (c *Compiler) generateReturnStatement(stmt *ast.ReturnStatement) {
	c.generate(stmt.Value) // TODO Check for null
	c.emitInstruction("  pop rax")
	c.emitInstruction("  jmp .L.return")
}

func (c *Compiler) generateExprStatement(statement *ast.ExprStatement) {
	c.generate(statement.Value)
}

func (c *Compiler) generateNumber(number *ast.NumberExpr) {
	c.emitInstruction(fmt.Sprintf("  push %d\n", int(number.Value))) // TODO Float
}

func (c *Compiler) generatePrefixExpr(expr *ast.PrefixExpr) {
	c.generate(expr.Right)

	c.emitInstruction("  pop rax")

	c.generateUnaryOperator(expr.Operator)

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

func (c *Compiler) generateBlockStatement(expr *ast.BlockStatement) {
	for _, e := range expr.Statements {
		c.generate(e)
	}
}

func (c *Compiler) generateIdentifierExpr(expr *ast.IdentifierExpr) {
	c.generateOffset(expr)
	c.emitInstruction("  pop rax")
	c.emitInstruction("  mov rax, [rax]")
	c.emitInstruction("  push rax")
}

func (c *Compiler) generateIfElseStatement(stmt *ast.IfElseStatement) {
	c.generate(stmt.Condition)
	c.emitInstruction("  cmp rax, 0")

	c.ifElseCount += 1

	c.emitInstruction(fmt.Sprintf("  je .L.else.%d", c.ifElseCount))
	c.generate(stmt.Then)
	c.emitInstruction(fmt.Sprintf("  jmp .L.end.%d", c.ifElseCount))
	c.emitInstruction(fmt.Sprintf(".L.else.%d:", c.ifElseCount))
	if stmt.Else != nil {
		c.generate(stmt.Else)
	}
	c.emitInstruction(fmt.Sprintf(".L.end.%d:", c.ifElseCount))
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
	case ast.EqualEqual:
		c.emitInstruction("  cmp rax, rdi")
		c.emitInstruction("  sete al")
	case ast.BangEqual:
		c.emitInstruction("  cmp rax, rdi")
		c.emitInstruction("  setne al")
	case ast.LessThanEqual:
		c.emitInstruction("  cmp rax, rdi")
		c.emitInstruction("  setle al")
	case ast.GreaterThanEqual:
		c.emitInstruction("  cmp rax, rdi")
		c.emitInstruction("  setge al")
	case ast.Greater:
		c.emitInstruction("  cmp rax, rdi")
		c.emitInstruction("  setg al")
	case ast.Less:
		c.emitInstruction("  cmp rax, rdi")
		c.emitInstruction("  setl al")
	default:
		panic("TODO")
	}
}

func (c *Compiler) generateUnaryOperator(operator ast.UnaryOperator) {
	switch operator {
	case ast.Negate:
		c.emitInstruction("  neg rax")
	case ast.Not:
		panic("TODO")
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