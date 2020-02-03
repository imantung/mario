package ast

// Loc represents the position of a parsed node in source file.
type Loc struct {
	Pos  int // Byte position
	Line int // Line number
}

// Location returns itself, and permits struct includers to satisfy that part of Node interface.
func (l Loc) Location() Loc {
	return l
}
