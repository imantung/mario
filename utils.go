package mario

import (
	"fmt"
	"path"
	"reflect"
	"strconv"
	"strings"
)

var htmlEscaper = strings.NewReplacer(
	`&`, "&amp;",
	`'`, "&apos;", // To stay in sync with JS implementation, and make mustache tests pass.
	`<`, "&lt;",
	`>`, "&gt;",
	`"`, "&quot;", // To stay in sync with JS implementation, and make mustache tests pass.
)

// Escape escapes special HTML characters.
//
// It can be used by helpers that return a SafeString and that need to escape some content by themselves.
func Escape(s string) string {
	return htmlEscaper.Replace(s)
}

// SafeString represents a string that must not be escaped.
//
// A SafeString can be returned by helpers to disable escaping.
type SafeString string

// isSafeString returns true if argument is a SafeString
func isSafeString(value interface{}) bool {
	if _, ok := value.(SafeString); ok {
		return true
	}
	return false
}

// Str returns string representation of any basic type value.
func Str(value interface{}) string {
	return strValue(reflect.ValueOf(value))
}

// strValue returns string representation of a reflect.Value
func strValue(value reflect.Value) string {
	result := ""

	ival, ok := printableValue(value)
	if !ok {
		panic(fmt.Errorf("Can't print value: %q", value))
	}

	val := reflect.ValueOf(ival)

	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			result += strValue(val.Index(i))
		}
	case reflect.Bool:
		result = "false"
		if val.Bool() {
			result = "true"
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		result = fmt.Sprintf("%d", ival)
	case reflect.Float32, reflect.Float64:
		result = strconv.FormatFloat(val.Float(), 'f', -1, 64)
	case reflect.Invalid:
		result = ""
	default:
		result = fmt.Sprintf("%s", ival)
	}

	return result
}

// printableValue returns the, possibly indirected, interface value inside v that
// is best for a call to formatted printer.
//
// NOTE: borrowed from https://github.com/golang/go/tree/master/src/text/template/exec.go
func printableValue(v reflect.Value) (interface{}, bool) {
	if v.Kind() == reflect.Ptr {
		v, _ = indirect(v) // fmt.Fprint handles nil.
	}
	if !v.IsValid() {
		return "", true
	}

	if !v.Type().Implements(errorType) && !v.Type().Implements(fmtStringerType) {
		if v.CanAddr() && (reflect.PtrTo(v.Type()).Implements(errorType) || reflect.PtrTo(v.Type()).Implements(fmtStringerType)) {
			v = v.Addr()
		} else {
			switch v.Kind() {
			case reflect.Chan, reflect.Func:
				return nil, false
			}
		}
	}
	return v.Interface(), true
}

// indirect returns the item at the end of indirection, and a bool to indicate if it's nil.
// We indirect through pointers and empty interfaces (only) because
// non-empty interfaces have methods we might need.
//
// NOTE: borrowed from https://github.com/golang/go/tree/master/src/text/template/exec.go
func indirect(v reflect.Value) (rv reflect.Value, isNil bool) {
	for ; v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface; v = v.Elem() {
		if v.IsNil() {
			return v, true
		}
		if v.Kind() == reflect.Interface && v.NumMethod() > 0 {
			break
		}
	}
	return v, false
}

// IsTrue returns true if obj is a truthy value.
func IsTrue(obj interface{}) bool {
	thruth, ok := isTrueValue(reflect.ValueOf(obj))
	if !ok {
		return false
	}
	return thruth
}

// isTrueValue reports whether the value is 'true', in the sense of not the zero of its type,
// and whether the value has a meaningful truth value
//
// NOTE: borrowed from https://github.com/golang/go/tree/master/src/text/template/exec.go
func isTrueValue(val reflect.Value) (truth, ok bool) {
	if !val.IsValid() {
		// Something like var x interface{}, never set. It's a form of nil.
		return false, true
	}
	switch val.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		truth = val.Len() > 0
	case reflect.Bool:
		truth = val.Bool()
	case reflect.Complex64, reflect.Complex128:
		truth = val.Complex() != 0
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Interface:
		truth = !val.IsNil()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		truth = val.Int() != 0
	case reflect.Float32, reflect.Float64:
		truth = val.Float() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		truth = val.Uint() != 0
	case reflect.Struct:
		truth = true // Struct values are always true.
	default:
		return
	}
	return truth, true
}

// canBeNil reports whether an untyped nil can be assigned to the type. See reflect.Zero.
//
// NOTE: borrowed from https://github.com/golang/go/tree/master/src/text/template/exec.go
func canBeNil(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return true
	}
	return false
}

// fileBase returns base file name
//
// example: /foo/bar/baz.png => baz
func fileBase(filePath string) string {
	fileName := path.Base(filePath)
	fileExt := path.Ext(filePath)

	return fileName[:len(fileName)-len(fileExt)]
}
