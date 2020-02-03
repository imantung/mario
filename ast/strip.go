package ast

import "fmt"

// Strip describes node whitespace management.
type Strip struct {
	Open  bool
	Close bool

	OpenStandalone   bool
	CloseStandalone  bool
	InlineStandalone bool
}

// NewStrip instanciates a Strip for given open and close mustaches.
func NewStrip(openStr, closeStr string) *Strip {
	return &Strip{
		Open:  (len(openStr) > 2) && openStr[2] == '~',
		Close: (len(closeStr) > 2) && closeStr[len(closeStr)-3] == '~',
	}
}

// NewStripForStr instanciates a Strip for given tag.
func NewStripForStr(str string) *Strip {
	return &Strip{
		Open:  (len(str) > 2) && str[2] == '~',
		Close: (len(str) > 2) && str[len(str)-3] == '~',
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (s *Strip) String() string {
	return fmt.Sprintf("Open: %t, Close: %t, OpenStandalone: %t, CloseStandalone: %t, InlineStandalone: %t", s.Open, s.Close, s.OpenStandalone, s.CloseStandalone, s.InlineStandalone)
}
