package ast

import "bytes"

type (
	// Node is for debugging because
	// it must return the literal value
	// of the token it's associated with
	Node interface {
		TokenLiteral() string
		String() string
	}
	// Statement represents a section of an Expression
	Statement interface {
		Node
		statementNode()
	}
	// Expression is 1 or more Statements
	Expression interface {
		Node
		expressionNode()
	}

	// Program is a collection of statements that
	// are contextually grouped for a single task
	Program struct {
		Statements []Statement
	}
)

// TokenLiteral finds and returns the string value
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// String allows us to print the string literal
// value for debugging purposes. It also allows
// for us to adhere to the ast.Node interface
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
