package lexer

import "monkey/token"

type Lexer struct {
	input         string
	position      int // current position in input (the current_char)
	read_position int // current reading position in input (after current_char)
	current_char  byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.read_char()
	return l
}

func new_token(TokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: TokenType, Literal: string(ch)}
}

func (lexer *Lexer) read_char() {
	if lexer.read_position >= len(lexer.input) {
		lexer.current_char = 0
	} else {
		lexer.current_char = lexer.input[lexer.read_position]
	}
	lexer.position = lexer.read_position
	lexer.read_position += 1
}

func (lexer *Lexer) peek_char() byte {
	if lexer.read_position >= len(lexer.input) {
		return 0
	} else {
		return lexer.input[lexer.read_position]
	}
}

func (lexer *Lexer) read_identifier() string {
	position := lexer.position
	for is_letter(lexer.current_char) {
		lexer.read_char()
	}
	return lexer.input[position:lexer.position]
}

func is_letter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (lexer *Lexer) skip_whitespace() {
	for lexer.current_char == ' ' || lexer.current_char == '\t' || lexer.current_char == '\n' || lexer.current_char == '\r' {
		lexer.read_char()
	}
}

func (lexer *Lexer) read_number() string {
	position := lexer.position
	for is_digit(lexer.current_char) {
		lexer.read_char()
	}
	return lexer.input[position:lexer.position]
}

func is_digit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (lexer *Lexer) NextToken() token.Token {
	var tok token.Token
	lexer.skip_whitespace()

	switch lexer.current_char {
	case '=':
		if lexer.peek_char() == '=' {
			ch := lexer.current_char
			lexer.read_char()
			literal := string(ch) + string(lexer.current_char)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = new_token(token.ASSIGN, lexer.current_char)
		}
	case '!':
		if lexer.peek_char() == '=' {
			ch := lexer.current_char
			lexer.read_char()
			literal := string(ch) + string(lexer.current_char)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = new_token(token.BANG, lexer.current_char)
		}
	case ';':
		tok = new_token(token.SEMICOLON, lexer.current_char)
	case '(':
		tok = new_token(token.LPAREN, lexer.current_char)
	case ')':
		tok = new_token(token.RPAREN, lexer.current_char)
	case '{':
		tok = new_token(token.LBRACE, lexer.current_char)
	case '}':
		tok = new_token(token.RBRACE, lexer.current_char)
	case ',':
		tok = new_token(token.COMMA, lexer.current_char)
	case '+':
		tok = new_token(token.PLUS, lexer.current_char)
	case '-':
		tok = new_token(token.MINUS, lexer.current_char)
	case '/':
		tok = new_token(token.SLASH, lexer.current_char)
	case '*':
		tok = new_token(token.ASTERISK, lexer.current_char)
	case '<':
		tok = new_token(token.LT, lexer.current_char)
	case '>':
		tok = new_token(token.GT, lexer.current_char)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	default:
		if is_letter(lexer.current_char) {
			tok.Literal = lexer.read_identifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
			return tok
		} else if is_digit(lexer.current_char) {
			tok.Type = token.INT
			tok.Literal = lexer.read_number()
			return tok
		} else {
			tok = new_token(token.ILLEGAL, lexer.current_char)
		}
	}
	lexer.read_char()
	return tok
}
