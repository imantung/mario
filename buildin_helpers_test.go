package mario_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/imantung/mario"
)

const (
	VERBOSE = false
)

//
// Helpers
//

func barHelper(options *mario.Options) string { return "bar" }

func echoHelper(str string, nb int) string {
	result := ""
	for i := 0; i < nb; i++ {
		result += str
	}

	return result
}

func boolHelper(b bool) string {
	if b {
		return "yes it is"
	}

	return "absolutely not"
}

func gnakHelper(nb int) string {
	result := ""
	for i := 0; i < nb; i++ {
		result += "GnAK!"
	}

	return result
}

func TestHelper(t *testing.T) {
	testcases := []testcase{
		{
			template: `{{foo}}`,
			helpers:  map[string]interface{}{"foo": barHelper},
			expected: `bar`,
		},
		{
			template: `{{echo "foo" 1}}`,
			helpers:  map[string]interface{}{"echo": echoHelper},
			expected: `foo`,
		},
		{
			template: `{{echo foo 1}}`,
			data:     map[string]interface{}{"foo": "bar"},
			helpers:  map[string]interface{}{"echo": echoHelper},
			expected: `bar`,
		},
		{
			template: `{{bool true}}`,
			helpers:  map[string]interface{}{"bool": boolHelper},
			expected: `yes it is`,
		},
		{
			template: `{{bool false}}`,
			helpers:  map[string]interface{}{"bool": boolHelper},
			expected: `absolutely not`,
		},
		{
			template: `{{gnak 5}}`,
			helpers:  map[string]interface{}{"gnak": gnakHelper},
			expected: `GnAK!GnAK!GnAK!GnAK!GnAK!`,
		},
		{
			template: `{{echo "GnAK!" 3}}`,
			helpers:  map[string]interface{}{"echo": echoHelper},
			expected: `GnAK!GnAK!GnAK!`,
		},
		{
			template: `{{#if true}}YES MAN{{/if}}`,
			expected: `YES MAN`,
		},
		{
			template: `{{#if false}}YES MAN{{/if}}`,
			expected: ``,
		},
		{
			template: `{{#if ok}}YES MAN{{/if}}`,
			data:     map[string]interface{}{"ok": true},
			expected: `YES MAN`,
		},
		{
			template: `{{#if ok}}YES MAN{{/if}}`,
			data:     map[string]interface{}{"ok": false},
			expected: ``,
		},
		{
			template: `{{#unless true}}YES MAN{{/unless}}`,
			expected: ``,
		},
		{
			template: `{{#unless false}}YES MAN{{/unless}}`,
			expected: `YES MAN`,
		},
		{
			template: `{{#unless ok}}YES MAN{{/unless}}`,
			data:     map[string]interface{}{"ok": true},
			expected: ``,
		},
		{
			template: `{{#unless ok}}YES MAN{{/unless}}`,
			data:     map[string]interface{}{"ok": false},
			expected: `YES MAN`,
		},
		{
			template: `{{#equal foo "bar"}}YES MAN{{/equal}}`,
			data:     map[string]interface{}{"foo": "bar"},
			expected: `YES MAN`,
		},
		{
			template: `{{#equal foo "baz"}}YES MAN{{/equal}}`,
			data:     map[string]interface{}{"foo": "bar"},
			expected: ``,
		},
		{
			template: `{{#equal foo bar}}YES MAN{{/equal}}`,
			data:     map[string]interface{}{"foo": "baz", "bar": "baz"},
			expected: `YES MAN`,
		},
		{
			template: `{{#equal foo bar}}YES MAN{{/equal}}`,
			data:     map[string]interface{}{"foo": "baz", "bar": "tag"},
			expected: ``,
		},
		{
			template: `{{#equal foo 1}}YES MAN{{/equal}}`,
			data:     map[string]interface{}{"foo": 1},
			expected: `YES MAN`,
		},
		{
			template: `{{#equal foo 0}}YES MAN{{/equal}}`,
			data:     map[string]interface{}{"foo": 1},
			expected: ``,
		},
		{
			template: `<option value="test" {{#equal value "test"}}selected{{/equal}}>Test</option>`,
			data:     map[string]interface{}{"value": "test"},
			expected: `<option value="test" selected>Test</option>`,
		},
		{
			template: `{{#equal foo "bar"}}foo is bar{{/equal}}
				{{#equal foo baz}}foo is the same as baz{{/equal}}
				{{#equal nb 0}}nothing{{/equal}}
				{{#equal nb 1}}there is one{{/equal}}
				{{#equal nb "1"}}everything is stringified before comparison{{/equal}}`,
			data: map[string]interface{}{
				"foo": "bar",
				"baz": "bar",
				"nb":  1,
			},
			expected: `foo is bar
				foo is the same as baz
				
				there is one
				everything is stringified before comparison`,
		},
	}

	t.Parallel()

	for i, tt := range testcases {
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

//
// Fixes: https://github.com/aymerick/raymond/issues/2
//

type Author struct {
	FirstName string
	LastName  string
}

func TestHelperCtx(t *testing.T) {
	templateHelper := func(name string, options *mario.Options) mario.SafeString {
		ctx := options.Ctx()
		template := name + " - {{ firstName }} {{ lastName }}"
		return mario.SafeString(compile(template, ctx))
	}

	tpl, err := mario.New().
		WithHelperFunc("template", templateHelper).
		Parse(`By {{ template "namefile" }}`)
	require.NoError(t, err)

	var b strings.Builder
	require.NoError(t, tpl.Execute(&b, Author{"Alan", "Johnson"}))
	require.Equal(t, "By namefile - Alan Johnson", b.String())
}
