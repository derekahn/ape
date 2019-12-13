package parser

import (
	"fmt"
	"strconv"

	"ape/ast"
	"ape/lexer"
	"ape/token"
)

type (
	// Parser is the data structure representing
	// the internal representation and position
	Parser struct {
		lex *lexer.Lexer

		curToken  token.Token
		peekToken token.Token

		errors []string

		prefixParsers map[token.Type]prefixParser
		infixParsers  map[token.Type]infixParser
	}

	prefixParser func() ast.Expression
	infixParser  func(ast.Expression) ast.Expression
)

// New is a factory function that produces
// a new initialized Parser{}; initializes
// with 'curToken' and 'peekToken' being set
func New(l *lexer.Lexer) *Parser {
	p := Parser{
		lex:    l,
		errors: []string{},
	}

	p.prefixParsers = make(map[token.Type]prefixParser)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)

	// Read two tokens, so curToken
	// and peekToken are both set
	p.nextToken()
	p.nextToken()

	return &p
}

// Errors is a get of the current parsing errors
func (p *Parser) Errors() []string {
	return p.errors
}

// ParseProgram constructs the root node of the *ast.Program.
// It then iterates over every token in the input until it finds
// an EOF token. Otherwise it appends it to ast.Statements
func (p *Parser) ParseProgram() *ast.Program {
	program := ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.currTokenIs(token.EOF) {
		if stmt := p.parseStatement(); stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return &program
}

func (p *Parser) currTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.peekError(t)
	return false
}

// small helper that advances both curToken and peekToken
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := ast.LetStatement{
		Token: p.curToken,
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO We're skipping the expressions until we
	// encounter a semicolon
	for !p.currTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return &stmt
}

// parseExpression checks whether we have a parsing fn associated with
// p.curToken.Type in the prefix position, if so return parsing fn
func (p *Parser) parseExpression(precedence Priority) ast.Expression {
	prefix := p.prefixParsers[p.curToken.Type]
	if prefix == nil {
		return nil
	}

	leftExp := prefix()
	return leftExp
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return &stmt
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	val, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	l := ast.IntegerLiteral{Token: p.curToken}
	l.Value = val

	return &l
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// TODO We're skipping the expressions until we
	// encounter a semicolon
	for !p.currTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return &stmt
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf(
		"expected next token to be %s, got %s instead",
		t, p.peekToken.Type,
	)

	p.errors = append(p.errors, msg)
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) registerPrefix(t token.Type, fn prefixParser) {
	p.prefixParsers[t] = fn
}

func (p *Parser) registerInfix(t token.Type, fn infixParser) {
	p.infixParsers[t] = fn
}
