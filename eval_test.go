package mario_test

import (
	"strings"
	"testing"

	"github.com/imantung/mario"
	"github.com/stretchr/testify/require"
)

func TestEval(t *testing.T) {
	testcases := []testcase{
		{
			template: "this is content",
			expected: "this is content",
		},
		{
			template: "{{#a}}{{one}}{{#b}}{{one}}{{two}}{{one}}{{/b}}{{/a}}",
			data: map[string]interface{}{
				"a": map[string]int{"one": 1},
				"b": map[string]int{"two": 2},
			},
			expected: "1121",
		},
		{
			template: "{{#foo as |bar|}}{{bar}}{{/foo}}{{bar}}",
			data: map[string]string{
				"foo": "baz",
				"bar": "bat",
			},
			expected: "bazbat",
		},
		{
			template: "{{#foo as |bar i|}}{{i}}.{{bar}} {{/foo}}",
			data: map[string][]string{
				"foo": {"baz", "bar", "bat"},
			},
			expected: "0.baz 1.bar 2.bat ",
		},
		{
			template: "{{#foos as |foo iFoo|}}{{#wats as |wat iWat|}}{{iFoo}}.{{iWat}}.{{foo}}-{{wat}} {{/wats}}{{/foos}}",
			data: map[string][]string{
				"foos": {"baz", "bar"},
				"wats": {"the", "phoque"},
			},
			expected: "0.0.baz-the 0.1.baz-phoque 1.0.bar-the 1.1.bar-phoque ",
		},
		{
			template: "{{#foo as |bar|}}{{bar.baz}}{{/foo}}",
			data: map[string]map[string]string{
				"foo": {"baz": "bat"},
			},
			expected: "bat",
		},
		{
			template: "{{#foo}}bar{{/foo}} baz",
			data:     map[string]interface{}{"foo": false},
			expected: " baz",
		},
		{
			template: "{{title}} - {{#bold}}{{body}}{{/bold}}",
			data: map[string]string{
				"title": "My new blog post",
				"body":  "I have so many things to say!",
			},
			helpers: map[string]interface{}{"bold": func(options *mario.Options) mario.SafeString {
				return mario.SafeString(`<div class="mybold">` + options.Fn() + "</div>")
			}},
			expected: `My new blog post - <div class="mybold">I have so many things to say!</div>`,
		},
		{
			template: "{{#if a}}A{{else if b}}B{{else}}C{{/if}}",
			data:     map[string]interface{}{"b": false},
			expected: "C",
		},
		{
			template: `{{foo "bar"}}`,
			data: map[string]interface{}{
				"foo": func(a string, b string) string {
					return "foo"
				},
			},
			expectedError: "Evaluation error: String{Value:'bar', Pos:7}: Helper 'foo' called with wrong number of arguments, needed 2 but got 1",
		},
		{
			template: "{{foo}}",
			data: map[string]interface{}{
				"foo": func() {},
			},
			expectedError: "Helper function must return a string or a SafeString: ",
		},
		{
			template: "{{foo}}",
			data: map[string]interface{}{
				"foo": func() (string, bool, string) {
					return "foo", true, "bar"
				},
			},
			expectedError: "Helper function must return a string or a SafeString: ",
		},

		// @todo Test with a "../../path" (depth 2 path) while context is only depth 1
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

func TestEvalStruct(t *testing.T) {
	t.Parallel()

	source := `<div class="post">
  <h1>By {{author.FirstName}} {{Author.lastName}}</h1>
  <div class="body">{{Body}}</div>

  <h1>Comments</h1>

  {{#each comments}}
  <h2>By {{Author.FirstName}} {{author.LastName}}</h2>
  <div class="body">{{body}}</div>
  {{/each}}
</div>`

	expected := `<div class="post">
  <h1>By Jean Valjean</h1>
  <div class="body">Life is difficult</div>

  <h1>Comments</h1>

  <h2>By Marcel Beliveau</h2>
  <div class="body">LOL!</div>
</div>`

	type Person struct {
		FirstName string
		LastName  string
	}

	type Comment struct {
		Author Person
		Body   string
	}

	type Post struct {
		Author   Person
		Body     string
		Comments []Comment
	}

	data := Post{
		Person{"Jean", "Valjean"},
		"Life is difficult",
		[]Comment{
			Comment{
				Person{"Marcel", "Beliveau"},
				"LOL!",
			},
		},
	}

	require.Equal(t, expected, compile(source, data))
}

func TestEvalStructTag(t *testing.T) {
	t.Parallel()

	source := `<div class="person">
	<h1>{{real-name}}</h1>
	<ul>
	  <li>City: {{info.location}}</li>
	  <li>Rug: {{info.[r.u.g]}}</li>
	  <li>Activity: {{info.activity}}</li>
	</ul>
	{{#each other-names}}
	<p>{{alias-name}}</p>
	{{/each}}
</div>`

	expected := `<div class="person">
	<h1>Lebowski</h1>
	<ul>
	  <li>City: Venice</li>
	  <li>Rug: Tied The Room Together</li>
	  <li>Activity: Bowling</li>
	</ul>
	<p>his dudeness</p>
	<p>el duderino</p>
</div>`

	type Alias struct {
		Name string `handlebars:"alias-name"`
	}

	type CharacterInfo struct {
		City     string `handlebars:"location"`
		Rug      string `handlebars:"r.u.g"`
		Activity string `handlebars:"not-activity"`
	}

	type Character struct {
		RealName string `handlebars:"real-name"`
		Info     CharacterInfo
		Aliases  []Alias `handlebars:"other-names"`
	}

	data := Character{
		"Lebowski",
		CharacterInfo{"Venice", "Tied The Room Together", "Bowling"},
		[]Alias{
			{"his dudeness"},
			{"el duderino"},
		},
	}

	require.Equal(t, expected, compile(source, data))
}

type TestFoo struct {
}

func (t *TestFoo) Subject() string {
	return "foo"
}

func TestEvalMethod(t *testing.T) {
	t.Parallel()

	source := `Subject is {{subject}}! YES I SAID {{Subject}}!`
	expected := `Subject is foo! YES I SAID foo!`
	data := &TestFoo{}

	require.Equal(t, expected, compile(source, data))
}

type TestBar struct {
}

func (t *TestBar) Subject() interface{} {
	return testBar
}

func testBar() string {
	return "bar"
}

func TestEvalMethodReturningFunc(t *testing.T) {
	t.Parallel()

	source := `Subject is {{subject}}! YES I SAID {{Subject}}!`
	expected := `Subject is bar! YES I SAID bar!`
	data := &TestBar{}

	require.Equal(t, expected, compile(source, data))

}
