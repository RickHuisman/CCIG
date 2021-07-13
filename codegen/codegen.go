package codegen

import (
	"CCIG/ast"
	"CCIG/parser"
	"fmt"
)

type Compiler struct {
	ifElseCount int
	currentFn   parser.Function
	asm         map[string][]string
}

func GenerateAsm(fn parser.Function) string {
	c := Compiler{
		ifElseCount: 0,
		asm:         map[string][]string{},
	}

	c.instruction("section .text")
	c.instruction(fmt.Sprintf("global %s", fn.Name))
	c.currentFn = fn

	c.labelInstruction(fn.Name, "  ; prologue")
	c.labelInstruction(fn.Name, "  push rbp")
	c.labelInstruction(fn.Name, "  mov rbp, rsp")
	c.labelInstruction(fn.Name, fmt.Sprintf("  sub rsp, %d", fn.StackSize))

	for _, node := range fn.Body.Statements {
		c.generate(node)

		//c.labelInstruction(c.currentFn.Name, "  pop rax")
	}

	c.labelInstruction(fn.Name, fmt.Sprintf(".L.return.%s:", fn.Name))
	c.labelInstruction(fn.Name, "  mov rsp, rbp")
	c.labelInstruction(fn.Name, "  pop rbp")
	c.labelInstruction(fn.Name, "  ret")
	return c.buildAsm()
}

func (c *Compiler) generate(node ast.Node) {
	switch node.(type) {
	case *ast.VarStatement:
		c.generateVarAssign(node.(*ast.VarStatement))
	case *ast.ReturnStatement:
		c.generateReturnStatement(node.(*ast.ReturnStatement))
	case *ast.IfElseStatement:
		c.generateIfElseStatement(node.(*ast.IfElseStatement))
	case *ast.BlockStatement:
		c.generateBlockStatement(node.(*ast.BlockStatement))
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
	case *ast.CallExpr:
		c.generateCallExpr(node.(*ast.CallExpr))
	case *ast.FunStatement:
		c.generateFunStatement(node.(*ast.FunStatement))
	default:
		panic("TODO")
	}
	c.labelInstruction(c.currentFn.Name, "  ;----- " + c.currentFn.Name)
}

func (c *Compiler) generateVarAssign(stmt *ast.VarStatement) {
	c.generateOffset(stmt)

	c.generate(stmt.Initializer)

	c.labelInstruction(c.currentFn.Name, "  pop rdi")
	c.labelInstruction(c.currentFn.Name, "  pop rax")
	c.labelInstruction(c.currentFn.Name, "  mov [rax], rdi")
	c.labelInstruction(c.currentFn.Name, "  push rdi")
}

func (c *Compiler) generateReturnStatement(stmt *ast.ReturnStatement) {
	c.generate(stmt.Value) // TODO Return null ???
	c.labelInstruction(c.currentFn.Name, "  pop rax")
	c.labelInstruction(
		c.currentFn.Name,
		fmt.Sprintf("  jmp .L.return.%s", c.currentFn.Name),
	)
}

func (c *Compiler) generateExprStatement(statement *ast.ExprStatement) {
	c.generate(statement.Value)
}

func (c *Compiler) generateNumber(number *ast.NumberExpr) {
	c.labelInstruction(
		c.currentFn.Name,
		fmt.Sprintf("  push %d\n", int(number.Value)),
	) // TODO Float
}

func (c *Compiler) generatePrefixExpr(expr *ast.PrefixExpr) {
	c.generate(expr.Right)
	c.labelInstruction(c.currentFn.Name, "  pop rax")
	c.generateUnaryOperator(expr.Operator)
	c.labelInstruction(c.currentFn.Name, "  push rax")
}

func (c *Compiler) generateInfixExpr(expr *ast.InfixExpr) {
	c.generate(expr.Left)
	c.generate(expr.Right)

	c.labelInstruction(c.currentFn.Name, "  pop rdi")
	c.labelInstruction(c.currentFn.Name, "  pop rax")

	c.generateBinaryOperator(expr.Operator)

	c.labelInstruction(c.currentFn.Name, "  push rax")
}

func (c *Compiler) generateBlockStatement(expr *ast.BlockStatement) {
	for _, e := range expr.Statements {
		c.generate(e)
	}
}

func (c *Compiler) generateFunStatement(stmt *ast.FunStatement) {
	oldFn := c.currentFn
	c.currentFn = parser.Function{
		Name:      stmt.Name,
		Body:      stmt.Body,
		StackSize: 0, // TODO
	}
	c.generate(stmt.Body)

	// Function epilogue
	c.labelInstruction(c.currentFn.Name, fmt.Sprintf(".L.return.%s:", c.currentFn.Name))

	c.labelInstruction(stmt.Name, "  ret")

	c.currentFn = oldFn
}

