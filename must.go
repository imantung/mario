package mario

// Must is a helper that wraps a call to a function returning (*Template, error)
// and panics if the error is non-nil. It is intended for use in variable
// initializations such as
func Must(t *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return t
}
