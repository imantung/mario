package ast

import "fmt"

// StringLiteral represents a string node.
type StringLiteral struct {
	NodeType
	Loc

	Value string
}

// NewStringLiteral instanciates a new string node.
func NewStringLiteral(pos int, line int, val string) *StringLiteral {
	return &StringLiteral{
		NodeType: NodeString,
		Loc:      Loc{pos, line},

		Value: val,
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *StringLiteral) String() string {
	return fmt.Sprintf("String{Value:'%s', Pos:%d}", node.Value, node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *StringLiteral) Accept(visitor Visitor) interface{} {
	return visitor.VisitString(node)
}
