package ast

type Node interface {
}

type Statement interface {
	Node
	statementNode()
}

type VarStatement struct {
	Name        string
	Initializer Expression
}

type BlockStatement struct {
	Statements []Statement
}

type FunStatement struct {
	Name   string
	Params []*IdentifierExpr
	Body   *BlockStatement
}

type ReturnStatement struct {
	ReturnValue Expression
}

type ExprStatement struct {
	Value Expression
}

func (*VarStatement) statementNode()    {}
func (*BlockStatement) statementNode()  {}
func (*FunStatement) statementNode()    {}
func (*ReturnStatement) statementNode() {}
func (*ExprStatement) statementNode()   {}

type Expression interface {
	Node
	expressionNode()
}

type IdentifierExpr struct {
	Value string
}

type NullExpr struct {
}

type NumberExpr struct {
	Value float64
}

type InfixExpr struct {
	Operator string
	Left     Expression
	Right    Expression
}

type CallExpr struct {
	Function Expression
	Args     []Expression
}

func (*IdentifierExpr) expressionNode() {}
func (*NullExpr) expressionNode()       {}
func (*NumberExpr) expressionNode()     {}
func (*InfixExpr) expressionNode()      {}
func (*CallExpr) expressionNode()       {}
