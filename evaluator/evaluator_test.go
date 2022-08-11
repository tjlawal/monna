package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestEvalIntegerExpression(l_test *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-10", -10},
		{"-5", -5},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"8 * (4 + 32 - (64 / 2) + 5) / (100 / 2 + (25 - 10 + 15) * 18)", 0},
	}

	for _, tt := range tests {
		evaluated := test_eval(tt.input)
		test_integer_object(l_test, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(l_test *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := test_eval(tt.input)
		test_boolean_object(l_test, evaluated, tt.expected)
	}
}

// Test to convert the '!' operator to boolean value and negate it
func TestBangOperator(l_test *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := test_eval(tt.input)
		test_boolean_object(l_test, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(l_test *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else {20}", 20},
		{"if (1 < 2) { 10 } else {20}", 10},
		{"if(10 > 1){if(10 > 1){ return 10;} return 1;}", 10},
	}
	for _, tt := range tests {
		evaluated := test_eval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			test_integer_object(l_test, evaluated, int64(integer))
		} else {
			test_null_object(l_test, evaluated)
		}
	}
}

func TestReturnStatements(l_test *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
	}

	for _, tt := range tests {
		evaluated := test_eval(tt.input)
		test_integer_object(l_test, evaluated, tt.expected)
	}
}

// Helpers
func test_eval(input string) object.Object {
	l_lexer := lexer.New(input)
	l_parser := parser.New(l_lexer)
	program := l_parser.ParseProgram()

	return Eval(program)
}

func test_integer_object(l_test *testing.T, l_object object.Object, expected int64) bool {
	result, ok := l_object.(*object.Integer)
	if !ok {
		l_test.Errorf("object is not integer, got=%T (%+v)", l_object, l_object)
		return false
	}
	if result.Value != expected {
		l_test.Errorf("object has wrong value, got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func test_boolean_object(l_test *testing.T, l_object object.Object, expected bool) bool {
	result, ok := l_object.(*object.Boolean)
	if !ok {
		l_test.Errorf("object is not Boolean, got=%T (%+v)", l_object, l_object)
		return false
	}

	if result.Value != expected {
		l_test.Errorf("object has wrong value, got=%T, want=%t", result.Value, expected)
		return false
	}
	return true
}

func test_null_object(l_test *testing.T, object object.Object) bool {
	if object != NULL {
		l_test.Errorf("object is not NULL, got=%T (%+v)", object, object)
		return false
	}
	return true
}
