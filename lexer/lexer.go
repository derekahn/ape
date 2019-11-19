package lexer

import "ape/token"

/*
In computer science, lexical analysis, lexing or tokenization
is the process of converting a sequence of characters into a sequence
of tokens (strings with an assigned and thus identified meaning).
*/

// Lexer represents a string of source code
type Lexer struct {
	char         byte // current char under examination
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	input        string
}

// New is a factory function to convert an
// input string into a Lexer and initializes
// it by placing the position +1
func New(input string) *Lexer {
	l := Lexer{input: input}
	l.readChar()
	return &l
}

// readChar gives the next char in the input string
// if the index goes past the length it'll reset it's
// position otherwise it increments through each char.
// It's important to note that it only supports ASCII
// for simplicity; otherwise Lexer.char would be a rune
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

// NextToken deciphers and determines the validity
// of a given character of a string of source code,
// if valid it will assign the representative "token"
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.char {
	case '=':
		tok = newToken(token.ASSIGN, l.char)
	case ';':
		tok = newToken(token.SEMICOLON, l.char)
	case '+':
		tok = newToken(token.PLUS, l.char)
	case '-':
		tok = newToken(token.MINUS, l.char)
	case '!':
			tok = newToken(token.BANG, l.char)
	case '/':
		tok = newToken(token.SLASH, l.char)
	case '*':
		tok = newToken(token.ASTERIX, l.char)
	case '<':
		tok = newToken(token.LT, l.char)
	case '>':
		tok = newToken(token.GT, l.char)
	case '(':
		tok = newToken(token.LPAREN, l.char)
	case ')':
		tok = newToken(token.RPAREN, l.char)
	case '{':
		tok = newToken(token.LBRACE, l.char)
	case '}':
		tok = newToken(token.RBRACE, l.char)
	case ',':
		tok = newToken(token.COMMA, l.char)
	case ';':
		tok = newToken(token.SEMICOLON, l.char)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.char) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		}
		if isDigit(l.char) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		}
		tok = newToken(token.ILLEGAL, l.char)
	}

	l.readChar()

	return tok
}

// readIdentifier finds a alphabetical character
// or a sequence of letters that make up a word
// and returns that specific section: ([a], [let])
func (l *Lexer) readIdentifier() string {
	begin := l.position
	for isLetter(l.char) {
		l.readChar()
	}
	end := l.position

	return l.input[begin:end]
}

// readNumber finds a single numeric character
// or a sequence of numbers that make up a big
// number and returns that specific section: ([1], [5523])
func (l *Lexer) readNumber() string {
	begin := l.position
	for isDigit(l.char) {
		l.readChar()
	}
	end := l.position

	return l.input[begin:end]
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readChar()
	}
}

func newToken(ty token.Type, char ...byte) token.Token {
	return token.Token{
		Type:    ty,
		Literal: string(char),
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
