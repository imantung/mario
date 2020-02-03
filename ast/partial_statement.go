package ast

import "fmt"

// PartialStatement represents a partial node.
type PartialStatement struct {
	NodeType
	Loc

	Name   Node   // PathExpression | SubExpression
	Params []Node // [ Expression ... ]
	Hash   *Hash

	// whitespace management
	Strip  *Strip
	Indent string
}

// NewPartialStatement instanciates a new partial node.
func NewPartialStatement(pos int, line int) *PartialStatement {
	return &PartialStatement{
		NodeType: NodePartial,
		Loc:      Loc{pos, line},
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *PartialStatement) String() string {
	return fmt.Sprintf("Partial{Name:%s, Pos:%d}", node.Name, node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *PartialStatement) Accept(visitor Visitor) interface{} {
	return visitor.VisitPartial(node)
}
