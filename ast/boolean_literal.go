package ast

import "fmt"

// BooleanLiteral represents a boolean node.
type BooleanLiteral struct {
	NodeType
	Loc

	Value    bool
	Original string
}

// NewBooleanLiteral instanciates a new boolean node.
func NewBooleanLiteral(pos int, line int, val bool, original string) *BooleanLiteral {
	return &BooleanLiteral{
		NodeType: NodeBoolean,
		Loc:      Loc{pos, line},

		Value:    val,
		Original: original,
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *BooleanLiteral) String() string {
	return fmt.Sprintf("Boolean{Value:%s, Pos:%d}", node.Canonical(), node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *BooleanLiteral) Accept(visitor Visitor) interface{} {
	return visitor.VisitBoolean(node)
}

// Canonical returns the canonical form of boolean node as a string (ie. "true" | "false").
func (node *BooleanLiteral) Canonical() string {
	if node.Value {
		return "true"
	}

	return "false"
}
