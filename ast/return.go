package ast

import (
	"ape/token"
	"bytes"
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

// String allows us to print the string literal
// value for debugging purposes. It also allows
// for us to adhere to the ast.Node interface
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + "  ")

	// TODO nil checks will be taken out when
	// we can fully build expressions
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}
