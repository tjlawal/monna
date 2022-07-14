package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

// Precedence of operations
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

// Precedence Table
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

func (l_parser *Parser) peek_precedence() int {
	if l_parser, ok := precedences[l_parser.peek_token.Type]; ok {
		return l_parser
	}
	return LOWEST
}

func (l_parser *Parser) current_precedence() int {
	if l_parser, ok := precedences[l_parser.current_token.Type]; ok {
		return l_parser
	}
	return LOWEST
}

// Pratt Parsing
type (
	prefix_parse_function func() ast.Expression
	infix_parse_function  func(ast.Expression) ast.Expression
)

func (l_parser *Parser) register_prefix(l_token_type token.TokenType, l_function prefix_parse_function) {
	l_parser.prefix_parse_functions[l_token_type] = l_function
}

func (l_parser *Parser) register_infix(l_token_type token.TokenType, l_function infix_parse_function) {
	l_parser.infix_parse_functions[l_token_type] = l_function
}

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

	// Prefix Operations
	l_parser.prefix_parse_functions = make(map[token.TokenType]prefix_parse_function)
	l_parser.register_prefix(token.IDENT, l_parser.parse_identifier)
	l_parser.register_prefix(token.INT, l_parser.parse_integer_literal)
	l_parser.register_prefix(token.BANG, l_parser.parse_prefix_expression)
	l_parser.register_prefix(token.MINUS, l_parser.parse_prefix_expression)

	// Infix Operation
	l_parser.infix_parse_functions = make(map[token.TokenType]infix_parse_function)
	l_parser.register_infix(token.PLUS, l_parser.parse_infix_expression)
	l_parser.register_infix(token.MINUS, l_parser.parse_infix_expression)
	l_parser.register_infix(token.SLASH, l_parser.parse_infix_expression)
	l_parser.register_infix(token.ASTERISK, l_parser.parse_infix_expression)
	l_parser.register_infix(token.EQ, l_parser.parse_infix_expression)
	l_parser.register_infix(token.NOT_EQ, l_parser.parse_infix_expression)
	l_parser.register_infix(token.LT, l_parser.parse_infix_expression)
	l_parser.register_infix(token.GT, l_parser.parse_infix_expression)

	// Boolean
	l_parser.register_prefix(token.TRUE, l_parser.parse_boolean)
	l_parser.register_prefix(token.FALSE, l_parser.parse_boolean)

	// Grouped Expression
	l_parser.register_prefix(token.LPAREN, l_parser.parse_grouped_expression)
	
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
	defer untrace(trace("parse_let_statement"))
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
	defer untrace(trace("parse_return_statement"))
	statement := &ast.ReturnStatement{Token: l_parser.current_token}
	l_parser.next_token()

	// TODO(tijani): Skipping the expression until there is semicolon
	for !l_parser.current_token_is(token.SEMICOLON) {
		l_parser.next_token()
	}
	return statement
}

func (l_parser *Parser) parse_expression_statement() *ast.ExpressionStatement {
	defer untrace(trace("parse_expression_statement"))
	statement := &ast.ExpressionStatement{Token: l_parser.current_token}
	statement.Expression = l_parser.parse_expression(LOWEST)
	if l_parser.peek_token_is(token.SEMICOLON) {
		l_parser.next_token()
	}
	return statement
}

func (l_parser *Parser) parse_identifier() ast.Expression {
	defer untrace(trace("parse_identifier"))
	return &ast.Identifier{Token: l_parser.current_token, Value: l_parser.current_token.Literal}
}

func (l_parser *Parser) parse_integer_literal() ast.Expression {
	defer untrace(trace("parse_integer_literal"))
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

// Here lies the heart of Pratt Parsing
func (l_parser *Parser) parse_expression(precedence int) ast.Expression {
	defer untrace(trace("parse_expression"))
	prefix := l_parser.prefix_parse_functions[l_parser.current_token.Type]
	if prefix == nil {
		l_parser.no_prefix_parse_function_error(l_parser.current_token.Type)
		return nil
	}
	left_expression := prefix()

	for !l_parser.peek_token_is(token.SEMICOLON) && precedence < l_parser.peek_precedence() {
		infix := l_parser.infix_parse_functions[l_parser.peek_token.Type]
		if infix == nil {
			return left_expression
		}
		l_parser.next_token()
		left_expression = infix(left_expression)
	}

	return left_expression
}

func (l_parser *Parser) parse_prefix_expression() ast.Expression {
	defer untrace(trace("parse_prefix_expression"))
	expression := &ast.PrefixExpression{
		Token:    l_parser.current_token,
		Operator: l_parser.current_token.Literal,
	}
	l_parser.next_token()
	expression.Right = l_parser.parse_expression(PREFIX)
	return expression
}

func (l_parser *Parser) parse_infix_expression(left ast.Expression) ast.Expression {
	defer untrace(trace("parse_prefix_expression"))
	expression := &ast.InfixExpression{
		Token:    l_parser.current_token,
		Operator: l_parser.current_token.Literal,
		Left:     left,
	}
	precedence := l_parser.current_precedence()
	l_parser.next_token()
	expression.Right = l_parser.parse_expression(precedence)

	return expression
}

func (l_parser *Parser) no_prefix_parse_function_error(l_token_type token.TokenType) {
	message := fmt.Sprintf("no prefix parse function for %s, found", l_token_type)
	l_parser.errors = append(l_parser.errors, message)
}

func (l_parser *Parser) parse_boolean() ast.Expression {
	defer untrace(trace("parse_boolean"))
	return &ast.Boolean{
		Token: l_parser.current_token,
		Value: l_parser.current_token_is(token.TRUE),
	}
}

func (l_parser *Parser) parse_grouped_expression() ast.Expression{
	l_parser.next_token()
	expression := l_parser.parse_expression(LOWEST)
	if !l_parser.expect_peek(token.RPAREN){
		return nil
	}
	return expression
}
