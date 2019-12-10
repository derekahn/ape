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
