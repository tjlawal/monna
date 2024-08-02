package ast

import (
	"monna/token"
	"testing"
)

func TestString(l_test *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "my_var"},
					Value: "my_var",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "another_var"},
					Value: "another_var",
				},
			},
		},
	}

	if program.String() != "let my_var = another_var;" {
		l_test.Errorf("program.String() wrong, got=%q", program.String())
	}
}
