package repl

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestStart(t *testing.T) {
	t.Run("it should produce a PROMPT", func(t *testing.T) {
		var (
			stdin  = strings.NewReader("let x = 5;")
			stdout = &bytes.Buffer{}
		)

		Start(stdin, stdout)
		b, err := ioutil.ReadAll(stdout)

	})
}
