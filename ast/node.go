// Package ast provides structures to represent a handlebars Abstract Syntax Tree, and a Visitor interface to visit that tree.
package ast

import (
	"io"
)

// References:
//   - https://github.com/wycats/handlebars.js/blob/master/lib/handlebars/compiler/ast.js
//   - https://github.com/wycats/handlebars.js/blob/master/docs/compiler-api.md
//   - https://github.com/golang/go/blob/master/src/text/template/parse/node.go

// Node is an element in the AST.
type Node interface {
	// node type
	Type() NodeType

	// location of node in original input string
	Location() Loc

	// string representation, used for debugging
	String() string

	// accepts visitor
	Accept(Visitor) interface{}
}

// Visitor is the interface to visit an AST.
type Visitor interface {
	VisitProgram(io.Writer, *Program) error

	// statements
	VisitMustache(*MustacheStatement) interface{}
	VisitBlock(*BlockStatement) interface{}
	VisitPartial(*PartialStatement) interface{}
	VisitContent(*ContentStatement) interface{}
	VisitComment(*CommentStatement) interface{}

	// expressions
	VisitExpression(*Expression) interface{}
	VisitSubExpression(*SubExpression) interface{}
	VisitPath(*PathExpression) interface{}

	// literals
	VisitString(*StringLiteral) interface{}
	VisitBoolean(*BooleanLiteral) interface{}
	VisitNumber(*NumberLiteral) interface{}

	// miscellaneous
	VisitHash(*Hash) interface{}
	VisitHashPair(*HashPair) interface{}
}

// HelperNameStr returns the string representation of a helper name, with a boolean set to false if this is not a valid helper name.
//
// helperName : path | dataName | STRING | NUMBER | BOOLEAN | UNDEFINED | NULL
func HelperNameStr(node Node) (string, bool) {
	// PathExpression
	if str, ok := PathExpressionStr(node); ok {
		return str, ok
	}

	// Literal
	if str, ok := LiteralStr(node); ok {
		return str, ok
	}

	return "", false
}

// PathExpressionStr returns the string representation of path expression value, with a boolean set to false if this is not a path expression.
func PathExpressionStr(node Node) (string, bool) {
	if path, ok := node.(*PathExpression); ok {
		result := path.Original

		// "[foo bar]"" => "foo bar"
		if (len(result) >= 2) && (result[0] == '[') && (result[len(result)-1] == ']') {
			result = result[1 : len(result)-1]
		}

		return result, true
	}

	return "", false
}

// LiteralStr returns the string representation of literal value, with a boolean set to false if this is not a literal.
func LiteralStr(node Node) (string, bool) {
	if lit, ok := node.(*StringLiteral); ok {
		return lit.Value, true
	}

	if lit, ok := node.(*BooleanLiteral); ok {
		return lit.Canonical(), true
	}

	if lit, ok := node.(*NumberLiteral); ok {
		return lit.Canonical(), true
	}

	return "", false
}
