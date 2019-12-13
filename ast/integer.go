package ast

import "ape/token"

// IntegerLiteral allows for integer expressions
// it needs to fulfills the ast.Expression interface
// IntegerLiteral.Value  is int64 for the actual value.
// ie. converting source code "5" into an int64
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

// TokenLiteral is a simple helper to retrieve
// the nested token.Literal string value
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

// String allows us to print the string literal
// value for debugging purposes. It also allows
// for us to adhere to the ast.Node interface
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}
