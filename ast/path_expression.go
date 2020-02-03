package ast

import "fmt"

// PathExpression represents a path expression node.
type PathExpression struct {
	NodeType
	Loc

	Original string
	Depth    int
	Parts    []string
	Data     bool
	Scoped   bool
}

// NewPathExpression instanciates a new path expression node.
func NewPathExpression(pos int, line int, data bool) *PathExpression {
	result := &PathExpression{
		NodeType: NodePath,
		Loc:      Loc{pos, line},

		Data: data,
	}

	if data {
		result.Original = "@"
	}

	return result
}

// String returns a string representation of receiver that can be used for debugging.
func (node *PathExpression) String() string {
	return fmt.Sprintf("Path{Original:'%s', Pos:%d}", node.Original, node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *PathExpression) Accept(visitor Visitor) interface{} {
	return visitor.VisitPath(node)
}

// Part adds path part.
func (node *PathExpression) Part(part string) {
	node.Original += part

	switch part {
	case "..":
		node.Depth++
		node.Scoped = true
	case ".", "this":
		node.Scoped = true
	default:
		node.Parts = append(node.Parts, part)
	}
}

// Sep adds path separator.
func (node *PathExpression) Sep(separator string) {
	node.Original += separator
}

// IsDataRoot returns true if path expression is @root.
func (node *PathExpression) IsDataRoot() bool {
	return node.Data && (node.Parts[0] == "root")
}
