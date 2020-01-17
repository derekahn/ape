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
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)

	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	p.infixParsers = make(map[token.Type]infixParser)
	for tokenType := range precedences {
		p.registerInfix(tokenType, p.parseInfixExpression)

		if tokenType == token.LPAREN {
			p.registerInfix(token.LPAREN, p.parseCallExpression)
		}
	}

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

func (p *Parser) currPrecedence() Priority {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
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

func (p *Parser) noPrefixParserError(t token.Type) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
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

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	// parses until it encounters a '}' signifying the end of the
	// or an EOF which tells us there's no more tokens left to parse
	for !p.currTokenIs(token.RBRACE) && !p.currTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return &block
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.curToken,
		Value: p.currTokenIs(token.TRUE),
	}
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	exp := ast.CallExpression{
		Token:    p.curToken,
		Function: fn,
	}

	exp.Arguments = p.parseCallArguments()

	return &exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	var args []ast.Expression

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

// parseExpression checks whether we have a parsing fn associated with
// p.curToken.Type in the prefix position, if so return parsing fn;
// The heart of our 'Prat Parser'
func (p *Parser) parseExpression(precedence Priority) ast.Expression {
	prefix := p.prefixParsers[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParserError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	// if it's not the end (denoted by ';') keep looping
	// until current priorty is lower than the next
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParsers[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

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

func (p *Parser) parseFunctionLiteral() ast.Expression {
	fn := ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	fn.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	fn.Body = p.parseBlockStatement()

	return &fn
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var identifiers []*ast.Identifier

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	identifiers = append(identifiers, &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	})

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		identifiers = append(identifiers, &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		})
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
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

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return &expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.currPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return &expression
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return &expression
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

func (p *Parser) peekPrecedence() Priority {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
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
