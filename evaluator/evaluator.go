package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return eval_statements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

		// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return native_bool_to_boolean_object(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return eval_prefix_expression(node.Operator, right)

	}

	return nil
}

func eval_statements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement)
	}
	return result
}

func native_bool_to_boolean_object(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func eval_prefix_expression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return eval_bang_operator_expression(right)
	case "-":
		return eval_minus_prefix_operator_expression(right)

	default:
		return NULL
	}
}

func eval_bang_operator_expression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE

	case FALSE:
		return TRUE

	case NULL:
		return TRUE

	default:
		return FALSE
	}
}

func eval_minus_prefix_operator_expression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJECT{
		return NULL
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}
