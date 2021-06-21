package tokenizer

type Tokenizer struct {
	source  string
	current int
}

func NewTokenizer(source string) *Tokenizer {
	return &Tokenizer{source: source, current: 0}
}

func (t *Tokenizer) Tokenize() []Token {
	var tokens []Token
	for !t.atEof() {
		tokens = append(tokens, t.readToken())
	}
	return tokens
}

func (t *Tokenizer) readToken() Token {
	t.skipWhitespace()

	if t.atEof() {
		return newToken(EOF, 0)
	}

	ch := t.advance()

	if isLetter(ch) {
		return t.identifier()
	}
	if isDigit(ch) {
		return t.number()
	}

	var tokenType TokenType
	switch ch {
	case '=':
		tokenType = Equal
	case '+':
		tokenType = Plus
	case '-':
		tokenType = Minus
	case '*':
		tokenType = Asterisk
	case '/':
		tokenType = Slash
	case ';':
		tokenType = Semicolon
	default:
		return newToken(Illegal, ch)
	}

	return newToken(tokenType, ch)
}

func (t *Tokenizer) identifier() Token {
	start := t.current - 1
	t.advanceWhile(func() bool {
		return isLetter(t.peek()) || isDigit(t.peek())
	})
	identifier := t.source[start:t.current]

	return Token{LookupIdentifierType(identifier), identifier}
}

func (t *Tokenizer) number() Token {
	start := t.current - 1

	t.advanceWhile(func() bool {
		return isDigit(t.peek())
	})

	// Look for a fractional part
	if t.peek() == '.' && isDigit(t.peekNext()) {
		// Consume the "."
		t.advance()

		t.advanceWhile(func() bool {
			return isDigit(t.peek())
		})
	}
	number := t.source[start:t.current]

	return Token{Number, number}
}

func newToken(tokenType TokenType, ch byte) Token {
	return Token{tokenType, string(ch)}
}

func (t *Tokenizer) skipWhitespace() {
	t.advanceWhile(func() bool {
		return isWhitespace(t.peek())
	})
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func (t *Tokenizer) advance() byte {
	t.current += 1
	return t.source[t.current-1]
}

func (t *Tokenizer) advanceWhile(condition func() bool) {
	func() {
		for condition() {
			t.advance()
		}
	}()
}

func (t *Tokenizer) peekNext() byte {
	if t.atEof() {
		return 0
	}
	return t.source[t.current+1]
}

func (t *Tokenizer) peek() byte {
	if t.atEof() {
		return 0
	}
	return t.source[t.current]
}

func (t *Tokenizer) atEof() bool {
	return t.current >= len(t.source)
}
