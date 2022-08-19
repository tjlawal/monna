package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return eval_program(node, env)

	case *ast.BlockStatement:
		return eval_block_statement(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if is_error(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if is_error(val) {
			return val
		}
		env.Set(node.Name.Value, val)

		// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return native_bool_to_boolean_object(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if is_error(right) {
			return right
		}
		return eval_prefix_expression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if is_error(left) {
			return left
		}

		right := Eval(node.Right, env)
		if is_error(right) {
			return right
		}
		return eval_infix_expression(node.Operator, left, right)

	case *ast.IfExpression:
		return eval_if_expression(node, env)

	case *ast.Identifier:
		return eval_identifier(node, env)
	}

	return nil
}

func eval_program(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value

		case *object.Error:
			return result
		}
	}
	return result
}

func eval_identifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return new_error("identifier not found: " + node.Value)
	}
	return val
}

func eval_block_statement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJECT || rt == object.ERROR_OBJECT {
				return result
			}
		}
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
		return new_error("unknown operator: %s%s", operator, right.Type())
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
	if right.Type() != object.INTEGER_OBJECT {
		return new_error("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func eval_infix_expression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJECT && right.Type() == object.INTEGER_OBJECT:
		return eval_integer_infix_expression(operator, left, right)

	case operator == "==":
		return native_bool_to_boolean_object(left == right)

	case operator == "!=":
		return native_bool_to_boolean_object(left != right)

	case left.Type() != right.Type():
		return new_error("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return new_error("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func eval_integer_infix_expression(operator string, left object.Object, right object.Object) object.Object {
	left_value := left.(*object.Integer).Value
	right_value := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: left_value + right_value}

	case "-":
		return &object.Integer{Value: left_value - right_value}

	case "*":
		return &object.Integer{Value: left_value * right_value}

	case "/":
		return &object.Integer{Value: left_value / right_value}

	case "<":
		return native_bool_to_boolean_object(left_value < right_value)

	case ">":
		return native_bool_to_boolean_object(left_value > right_value)

	case "==":
		return native_bool_to_boolean_object(left_value == right_value)

	case "!=":
		return native_bool_to_boolean_object(left_value != right_value)

	default:
		return new_error("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func eval_if_expression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)

	if is_error(condition) {
		return condition
	}

	if is_truthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func is_truthy(object object.Object) bool {
	switch object {
	case NULL:
		return false

	case TRUE:
		return true

	case FALSE:
		return false

	default:
		return true
	}
}

func new_error(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func is_error(l_object object.Object) bool {
	if l_object != nil {
		return l_object.Type() == object.ERROR_OBJECT
	}
	return false
}
