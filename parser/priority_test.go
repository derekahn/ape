package parser

import "testing"

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
