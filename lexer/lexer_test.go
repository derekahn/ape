package lexer

import (
	"ape/token"
	"testing"
)

type tokenTest struct {
	expectedType    token.Type
	expectedLiteral string
}

func TestNextToken(t *testing.T) {
	t.Run("it should perform basic lexical analysis", func(t *testing.T) {
		input := `=+(){},;`

		tests := []tokenTest{
			{token.ASSIGN, "="},
			{token.PLUS, "+"},
			{token.LPAREN, "("},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.RBRACE, "}"},
			{token.COMMA, ","},
			{token.SEMICOLON, ";"},
			{token.EOF, ""},
		}

		run(t, input, tests)
	})

	t.Run("it should perform more advanced lexical analysis", func(t *testing.T) {
		input := `
			let five = 5;
			let ten = 10;
			let add = fn(x, y) {
				x + y;
			};

			let result = add(five, ten);
		`

		tests := []tokenTest{
			{token.LET, "let"},
			{token.IDENT, "five"},
			{token.ASSIGN, "="},
			{token.INT, "5"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "ten"},
			{token.ASSIGN, "="},
			{token.INT, "10"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "add"},
			{token.ASSIGN, "="},
			{token.FUNCTION, "fn"},
			{token.LPAREN, "("},
			{token.IDENT, "x"},
			{token.COMMA, ","},
			{token.IDENT, "y"},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.IDENT, "x"},
			{token.PLUS, "+"},
			{token.IDENT, "y"},
			{token.SEMICOLON, ";"},
			{token.RBRACE, "}"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "result"},
			{token.ASSIGN, "="},
			{token.IDENT, "add"},
			{token.LPAREN, "("},
			{token.IDENT, "five"},
			{token.COMMA, ","},
			{token.IDENT, "ten"},
			{token.RPAREN, ")"},
			{token.SEMICOLON, ";"},
			{token.EOF, ""},
		}

		run(t, input, tests)
	})

	t.Run("it should return ILLEGAL for unknown character", func(t *testing.T) {
		input := `let & 5 [];`

		tests := []tokenTest{
			{token.LET, "let"},
			{token.ILLEGAL, "&"},
			{token.INT, "5"},
			{token.ILLEGAL, "["},
			{token.ILLEGAL, "]"},
			{token.SEMICOLON, ";"},
			{token.EOF, ""},
		}

		run(t, input, tests)
	})

	t.Run("it should handle other operators: !-/*<> ", func(t *testing.T) {
		input := `
			!-/*5;
			5 < 10 > 5;
		`

		tests := []tokenTest{
			{token.BANG, "!"},
			{token.MINUS, "-"},
			{token.SLASH, "/"},
			{token.ASTERIX, "*"},
			{token.INT, "5"},
			{token.SEMICOLON, ";"},
			{token.INT, "5"},
			{token.LT, "<"},
			{token.INT, "10"},
			{token.GT, ">"},
			{token.INT, "5"},
			{token.SEMICOLON, ";"},
		}

		run(t, input, tests)
	})

	t.Run("it should handle conditional statement", func(t *testing.T) {
		input := `
			if (5 < 10) {
				return true;
			} else {
				return false;
			}
		`

		tests := []tokenTest{
			{token.IF, "if"},
			{token.LPAREN, "("},
			{token.INT, "5"},
			{token.LT, "<"},
			{token.INT, "10"},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.RETURN, "return"},
			{token.TRUE, "true"},
			{token.SEMICOLON, ";"},
			{token.RBRACE, "}"},
			{token.ELSE, "else"},
			{token.LBRACE, "{"},
			{token.RETURN, "return"},
			{token.FALSE, "false"},
			{token.SEMICOLON, ";"},
			{token.RBRACE, "}"},
			{token.EOF, ""},
		}

		run(t, input, tests)
	})

		}

		run(t, input, tests)
	})
}

func run(t *testing.T, input string, ts []tokenTest) {
	l := New(input)

	for i, tt := range ts {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf(
				"tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type,
			)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf(
				"tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal,
			)
		}
	}
}
