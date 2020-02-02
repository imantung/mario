package mario

import (
	"io/ioutil"
	"reflect"
	"runtime"
	"sync"

	"github.com/imantung/mario/ast"
	"github.com/imantung/mario/parser"
)

// Template represents a handlebars template.
type Template struct {
	program  *ast.Program
	helpers  map[string]*Helper
	partials map[string]*Template
	mutex    sync.RWMutex // protects helpers and partials
}

// New mustache handlebars template
func New() *Template {
	return &Template{
		helpers: map[string]*Helper{
			"if":     ifHelper,
			"unless": unlessHelper,
			"with":   withHelper,
			"each":   eachHelper,
			"log":    logHelper,
			"lookup": lookupHelper,
			"equal":  equalHelper,
		},
		partials: make(map[string]*Template),
	}
}

// Must is a helper that wraps a call to a function returning (*Template, error)
// and panics if the error is non-nil. It is intended for use in variable
// initializations such as
func Must(t *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return t
}

// Parse the template
func (tpl *Template) Parse(source string) (*Template, error) {
	var program *ast.Program
	var err error
	if program, err = parser.Parse(source); err != nil {
		return nil, err
	}
	tpl.program = program
	return tpl, nil
}

// ParseFile reads given file and returns parsed template.
func ParseFile(filePath string) (*Template, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return New().Parse(string(b))
}

// Execute evaluates template with given context.
func (tpl *Template) Execute(ctx interface{}) (result string, err error) {
	return tpl.ExecuteWith(ctx, nil)
}

// ExecuteWith evaluates template with given context and private data frame.
func (tpl *Template) ExecuteWith(ctx interface{}, frame *DataFrame) (result string, err error) {
	defer errRecover(&err)
	if frame == nil {
		frame = NewDataFrame()
	}
	eval := &evaluator{
		helpers:   tpl.helpers,
		partials:  tpl.partials,
		ctx:       []reflect.Value{reflect.ValueOf(ctx)},
		dataFrame: frame,
		exprFunc:  make(map[*ast.Expression]bool),
	}
	result, _ = tpl.program.Accept(eval).(string)
	return
}

// WithHelperFunc to create and set helper
func (tpl *Template) WithHelperFunc(name string, fn interface{}) *Template {
	return tpl.WithHelper(name, CreateHelper(fn))
}

// WithHelper to set helper
func (tpl *Template) WithHelper(name string, helper *Helper) *Template {
	tpl.helpers[name] = helper
	return tpl
}

// WithPartial registers an already parsed partial for that template.
func (tpl *Template) WithPartial(name string, template *Template) *Template {
	tpl.partials[name] = template
	return tpl
}

// Program return program
func (tpl *Template) Program() *ast.Program {
	return tpl.program
}

func errRecover(errp *error) {
	e := recover()
	if e != nil {
		switch err := e.(type) {
		case runtime.Error:
			panic(e)
		case error:
			*errp = err
		default:
			panic(e)
		}
	}
}
