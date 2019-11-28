package repl

// Read, Eval, Print, Loop

import (
	"bufio"
	"fmt"
	"io"

	"ape/lexer"
	"ape/token"
)

// PROMPT is the console symbol indicator
const PROMPT = ">> "

// Start creates an interactive session
// to interpret statements of ApeScript
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n\n", tok)
		}
	}
}
