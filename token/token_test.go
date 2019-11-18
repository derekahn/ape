package token

import "testing"

func TestLookupIdent(t *testing.T) {
	t.Run("it should find the proper identifier for a given word", func(t *testing.T) {
		tests := []struct {
			input    string
			expected Type
		}{
			{"fn", FUNCTION},
			{"let", LET},
			{"foo", IDENT},
			{"random_word", IDENT},
		}

		for i, tt := range tests {
			got := LookupIdent(tt.input)
			if got != tt.expected {
				t.Fatalf(
					"tests[%d] - tokentype wrong. expected=%q, got=%q",
					i, tt.expected, got,
				)
			}
		}
	})
}
