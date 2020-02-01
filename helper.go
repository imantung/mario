package mario

import (
	"fmt"
	"log"
	"reflect"
	"sync"
)

// helpers stores all globally registered helpers
var helpers = make(map[string]reflect.Value)

// protects global helpers
var helpersMutex sync.RWMutex

func init() {
	// register builtin helpers
	RegisterHelper("if", ifHelper)
	RegisterHelper("unless", unlessHelper)
	RegisterHelper("with", withHelper)
	RegisterHelper("each", eachHelper)
	RegisterHelper("log", logHelper)
	RegisterHelper("lookup", lookupHelper)
	RegisterHelper("equal", equalHelper)
}

// Helpers to return helpers
func Helpers() map[string]reflect.Value {
	return helpers
}

// RegisterHelper registers a global helper. That helper will be available to all templates.
func RegisterHelper(name string, helper interface{}) {
	helpersMutex.Lock()
	defer helpersMutex.Unlock()

	if helpers[name] != zero {
		panic(fmt.Errorf("Helper already registered: %s", name))
	}

	val := reflect.ValueOf(helper)
	ensureValidHelper(name, val)

	helpers[name] = val
}

// RegisterHelpers registers several global helpers. Those helpers will be available to all templates.
func RegisterHelpers(helpers map[string]interface{}) {
	for name, helper := range helpers {
		RegisterHelper(name, helper)
	}
}

// RemoveHelper unregisters a global helper
func RemoveHelper(name string) {
	helpersMutex.Lock()
	defer helpersMutex.Unlock()

	delete(helpers, name)
}

// RemoveAllHelpers unregisters all global helpers
func RemoveAllHelpers() {
	helpersMutex.Lock()
	defer helpersMutex.Unlock()

	helpers = make(map[string]reflect.Value)
}

// ensureValidHelper panics if given helper is not valid
func ensureValidHelper(name string, funcValue reflect.Value) {
	if funcValue.Kind() != reflect.Func {
		panic(fmt.Errorf("Helper must be a function: %s", name))
	}

	funcType := funcValue.Type()

	if funcType.NumOut() != 1 {
		panic(fmt.Errorf("Helper function must return a string or a SafeString: %s", name))
	}

	// @todo Check if first returned value is a string, SafeString or interface{} ?
}

// findHelper finds a globally registered helper
func findHelper(name string) reflect.Value {
	helpersMutex.RLock()
	defer helpersMutex.RUnlock()

	return helpers[name]
}

//
// Builtin helpers
//

// #if block helper
func ifHelper(conditional interface{}, options *Options) interface{} {
	if options.isIncludableZero() || IsTrue(conditional) {
		return options.Fn()
	}

	return options.Inverse()
}

// #unless block helper
func unlessHelper(conditional interface{}, options *Options) interface{} {
	if options.isIncludableZero() || IsTrue(conditional) {
		return options.Inverse()
	}

	return options.Fn()
}

// #with block helper
func withHelper(context interface{}, options *Options) interface{} {
	if IsTrue(context) {
		return options.FnWith(context)
	}

	return options.Inverse()
}

// #each block helper
func eachHelper(context interface{}, options *Options) interface{} {
	if !IsTrue(context) {
		return options.Inverse()
	}

	result := ""

	val := reflect.ValueOf(context)
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			// computes private data
			data := options.newIterDataFrame(val.Len(), i, nil)

			// evaluates block
			result += options.evalBlock(val.Index(i).Interface(), data, i)
		}
	case reflect.Map:
		// note: a go hash is not ordered, so result may vary, this behaviour differs from the JS implementation
		keys := val.MapKeys()
		for i := 0; i < len(keys); i++ {
			key := keys[i].Interface()
			ctx := val.MapIndex(keys[i]).Interface()

			// computes private data
			data := options.newIterDataFrame(len(keys), i, key)

			// evaluates block
			result += options.evalBlock(ctx, data, key)
		}
	case reflect.Struct:
		var exportedFields []int

		// collect exported fields only
		for i := 0; i < val.NumField(); i++ {
			if tField := val.Type().Field(i); tField.PkgPath == "" {
				exportedFields = append(exportedFields, i)
			}
		}

		for i, fieldIndex := range exportedFields {
			key := val.Type().Field(fieldIndex).Name
			ctx := val.Field(fieldIndex).Interface()

			// computes private data
			data := options.newIterDataFrame(len(exportedFields), i, key)

			// evaluates block
			result += options.evalBlock(ctx, data, key)
		}
	}

	return result
}

// #log helper
func logHelper(message string) interface{} {
	log.Print(message)
	return ""
}

// #lookup helper
func lookupHelper(obj interface{}, field string, options *Options) interface{} {
	return Str(options.Eval(obj, field))
}

// #equal helper
// Ref: https://github.com/imantung/mario/issues/7
func equalHelper(a interface{}, b interface{}, options *Options) interface{} {
	if Str(a) == Str(b) {
		return options.Fn()
	}

	return ""
}
