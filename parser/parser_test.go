package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatement(l_test *testing.T) {
	input := `
   let x = 4;
   let y = 19;
   let foobar = 8948398493;
  `

	l_lexer := lexer.New(input)
	l_parser := New(l_lexer)

	program := l_parser.ParseProgram()
	if program == nil {
		l_test.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		l_test.Fatalf("program.Statements does not contain 3 statements, got=%d", len(program.Statements))
	}

	tests := []struct {
		expected_identifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]
		if !testLetStatement(l_test, statement, tt.expected_identifier) {
			return
		}
	}
}

func testLetStatement(l_test *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "let" {
		l_test.Errorf("statement.TokenLiteral not let, got=%q", statement.TokenLiteral())
		return false
	}

	let_statement, ok := statement.(*ast.LetStatement)
	if !ok {
		l_test.Errorf("statement not *ast.LetStatement, got=%T", statement)
		return false
	}

	if let_statement.Name.Value != name {
		l_test.Errorf("let_statement.name.Value not %s, got=%s", name, let_statement.Name.Value)
		return false
	}

	if let_statement.Name.TokenLiteral() != name {
		l_test.Errorf("let_statement.name.TokenLiteral() not %s, got=%s", name, let_statement.Name.TokenLiteral())
		return false

	}
	return true
}
