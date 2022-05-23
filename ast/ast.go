package ast

import "monkey/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statement_node()
}

type Expression interface {
	Node
	expression_node()
}

type Program struct {
	Statements []Statement
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statement_node() {}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (i *Identifier) expression_node() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
