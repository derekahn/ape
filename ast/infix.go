package ast

import (
	"ape/token"
	"bytes"
)

// InfixExpression as the formula of:
// 2 operands and 1 operator in between
// ie. 5 + 5;
type InfixExpression struct {
	Token    token.Token // The operator token; ie. '+'
	Left     Expression
	Right    Expression
	Operator string
}

func (ie *InfixExpression) expressionNode() {}

// TokenLiteral is a simple helper to retrieve
// the nested token.Literal string value
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// String formats the expression to be grouped by
// parenthesis so we know which operand belongs to
// which operator. ie (5 + 5)
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}
