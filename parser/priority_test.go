package parser

import (
	"fmt"
	"testing"
)

func TestPriority(t *testing.T) {
	tests := []struct {
		constant Priority
		expected string
	}{
		{LOWEST, "LOWEST"},
		{EQUALS, "EQUALS"},
		{LESSGREATER, "LESSGREATER"},
		{SUM, "SUM"},
		{PRODUCT, "PRODUCT"},
		{PREFIX, "PREFIX"},
		{CALL, "CALL"},
		{Priority(0), "UNKNOWN"},
	}

	t.Run("it should have a string method for each Priority", func(t *testing.T) {
		for i, tt := range tests {
			got := tt.constant.String()
			if got != tt.expected {
				t.Fatalf(
					"tests[%d] - Priority wrong. expected=%q, got=%q",
					i, tt.expected, got,
				)
			}
		}
	})
}

func TestOperatorPrcedenceParsing(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("OperatorPrecedence[%d]: it should evaluate the Priority of '%s'", i, tt.input), func(t *testing.T) {
			_, program := initProgram(t, tt.input)
			got := program.String()
			if got != tt.expected {
				t.Errorf("expected=%q, got=%q\n", tt.expected, got)
			}
		})
	}
}
