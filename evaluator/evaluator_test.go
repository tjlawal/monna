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
	}

	for _, tt := range tests {
		evaluated := test_eval(tt.input)
		test_integer_object(l_test, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(l_test *testing.T){
	tests := []struct {
		input string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests{
			evaluated := test_eval(tt.input)
			test_boolean_object(l_test, evaluated, tt.expected)
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

func test_boolean_object(l_test *testing.T, l_object object.Object, expected bool)bool {
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
