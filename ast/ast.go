package ast

type (
	// Node is for debugging because
	// it msut return the literal value
	// of the token it's associated with
	Node interface {
		TokenLiteral() string
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
