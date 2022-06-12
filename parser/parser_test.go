package parser

import (
	"fmt"

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
	check_parser_errors(l_test, l_parser)

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

func TestReturnStatement(l_test *testing.T) {
	input := `
   return 6;
   return 10;
   return 8419849;
  `

	l_lexer := lexer.New(input)
	l_parser := New(l_lexer)

	program := l_parser.ParseProgram()
	check_parser_errors(l_test, l_parser)

	if len(program.Statements) != 3 {
		l_test.Fatalf("program.Statements does not contain 3 statements, got=%d", len(program.Statements))
	}

	for _, statement := range program.Statements {
		return_statement, ok := statement.(*ast.ReturnStatement)

		if !ok {
			l_test.Errorf("statment not *ast.ReturnStatement, got =%T", statement)
			continue
		}

		if return_statement.TokenLiteral() != "return" {
			l_test.Errorf("return_statement.TokenLiteral() not 'return', got %q", return_statement.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(l_test *testing.T) {
	input := "foobar;"

	l_lexer := lexer.New(input)
	l_parser := New(l_lexer)
	program := l_parser.ParseProgram()
	check_parser_errors(l_test, l_parser)

	if len(program.Statements) != 1 {
		l_test.Fatalf("program does not have enough staments, got=%d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		l_test.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	identifier, ok := statement.Expression.(*ast.Identifier)
	if !ok {
		l_test.Fatalf("expression not *ast.Identifier, got=%T", statement.Expression)
	}

	if identifier.Value != "foobar" {
		l_test.Errorf("identifier.Value not %s, got=%s", "foobar", identifier.Value)
	}

	if identifier.TokenLiteral() != "foobar" {
		l_test.Errorf("identifier.TokenLiteral not %s, got=%s", "foobar", identifier.TokenLiteral())
	}
}
 
func TestIntegerLiteralExpressions(l_test *testing.T) {
	input := "5;"

	l_lexer := lexer.New(input)
	l_parser := New(l_lexer)
	program := l_parser.ParseProgram()

	check_parser_errors(l_test, l_parser)

	if len(program.Statements) != 1 {
		l_test.Fatalf("program does not have enough statements, got=%d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		l_test.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	literal, ok := statement.Expression.(*ast.IntegerLiteral)
	if !ok {
		l_test.Fatalf("expression not *ast.IntegerLiteral, got=%T", statement.Expression)
	}

	if literal.Value != 5 {
		l_test.Errorf("literal.Value not %d, got=%d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		l_test.Errorf("literal.TokenLiteral not %s, got=%s", "5", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpression(l_test *testing.T) {
	prefix_tests := []struct {
		input         string
		operator      string
		integer_value int64
	}{
		{"!5;", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefix_tests {
		l_lexer := lexer.New(tt.input)
		l_parser := New(l_lexer)
		program := l_parser.ParseProgram()
		check_parser_errors(l_test, l_parser)

		if len(program.Statements) != 1 {
			l_test.Fatalf("program.Statements does not contain %d statements, got=%d\n", 1, len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			l_test.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
		}

		expression, ok := statement.Expression.(*ast.PrefixExpression)
		if !ok {
			l_test.Fatalf("program.Statements[0] is not ast.PrefixEXpression, got=%T", statement.Expression)
		}
		if expression.Operator != tt.operator {
			l_test.Fatalf("exp.Operator is not '%s', got %s", tt.operator, expression.Operator)
		}

		if !testIntegerLiteral(l_test, expression.Right, tt.integer_value) {
			return
		}
	}
}

func TestParsingInfixExpressions(l_test *testing.T) {
	infix_tests := []struct {
		input       string
		left_value  int64
		operator    string
		right_value int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infix_tests {
		l_lexer := lexer.New(tt.input)
		l_parser := New(l_lexer)
		program := l_parser.ParseProgram()
		check_parser_errors(l_test, l_parser)

		if len(program.Statements) != 1 {
			l_test.Fatalf("program.Statements does not contain %d statements, got=%d\n", 1, len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			l_test.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
		}

		expression, ok := statement.Expression.(*ast.InfixExpression)
		if !ok {
			l_test.Fatalf("expression is not ast.InfixExpression, got=%T", statement.Expression)
		}

		if !testIntegerLiteral(l_test, expression.Left, tt.left_value) {
			return
		}

		if expression.Operator != tt.operator {
			l_test.Fatalf("expression.Operator is not '%s', got=%s", tt.operator, expression.Operator)
		}

		if !testIntegerLiteral(l_test, expression.Right, tt.right_value) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(l_test *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}
	for _, tt := range tests {
		l_lexer := lexer.New(tt.input)
		l_parser := New(l_lexer)
		program := l_parser.ParseProgram()
		check_parser_errors(l_test, l_parser)

		actual := program.String()
		if actual != tt.expected {
			l_test.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}


// Helpers

func check_parser_errors(l_test *testing.T, l_parser *Parser) {
	errors := l_parser.Errors()
	if len(errors) == 0 {
		return
	}

	l_test.Errorf("parser has %d errors", len(errors))

	for _, message := range errors {
		l_test.Errorf("parser error: %q", message)
	}
	l_test.FailNow()
}

func testIntegerLiteral(l_test *testing.T, il ast.Expression, value int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		l_test.Errorf("il not *ast.IntegerLiteral, got=%T", il)
		return false
	}

	if integer.Value != value {
		l_test.Errorf("integer.Value not %d, got=%d", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		l_test.Errorf("integer.TokenLiteral not %d, got=%s", value, integer.TokenLiteral())
		return false
	}
	return true
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

