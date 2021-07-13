package parser

import (
	"CCIG/ast"
	"CCIG/tokenizer"
)

func (p *Parser) parseVarStatement() ast.Statement {
	name := p.expect(tokenizer.Identifier, "Expect identifier.").Source

	// Initialize let with null
	var initializer ast.Expression = &ast.NullExpr{}

	// Check for assignment
	if p.check(tokenizer.Equal) {
		p.consume() // Consume '='
		initializer = p.parseExpression(None)
	}

	p.expect(tokenizer.Semicolon, "Expect ';' after expression.")

	return &ast.VarStatement{
		Name:        name,
		Initializer: initializer,
	}
}

func (p *Parser) parseFunction() ast.Statement {
	lit := &ast.FunStatement{}

	lit.Name = p.expect(tokenizer.Identifier, "Expect identifier after fn keyword.").Source

	p.expect(tokenizer.LeftParen, "Expect '(' after function identifier.")

	lit.Params = p.parseFunctionParameters()

	p.expect(tokenizer.LeftBrace, "Expect '{' after function block.")
	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.IdentifierExpr {
	var identifiers []*ast.IdentifierExpr

	if p.check(tokenizer.RightParen) {
		p.consume()
		return identifiers
	}

	ident := &ast.IdentifierExpr{Value: p.consume().Source}
	identifiers = append(identifiers, ident)

	for p.check(tokenizer.Comma) {
		p.consume() // Pop ','
		ident := &ast.IdentifierExpr{Value: p.consume().Source}
		identifiers = append(identifiers, ident)
	}

	p.expect(tokenizer.RightParen, "") // TODO Messsage

	return identifiers
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	var statements []ast.Statement

	for !p.check(tokenizer.RightBrace) && !p.check(tokenizer.EOF) {
		stmt := p.statement()
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}

	p.expect(tokenizer.RightBrace, "Expect '}' after block statement.")

	return &ast.BlockStatement{Statements: statements}
}

func (p *Parser) parseReturn() ast.Statement {
	stmt := &ast.ReturnStatement{
		Value: p.parseExpression(None),
	}

	p.expect(tokenizer.Semicolon, "Expect ';' after return statement.")

	return stmt
}

func (p *Parser) parseIf() ast.Statement {
	p.expect(tokenizer.LeftParen, "Expect '(' before if condition.")
	stmt := &ast.IfElseStatement{
		Condition: p.parseExpression(None),
	}
	p.expect(tokenizer.RightParen, "Expect ')' after if condition.")

	p.expect(tokenizer.LeftBrace, "Expect '{' before block statement.")
	stmt.Then = p.parseBlockStatement()

	if p.match(tokenizer.Else) {
		p.expect(tokenizer.LeftBrace, "Expect '{' before block statement.")

		stmt.Else = p.parseBlockStatement()
	}

	return stmt
}

func (p *Parser) parseExprStatement() ast.Statement {
	expr := &ast.ExprStatement{Value: p.parseExpression(Assignment)} // TODO Precedence
	p.expect(tokenizer.Semicolon, "Expect ';' after expression.")
	return expr
}
