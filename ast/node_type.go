package ast

const (
	// NodeProgram is the program node
	NodeProgram NodeType = iota

	// NodeMustache is the mustache statement node
	NodeMustache

	// NodeBlock is the block statement node
	NodeBlock

	// NodePartial is the partial statement node
	NodePartial

	// NodeContent is the content statement node
	NodeContent

	// NodeComment is the comment statement node
	NodeComment

	// NodeExpression is the expression node
	NodeExpression

	// NodeSubExpression is the subexpression node
	NodeSubExpression

	// NodePath is the expression path node
	NodePath

	// NodeBoolean is the literal boolean node
	NodeBoolean

	// NodeNumber is the literal number node
	NodeNumber

	// NodeString is the literal string node
	NodeString

	// NodeHash is the hash node
	NodeHash

	// NodeHashPair is the hash pair node
	NodeHashPair
)

// NodeType represents an AST Node type.
type NodeType int

// Type returns itself, and permits struct includers to satisfy that part of Node interface.
func (t NodeType) Type() NodeType {
	return t
}
