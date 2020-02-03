package ast

import "fmt"

// MustacheStatement represents a mustache node.
type MustacheStatement struct {
	NodeType
	Loc

	Unescaped  bool
	Expression *Expression

	// whitespace management
	Strip *Strip
}

// NewMustacheStatement instanciates a new mustache node.
func NewMustacheStatement(pos int, line int, unescaped bool) *MustacheStatement {
	return &MustacheStatement{
		NodeType:  NodeMustache,
		Loc:       Loc{pos, line},
		Unescaped: unescaped,
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *MustacheStatement) String() string {
	return fmt.Sprintf("Mustache{Pos: %d}", node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *MustacheStatement) Accept(visitor Visitor) interface{} {
	return visitor.VisitMustache(node)
}
