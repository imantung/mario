package mario

import (
	"fmt"
	"reflect"
)

// Helper implement functionality that is not part of the Handlebars language itself
// https://handlebarsjs.com/guide/expressions.html#helpers
type Helper struct {
	reflect.Value
}

// CreateHelper from function
func CreateHelper(fn interface{}) *Helper {
	val := reflect.ValueOf(fn)
	if err := isValidFunction(val); err != nil {
		panic(err)
	}
	return &Helper{
		Value: val,
	}
}

// panicIfInvalidHelper panics if given helper is not valid
func isValidFunction(fnVal reflect.Value) error {
	fnType := fnVal.Type()
	name := fnType.Name()
	if fnVal.Kind() != reflect.Func {
		return fmt.Errorf("Helper must be a function: %s", name)
	}
	if fnType.NumOut() != 1 {
		return fmt.Errorf("Helper function must return a string or a SafeString: %s", name)
	}
	// TODO: Check if first returned value is a string, SafeString or interface{} ?
	return nil
}
