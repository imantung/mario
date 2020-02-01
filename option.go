package mario

import "reflect"

// Options represents the options argument provided to helpers and context functions.
type Options struct {
	// evaluation visitor
	eval *evaluator

	// params
	params []interface{}
	hash   map[string]interface{}
}

// newOptions instanciates a new Options
func newOptions(eval *evaluator, params []interface{}, hash map[string]interface{}) *Options {
	return &Options{
		eval:   eval,
		params: params,
		hash:   hash,
	}
}

// newEmptyOptions instanciates a new empty Options
func newEmptyOptions(eval *evaluator) *Options {
	return &Options{
		eval: eval,
		hash: make(map[string]interface{}),
	}
}

//
// Context Values
//

// Value returns field value from current context.
func (options *Options) Value(name string) interface{} {
	value := options.eval.evalField(options.eval.curCtx(), name, false)
	if !value.IsValid() {
		return nil
	}

	return value.Interface()
}

// ValueStr returns string representation of field value from current context.
func (options *Options) ValueStr(name string) string {
	return Str(options.Value(name))
}

// Ctx returns current evaluation context.
func (options *Options) Ctx() interface{} {
	return options.eval.curCtx().Interface()
}

//
// Hash Arguments
//

// HashProp returns hash property.
func (options *Options) HashProp(name string) interface{} {
	return options.hash[name]
}

// HashStr returns string representation of hash property.
func (options *Options) HashStr(name string) string {
	return Str(options.hash[name])
}

// Hash returns entire hash.
func (options *Options) Hash() map[string]interface{} {
	return options.hash
}

//
// Parameters
//

// Param returns parameter at given position.
func (options *Options) Param(pos int) interface{} {
	if len(options.params) > pos {
		return options.params[pos]
	}

	return nil
}

// ParamStr returns string representation of parameter at given position.
func (options *Options) ParamStr(pos int) string {
	return Str(options.Param(pos))
}

// Params returns all parameters.
func (options *Options) Params() []interface{} {
	return options.params
}

//
// Private data
//

// Data returns private data value.
func (options *Options) Data(name string) interface{} {
	return options.eval.dataFrame.Get(name)
}

// DataStr returns string representation of private data value.
func (options *Options) DataStr(name string) string {
	return Str(options.eval.dataFrame.Get(name))
}

// DataFrame returns current private data frame.
func (options *Options) DataFrame() *DataFrame {
	return options.eval.dataFrame
}

// NewDataFrame instanciates a new data frame that is a copy of current evaluation data frame.
//
// Parent of returned data frame is set to current evaluation data frame.
func (options *Options) NewDataFrame() *DataFrame {
	return options.eval.dataFrame.Copy()
}

// newIterDataFrame instanciates a new data frame and set iteration specific vars
func (options *Options) newIterDataFrame(length int, i int, key interface{}) *DataFrame {
	return options.eval.dataFrame.newIterDataFrame(length, i, key)
}

//
// Evaluation
//

// evalBlock evaluates block with given context, private data and iteration key
func (options *Options) evalBlock(ctx interface{}, data *DataFrame, key interface{}) string {
	result := ""

	if block := options.eval.curBlock(); (block != nil) && (block.Program != nil) {
		result = options.eval.evalProgram(block.Program, ctx, data, key)
	}

	return result
}

// Fn evaluates block with current evaluation context.
func (options *Options) Fn() string {
	return options.evalBlock(nil, nil, nil)
}

// FnCtxData evaluates block with given context and private data frame.
func (options *Options) FnCtxData(ctx interface{}, data *DataFrame) string {
	return options.evalBlock(ctx, data, nil)
}

// FnWith evaluates block with given context.
func (options *Options) FnWith(ctx interface{}) string {
	return options.evalBlock(ctx, nil, nil)
}

// FnData evaluates block with given private data frame.
func (options *Options) FnData(data *DataFrame) string {
	return options.evalBlock(nil, data, nil)
}

// Inverse evaluates "else block".
func (options *Options) Inverse() string {
	result := ""
	if block := options.eval.curBlock(); (block != nil) && (block.Inverse != nil) {
		result, _ = block.Inverse.Accept(options.eval).(string)
	}

	return result
}

// Eval evaluates field for given context.
func (options *Options) Eval(ctx interface{}, field string) interface{} {
	if ctx == nil {
		return nil
	}

	if field == "" {
		return nil
	}

	val := options.eval.evalField(reflect.ValueOf(ctx), field, false)
	if !val.IsValid() {
		return nil
	}

	return val.Interface()
}

//
// Misc
//

// isIncludableZero returns true if 'includeZero' option is set and first param is the number 0
func (options *Options) isIncludableZero() bool {
	b, ok := options.HashProp("includeZero").(bool)
	if ok && b {
		nb, ok := options.Param(0).(int)
		if ok && nb == 0 {
			return true
		}
	}

	return false
}
