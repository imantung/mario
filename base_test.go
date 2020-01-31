package mario_test

import (
	"fmt"

	"regexp"
	"testing"

	"github.com/imantung/mario"
	"github.com/imantung/mario/ast"
)

type Test struct {
	name     string
	input    string
	data     interface{}
	privData map[string]interface{}
	helpers  map[string]interface{}
	partials map[string]string
	output   interface{}
}

func launchTests(t *testing.T, tests []Test) {
	// NOTE: TestMustache() makes Parallel testing fail
	// t.Parallel()

	for _, test := range tests {
		var (
			err error
			tpl *mario.Template
		)

		// parse template
		tpl, err = mario.New().Parse(test.input)
		if err != nil {
			t.Errorf("Test '%s' failed - Failed to parse template\ninput:\n\t'%s'\nerror:\n\t%s", test.name, test.input, err)
		} else {
			if len(test.helpers) > 0 {
				// register helpers
				tpl.RegisterHelpers(test.helpers)
			}

			if len(test.partials) > 0 {
				// register partials
				tpl.RegisterPartials(test.partials)
			}

			// setup private data frame
			var privData *mario.DataFrame
			if test.privData != nil {
				privData = mario.NewDataFrame()
				for k, v := range test.privData {
					privData.Set(k, v)
				}
			}

			// render template
			output, err := tpl.ExecuteWith(test.data, privData)
			if err != nil {
				t.Errorf("Test '%s' failed\ninput:\n\t'%s'\ndata:\n\t%s\nerror:\n\t%s\nAST:\n\t%s", test.name, test.input, mario.Str(test.data), err, ast.Print(tpl.Program()))
			} else {
				// check output
				var expectedArr []string
				expectedArr, ok := test.output.([]string)
				if ok {
					match := false
					for _, expectedStr := range expectedArr {
						if expectedStr == output {
							match = true
							break
						}
					}

					if !match {
						t.Errorf("Test '%s' failed\ninput:\n\t'%s'\ndata:\n\t%s\npartials:\n\t%s\nexpected\n\t%q\ngot\n\t%q\nAST:\n%s", test.name, test.input, mario.Str(test.data), mario.Str(test.partials), expectedArr, output, ast.Print(tpl.Program()))
					}
				} else {
					expectedStr, ok := test.output.(string)
					if !ok {
						panic(fmt.Errorf("Erroneous test output description: %q", test.output))
					}

					if expectedStr != output {
						t.Errorf("Test '%s' failed\ninput:\n\t'%s'\ndata:\n\t%s\npartials:\n\t%s\nexpected\n\t%q\ngot\n\t%q\nAST:\n%s", test.name, test.input, mario.Str(test.data), mario.Str(test.partials), expectedStr, output, ast.Print(tpl.Program()))
					}
				}
			}
		}
	}
}

func launchErrorTests(t *testing.T, tests []Test) {
	t.Parallel()

	for _, test := range tests {
		var (
			err error
			tpl *mario.Template
		)

		// parse template
		tpl, err = mario.New().Parse(test.input)
		if err != nil {
			t.Errorf("Test '%s' failed - Failed to parse template\ninput:\n\t'%s'\nerror:\n\t%s", test.name, test.input, err)
		} else {
			if len(test.helpers) > 0 {
				// register helpers
				tpl.RegisterHelpers(test.helpers)
			}

			if len(test.partials) > 0 {
				// register partials
				tpl.RegisterPartials(test.partials)
			}

			// setup private data frame
			var privData *mario.DataFrame
			if test.privData != nil {
				privData := mario.NewDataFrame()
				for k, v := range test.privData {
					privData.Set(k, v)
				}
			}

			// render template
			output, err := tpl.ExecuteWith(test.data, privData)
			if err == nil {
				t.Errorf("Test '%s' failed - Error expected\ninput:\n\t'%s'\ngot\n\t%q\nAST:\n%q", test.name, test.input, output, ast.Print(tpl.Program()))
			} else {
				var errMatch error
				match := false

				// check output
				var expectedArr []string
				expectedArr, ok := test.output.([]string)
				if ok {
					if len(expectedArr) > 0 {
						for _, expectedStr := range expectedArr {
							match, errMatch = regexp.MatchString(regexp.QuoteMeta(expectedStr), fmt.Sprint(err))
							if errMatch != nil {
								panic("Failed to match regexp")
							}

							if match {
								break
							}
						}
					} else {
						// nothing to test
						match = true
					}
				} else {
					expectedStr, ok := test.output.(string)
					if !ok {
						panic(fmt.Errorf("Erroneous test output description: %q", test.output))
					}

					if expectedStr != "" {
						match, errMatch = regexp.MatchString(regexp.QuoteMeta(expectedStr), fmt.Sprint(err))
						if errMatch != nil {
							panic("Failed to match regexp")
						}
					} else {
						// nothing to test
						match = true
					}
				}

				if !match {
					t.Errorf("Test '%s' failed - Incorrect error returned\ninput:\n\t'%s'\ndata:\n\t%s\nexpected\n\t%q\ngot\n\t%q", test.name, test.input, mario.Str(test.data), test.output, err)
				}
			}
		}
	}
}
