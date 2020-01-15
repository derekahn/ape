package parser

import (
	"fmt"
	"strings"
	"testing"

	"ape/ast"
	"ape/lexer"
)

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	desc := "BooleanExpression[%d]: should parse '%s' as a boolean"
	for i, tt := range tests {
		t.Run(fmt.Sprintf(desc, i, tt.input), func(t *testing.T) {
			_, program := initProgram(t, tt.input)
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

			boolean, ok := stmt.Expression.(*ast.Boolean)
			if !ok {
				t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
			}
			if boolean.Value != tt.expectedBoolean {
				t.Errorf(
					"boolean.Value not %t. got=%t",
					tt.expectedBoolean,
					boolean.Value,
				)
			}
		})
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := "fn(x, y) { x + y; }"

	_, program := initProgram(t, input)
	if len(program.Statements) != 1 {
		t.Fatalf(
			"program.Statements does not contain %d statements. got=%d",
			1,
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

	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf(
			"stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression,
		)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf(
			"function literal parameters wrong. want 2, got=%d\n",
			len(fn.Parameters),
		)
	}

	testLiteralExpression(t, fn.Parameters[0], "x")
	testLiteralExpression(t, fn.Parameters[1], "y")

	if len(fn.Body.Statements) != 1 {
		t.Fatalf(
			"function.Body.Statements has not 1 statements. got=%d\n",
			len(fn.Body.Statements),
		)
	}

	bodyStmt, ok := fn.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"function body stmt is not ast.ExpressionStatement. got=%T",
			fn.Body.Statements[0],
		)
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	desc := "FunctionParameterParsing[%d]: should parse input '%s'"
	for i, tt := range tests {
		t.Run(fmt.Sprintf(desc, i, tt.input), func(t *testing.T) {

			_, program := initProgram(t, tt.input)
			stmt := program.Statements[0].(*ast.ExpressionStatement)
			fn := stmt.Expression.(*ast.FunctionLiteral)

			if len(fn.Parameters) != len(tt.expectedParams) {
				t.Errorf(
					"length parameters wrong. want %d, got=%d\n",
					len(tt.expectedParams),
					len(fn.Parameters),
				)
			}

			for j, ident := range tt.expectedParams {
				testLiteralExpression(t, fn.Parameters[j], ident)
			}
		})
	}

}

func TestIfExpression(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"if (x < y) { x }"},
		{"if (x < y) { x } else { y }"},
	}
	minStmt := 1

	desc := "IfExpression[%d]: it should parse '%s'"
	for i, tt := range tests {
		t.Run(fmt.Sprintf(desc, i, tt.input), func(t *testing.T) {
			_, program := initProgram(t, tt.input)
			if len(program.Statements) != minStmt {
				t.Fatalf(
					"program.Statements does not contain %d statements. got=%d\n",
					minStmt,
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

			exp, ok := stmt.Expression.(*ast.IfExpression)
			if !ok {
				t.Fatalf(
					"stmt.Expression is not ast.IfExpression. got=%T",
					stmt.Expression,
				)
			}

			if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
				return
			}

			if len(exp.Consequence.Statements) != minStmt {
				t.Errorf(
					"'consequence' is not %d statements. got=%d\n",
					minStmt,
					len(exp.Consequence.Statements),
				)
			}

			consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf(
					"Statements[0] is not ast.ExpressionStatement. got=%T",
					exp.Consequence.Statements[0],
				)
			}

			if !testIdentifier(t, consequence.Expression, "x") {
				return
			}

			// Seperate test logic fo IfElse statements
			if strings.Contains(tt.input, "else") {
				if len(exp.Alternative.Statements) != minStmt {
					t.Errorf(
						"exp.Alternative.Statements does not contain 1 statements. got=%d\n",
						len(exp.Alternative.Statements),
					)
				}

				alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
				if !ok {
					t.Fatalf(
						"Statements[0] is not ast.ExpressionStatement. got=%T",
						exp.Alternative.Statements[0],
					)
				}

				if !testIdentifier(t, alternative.Expression, "y") {
					return
				}
			} else {
				if exp.Alternative != nil {
					t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
				}
			}
		})
	}
}

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
		p := New(lexer.New(input))
		p.ParseProgram()

		errors := p.Errors()
		if len(errors) < 3 {
			t.Error("expecting to have errors for a bad let expression")
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

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
	}

	desc := "OperatorPrecedenceParsing[%d]: it should parse '%s' with correct precedence"
	for i, tt := range tests {
		t.Run(fmt.Sprintf(desc, i, tt.input), func(t *testing.T) {

		})

	}

}

func TestParsingPrefixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	desc := "PrefixExpression[%d]: it should parse '%s' with a prefix operator"
	for i, tt := range tests {
		t.Run(fmt.Sprintf(desc, i, tt.input), func(t *testing.T) {
			_, program := initProgram(t, tt.input)
			if len(program.Statements) != 1 {
				t.Fatalf(
					"program.Statements does not contain %d Statements. got=%d",
					1,
					len(program.Statements),
				)
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf(
					"program.Statements[0] is not an ast.ExpressionStatement. got=%T",
					program.Statements[0],
				)
			}

			exp, ok := stmt.Expression.(*ast.PrefixExpression)
			if !ok {
				t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
			}
			if exp.Operator != tt.operator {
				t.Fatalf("exp.Operator is not '%s'. got=%s",
					tt.operator, exp.Operator)
			}
			if !testLiteralExpression(t, exp.Right, tt.value) {
				return
			}

		})
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		leftVal  interface{}
		operator string
		rightVal interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	desc := "ParseInfixExpression[%d]: it should parse '%s'"
	for i, tt := range tests {
		t.Run(fmt.Sprintf(desc, i, tt.input), func(t *testing.T) {
			_, program := initProgram(t, tt.input)
			if len(program.Statements) != 1 {
				t.Fatalf(
					"program.Statements does not contain %d statements. got=%d\n",
					1,
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
			if !testInfixExpression(t, stmt.Expression, tt.leftVal, tt.operator, tt.rightVal) {
				return
			}
		})
	}
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

/*******************
			HELPERS
*******************/

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

func initProgram(t *testing.T, input string) (*Parser, *ast.Program) {
	p := New(lexer.New(input))
	program := p.ParseProgram()

	checkParserErrors(t, p)

	return p, program
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}
	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s", value, bo.TokenLiteral())
		return false
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}
	return true
}

func testInfixExpression(
	t *testing.T,
	exp ast.Expression,
	left interface{},
	operator string,
	right interface{},
) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		return true
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}
	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	i, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il is not an *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if i.Value != value {
		t.Errorf("i.Value is not %d. got=%d", value, i.Value)
		return false
	}
	if i.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("i.TokenLiteral is not %d. got=%s", value, i.TokenLiteral())
		return false
	}
	return true
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

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case bool:
		return testBooleanLiteral(t, exp, v)
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}
