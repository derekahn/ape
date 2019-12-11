package ast

import "ape/token"

// Identifier is for special symbols/characters
type Identifier struct {
	Token token.Token // the token.Ident token
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral is a simple helper to retrieve
// the nested token.Literal string value which
// is a special symbol or character
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// String allows us to print the string literal
// value for debugging purposes. Although it isn't
// required for an interace, it's for consistency
func (i *Identifier) String() string {
	return i.Value
}
