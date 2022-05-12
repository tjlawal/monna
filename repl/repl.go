package repl

import (
	"bufio"
	"fmt"
	"io"

	"monkey/lexer"
	"monkey/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer){
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned{
			return
		}

		line := scanner.Text()
		l := lexer.new(line)

		for tok := l.next_token(); tok.Type != token.EOF; tok = l.next_token(){
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
