package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

// PRECEDENCE of operations
const (
	_ int = iota // iota means start from 0, hence _ starts from 0
	LOWEST
	EQUALS      // ==
	LESSGREATER // > OR <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -x OR !x
	CALL        // simple_function(x)
)

// Pratt Parsing
type (
	prefix_parse_function func() ast.Expression
	infix_parse_function  func(ast.Expression) ast.Expression
)

type Parser struct {
	lexer         *lexer.Lexer
	current_token token.Token
	peek_token    token.Token

	errors []string

	prefix_parse_functions map[token.TokenType]prefix_parse_function
	infix_parse_functions  map[token.TokenType]infix_parse_function
}

func New(l_lexer *lexer.Lexer) *Parser {
	l_parser := &Parser{lexer: l_lexer, errors: []string{}}

	// Read two tokens so current_token and peek_token are both set
	l_parser.next_token()
	l_parser.next_token()

	l_parser.prefix_parse_functions = make(map[token.TokenType]prefix_parse_function)
	l_parser.register_prefix(token.IDENT, l_parser.parse_identifier)
	l_parser.register_prefix(token.INT, l_parser.parse_integer_literal)
	l_parser.register_prefix(token.BANG, l_parser.parse_prefix_expression)
	l_parser.register_prefix(token.MINUS, l_parser.parse_prefix_expression)

	return l_parser
}

func (l_parser *Parser) Errors() []string {
	return l_parser.errors
}

func (l_parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !l_parser.current_token_is(token.EOF) {
		statement := l_parser.parse_statement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		l_parser.next_token()
	}
	return program
}

func (l_parser *Parser) peek_error(l_token token.TokenType) {
	message := fmt.Sprintf("expected next token to be %s, got %s", l_token, l_parser.peek_token.Type)
	l_parser.errors = append(l_parser.errors, message)
}

func (l_parser *Parser) current_token_is(l_token token.TokenType) bool {
	return l_parser.current_token.Type == l_token
}

func (l_parser *Parser) peek_token_is(l_token token.TokenType) bool {
	return l_parser.peek_token.Type == l_token
}

func (l_parser *Parser) next_token() {
	l_parser.current_token = l_parser.peek_token
	l_parser.peek_token = l_parser.lexer.NextToken()
}

func (l_parser *Parser) expect_peek(l_token token.TokenType) bool {
	if l_parser.peek_token_is(l_token) {
		l_parser.next_token()
		return true
	} else {
		l_parser.peek_error(l_token)
		return false
	}
}

func (l_parser *Parser) parse_statement() ast.Statement {
	switch l_parser.current_token.Type {
	case token.LET:
		return l_parser.parse_let_statement()
	case token.RETURN:
		return l_parser.parse_return_statement()
	default:
		return l_parser.parse_expression_statement()
	}
}

func (l_parser *Parser) parse_let_statement() *ast.LetStatement {
	statement := &ast.LetStatement{Token: l_parser.current_token}
	if !l_parser.expect_peek(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: l_parser.current_token, Value: l_parser.current_token.Literal}
	if !l_parser.expect_peek(token.ASSIGN) {
		return nil
	}

	// TODO(tijani): Skipping the expressins until there is a semicolon
	for !l_parser.current_token_is(token.SEMICOLON) {
		l_parser.next_token()
	}
	return statement
}

func (l_parser *Parser) parse_return_statement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: l_parser.current_token}
	l_parser.next_token()

	// TODO(tijani): Skipping the expression until there is semicolon
	for !l_parser.current_token_is(token.SEMICOLON) {
		l_parser.next_token()
	}
	return statement
}

func (l_parser *Parser) parse_expression_statement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: l_parser.current_token}
	statement.Expression = l_parser.parse_expression(LOWEST)
	if l_parser.peek_token_is(token.SEMICOLON) {
		l_parser.next_token()
	}
	return statement
}

func (l_parser *Parser) register_prefix(l_token_type token.TokenType, l_function prefix_parse_function) {
	l_parser.prefix_parse_functions[l_token_type] = l_function
}

func (l_parser *Parser) parse_identifier() ast.Expression {
	return &ast.Identifier{Token: l_parser.current_token, Value: l_parser.current_token.Literal}
}

func (l_parser *Parser) register_infix(l_token_type token.TokenType, l_function infix_parse_function) {
	l_parser.infix_parse_functions[l_token_type] = l_function
}

func (l_parser *Parser) parse_integer_literal() ast.Expression {
	literal := &ast.IntegerLiteral{Token: l_parser.current_token}
	value, error := strconv.ParseInt(l_parser.current_token.Literal, 0, 64)
	if error != nil {
		message := fmt.Sprintf("could not parse %q as integer", l_parser.current_token.Literal)
		l_parser.errors = append(l_parser.errors, message)
		return nil
	}
	literal.Value = value
	return literal
}

func (l_parser *Parser) parse_expression(precedence int) ast.Expression {
	prefix := l_parser.prefix_parse_functions[l_parser.current_token.Type]
	if prefix == nil {
		l_parser.no_prefix_parse_function_error(l_parser.current_token.Type)
		return nil
	}
	left_expression := prefix()
	return left_expression
}

func (l_parser *Parser) no_prefix_parse_function_error(l_token_type token.TokenType) {
	message := fmt.Sprintf("no prefix parse function for %s, found", l_token_type)
	l_parser.errors = append(l_parser.errors, message)
}

func (l_parser *Parser) parse_prefix_expression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    l_parser.current_token,
		Operator: l_parser.current_token.Literal,
	}
	l_parser.next_token()
	expression.Right = l_parser.parse_expression(PREFIX)
	return expression
}
