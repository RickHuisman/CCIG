package ast

type BinaryOperator string
type UnaryOperator string

const (
	// Binary operators
	Add      = "+"
	Subtract = "-"
	Multiply = "*"
	Divide   = "/"

	// Unary operators
	Negate = "-"
	Not    = "!"
)

var binaryOperators = map[string]BinaryOperator{
	"+": Add,
	"-": Subtract,
	"*": Multiply,
	"/": Divide,
}

var unaryOperators = map[string]UnaryOperator{
	"-": Negate,
	"!": Not,
}

func LookupBinaryOperator(operator string) BinaryOperator {
	if op, ok := binaryOperators[operator]; ok {
		return op
	}
	panic("TODO") // TODO
}

func LookupUnaryOperator(operator string) UnaryOperator {
	if op, ok := unaryOperators[operator]; ok {
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
	Value Expression
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
	Operator UnaryOperator
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
