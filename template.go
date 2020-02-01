package mario

import (
	"fmt"
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
	partials map[string]*partial
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
		partials: make(map[string]*partial),
	}
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

// RegisterPartial registers a partial for that template.
func (tpl *Template) RegisterPartial(name string, source string) {
	tpl.addPartial(name, source, nil)
}

// RegisterPartials registers several partials for that template.
func (tpl *Template) RegisterPartials(partials map[string]string) {
	for name, partial := range partials {
		tpl.RegisterPartial(name, partial)
	}
}

// RegisterPartialFile reads given file and registers its content as a partial with given name.
func (tpl *Template) RegisterPartialFile(filePath string, name string) (err error) {
	var b []byte
	if b, err = ioutil.ReadFile(filePath); err != nil {
		return err
	}
	tpl.RegisterPartial(name, string(b))
	return nil
}

// RegisterPartialFiles reads several files and registers them as partials, the filename base is used as the partial name.
func (tpl *Template) RegisterPartialFiles(filePaths ...string) error {
	if len(filePaths) == 0 {
		return nil
	}
	for _, filePath := range filePaths {
		name := fileBase(filePath)
		if err := tpl.RegisterPartialFile(filePath, name); err != nil {
			return err
		}
	}
	return nil
}

// RegisterPartialTemplate registers an already parsed partial for that template.
func (tpl *Template) RegisterPartialTemplate(name string, template *Template) {
	tpl.addPartial(name, "", template)
}

// Program return program
func (tpl *Template) Program() *ast.Program {
	return tpl.program
}

func (tpl *Template) addPartial(name string, source string, template *Template) {
	tpl.mutex.Lock()
	defer tpl.mutex.Unlock()

	if tpl.partials[name] != nil {
		panic(fmt.Sprintf("Partial %s already registered", name))
	}

	tpl.partials[name] = newPartial(name, source, template)
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
