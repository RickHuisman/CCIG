package ast

type BinaryOperator string

const (
	Add      = "+"
	Subtract = "-"
	Multiply = "*"
	Divide   = "/"
)

var operators = map[string]BinaryOperator{
	"+": Add,
	"-": Subtract,
	"*": Multiply,
	"/": Divide,
}

func LookupBinaryOperator(operator string) BinaryOperator {
	if op, ok := operators[operator]; ok {
		return op
	}
	panic("TODO") // TODO
}

type Node interface {
}

type Statement interface {
	Node
	statementNode()
}

type VarStatement struct {
	Name        string
	Initializer Expression
	Offset      int
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
	Value  string
	Offset int
}

type NullExpr struct {
}

type NumberExpr struct {
	Value float64
}

type PrefixExpr struct {
	Operator BinaryOperator
	Right    Expression
}

type InfixExpr struct {
	Operator BinaryOperator
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
func (*PrefixExpr) expressionNode()     {}
func (*InfixExpr) expressionNode()      {}
func (*CallExpr) expressionNode()       {}
