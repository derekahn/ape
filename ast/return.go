package ast

import (
	"ape/token"
)

// ReturnStatement is for the reserved word 'return'
type ReturnStatement struct {
	Token       token.Token // the 'token.Return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral is a simple helper to retrieve
// the nested token.Literal string value which
// is the reserved word 'return'
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
