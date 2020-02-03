package ast

import "fmt"

// Hash represents a hash node.
type Hash struct {
	NodeType
	Loc

	Pairs []*HashPair
}

// NewHash instanciates a new hash node.
func NewHash(pos int, line int) *Hash {
	return &Hash{
		NodeType: NodeHash,
		Loc:      Loc{pos, line},
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *Hash) String() string {
	result := fmt.Sprintf("Hash{[%d", node.Loc.Pos)

	for i, p := range node.Pairs {
		if i > 0 {
			result += ", "
		}
		result += p.String()
	}

	return result + fmt.Sprintf("], Pos:%d}", node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *Hash) Accept(visitor Visitor) interface{} {
	return visitor.VisitHash(node)
}
