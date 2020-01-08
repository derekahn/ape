package parser

import "testing"

func TestTrace(t *testing.T) {
	t.Run("it should return the string it was given", func(t *testing.T) {
		expected := "1 + 2;"

		got := trace(expected)
		if got != expected {
			t.Fatalf("got='%s'", got)
		}
	})
}

func TestUntrace(t *testing.T) {
	t.Run("This is an example of using 'parser/tracing.go'", func(t *testing.T) {
		defer untrace(trace("parseExpressionStatement"))
	})
}
