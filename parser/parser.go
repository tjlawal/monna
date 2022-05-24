package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	lexer *lexer.Lexer

	current_token token.Token
	peek_token    token.Token
}

func New(l_lexer *lexer.Lexer) *Parser {
	l_parser := &Parser{lexer: l_lexer}

	// Read two tokens, one for curren_token and another for peek_token
	l_parser.next_token()
	l_parser.next_token()

	return l_parser
}

// When this is first called, current_token and peek_token should point to the same
// token
func (l_parser *Parser) next_token() {
	l_parser.current_token = l_parser.peek_token
	l_parser.peek_token = l_parser.lexer.NextToken()
}

func (l_parser *Parser) ParseProgram() *ast.Program {
	return nil
}