func (c *Compiler) generateIdentifierExpr(expr *ast.IdentifierExpr) {
	c.generateOffset(expr)
	c.labelInstruction(c.currentFn.Name, "  pop rax")
	c.labelInstruction(c.currentFn.Name, "  mov rax, [rax]")
	c.labelInstruction(c.currentFn.Name, "  push rax")
}

func (c *Compiler) generateCallExpr(expr *ast.CallExpr) {
	c.labelInstruction(c.currentFn.Name, "  mov rax, 0")
	c.labelInstruction(c.currentFn.Name, fmt.Sprintf("  call %s", expr.Function))
	c.labelInstruction(c.currentFn.Name, "  push rax")
}

func (c *Compiler) generateIfElseStatement(stmt *ast.IfElseStatement) {
	c.generate(stmt.Condition)
	c.labelInstruction(c.currentFn.Name, "  cmp rax, 0")

	c.ifElseCount += 1

	c.labelInstruction(c.currentFn.Name, fmt.Sprintf("  je .L.else.%d", c.ifElseCount))
	c.generate(stmt.Then)
	c.labelInstruction(c.currentFn.Name, fmt.Sprintf("  jmp .L.end.%d", c.ifElseCount))
	c.labelInstruction(c.currentFn.Name, fmt.Sprintf(".L.else.%d:", c.ifElseCount))
	if stmt.Else != nil {
		c.generate(stmt.Else)
	}
	c.labelInstruction(c.currentFn.Name, fmt.Sprintf(".L.end.%d:", c.ifElseCount))
}

func (c *Compiler) generateBinaryOperator(operator ast.BinaryOperator) {
	switch operator {
	case ast.Add:
		c.labelInstruction(c.currentFn.Name, "  add rax, rdi")
	case ast.Subtract:
		c.labelInstruction(c.currentFn.Name, "  sub rax, rdi")
	case ast.Multiply:
		c.labelInstruction(c.currentFn.Name, "  imul rax, rdi")
	case ast.Divide:
		c.labelInstruction(c.currentFn.Name, "  cqo")
		c.labelInstruction(c.currentFn.Name, "  idiv rdi")
	case ast.EqualEqual:
		c.labelInstruction(c.currentFn.Name, "  cmp rax, rdi")
		c.labelInstruction(c.currentFn.Name, "  sete al")
	case ast.BangEqual:
		c.labelInstruction(c.currentFn.Name, "  cmp rax, rdi")
		c.labelInstruction(c.currentFn.Name, "  setne al")
	case ast.LessThanEqual:
		c.labelInstruction(c.currentFn.Name, "  cmp rax, rdi")
		c.labelInstruction(c.currentFn.Name, "  setle al")
	case ast.GreaterThanEqual:
		c.labelInstruction(c.currentFn.Name, "  cmp rax, rdi")
		c.labelInstruction(c.currentFn.Name, "  setge al")
	case ast.Greater:
		c.labelInstruction(c.currentFn.Name, "  cmp rax, rdi")
		c.labelInstruction(c.currentFn.Name, "  setg al")
	case ast.Less:
		c.labelInstruction(c.currentFn.Name, "  cmp rax, rdi")
		c.labelInstruction(c.currentFn.Name, "  setl al")
	default:
		panic("TODO")
	}
}

func (c *Compiler) generateUnaryOperator(operator ast.UnaryOperator) {
	switch operator {
	case ast.Negate:
		c.labelInstruction(c.currentFn.Name, "  neg rax")
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

		c.labelInstruction(c.currentFn.Name, "  mov rax, rbp")
		c.labelInstruction(c.currentFn.Name, fmt.Sprintf("  sub rax, %d", offset))
		c.labelInstruction(c.currentFn.Name, "  push rax")
		return
	}
	if _, ok := node.(*ast.IdentifierExpr); ok {
		offset := node.(*ast.IdentifierExpr).Offset

		c.labelInstruction(c.currentFn.Name, "  mov rax, rbp")
		c.labelInstruction(c.currentFn.Name, fmt.Sprintf("  sub rax, %d", offset))
		c.labelInstruction(c.currentFn.Name, "  push rax")
		return
	}
	panic("TODO")
}

func (c *Compiler) buildAsm() string {
	var asm string

	for label, items := range c.asm {
		if label != "" {
			asm += label + ":" + "\n"
		}
		for _, item := range items {
			asm += item + "\n"
		}
	}

	return asm
}

func (c *Compiler) instruction(instr string) {
	c.asm[""] = append(c.asm[""], instr)
}

func (c *Compiler) labelInstruction(label string, instr string) {
	c.asm[label] = append(c.asm[label], instr)
}
