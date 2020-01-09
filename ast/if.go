package ast

import (
	"ape/token"
	"bytes"
)

// IfExpression does something
// ie. if (<condition>) <consequence> else <alternative>
type IfExpression struct {
	Token       token.Token     // the 'if' token
	Condition   Expression      // A single statement
	Consequence *BlockStatement // truthy block
	Alternative *BlockStatement // falsey block
}

func (ie *IfExpression) expressionNode() {}

// TokenLiteral is a simple helper to retrieve
// the nested token.Literal string value
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}
