package ast

// HashPair represents a hash pair node.
type HashPair struct {
	NodeType
	Loc

	Key string
	Val Node // Expression
}

// NewHashPair instanciates a new hash pair node.
func NewHashPair(pos int, line int) *HashPair {
	return &HashPair{
		NodeType: NodeHashPair,
		Loc:      Loc{pos, line},
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *HashPair) String() string {
	return node.Key + "=" + node.Val.String()
}

// Accept is the receiver entry point for visitors.
func (node *HashPair) Accept(visitor Visitor) interface{} {
	return visitor.VisitHashPair(node)
}
