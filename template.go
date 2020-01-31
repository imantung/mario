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
	source   string
	program  *ast.Program
	helpers  map[string]reflect.Value
	Partials map[string]*partial
	mutex    sync.RWMutex // protects helpers and partials
}

// Parse instanciates a template by parsing given source.
func Parse(source string) (*Template, error) {
	tpl := newTemplate(source)
	err := tpl.parse()
	return tpl, err
}

// ParseFile reads given file and returns parsed template.
func ParseFile(filePath string) (*Template, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return Parse(string(b))
}

// Exec evaluates template with given context.
func (tpl *Template) Exec(ctx interface{}) (result string, err error) {
	return tpl.ExecWith(ctx, nil)
}

// ExecWith evaluates template with given context and private data frame.
func (tpl *Template) ExecWith(ctx interface{}, privData *DataFrame) (result string, err error) {
	defer errRecover(&err)

	// parses template if necessary
	if err = tpl.parse(); err != nil {
		return
	}

	// setup visitor
	v := newEvalVisitor(tpl, ctx, privData)

	// visit AST
	result, _ = tpl.program.Accept(v).(string)

	// named return values
	return
}

// Clone returns a copy of that template.
func (tpl *Template) Clone() *Template {
	result := newTemplate(tpl.source)

	result.program = tpl.program

	tpl.mutex.RLock()
	defer tpl.mutex.RUnlock()

	for name, helper := range tpl.helpers {
		result.RegisterHelper(name, helper.Interface())
	}

	for name, partial := range tpl.Partials {
		result.addPartial(name, partial.source, partial.tpl)
	}

	return result
}

// RegisterHelper registers a helper for that template.
func (tpl *Template) RegisterHelper(name string, helper interface{}) {
	tpl.mutex.Lock()
	defer tpl.mutex.Unlock()

	if tpl.helpers[name] != zero {
		panic(fmt.Sprintf("Helper %s already registered", name))
	}

	val := reflect.ValueOf(helper)
	ensureValidHelper(name, val)

	tpl.helpers[name] = val
}

// RegisterHelpers registers several helpers for that template.
func (tpl *Template) RegisterHelpers(helpers map[string]interface{}) {
	for name, helper := range helpers {
		tpl.RegisterHelper(name, helper)
	}
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

// PrintAST returns string representation of parsed template.
func (tpl *Template) PrintAST() string {
	if err := tpl.parse(); err != nil {
		return fmt.Sprintf("PARSER ERROR: %s", err)
	}
	return ast.Print(tpl.program)
}

func newTemplate(source string) *Template {
	return &Template{
		source:   source,
		helpers:  make(map[string]reflect.Value),
		Partials: make(map[string]*partial),
	}
}

func (tpl *Template) parse() (err error) {
	if tpl.program == nil {
		if tpl.program, err = parser.Parse(tpl.source); err != nil {
			return
		}
	}
	return
}

func (tpl *Template) addPartial(name string, source string, template *Template) {
	tpl.mutex.Lock()
	defer tpl.mutex.Unlock()

	if tpl.Partials[name] != nil {
		panic(fmt.Sprintf("Partial %s already registered", name))
	}

	tpl.Partials[name] = newPartial(name, source, template)
}

func (tpl *Template) findPartial(name string) *partial {
	tpl.mutex.RLock()
	defer tpl.mutex.RUnlock()

	return tpl.Partials[name]
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

func (tpl *Template) findHelper(name string) reflect.Value {
	tpl.mutex.RLock()
	defer tpl.mutex.RUnlock()

	return tpl.helpers[name]
}
