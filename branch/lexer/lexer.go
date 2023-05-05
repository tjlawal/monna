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

func (l_lexer *Lexer) NextToken() token.Token {
	var tok token.Token
	l_lexer.skip_whitespace()

	switch l_lexer.current_char {
	case '=':
		if l_lexer.peek_char() == '=' {
			ch := l_lexer.current_char
			l_lexer.read_char()
			literal := string(ch) + string(l_lexer.current_char)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = new_token(token.ASSIGN, l_lexer.current_char)
		}
	case '!':
		if l_lexer.peek_char() == '=' {
			ch := l_lexer.current_char
			l_lexer.read_char()
			literal := string(ch) + string(l_lexer.current_char)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = new_token(token.BANG, l_lexer.current_char)
		}
	case ';':
		tok = new_token(token.SEMICOLON, l_lexer.current_char)
	case '(':
		tok = new_token(token.LPAREN, l_lexer.current_char)
	case ')':
		tok = new_token(token.RPAREN, l_lexer.current_char)
	case '{':
		tok = new_token(token.LBRACE, l_lexer.current_char)
	case '}':
		tok = new_token(token.RBRACE, l_lexer.current_char)
	case ',':
		tok = new_token(token.COMMA, l_lexer.current_char)
	case '+':
		tok = new_token(token.PLUS, l_lexer.current_char)
	case '-':
		tok = new_token(token.MINUS, l_lexer.current_char)
	case '/':
		tok = new_token(token.SLASH, l_lexer.current_char)
	case '*':
		tok = new_token(token.ASTERISK, l_lexer.current_char)
	case '<':
		tok = new_token(token.LT, l_lexer.current_char)
	case '>':
		tok = new_token(token.GT, l_lexer.current_char)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	default:
		if is_letter(l_lexer.current_char) {
			tok.Literal = l_lexer.read_identifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
			return tok
		} else if is_digit(l_lexer.current_char) {
			tok.Type = token.INT
			tok.Literal = l_lexer.read_number()
			return tok
		} else {
			tok = new_token(token.ILLEGAL, l_lexer.current_char)
		}
	}
	l_lexer.read_char()
	return tok
}

func new_token(TokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: TokenType, Literal: string(ch)}
}

func (l_lexer *Lexer) read_char() {
	if l_lexer.read_position >= len(l_lexer.input) {
		l_lexer.current_char = 0
	} else {
		l_lexer.current_char = l_lexer.input[l_lexer.read_position]
	}
	l_lexer.position = l_lexer.read_position
	l_lexer.read_position += 1
}

func (l_lexer *Lexer) peek_char() byte {
	if l_lexer.read_position >= len(l_lexer.input) {
		return 0
	} else {
		return l_lexer.input[l_lexer.read_position]
	}
}

func (l_lexer *Lexer) read_identifier() string {
	position := l_lexer.position
	for is_letter(l_lexer.current_char) {
		l_lexer.read_char()
	}
	return l_lexer.input[position:l_lexer.position]
}

func is_letter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l_lexer *Lexer) skip_whitespace() {
	for l_lexer.current_char == ' ' || l_lexer.current_char == '\t' || l_lexer.current_char == '\n' || l_lexer.current_char == '\r' {
		l_lexer.read_char()
	}
}

func (l_lexer *Lexer) read_number() string {
	position := l_lexer.position
	for is_digit(l_lexer.current_char) {
		l_lexer.read_char()
	}
	return l_lexer.input[position:l_lexer.position]
}

func is_digit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
