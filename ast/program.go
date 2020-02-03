package ast

import (
	"fmt"
	"strings"
)

// Program represents a program node.
type Program struct {
	NodeType
	Loc

	Body        []Node // [ Statement ... ]
	BlockParams []string
	Chained     bool

	// whitespace management
	Strip *Strip
}

// NewProgram instanciates a new program node.
func NewProgram(pos int, line int) *Program {
	return &Program{
		NodeType: NodeProgram,
		Loc:      Loc{pos, line},
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *Program) String() string {
	return fmt.Sprintf("Program{Pos: %d}", node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *Program) Accept(visitor Visitor) interface{} {
	var b strings.Builder
	visitor.VisitProgram(&b, node)
	return b.String()
}

// AddStatement adds given statement to program.
func (node *Program) AddStatement(statement Node) {
	node.Body = append(node.Body, statement)
}
