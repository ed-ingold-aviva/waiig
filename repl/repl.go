package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		_, _ = fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		if len(line) == 0 {
			return
		}

		for tok := range lexer.Lex(line) {
			_, _ = fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
