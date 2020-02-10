# Mario 

[![Build Status](https://secure.travis-ci.org/aymerick/mario.svg?branch=master)](http://travis-ci.org/aymerick/mario) [![GoDoc](https://godoc.org/github.com/imantung/mario?status.svg)](http://godoc.org/github.com/imantung/mario)


Continues work of [raymond](https://github.com/aymerick/raymond) to support [Handlebars](https://handlebarsjs.com/) for [golang](https://golang.org) 



<img src="mario.jpeg?raw=true">

## Handlebars 

Handlebars is a superset of [mustache](https://mustache.github.io) but it differs on those points:

- Alternative delimiters are not supported
- There is no recursive lookup

## Install

```bash
go get -u github.com/imantung/mario
```

## Usage

Mario's function implementation mimic [go template](https://golang.org/pkg/text/template/) for convenience.

```go
source := `Hello {{name}}`

ctx := map[string]string{
  "name": "World",
}

tpl, err := mario.New().Parse(source)
if err != nil {
  panic(err)
}

var b strings.Builder
if err := tpl.Execute(&b, ctx); err != nil {
  panic(err)
}

fmt.Println(b.String())

// Output: 
// Hello World
```

For shortcut, using `Must` to create template object without error
```go
var t = mario.Must(mario.New().Parse(source))
```

## Rendering Context

The rendering context can contain any type of values, including `array`, `slice`, `map`, `struct` and `func`.

### Map

```go
source := `{ "name": "{{name}}", "age": {{age}} }`
data := map[string]interface{}{
	"name": "Mario",
	"age":  12,
}
```


### Struct

```go
source := `<div class="post">
  <h1>By {{author.firstName}} {{author.lastName}}</h1>
  <div class="body">{{body}}</div>

  <h1>Comments</h1>

  {{#each comments}}
  <h2>By {{author.firstName}} {{author.lastName}}</h2>
  <div class="body">{{content}}</div>
  {{/each}}
</div>`

ctx := Post{
  Person{"Jean", "Valjean"},
  "Life is difficult",
  []Comment{
    Comment{
      Person{"Marcel", "Beliveau"},
      "LOL!",
    },
  },
}
```

Note:
- When using structs, be warned that only exported fields are accessible. 
  ```go
  type Person struct {
    FirstName  string
    LastName   string
    motherName string // NOTE: Unexported and can't be access
  }
  ```

- However you can access exported fields in template with their lowercase names. For example, both `{{author.firstName}}` and `{{Author.FirstName}}` references give the same result, as long as `Author` and `FirstName` are exported struct fields.

- More, you can use the `handlebars` struct tag to specify a template variable name different from the struct field name.

### Function

Randomly renders `I hate you` or `I love you`.
```go
source := "I {{feeling}} you"

ctx := map[string]interface{}{
    "feeling": func() string {
        rand.Seed(time.Now().UTC().UnixNano())

        feelings := []string{"hate", "love"}
        return feelings[rand.Intn(len(feelings))]
    },
}
```

Note:
- Those context functions behave like helper functions: they can be called with parameters and they can have an `Options` argument.



## Helpers

Handlebarjs [built-in helpers](https://handlebarsjs.com/guide/builtin-helpers.html):

- `if` to conditionally render a block. 

  ```html
  <div class="entry">
    {{#if author}}
      <h1>{{firstName}} {{lastName}}</h1>
    {{/if}}
  </div>
  ```

- `unless` inverse of the `if` helper. 

  ```html
  <div class="entry">
    {{#unless license}}
    <h3 class="warning">WARNING: This entry does not have a license!</h3>
    {{/unless}}
  </div>
  ```

- `each`: iterate over an array, a slice, a map or a struct instance using this built-in `each` helper. Inside the block, you can use `this` to reference the element being iterated over.

  ```html
  <ul class="people">
    {{#each people}}
      <li>{{this}}</li>
    {{/each}}
  </ul>
  ```

- `with`: shift the context for a section of a template by using the built-in `with` block helper.

  ```html
  <div class="entry">
    <h1>{{title}}</h1>

    {{#with author}}
    <h2>By {{firstName}} {{lastName}}</h2>
    {{/with}}
  </div>
  ```

- `lookup`: allows for dynamic parameter resolution using handlebars variables.

  ```html
  {{#each bar}}
    {{lookup ../foo @index}}
  {{/each}}
  ```

- `log`: allows for logging while rendering a template. 

  ```html
  {{log "Look at me!"}}
  ```

Additional helper:
- `equal`: renders a block if the string version of both arguments are equals.

  ```html
  {{#equal foo "bar"}}foo is bar{{/equal}}
  {{#equal foo baz}}foo is the same as baz{{/equal}}
  {{#equal nb 0}}nothing{{/equal}}
  {{#equal nb 1}}there is one{{/equal}}
  {{#equal nb "1"}}everything is stringified before comparison{{/equal}}
  ```

## Custom Helper

_TODO: Implementation of custom helper_

## Language Features

_TODO: <https://handlebarsjs.com/guide/#language-features>_

## Limitations

These handlebars options are currently NOT implemented:

- `compat` - enables recursive field lookup
- `knownHelpers` - list of helpers that are known to exist (truthy) at template execution time
- `knownHelpersOnly` - allows further optimizations based on the known helpers list
- `trackIds` - include the id names used to resolve parameters for helpers
- `noEscape` - disables HTML escaping globally
- `strict` - templates will throw rather than silently ignore missing fields
- `assumeObjects` - removes object existence checks when traversing paths
- `preventIndent` - disables the auto-indententation of nested partials
- `stringParams` - resolves a parameter to it's name if the value isn't present in the context stack

These handlebars features are currently NOT implemented:

- raw block content is not passed as a parameter to helper
- `blockHelperMissing` - helper called when a helper can not be directly resolved
- `helperMissing` - helper called when a potential helper expression was not found
- `@contextPath` - value set in `trackIds` mode that records the lookup path for the current context
- `@level` - log level


## References

  - <http://handlebarsjs.com/>
  - <https://mustache.github.io/mustache.5.html>
  - <https://github.com/golang/go/tree/master/src/text/template>
  - <https://www.youtube.com/watch?v=HxaD_trXwRE>
  - <https://github.com/aymerick/raymond>


## Others Implementations

- [raymond](https://github.com/aymerick/raymond) - golang (OG)
- [handlebars.js](http://handlebarsjs.com) - javascript
- [handlebars.java](https://github.com/jknack/handlebars.java) - java
- [handlebars.rb](https://github.com/cowboyd/handlebars.rb) - ruby
- [handlebars.php](https://github.com/XaminProject/handlebars.php) - php
- [handlebars-objc](https://github.com/Bertrand/handlebars-objc) - Objective C
- [rumblebars](https://github.com/nicolas-cherel/rumblebars) - rust
