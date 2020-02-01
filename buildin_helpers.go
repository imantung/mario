package mario

import (
	"log"
	"reflect"
)

// buildinHelpers stores buildin helpers
var buildinHelpers = map[string]*Helper{
	"if":     CreateHelper(ifHelper),
	"unless": CreateHelper(unlessHelper),
	"with":   CreateHelper(withHelper),
	"each":   CreateHelper(eachHelper),
	"log":    CreateHelper(logHelper),
	"lookup": CreateHelper(lookupHelper),
	"equal":  CreateHelper(equalHelper),
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
// Ref: https://github.com/aymerick/raymond/issues/7
func equalHelper(a interface{}, b interface{}, options *Options) interface{} {
	if Str(a) == Str(b) {
		return options.Fn()
	}

	return ""
}
