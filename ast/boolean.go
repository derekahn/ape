package ast

import "ape/token"

// Boolean AST; In ApeScript we can use
// booleans in place of any other expression
// ie: true;, false;, let foobar = true;
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}

// TokenLiteral is a simple helper to retrieve
// the nested token.Literal string value
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

// String allows us to print the string literal
// value for debugging purposes. It also allows
// for us to adhere to the ast.Node interface
func (b *Boolean) String() string {
	return b.Token.Literal
}
