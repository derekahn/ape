package ast

import (
	"ape/token"
)

// LetStatement is for the reserved word 'let'
type LetStatement struct {
	Token token.Token // the token.Let token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

// TokenLiteral is a simple helper to retrieve
// the nested token.Literal string value which
// is the reserved word 'let'
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
