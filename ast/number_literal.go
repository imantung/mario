package ast

import (
	"fmt"
	"strconv"
)

// NumberLiteral represents a number node.
type NumberLiteral struct {
	NodeType
	Loc

	Value    float64
	IsInt    bool
	Original string
}

// NewNumberLiteral instanciates a new number node.
func NewNumberLiteral(pos int, line int, val float64, isInt bool, original string) *NumberLiteral {
	return &NumberLiteral{
		NodeType: NodeNumber,
		Loc:      Loc{pos, line},

		Value:    val,
		IsInt:    isInt,
		Original: original,
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *NumberLiteral) String() string {
	return fmt.Sprintf("Number{Value:%s, Pos:%d}", node.Canonical(), node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *NumberLiteral) Accept(visitor Visitor) interface{} {
	return visitor.VisitNumber(node)
}

// Canonical returns the canonical form of number node as a string (eg: "12", "-1.51").
func (node *NumberLiteral) Canonical() string {
	prec := -1
	if node.IsInt {
		prec = 0
	}
	return strconv.FormatFloat(node.Value, 'f', prec, 64)
}

// Number returns an integer or a float.
func (node *NumberLiteral) Number() interface{} {
	if node.IsInt {
		return int(node.Value)
	}

	return node.Value
}
