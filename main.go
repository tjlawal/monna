/* TODO(tijani):

- add quit command to the repl.
- remove semicolons from the language to mark the end of a statement.
- add line and position reporting in the lexer when a lexing or parsing error occurs.
*/

package main

import (
	"fmt"
	"os"

	"monna/repl"
)

func main() {
	fmt.Printf("Hello human, type some commands: \n")
	repl.Start(os.Stdin, os.Stdout)
}
