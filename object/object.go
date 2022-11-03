package object

import (
	"bytes"
	"fmt"
	"monkey/ast"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJECT      = "INTEGER"
	BOOLEAN_OBJECT      = "BOOLEAN"
	NULL_OBJECT         = "NULL"
	RETURN_VALUE_OBJECT = "RETURN_VALUE"
	ERROR_OBJECT        = "ERROR"
	FUNCTION_OBJECT     = "FUNCTION"
	STRING_OBJECT       = "STRING"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJECT
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// Booleans
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJECT
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

// Null
type Null struct{}

func (n *Null) Type() ObjectType {
	return NULL_OBJECT
}

func (n *Null) Inspect() string {
	return "null"
}

// Return
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJECT }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// Error
type Error struct {
	Message string
}

func (err *Error) Type() ObjectType { return ERROR_OBJECT }
func (err *Error) Inspect() string  { return "ERROR: " + err.Message }

// Function
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJECT }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// String
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJECT }
func (s *String) Inspec() string   { return s.Value }
