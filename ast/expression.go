package ast

import "fmt"

// Expression represents an expression node.
type Expression struct {
	NodeType
	Loc

	Path   Node   // PathExpression | StringLiteral | BooleanLiteral | NumberLiteral
	Params []Node // [ Expression ... ]
	Hash   *Hash
}

// NewExpression instanciates a new expression node.
func NewExpression(pos int, line int) *Expression {
	return &Expression{
		NodeType: NodeExpression,
		Loc:      Loc{pos, line},
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *Expression) String() string {
	return fmt.Sprintf("Expr{Path:%s, Pos:%d}", node.Path, node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *Expression) Accept(visitor Visitor) interface{} {
	return visitor.VisitExpression(node)
}

// HelperName returns helper name, or an empty string if this expression can't be a helper.
func (node *Expression) HelperName() string {
	path, ok := node.Path.(*PathExpression)
	if !ok {
		return ""
	}

	if path.Data || (len(path.Parts) != 1) || (path.Depth > 0) || path.Scoped {
		return ""
	}

	return path.Parts[0]
}

// FieldPath returns path expression representing a field path, or nil if this is not a field path.
func (node *Expression) FieldPath() *PathExpression {
	path, ok := node.Path.(*PathExpression)
	if !ok {
		return nil
	}

	return path
}

// LiteralStr returns the string representation of literal value, with a boolean set to false if this is not a literal.
func (node *Expression) LiteralStr() (string, bool) {
	return LiteralStr(node.Path)
}

// Canonical returns the canonical form of expression node as a string.
func (node *Expression) Canonical() string {
	if str, ok := HelperNameStr(node.Path); ok {
		return str
	}

	return ""
}
