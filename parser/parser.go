package parser

import (
	"CCIG/ast"
	"CCIG/tokenizer"
	"fmt"
	"os"
)

const (
	_ int = iota
	None
	Assignment    // =
	LessOrGreater // < or >
	Equals        // == or !=
	Sum           // +
	Product       // *
	Prefix        // -X or !X
	Call          // myFunction(X)
	Index         // array[index]
)

type Function struct {
	Name string
	Body *ast.BlockStatement
	//Locals    []Local TODO
	StackSize int
}

//// Local variable
//type Local struct {
//	name   string
//	offset int
//}

var precedences = map[tokenizer.TokenType]int{
	tokenizer.Equal:            Assignment,
	tokenizer.Plus:             Sum,
	tokenizer.Minus:            Sum,
	tokenizer.Slash:            Product,
	tokenizer.Asterisk:         Product,
	tokenizer.EqualEqual:       Equals,
	tokenizer.BangEqual:        Equals,
	tokenizer.LessThanEqual:    Equals,
	tokenizer.GreaterThanEqual: Equals,
	tokenizer.Less:             Equals,
	tokenizer.Greater:          Equals,
	tokenizer.LeftParen:        Call,
}

type (
	prefixParseFn func(tokenizer.Token) ast.Expression
	infixParseFn  func(tokenizer.Token, ast.Expression) ast.Expression
)

type Parser struct {
	tokens         []tokenizer.Token
	prefixParseFns map[tokenizer.TokenType]prefixParseFn
	infixParseFns  map[tokenizer.TokenType]infixParseFn
}

func NewParser(tokens []tokenizer.Token) *Parser {
	p := Parser{tokens: tokens}

	p.prefixParseFns = make(map[tokenizer.TokenType]prefixParseFn)
	p.registerPrefix(tokenizer.Identifier, p.parseIdentifier)
	p.registerPrefix(tokenizer.Number, p.parseNumber)
	p.registerPrefix(tokenizer.Minus, p.parsePrefixExpr)
	p.registerPrefix(tokenizer.Bang, p.parsePrefixExpr)

	p.infixParseFns = make(map[tokenizer.TokenType]infixParseFn)
	p.registerInfix(tokenizer.LeftParen, p.parseCallExpression)
	p.registerInfix(tokenizer.Plus, p.parseInfixExpr)
	p.registerInfix(tokenizer.Minus, p.parseInfixExpr)
	p.registerInfix(tokenizer.Slash, p.parseInfixExpr)
	p.registerInfix(tokenizer.Asterisk, p.parseInfixExpr)
	p.registerInfix(tokenizer.EqualEqual, p.parseInfixExpr)
	p.registerInfix(tokenizer.BangEqual, p.parseInfixExpr)
	p.registerInfix(tokenizer.LessThanEqual, p.parseInfixExpr)
	p.registerInfix(tokenizer.GreaterThanEqual, p.parseInfixExpr)
	p.registerInfix(tokenizer.Less, p.parseInfixExpr)
	p.registerInfix(tokenizer.Greater, p.parseInfixExpr)

	return &p
}

func (p *Parser) Parse() Function {
	var stmts []ast.Statement
	for !p.check(tokenizer.EOF) {
		stmts = append(stmts, p.statement())
	}
	stackSize := assignOffsets(stmts)

	return Function{
		Body:      &ast.BlockStatement{Statements: stmts},
		StackSize: stackSize,
	}
}

func (p *Parser) statement() ast.Statement {
	switch p.peekType() {
	case tokenizer.Var:
		p.consume() // Consume "let"
		return p.parseVarStatement()
	case tokenizer.Function:
		p.consume() // Consume "fn"
		return p.parseFunction()
	case tokenizer.Return:
		p.consume() // Consume "return"
		return p.parseReturn()
	case tokenizer.If:
		p.consume() // Consume "if"
		return p.parseIf()
	default:
		return p.parseExprStatement()
	}
}

func (p *Parser) registerPrefix(tokenType tokenizer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType tokenizer.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) expect(tokenType tokenizer.TokenType, message string) tokenizer.Token {
	if p.check(tokenType) {
		return p.consume()
	}
	reportError(message)
	return p.tokens[0] // TODO
}

func reportError(error string) {
	fmt.Println(error)
	os.Exit(0)
}

func (p *Parser) consume() tokenizer.Token {
	popped := p.tokens[0]
	p.tokens = p.tokens[1:]
	return popped
}

func (p *Parser) match(tokenType tokenizer.TokenType) bool {
	if !p.check(tokenType) {
		return false
	}
	p.consume()
	return true
}

func (p *Parser) check(tokenType tokenizer.TokenType) bool {
	return p.peekType() == tokenType
}

func (p *Parser) peekType() tokenizer.TokenType {
	if !p.hasNext() {
		return tokenizer.EOF
	}
	return p.tokens[0].TokenType
}

func (p *Parser) hasNext() bool {
	return len(p.tokens) != 0
}
