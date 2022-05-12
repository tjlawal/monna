package main

import (
	"fmt"
	"os"
	"os/user"

	"monkey/repl"
)

// TODO(tijani): add quit command.
func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! Welcome to the Monkey programming language!\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
