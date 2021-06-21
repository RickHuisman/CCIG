package parser

import (
	"CCIG/ast"
	"CCIG/tokenizer"
	"strconv"
)

func (p *Parser) parseExpression(precedence int) ast.Expression {
	lhs := p.consume()
	prefix := p.prefixParseFns[lhs.TokenType]
	if prefix == nil {
		panic("TODO") // TODO
		// p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}

	leftExp := prefix(lhs)

	for !p.check(tokenizer.Semicolon) && precedence < precedences[p.peekType()] {
		infix := p.infixParseFns[p.peekType()]
		if infix == nil {
			return leftExp
		}

		foobar := p.consume() // TODO
		leftExp = infix(foobar, leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier(token tokenizer.Token) ast.Expression {
	return &ast.IdentifierExpr{Value: token.Source}
}

func (p *Parser) parseNumber(token tokenizer.Token) ast.Expression {
	value, err := strconv.ParseFloat(token.Source, 64)
	if err != nil {
		//msg := fmt.Sprintf("could not parse %q as integer", p.curtokenizer.Literal)
		//msg := fmt.Sprintf("could not parse %q as integer", "10")
		//p.errors = append(p.errors, msg)
		panic("TODO") // TODO
	}

	return &ast.NumberExpr{
		Value: value,
	}
}

func (p *Parser) parseInfixExpression(token tokenizer.Token, left ast.Expression) ast.Expression {
	expr := &ast.InfixExpr{
		Operator: token.Source,
		Left:     left,
	}

	precedence := precedences[token.TokenType]
	expr.Right = p.parseExpression(precedence)

	return expr
}

func (p *Parser) parseCallExpression(_ tokenizer.Token, function ast.Expression) ast.Expression {
	exp := &ast.CallExpr{Function: function}
	exp.Args = p.parseExpressionList(tokenizer.RightParen)
	return exp
}

func (p *Parser) parseExpressionList(end tokenizer.TokenType) []ast.Expression {
	var list []ast.Expression

	if p.check(end) {
		p.consume()
		return list
	}

	list = append(list, p.parseExpression(None))

	for p.check(tokenizer.Comma) {
		p.consume()
		p.consume()
		list = append(list, p.parseExpression(None))
	}

	p.expect(end, "") // TODO Message

	return list
}
