package token

// Type allows many types and
// allows us to distinguish between them
type Type string

// Token for our lexical analysis
type Token struct {
	Type    Type
	Literal string
}

const (
	// ILLEGAL represents an unknown token /character
	ILLEGAL = "ILLEGAL"
	// EOF stands for End Of File
	EOF = "EOF"

	/* Identifiers + Literals */

	// IDENT = indetifier
	IDENT = "IDENT" // add, foobar, x, y, ...
	// INT type
	INT = "INT" // 1343456

	/* Operators */

	// ASSIGN is to attach an IDENT with a value
	ASSIGN = "="
	// PLUS is for mathematical addition
	PLUS = "+"
	// MINUS is for mathematical subtraction
	MINUS = "-"
	// BANG is for inverted boolean logic
	BANG = "!"
	// ASTERIX is for mathematical multiplication
	ASTERIX = "*"
	// SLASH is for mathematical division
	SLASH = "/"
	// LT is for "less than" evaluation
	LT = "<"
	// GT is for "greater than" evaluation
	GT = ">"
	// EQ equal then
	EQ = "=="
	// NEQ is not equal to
	NEQ = "!="

	/* Delimiters */

	// COMMA are argument delimiters
	COMMA = ","
	// SEMICOLON marks the end of an expression
	SEMICOLON = ";"
	// LPAREN = left parenthesis
	LPAREN = "("
	// RPAREN = right parenthesis
	RPAREN = ")"
	// LBRACE = left curly bracket
	LBRACE = "{"
	// RBRACE = right curly bracket
	RBRACE = "}"

	/* Keywords */

	// FUNCTION indicates an expression involving 'n' variables
	FUNCTION = "FUNCTION"
	// LET indicates an assignment expression
	LET = "LET"
	// TRUE is a primitive boolean
	TRUE = "TRUE"
	// FALSE is a primitive boolean
	FALSE = "FALSE"
	// IF is the conditional indicator
	IF = "IF"
	// ELSE is the default conditional
	ELSE = "ELSE"
	// RETURN exits/escapes a function
	RETURN = "RETURN"
)

var keywords = map[string]Type{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// LookupIdent validates if a keyword exists,
// if so it returns it's token.Type. Otherwise
// it'll return the default token.IDENT
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
