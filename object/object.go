package object
import (
	"fmt"
)
type ObjectType string
const (
	INTEGER_OBJ = "INTEGER",
	BOOLEAN_OBJ = "BOOLEAN",
	STRING_OBJ = "STRING",
)

type Object interface {
	Type() ObjectType
	Inspect() string
}	

type Integer struct {
	Value int64
}
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }


type Boolean struct {
	Value bool
}
func (b *Boolean) Inspect() string { return fmt.Sprintf("%d", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }


type Null struct{}
func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string { return "null" }


type String struct {
	Value string
}

func (s *String) Inspect() string { return fmt.Sprintf("%d", s.Value) }
func (s *String) Type() ObjectType { return STRING_OBJ}