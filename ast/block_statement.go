package ast

import "fmt"

// BlockStatement represents a block node.
type BlockStatement struct {
	NodeType
	Loc

	Expression *Expression

	Program *Program
	Inverse *Program

	// whitespace management
	OpenStrip    *Strip
	InverseStrip *Strip
	CloseStrip   *Strip
}

// NewBlockStatement instanciates a new block node.
func NewBlockStatement(pos int, line int) *BlockStatement {
	return &BlockStatement{
		NodeType: NodeBlock,
		Loc:      Loc{pos, line},
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *BlockStatement) String() string {
	return fmt.Sprintf("Block{Pos: %d}", node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *BlockStatement) Accept(visitor Visitor) interface{} {
	return visitor.VisitBlock(node)
}
