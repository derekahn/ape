package ast

import (
	"ape/token"
	"bytes"
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

// String allows us to print the string literal
// value for debugging purposes. It also allows
// for us to adhere to the ast.Node interface
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	// TODO nil checks will be taken out when
	// we can fully build expressions
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}
