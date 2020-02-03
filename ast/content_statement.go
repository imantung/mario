package ast

import "fmt"

// ContentStatement represents a content node.
type ContentStatement struct {
	NodeType
	Loc

	Value    string
	Original string

	// whitespace management
	RightStripped bool
	LeftStripped  bool
}

// NewContentStatement instanciates a new content node.
func NewContentStatement(pos int, line int, val string) *ContentStatement {
	return &ContentStatement{
		NodeType: NodeContent,
		Loc:      Loc{pos, line},

		Value:    val,
		Original: val,
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *ContentStatement) String() string {
	return fmt.Sprintf("Content{Value:'%s', Pos:%d}", node.Value, node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *ContentStatement) Accept(visitor Visitor) interface{} {
	return visitor.VisitContent(node)
}
