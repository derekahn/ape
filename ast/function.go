package ast

import (
	"bytes"
	"strings"

	"ape/token"
)

// FunctionLiteral is how we define functions
// fn <parameters> <block statement>
type FunctionLiteral struct {
	Token      token.Token // the 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

// TokenLiteral is a simple helper to retrieve
// the nested token.Literal string value
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	var params []string
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("( ")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}
