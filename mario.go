package mario

var (
	helpers map[string]*Helper
)

func init() {
	ResetHelpers()
}

// ResetHelpers to return current build-in helpers
func ResetHelpers() {
	helpers = map[string]*Helper{
		// Build-in: https://handlebarsjs.com/guide/builtin-helpers.html
		"if":     CreateHelper(ifHelper),
		"unless": CreateHelper(unlessHelper),
		"with":   CreateHelper(withHelper),
		"each":   CreateHelper(eachHelper),
		"log":    CreateHelper(logHelper),
		"lookup": CreateHelper(lookupHelper),

		// Common helper
		"equal": CreateHelper(equalHelper),
	}
}

// RegisterHelper to register new build-in helpers
func RegisterHelper(name string, fn interface{}) {
	helpers[name] = CreateHelper(fn)
}
