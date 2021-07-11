package tokenizer

type TokenType string

type Token struct {
	TokenType TokenType
	Source    string
}

const (
	Illegal = "Illegal"
	EOF     = "EOF"

	// Identifiers + Literals
	Identifier = "Identifier"
	Number     = "Number"

	// Operators
	Equal    = "="
	Plus     = "+"
	Minus    = "-"
	Asterisk = "*"
	Slash    = "/"
	Bang     = "!"

	// Logical operators
	EqualEqual       = "=="
	BangEqual        = "!="
	LessThanEqual    = "<="
	GreaterThanEqual = ">="
	Less             = "<"
	Greater          = ">"

	// Delimiters
	Comma     = ","
	Semicolon = ";"

	LeftParen  = "("
	RightParen = ")"
	LeftBrace  = "{"
	RightBrace = "}"

	// Keywords
	Function = "Function"
	Var      = "var"
	True     = "True"
	False    = "False"
	If       = "If"
	Else     = "Else"
	Return   = "Return"
)

var keywords = map[string]TokenType{
	"fn":     Function,
	"var":    Var,
	"true":   True,
	"false":  False,
	"if":     If,
	"else":   Else,
	"return": Return,
}

func LookupIdentifierType(identifier string) TokenType {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return Identifier
}
