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

func TestErrorHandling(l_test *testing.T) {
	tests := []struct {
		input            string
		expected_message string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true;", "unknown operator: -BOOLEAN"},
		{"true + false;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) {true + false; }", "unknown operator: BOOLEAN + BOOLEAN"},
		{`
       if (10 > 1){
        if (10 > 1){
         return true + false;
        }
        return 1;
       }
      `, "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found: foobar"},

		// Test to only make sure there is only support for + operator in string concatenation, anything else would be wrong.
		{`"Hello" - "World"`, "unknown operator: STRING - STRING"},
	}

	for _, tt := range tests {
		evaluated := test_eval(tt.input)

		error_object, ok := evaluated.(*object.Error)
		if !ok {
			l_test.Errorf("no error object returned, got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if error_object.Message != tt.expected_message {
			l_test.Errorf("wrong error message, expected=%q, got=%q", tt.expected_message, error_object.Message)
		}
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

func TestLetStatements(l_test *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		test_integer_object(l_test, test_eval(tt.input), tt.expected)
	}
}

func TestFunctionObject(l_test *testing.T) {
	input := "fn(x) { x + 2;};"
	evaluated := test_eval(input)

	fn, ok := evaluated.(*object.Function)
	if !ok {
		l_test.Fatalf("object is not Function, got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		l_test.Fatalf("function has wrong parameters, Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		l_test.Fatalf("parameter is not 'x', got=%q", fn.Parameters[0])
	}

	expected_body := "(x + 2)"

	if fn.Body.String() != expected_body {
		l_test.Fatalf("body is not %q, got=%q", expected_body, fn.Body.String())
	}
}

func TestFunctionApplication(l_test *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		test_integer_object(l_test, test_eval(tt.input), tt.expected)
	}
}

func TestClosures(l_test *testing.T) {
	input := `
    let newAdder = fn(x) {
        fn(y) { x + y };
    };
    let addTwo = newAdder(2);
    addTwo(2);
`
	test_integer_object(l_test, test_eval(input), 4)
}

func TestStringLiteral(l_test *testing.T) {
	input := `"Hello, world!"`

	evaluated := test_eval(input)
	string, ok := evaluated.(*object.String)
	if !ok {
		l_test.Fatalf("object is not String, got=%T (%+v)", evaluated, evaluated)
	}

	if string.Value != "Hello, world!" {
		l_test.Errorf("String has wrong value, got=%q", string.Value)
	}
}

func TestStringConcatenation(l_test *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := test_eval(input)
	string, ok := evaluated.(*object.String)
	if !ok {
		l_test.Fatalf("object is not String, got=%T (%+v)", evaluated, evaluated)
	}

	if string.Value != "Hello World!" {
		l_test.Errorf("String has wrong value, got=%q", string.Value)
	}
}

func TestBuiltinFunctions(l_test *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments, got=2, want=1"},
	}

	for _, tt := range tests {
		evaluated := test_eval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			test_integer_object(l_test, evaluated, int64(expected))
		case string:
			error_object, ok := evaluated.(*object.Error)
			if !ok {
				l_test.Errorf("object is not Error, got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if error_object.Message != expected {
				l_test.Errorf("wrong error message, expected=%q, got=%q", expected, error_object.Message)
			}
		}
	}
}

// Helpers
func test_eval(input string) object.Object {
	l_lexer := lexer.New(input)
	l_parser := parser.New(l_lexer)
	program := l_parser.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
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
