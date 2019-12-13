package parser

import (
	"testing"

	"ape/ast"
	"ape/lexer"
)

func TestLetStatement(t *testing.T) {
	t.Run("it should parse let assignment IDENT portion", func(t *testing.T) {
		input := `
			let x = 5;
			let y = 10;
			let foobar = 838383;
		`

		_, program := initProgram(t, input)
		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}
		if len(program.Statements) != 3 {
			t.Fatalf(
				"program.Statements does not contain 3 statements. got=%d",
				len(program.Statements),
			)
		}

		tests := []struct {
			expectedIdentifier string
		}{
			{"x"},
			{"y"},
			{"foobar"},
		}

		for i, tt := range tests {
			stmt := program.Statements[i]
			if !testingLet(t, stmt, tt.expectedIdentifier) {
				return
			}
		}
	})

	t.Run("it should parse and handle errors", func(t *testing.T) {
		input := `
			let x = 5;
			let = 10;
			let 838383;
		`
		p, _ := initProgram(t, input)
		errors := p.Errors()
		if len(errors) < 1 {
			t.Error("expecting to have errors for a bad let expression")
		}
	})
}

func testingLet(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf(
			"letStmt.Name.Value not '%s'. got=%s",
			name,
			letStmt.Name.TokenLiteral(),
		)
		return false
	}
	return true
}

func TestReturnStatement(t *testing.T) {
	t.Run("it should parse a 'return' expression", func(t *testing.T) {
		input := `
			return 5;
			return 10;
			return 993322;
		`

		_, program := initProgram(t, input)
		if len(program.Statements) != 3 {
			t.Fatalf(
				"program.Statements does not contain 3 statements. got=%d",
				len(program.Statements),
			)
		}

		for _, stmt := range program.Statements {
			returnStmt, ok := stmt.(*ast.ReturnStatement)
			if !ok {
				t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
				continue
			}
			if returnStmt.TokenLiteral() != "return" {
				t.Errorf(
					"returnStmt.Token.TokenLiteral not 'return', got %q",
					returnStmt.TokenLiteral(),
				)
			}
		}
	})
}

func TestIdentifierExpression(t *testing.T) {
	t.Run("it should parse an identifer expression", func(t *testing.T) {
		input := "foobar"

		_, program := initProgram(t, input)
		if len(program.Statements) != 1 {
			t.Fatalf(
				"program doesn't have enough statements. got=%d",
				len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0].(*ast.ExpressionStatement). got=%T",
				program.Statements[0],
			)
		}

		ident, ok := stmt.Expression.(*ast.Identifier)
		if !ok {
			t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
		}
		if ident.Value != input {
			t.Errorf("ident.Value not %s. got=%s", input, ident.Value)
		}
		if ident.TokenLiteral() != input {
			t.Errorf("ident.TokenLiteral not %s. got=%s", input, ident.TokenLiteral())
		}
	})
}

func TestIntegerLiteralExpression(t *testing.T) {
	t.Run("it should parse integers as expressions", func(t *testing.T) {
		input := "5;"
		_, program := initProgram(t, input)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program has not enough statements. got=%d",
				len(program.Statements),
			)
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0],
			)
		}

		literal, ok := stmt.Expression.(*ast.IntegerLiteral)
		if !ok {
			t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
		}
		if literal.Value != 5 {
			t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
		}
		if literal.TokenLiteral() != "5" {
			t.Errorf(
				"literal.TokenLiteral not %s. got=%s",
				"5",
				literal.TokenLiteral(),
			)
		}
	})
}

func initProgram(t *testing.T, input string) (*Parser, *ast.Program) {
	p := New(lexer.New(input))
	checkParserErrors(t, p)

	return p, p.ParseProgram()
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}
