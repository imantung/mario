package mario_test

import (
	"testing"

	"github.com/imantung/mario"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	source := "<h1>{{title}}</h1><p>{{body.content}}</p>"

	ctx := map[string]interface{}{
		"title": "foo",
		"body":  map[string]string{"content": "bar"},
	}

	// parse template
	tpl, err := mario.New().Parse(source)
	require.NoError(t, err)

	// evaluate template with context
	output, err := tpl.Exec(ctx)
	require.NoError(t, err)
	require.Equal(t, `<h1>foo</h1><p>bar</p>`, output)
}

func Test_struct(t *testing.T) {

	type Person struct {
		FirstName string
		LastName  string
	}

	type Comment struct {
		Author Person
		Body   string `handlebars:"content"`
	}

	type Post struct {
		Author   Person
		Body     string
		Comments []Comment
	}

	mario.RegisterHelper("fullName", func(person Person) string {
		return person.FirstName + " " + person.LastName
	})

	require.Equal(t,
		"<div class=\"post\">\n  <h1>By Jean Valjean</h1>\n  <div class=\"body\">Life is difficult</div>\n\n  <h1>Comments</h1>\n\n  <h2>By Marcel Beliveau</h2>\n  <div class=\"body\">LOL!</div>\n</div>",
		mario.MustRender(
			`<div class="post">
  <h1>By {{fullName author}}</h1>
  <div class="body">{{body}}</div>

  <h1>Comments</h1>

  {{#each comments}}
  <h2>By {{fullName author}}</h2>
  <div class="body">{{content}}</div>
  {{/each}}
</div>`,
			Post{
				Person{"Jean", "Valjean"},
				"Life is difficult",
				[]Comment{
					Comment{
						Person{"Marcel", "Beliveau"},
						"LOL!",
					},
				},
			},
		),
	)

}

func TestRender(t *testing.T) {
	// render template with context
	output, err := mario.Render(
		"<h1>{{title}}</h1><p>{{body.content}}</p>",
		map[string]interface{}{
			"title": "foo",
			"body":  map[string]string{"content": "bar"},
		})
	require.NoError(t, err)
	require.Equal(t, `<h1>foo</h1><p>bar</p>`, output)

}

func TestMustRender(t *testing.T) {
	require.Equal(t,
		`<h1>foo</h1><p>bar</p>`,
		mario.MustRender(
			"<h1>{{title}}</h1><p>{{body.content}}</p>",
			map[string]interface{}{
				"title": "foo",
				"body":  map[string]string{"content": "bar"},
			},
		),
	)
}
