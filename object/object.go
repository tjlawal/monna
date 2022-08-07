package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJECT = "INTEGER"
	BOOLEAN_OBJECT = "BOOLEAN"
	NULL_OBJECT    = "NULL"
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

func (b *Boolean) Inpect() string {
	return fmt.Sprintf("%t", b.Value)
}

// Null
type Null struct{}

func (n *Null) Type() ObjectType {
	return NULL_OBJECT
}

func (n *Null) Inpect() string {
	return "null"
}
