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

func TestParsingPrefixExpressions(l_test *testing.T) {
	prefix_tests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
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

		if !testLiteralExpression(l_test, expression.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(l_test *testing.T) {
	infix_tests := []struct {
		input       string
		left_value  interface{}
		operator    string
		right_value interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
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

		if !testInfixExpression(l_test, statement.Expression, tt.left_value, tt.operator, tt.right_value) {
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
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
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

func TestBooleanExpression(l_test *testing.T) {
	tests := []struct {
		input            string
		expected_boolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l_lexer := lexer.New(tt.input)
		l_parser := New(l_lexer)

		program := l_parser.ParseProgram()
		check_parser_errors(l_test, l_parser)

		if len(program.Statements) != 1 {
			l_test.Fatalf("program.Statements does not have enough statements, got=%d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			l_test.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
		}

		boolean, ok := statement.Expression.(*ast.Boolean)
		if !ok {
			l_test.Fatalf("exp not *ast.Boolean, got=%T", statement.Expression)
		}
		if boolean.Value != tt.expected_boolean {
			l_test.Errorf("boolean.Value not %t, got=%t", tt.expected_boolean, boolean.Value)
		}
	}
}

func TestIfExpression(l_test *testing.T) {
	input := `if (x < y) { x }`

	l_lexer := lexer.New(input)
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
	expression, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		l_test.Fatalf("statement.Expression is not ast.IfExpression, got=%T", statement.Expression)
	}

	if !testInfixExpression(l_test, expression.Condition, "x", "<", "y") {
		return
	}
	if len(expression.Consequence.Statements) != 1 {
		l_test.Errorf("consequence is not 1 statements, got=%d\n", len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		l_test.Fatalf("Statements[0] is not ast.ExpressionStatement, got=%T", expression.Consequence.Statements[0])
	}

	if !testIdentifier(l_test, consequence.Expression, "x") {
		return
	}

	if expression.Alternative != nil {
		l_test.Errorf("expression.Alternative.Statements was not nil, got=%+v", expression.Alternative)
	}
}

func TestIfElseExpression(l_test *testing.T) {
	input := `if (x < y) { x } else { y }`

	l_lexer := lexer.New(input)
	l_parser := New(l_lexer)
	program := l_parser.ParseProgram()
	check_parser_errors(l_test, l_parser)

	if len(program.Statements) != 1 {
		l_test.Fatalf("program.Statements does not contain %d statements, got=%d\n", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		l_test.Fatalf("program.Statements[0] is not an ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		l_test.Fatalf("statement.Expression is not ast.IfExpression, got=%T", statement.Expression)
	}

	if !testInfixExpression(l_test, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		l_test.Errorf("consequence is not 1 statements, got=%d\n", len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		l_test.Fatalf("Statements[0] is not ast.ExpressionStatement, got=%T", expression.Consequence.Statements[0])
	}
	if !testIdentifier(l_test, consequence.Expression, "x") {
		return
	}

	if len(expression.Alternative.Statements) != 1 {
		l_test.Errorf("expression.Alterative.Statements does not contain 1 statement, got=%d\n", len(expression.Alternative.Statements))
	}

	alternative, ok := expression.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		l_test.Fatalf("Statements[0] is not ast.ExpressionStatement, got=%T", expression.Alternative.Statements[0])
	}

	if !testIdentifier(l_test, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(l_test *testing.T) {
	input := `fn(x, y) { x + y; }`

	l_lexer := lexer.New(input)
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

	function, ok := statement.Expression.(*ast.FunctionLiteral)
	if !ok {
		l_test.Fatalf("statement.Expression is not ast.FunctionLiteral, got=%T", statement.Expression)
	}

	if len(function.Parameters) != 2 {
		l_test.Fatalf("function literal parameters wrong, want 2, got=%d\n", len(function.Parameters))
	}
	testLiteralExpression(l_test, function.Parameters[0], "x")
	testLiteralExpression(l_test, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		l_test.Fatalf("function.Body.Statements does not have 1 statement, got=%d\n", len(function.Body.Statements))
	}

	body_statement, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		l_test.Fatalf("function body statement is not ast.ExpressionStatemes, got=%T", function.Body.Statements[0])
	}
	testInfixExpression(l_test, body_statement.Expression, "x", "+", "y")

}

func TestFunctionParameterParsing(l_test *testing.T) {
	tests := []struct {
		input           string
		expected_params []string
	}{
		{input: "fn() {};", expected_params: []string{}},
		{input: "fn(x) {};", expected_params: []string{"x"}},
		{input: "fn(x, y, z) {};", expected_params: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l_lexer := lexer.New(tt.input)
		l_parser := New(l_lexer)
		program := l_parser.ParseProgram()
		check_parser_errors(l_test, l_parser)

		statement := program.Statements[0].(*ast.ExpressionStatement)
		function := statement.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expected_params) {
			l_test.Errorf("length of parameters is wrong, want %d, got=%d\n",
				len(tt.expected_params), len(function.Parameters))
		}

		for i, identifier := range tt.expected_params {
			testLiteralExpression(l_test, function.Parameters[i], identifier)
		}
	}
}

func TestCallExpressionParsing(l_test *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l_lexer := lexer.New(input)
	l_parser := New(l_lexer)
	program := l_parser.ParseProgram()
	check_parser_errors(l_test, l_parser)

	if len(program.Statements) != 1 {
		l_test.Fatalf("program.Statements does not contain %d statements, got=%d\n", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		l_test.Fatalf("statement is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.CallExpression)
	if !ok {
		l_test.Fatalf("statemnt.Expression is not ast.CallExpression, got=%T", statement.Expression)
	}

	if !testIdentifier(l_test, expression.Function, "add") {
		return
	}

	if len(expression.Arguments) != 3 {
		l_test.Fatalf("wrong length of arguments, got=%d", len(expression.Arguments))
	}

	testLiteralExpression(l_test, expression.Arguments[0], 1)
	testInfixExpression(l_test, expression.Arguments[1], 2, "*", 3)
	testInfixExpression(l_test, expression.Arguments[2], 4, "+", 5)
}

func TestCallExpressionParameterParsing(l_test *testing.T) {
	tests := []struct {
		input          string
		expected_ident string
		expected_args  []string
	}{
		{
			input:          "add();",
			expected_ident: "add",
			expected_args:  []string{},
		},
		{
			input:          "add(1);",
			expected_ident: "add",
			expected_args:  []string{"1"},
		},
		{
			input:          "add(1, 2 * 3, 4 + 5);",
			expected_ident: "add",
			expected_args:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, tt := range tests {
		l_lexer := lexer.New(tt.input)
		l_parser := New(l_lexer)
		program := l_parser.ParseProgram()
		check_parser_errors(l_test, l_parser)

		statement := program.Statements[0].(*ast.ExpressionStatement)
		expression, ok := statement.Expression.(*ast.CallExpression)
		if !ok {
			l_test.Fatalf("statement.Expression is not ast.CallExpression, got=%T",
				statement.Expression)
		}

		if !testIdentifier(l_test, expression.Function, tt.expected_ident) {
			return
		}

		if len(expression.Arguments) != len(tt.expected_args) {
			l_test.Fatalf("wrong number of arguments, want=%d, got=%d",
				len(tt.expected_args), len(expression.Arguments))
		}

		for i, arg := range tt.expected_args {
			if expression.Arguments[i].String() != arg {
				l_test.Errorf("argument %d wrong. want=%q, got=%q", i,
					arg, expression.Arguments[i].String())
			}
		}
	}
}

func TestLetStatements(l_test *testing.T) {
	tests := []struct {
		input               string
		expected_identifier string
		expected_value      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l_lexer := lexer.New(tt.input)
		l_parser := New(l_lexer)
		program := l_parser.ParseProgram()
		check_parser_errors(l_test, l_parser)

		if len(program.Statements) != 1 {
			l_test.Fatalf("program.Statements does not contain 1 statements, got=%d",
				len(program.Statements))
		}

		statement := program.Statements[0]
		if !testLetStatement(l_test, statement, tt.expected_identifier) {
			return
		}

		val := statement.(*ast.LetStatement).Value
		if !testLiteralExpression(l_test, val, tt.expected_value) {
			return
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

func testIdentifier(l_test *testing.T, exp ast.Expression, value string) bool {
	identifier, ok := exp.(*ast.Identifier)
	if !ok {
		l_test.Errorf("exp not *ast.Identifier, got=%T", exp)
		return false
	}

	if identifier.Value != value {
		l_test.Errorf("identifier.Value not %s, got=%s", value, identifier.Value)
		return false
	}

	if identifier.TokenLiteral() != value {
		l_test.Errorf("identifier.TokenLiteral not %s, got=%s", value, identifier.TokenLiteral())
		return false
	}

	return true
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

func testLiteralExpression(l_test *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(l_test, exp, int64(v))
	case int64:
		return testIntegerLiteral(l_test, exp, v)
	case string:
		return testIdentifier(l_test, exp, v)
	case bool:
		return testBooleanLiteral(l_test, exp, v)
	}

	l_test.Errorf("type of exp not handled, got=%T", exp)
	return false
}

func testInfixExpression(l_test *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	operator_expression, ok := exp.(*ast.InfixExpression)
	if !ok {
		l_test.Errorf("exp is not ast.InfixExpression, got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(l_test, operator_expression.Left, left) {
		return false
	}

	if operator_expression.Operator != operator {
		l_test.Errorf("exp.Operator is not '%s', got=%q", operator, operator_expression.Operator)
		return false
	}

	if !testLiteralExpression(l_test, operator_expression.Right, right) {
		return false
	}
	return true
}

func testBooleanLiteral(l_test *testing.T, exp ast.Expression, value bool) bool {
	boolean, ok := exp.(*ast.Boolean)
	if !ok {
		l_test.Errorf("exp not *ast.Boolean, got=%T", exp)
		return false
	}

	if boolean.Value != value {
		l_test.Errorf("boolean.Value is not %t, got=%t", value, boolean.Value)
		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		l_test.Errorf("boolean.TokenLiteral is not %t, got=%s", value, boolean.TokenLiteral())
		return false
	}

	return true
}
