package ast

import "fmt"

// CommentStatement represents a comment node.
type CommentStatement struct {
	NodeType
	Loc

	Value string

	// whitespace management
	Strip *Strip
}

// NewCommentStatement instanciates a new comment node.
func NewCommentStatement(pos int, line int, val string) *CommentStatement {
	return &CommentStatement{
		NodeType: NodeComment,
		Loc:      Loc{pos, line},

		Value: val,
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *CommentStatement) String() string {
	return fmt.Sprintf("Comment{Value:'%s', Pos:%d}", node.Value, node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *CommentStatement) Accept(visitor Visitor) interface{} {
	return visitor.VisitComment(node)
}
