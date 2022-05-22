/* NOTE(tijani): 
 * l_ in variables names stand for local_ identifiers. The reason is because I do not like
 * single letter variable names, but I also do not want to confuse the myself 
 * when I come back to reading this code after a while hence l_ to signify that it is local to that function.
*/

package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	l_lexer *lexer.Lexer
	current_token token.Token
	peek_token token.Token
}

func (l_parser *Parser) next_token(){
	l_parser.current_token	 = l_parser.peek_token 
	l_parser.peek_token = l_parser.l_lexer.NextToken() 
}

func (p *Parser) ParseProgram() *ast.Program {
	return nil
}

func New(lexer *lexer.Lexer) *Parser {
	l_parser := &Parser {l_lexer: lexer}
	// Read two tokens so current_token and peek_token are set.
	// NOTE(tijani): the first time l_parser.next_token() is called, current_token and peek_token will be pointing to the same token.
	l_parser.next_token()
	l_parser.next_token()

	return l_parser
}

