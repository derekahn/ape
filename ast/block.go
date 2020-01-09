package ast

import (
	"ape/token"
	"bytes"
)

// BlockStatement represents a series of statements
type BlockStatement struct {
	Token      token.Token // the '{' token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

// TokenLiteral is a simple helper to retrieve
// the nested token.Literal string value
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, stmt := range bs.Statements {
		out.WriteString(stmt.String())
	}
	return out.String()
}
