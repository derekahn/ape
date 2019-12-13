package ast

import (
	"ape/token"
	"bytes"
)

// PrefixExpression is to identify expressions
// with a preceding token.BANG or token.MINUS
type PrefixExpression struct {
	Token    token.Token // The prefix token, ie. !
	Operator string      // Will only be: '!' or '-'
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral is a simple helper to retrieve
// the nested token.Literal string value which
// is the token.BANG or token.MINUS
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

// String formats the expression to be grouped by
// parenthesis so we know which operand belongs to
// which operator. ie (!thing)
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}
