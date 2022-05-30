/* TODO(tijani):

- add quit command to the repl.
- remove semicolons from the language to mark the end of a statement.
- add line and position reporting in the lexer when a lexing or parsing error occurs.
*/

package main

import (
	"fmt"
	"os"

	"monkey/repl"
)

func main() {
	// user, err := user.Current()
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Printf("Welcome to the Monk programming language!\n")
	repl.Start(os.Stdin, os.Stdout)
}
