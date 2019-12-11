package ast

import (
	"ape/token"
)

// ExpressionStatement allows for a single expression
// most scripting languages allow for it; ie. 'x + 10;'
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// TokenLiteral is a simple helper to retrieve
// the nested token.Literal string value which
// is evaluation of a single statement/line
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

// String allows us to print the string literal
// value for debugging purposes. It also allows
// for us to adhere to the ast.Node interface
func (es *ExpressionStatement) String() string {

	// TODO nil checks will be taken out when
	// we can fully build expressions
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}
