package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	lexer         *lexer.Lexer
	current_token token.Token
	peek_token    token.Token

	errors []string
}

func New(l_lexer *lexer.Lexer) *Parser {
	l_parser := &Parser{lexer: l_lexer, errors: []string{}}

	// Read two tokens so current_token and peek_token are both set
	l_parser.next_token()
	l_parser.next_token()

	return l_parser
}

func (l_parser *Parser) Errors() []string {
	return l_parser.errors
}

func (l_parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for l_parser.current_token.Type != token.EOF {
		statement := l_parser.parse_statement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		l_parser.next_token()
	}
	return program
}

func (l_parser *Parser) next_token() {
	l_parser.current_token = l_parser.peek_token
	l_parser.peek_token = l_parser.lexer.NextToken()
}

func (l_parser *Parser) parse_statement() ast.Statement {
	switch l_parser.current_token.Type {
	case token.LET:
		return l_parser.parse_let_statement()
	default:
		return nil
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

func (l_parser *Parser) expect_peek(l_token token.TokenType) bool {
	if l_parser.peek_token_is(l_token) {
		l_parser.next_token()
		return true
	} else {
		l_parser.peek_error(l_token)
		return false
	}
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
