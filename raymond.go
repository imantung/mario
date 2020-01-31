// Package mario provides handlebars evaluation
package mario

// Render parses a template and evaluates it with given context
//
// Note that this function call is not optimal as your template is parsed everytime you call it. You should use Parse() function instead.
func Render(source string, data interface{}) (s string, err error) {
	var tpl *Template
	if tpl, err = Parse(source); err != nil {
		return "", err
	}
	return tpl.Exec(data)
}

// MustRender parses a template and evaluates it with given context. It panics on error.
//
// Note that this function call is not optimal as your template is parsed everytime you call it. You should use Parse() function instead.
func MustRender(source string, data interface{}) string {
	result, err := Render(source, data)
	if err != nil {
		panic(err)
	}
	return result
}
