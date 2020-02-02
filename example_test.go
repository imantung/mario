package mario_test

import (
	"fmt"
	"strings"

	"github.com/imantung/mario"
)

func Example() {
	source := `Hello {{name}}`
	data := map[string]string{
		"name": "World",
	}

	tpl, err := mario.New().Parse(source)
	if err != nil {
		panic(err)
	}

	var b strings.Builder
	if err := tpl.Execute(&b, data); err != nil {
		panic(err)
	}

	fmt.Println(b.String())

	// Output:
	// Hello World
}

func ExampleMust() {
	source := `Hello {{name}}`
	data := map[string]string{
		"name": "World",
	}

	var b strings.Builder
	if err := mario.Must(mario.New().Parse(source)).Execute(&b, data); err != nil {
		panic(err)
	}

	fmt.Println(b.String())

	// Output:
	// Hello World
}

func Example_html_escape() {
	source := "{{text}}\n{{{text}}}"

	data := map[string]string{
		"text": "This is <html> Tags",
	}

	var b strings.Builder
	if err := mario.Must(mario.New().Parse(source)).Execute(&b, data); err != nil {
		panic(err)
	}

	fmt.Println(b.String())

	// Output:
	// This is &lt;html&gt; Tags
	// This is <html> Tags

}

func Example_rendering_context() {
	source := `<div class="post">
  <h1>By {{author.firstName}} {{author.lastName}}</h1>
  <div class="body">{{body}}</div>

  <h1>Comments</h1>

  {{#each comments}}
  <h2>By {{author.firstName}} {{author.lastName}}</h2>
  <div class="body">{{content}}</div>
  {{/each}}
</div>`

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

	var b strings.Builder
	if err := mario.Must(mario.New().Parse(source)).Execute(&b, data); err != nil {
		panic(err)
	}

	fmt.Println(b.String())

	// Output:
	// <div class="post">
	//   <h1>By Jean Valjean</h1>
	//   <div class="body">Life is difficult</div>
	//
	//   <h1>Comments</h1>
	//
	//   <h2>By Marcel Beliveau</h2>
	//   <div class="body">LOL!</div>
	// </div>
}

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

func ExampleHelper_safe_string() {
	tpl, err := mario.New().
		WithHelperFunc("link", func(url, text string) mario.SafeString {
			return mario.SafeString("<a href='" + mario.Escape(url) + "'>" + mario.Escape(text) + "</a>")
		}).
		Parse("{{link url text}}")
	if err != nil {
		panic(err)
	}

	data := map[string]string{
		"url":  "http://www.aymerick.com/",
		"text": "This is a <em>cool</em> website",
	}

	var b strings.Builder
	if err := tpl.Execute(&b, data); err != nil {
		panic(err)
	}

	fmt.Println(b.String())

	// Output:
	// <a href='http://www.aymerick.com/'>This is a &lt;em&gt;cool&lt;/em&gt; website</a>
}
