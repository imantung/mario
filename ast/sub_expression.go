package ast

import "fmt"

// SubExpression represents a subexpression node.
type SubExpression struct {
	NodeType
	Loc

	Expression *Expression
}

// NewSubExpression instanciates a new subexpression node.
func NewSubExpression(pos int, line int) *SubExpression {
	return &SubExpression{
		NodeType: NodeSubExpression,
		Loc:      Loc{pos, line},
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *SubExpression) String() string {
	return fmt.Sprintf("Sexp{Path:%s, Pos:%d}", node.Expression.Path, node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *SubExpression) Accept(visitor Visitor) interface{} {
	return visitor.VisitSubExpression(node)
}
