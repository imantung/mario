package mario_test

import (
	"io/ioutil"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/imantung/mario"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

//
// Note, as the JS implementation, the divergences from mustache spec:
//   - we don't support alternative delimeters
//   - the mustache lambda spec differs
//

type mustacheTest struct {
	Name     string
	Desc     string
	Data     interface{}
	Template string
	Expected string
	Partials map[string]string
}

type mustacheTestFile struct {
	Overview string
	Tests    []mustacheTest
}

var (
	rAltDelim = regexp.MustCompile(regexp.QuoteMeta("{{="))
)

var (
	musTestLambdaInterMult = 0
)

func TestMustache(t *testing.T) {
	skipFiles := map[string]bool{
		// mustache lambdas differ from handlebars lambdas
		"~lambdas.yml": true,
	}

	for _, fileName := range mustacheTestFiles() {
		if skipFiles[fileName] {
			// fmt.Printf("Skipped file: %s\n", fileName)
			continue
		}

		for i, tt := range testsFromMustacheFile(fileName) {
			tpl := mario.New()
			for name, fn := range tt.helpers {
				tpl.WithHelperFunc(name, fn)
			}
			for name, source := range tt.partials {
				tpl.WithPartial(name, mario.Must(mario.New().Parse(source)))
			}

			var b strings.Builder
			if err := mario.Must(tpl.Parse(tt.template)).Execute(&b, tt.data); err != nil {
				require.EqualError(t, err, tt.expectedError, i)
			} else {
				require.Equal(t, tt.expected, b.String(), i)
			}
		}

	}
}

func testsFromMustacheFile(fileName string) []testcase {
	result := []testcase{}

	fileData, err := ioutil.ReadFile(path.Join("mustache", "specs", fileName))
	if err != nil {
		panic(err)
	}

	var testFile mustacheTestFile
	if err := yaml.Unmarshal(fileData, &testFile); err != nil {
		panic(err)
	}

	for _, mustacheTest := range testFile.Tests {
		if mustBeSkipped(mustacheTest, fileName) {
			// fmt.Printf("Skipped test: %s\n", mustacheTest.Name)
			continue
		}

		test := testcase{
			template: mustacheTest.Template,
			data:     mustacheTest.Data,
			partials: mustacheTest.Partials,
			expected: mustacheTest.Expected,
		}

		result = append(result, test)
	}

	return result
}

// returns true if test must be skipped
func mustBeSkipped(test mustacheTest, fileName string) bool {
	// handlebars does not support alternative delimiters
	return haveAltDelimiter(test) ||
		// the JS implementation skips those tests
		fileName == "partials.yml" && (test.Name == "Failed Lookup" || test.Name == "Standalone Indentation")
}

// returns true if test have alternative delimeter in template or in partials
func haveAltDelimiter(test mustacheTest) bool {
	// check template
	if rAltDelim.MatchString(test.Template) {
		return true
	}

	// check partials
	for _, partial := range test.Partials {
		if rAltDelim.MatchString(partial) {
			return true
		}
	}

	return false
}

func mustacheTestFiles() []string {
	var result []string

	files, err := ioutil.ReadDir(path.Join("mustache", "specs"))
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fileName := file.Name()

		if !file.IsDir() && strings.HasSuffix(fileName, ".yml") {
			result = append(result, fileName)
		}
	}

	return result
}

//
// Following tests come fron ~lambdas.yml
//

var mustacheLambdasTests = []testcase{
	{
		template: "Hello, {{lambda}}!",
		data: map[string]interface{}{
			"lambda": func() string {
				return "world"
			}},
		expected: "Hello, world!",
	},

	// // SKIP: lambda return value is not parsed
	// {
	// 	"Interpolation - Expansion",
	// 	"Hello, {{lambda}}!",
	// 	map[string]interface{}{"lambda": func() string { return "{{planet}}" }},
	// 	nil, nil, nil,
	// 	"Hello, world!",
	// },

	// SKIP "Interpolation - Alternate Delimiters"

	{
		template: "{{lambda}} == {{{lambda}}} == {{lambda}}",
		data: map[string]interface{}{
			"lambda": func() string {
				musTestLambdaInterMult++
				return mario.Str(musTestLambdaInterMult)
			}},
		expected: "1 == 2 == 3",
	},

	{
		template: "<{{lambda}}{{{lambda}}}",
		data: map[string]interface{}{
			"lambda": func() string {
				return ">"
			}},
		expected: "<&gt;>",
	},

	// // SKIP: "Lambdas used for sections should receive the raw section string."
	// {
	// 	"Section",
	// 	"<{{#lambda}}{{x}}{{/lambda}}>",
	// 	map[string]interface{}{"lambda": func(param string) string {
	// 		if param == "{{x}}" {
	// 			return "yes"
	// 		}

	// 		return "false"
	// 	}, "x": "Error!"},
	// 	nil, nil, nil,
	// 	"<yes>",
	// },

	// // SKIP: lambda return value is not parsed
	// {
	// 	"Section - Expansion",
	// 	"<{{#lambda}}-{{/lambda}}>",
	// 	map[string]interface{}{"lambda": func(param string) string {
	// 		return param + "{{planet}}" + param
	// 	}, "planet": "Earth"},
	// 	nil, nil, nil,
	// 	"<-Earth->",
	// },

	// SKIP: "Section - Alternate Delimiters"

	{
		template: "{{#lambda}}FILE{{/lambda}} != {{#lambda}}LINE{{/lambda}}",
		data: map[string]interface{}{
			"lambda": func(options *mario.Options) string {
				return "__" + options.Fn() + "__"
			}},
		expected: "__FILE__ != __LINE__",
	},

	// // SKIP: "Lambdas used for inverted sections should be considered truthy."
	// {
	// 	"Inverted Section",
	// 	"<{{^lambda}}{{static}}{{/lambda}}>",
	// 	map[string]interface{}{
	// 		"lambda": func() interface{} {
	// 			return false
	// 		},
	// 		"static": "static",
	// 	},
	// 	nil, nil, nil,
	// 	"<>",
	// },
}

func TestMustacheLambdas(t *testing.T) {
	t.Parallel()

	for i, tt := range mustacheLambdasTests {
		tpl := mario.New()
		for name, fn := range tt.helpers {
			tpl.WithHelperFunc(name, fn)
		}
		for name, source := range tt.partials {
			tpl.WithPartial(name, mario.Must(mario.New().Parse(source)))
		}

		var b strings.Builder
		if err := mario.Must(tpl.Parse(tt.template)).Execute(&b, tt.data); err != nil {
			require.EqualError(t, err, tt.expectedError, i)
		} else {
			require.Equal(t, tt.expected, b.String(), i)
		}
	}

}
