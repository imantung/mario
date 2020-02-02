# Mario 

[![Build Status](https://secure.travis-ci.org/aymerick/mario.svg?branch=master)](http://travis-ci.org/aymerick/mario) [![GoDoc](https://godoc.org/github.com/imantung/mario?status.svg)](http://godoc.org/github.com/imantung/mario)


Continues work of [raymond](https://github.com/aymerick/raymond) to support [Handlebars](https://handlebarsjs.com/) for [golang](https://golang.org) 



<img src="mario.jpeg?raw=true">

## Handlebars 

Handlebars is a superset of [mustache](https://mustache.github.io) but it differs on those points:

- Alternative delimiters are not supported
- There is no recursive lookup



## Installation

```bash
go get -u github.com/imantung/mario
```

## Usages

Mario's function implement mimic [go template](https://golang.org/pkg/text/template/) for convenience.

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

var b strings.Builder
if err := mario.Must(mario.New().Parse(source)).Execute(&b, ctx); err != nil {
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
```

When using structs, be warned that only exported fields are accessible. 
```go
type Person struct {
	FirstName  string
	LastName   string
	motherName string // NOTE: Unexported and can't be access
}
```

However you can access exported fields in template with their lowercase names. For example, both `{{author.firstName}}` and `{{Author.FirstName}}` references give the same result, as long as `Author` and `FirstName` are exported struct fields.

More, you can use the `handlebars` struct tag to specify a template variable name different from the struct field name.

### Function

In addition to helpers, lambdas found in context are evaluated.

For example, that template and context:

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

Randomly renders `I hate you` or `I love you`.

Those context functions behave like helper functions: they can be called with parameters and they can have an `Options` argument.


## HTML Escaping

By default, the result of a mustache expression is HTML escaped. Use the triple mustache `{{{` to output unescaped values.

```go
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
```

For custom helper implementation, you should return a `SafeString` if you don't want it to be escaped by default. When using `SafeString` all unknown or unsafe data should be manually escaped with the `Escape` method.

```go
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
```


## Helpers

### Built-In

Those built-in helpers are available to all templates.


#### The `if` block helper

You can use the `if` helper to conditionally render a block. If its argument returns `false`, `nil`, `0`, `""`, an empty array, an empty slice or an empty map, then mario will not render the block.

```html
<div class="entry">
  {{#if author}}
    <h1>{{firstName}} {{lastName}}</h1>
  {{/if}}
</div>
```

When using a block expression, you can specify a template section to run if the expression returns a falsy value. That section, marked by `{{else}}` is called an "else section".

```html
<div class="entry">
  {{#if author}}
    <h1>{{firstName}} {{lastName}}</h1>
  {{else}}
    <h1>Unknown Author</h1>
  {{/if}}
</div>
```

You can chain several blocks. For example that template:

```html
{{#if isActive}}
  <img src="star.gif" alt="Active">
{{else if isInactive}}
  <img src="cry.gif" alt="Inactive">
{{else}}
  <img src="wat.gif" alt="Unknown">
{{/if}}
```


#### The `unless` block helper

You can use the `unless` helper as the inverse of the `if` helper. Its block will be rendered if the expression returns a falsy value.

```html
<div class="entry">
  {{#unless license}}
  <h3 class="warning">WARNING: This entry does not have a license!</h3>
  {{/unless}}
</div>
```


#### The `each` block helper

You can iterate over an array, a slice, a map or a struct instance using this built-in `each` helper. Inside the block, you can use `this` to reference the element being iterated over.

For example:

```html
<ul class="people">
  {{#each people}}
    <li>{{this}}</li>
  {{/each}}
</ul>
```

With this context:

```go
map[string]interface{}{
    "people": []string{
        "Marcel", "Jean-Claude", "Yvette",
    },
}
```

Outputs:

```html
<ul class="people">
  <li>Marcel</li>
  <li>Jean-Claude</li>
  <li>Yvette</li>
</ul>
```

You can optionally provide an `{{else}}` section which will display only when the passed argument is an empty array, an empty slice or an empty map (a `struct` instance is never considered empty).

```html
{{#each paragraphs}}
  <p>{{this}}</p>
{{else}}
  <p class="empty">No content</p>
{{/each}}
```

When looping through items in `each`, you can optionally reference the current loop index via `{{@index}}`.

```html
{{#each array}}
  {{@index}}: {{this}}
{{/each}}
```

Additionally for map and struct instance iteration, `{{@key}}` references the current map key or struct field name:

```html
{{#each map}}
  {{@key}}: {{this}}
{{/each}}
```

The first and last steps of iteration are noted via the `@first` and `@last` variables.


#### The `with` block helper

You can shift the context for a section of a template by using the built-in `with` block helper.

```html
<div class="entry">
  <h1>{{title}}</h1>

  {{#with author}}
  <h2>By {{firstName}} {{lastName}}</h2>
  {{/with}}
</div>
```

With this context:

```go
map[string]interface{}{
    "title": "My first post!",
    "author": map[string]string{
        "firstName": "Jean",
        "lastName":  "Valjean",
    },
}
```

Outputs:

```html
<div class="entry">
  <h1>My first post!</h1>

  <h2>By Jean Valjean</h2>
</div>
```

You can optionally provide an `{{else}}` section which will display only when the passed argument is falsy.

```html
{{#with author}}
  <p>{{name}}</p>
{{else}}
  <p class="empty">No content</p>
{{/with}}
```


#### The `lookup` helper

The `lookup` helper allows for dynamic parameter resolution using handlebars variables.

```html
{{#each bar}}
  {{lookup ../foo @index}}
{{/each}}
```


#### The `log` helper

The `log` helper allows for logging while rendering a template.

```html
{{log "Look at me!"}}
```

Note that the handlebars.js `@level` variable is not supported.


#### The `equal` helper

The `equal` helper renders a block if the string version of both arguments are equals.

For example that template:

```html
{{#equal foo "bar"}}foo is bar{{/equal}}
{{#equal foo baz}}foo is the same as baz{{/equal}}
{{#equal nb 0}}nothing{{/equal}}
{{#equal nb 1}}there is one{{/equal}}
{{#equal nb "1"}}everything is stringified before comparison{{/equal}}
```

With that context:

```go
ctx := map[string]interface{}{
    "foo": "bar",
    "baz": "bar",
    "nb":  1,
}
```

Outputs:

```html
foo is bar
foo is the same as baz

there is one
everything is stringified before comparison
```


### Block Helpers

Block helpers make it possible to define custom iterators and other functionality that can invoke the passed block with a new context.


#### Block Evaluation

As an example, let's define a block helper that adds some markup to the wrapped text.

```html
<div class="entry">
  <h1>{{title}}</h1>
  <div class="body">
    {{#bold}}{{body}}{{/bold}}
  </div>
</div>
```

The `bold` helper will add markup to make its text bold.

```go
mario.RegisterHelper("bold", func(options *mario.Options) mario.SafeString {
    return mario.SafeString(`<div class="mybold">` + options.Fn() + "</div>")
})
```

A helper evaluates the block content with current context by calling `options.Fn()`.

If you want to evaluate the block with another context, then use `options.FnWith(ctx)`, like this french version of built-in `with` block helper:

```go
mario.RegisterHelper("avec", func(context interface{}, options *mario.Options) string {
    return options.FnWith(context)
})
```

With that template:

```html
{{#avec obj.text}}{{this}}{{/avec}}
```


#### Conditional

Let's write a french version of `if` block helper:

```go
source := `{{#si yep}}YEP !{{/si}}`

ctx := map[string]interface{}{"yep": true}

mario.RegisterHelper("si", func(conditional bool, options *mario.Options) string {
    if conditional {
        return options.Fn()
    }
    return ""
})
```

Note that as the first parameter of the helper is typed as `bool` an automatic conversion is made if corresponding context value is not a boolean. So this helper works with that context too:

```go
ctx := map[string]interface{}{"yep": "message"}
```

Here, `"message"` is converted to `true` because it is an non-empty string. See `IsTrue()` function for more informations on boolean conversion.


#### Else Block Evaluation

We can enhance the `si` block helper to evaluate the `else block` by calling `options.Inverse()` if conditional is false:

```go
source := `{{#si yep}}YEP !{{else}}NOP !{{/si}}`

ctx := map[string]interface{}{"yep": false}

mario.RegisterHelper("si", func(conditional bool, options *mario.Options) string {
    if conditional {
        return options.Fn()
    }
    return options.Inverse()
})
```

Outputs:
```
NOP !
```


#### Block Parameters

It's possible to receive named parameters from supporting helpers.

```html
{{#each users as |user userId|}}
  Id: {{userId}} Name: {{user.name}}
{{/each}}
```

In this particular example, `user` will have the same value as the current context and `userId` will have the index/key value for the iteration.

This allows for nested helpers to avoid name conflicts.

For example:

```html
{{#each users as |user userId|}}
  {{#each user.books as |book bookId|}}
    User: {{userId}} Book: {{bookId}}
  {{/each}}
{{/each}}
```

With this context:

```go
ctx := map[string]interface{}{
    "users": map[string]interface{}{
        "marcel": map[string]interface{}{
            "books": map[string]interface{}{
                "book1": "My first book",
                "book2": "My second book",
            },
        },
        "didier": map[string]interface{}{
            "books": map[string]interface{}{
                "bookA": "Good book",
                "bookB": "Bad book",
            },
        },
    },
}
```

Outputs:

```html
  User: marcel Book: book1
  User: marcel Book: book2
  User: didier Book: bookA
  User: didier Book: bookB
```

As you can see, the second block parameter is the map key. When using structs, it is the struct field name.

When using arrays and slices, the second parameter is element index:

```go
ctx := map[string]interface{}{
    "users": []map[string]interface{}{
        {
            "id": "marcel",
            "books": []map[string]interface{}{
                {"id": "book1", "title": "My first book"},
                {"id": "book2", "title": "My second book"},
            },
        },
        {
            "id": "didier",
            "books": []map[string]interface{}{
                {"id": "bookA", "title": "Good book"},
                {"id": "bookB", "title": "Bad book"},
            },
        },
    },
}
```

Outputs:

```html
    User: 0 Book: 0
    User: 0 Book: 1
    User: 1 Book: 0
    User: 1 Book: 1
```


### Helper Parameters

When calling a helper in a template, mario expects the same number of arguments as the number of helper function parameters.

So this template:

```html
{{add a}}
```

With this helper:

```go
mario.RegisterHelper("add", func(val1, val2 int) string {
    return strconv.Itoa(val1 + val2)
})
```

Will simply panics, because we call the helper with one argument whereas it expects two.


#### Automatic conversion

Let's create a `concat` helper that expects two strings and concat them:

```go
source := `{{concat a b}}`

ctx := map[string]interface{}{
    "a": "Jean",
    "b": "Valjean",
}

mario.RegisterHelper("concat", func(val1, val2 string) string {
    return val1 + " " + val2
})
```

Everything goes well, two strings are passed as arguments to the helper that outputs:

```html
Jean VALJEAN
```

But what happens if there is another type than `string` in the context ? For example:

```go
ctx := map[string]interface{}{
    "a": 10,
    "b": "Valjean",
}
```

Actually, mario perfoms automatic string conversion. So because the first parameter of the helper is typed as `string`, the first argument will be converted from the `10` integer to `"10"`, and the helper outputs:

```html
10 VALJEAN
```

Note that this kind of automatic conversion is done with `bool` type too, thanks to the `IsTrue()` function.


### Options Argument

If a helper needs the `Options` argument, just add it at the end of helper parameters:

```go
mario.RegisterHelper("add", func(val1, val2 int, options *mario.Options) string {
    return strconv.Itoa(val1 + val2) + " " + options.ValueStr("bananas")
})
```

Thanks to the `options` argument, helpers have access to the current evaluation context, to the `Hash` arguments, and they can manipulate the private data variables.

The `Options` argument is even necessary for Block Helpers to evaluate block and "else block".


#### Context Values

Helpers fetch current context values with `options.Value()` and `options.ValuesStr()`.

`Value()` returns an `interface{}` and lets the helper do the type assertions whereas `ValueStr()` automatically converts the value to a `string`.

For example:

```go
source := `{{concat a b}}`

ctx := map[string]interface{}{
    "a":      "Marcel",
    "b":      "Beliveau",
    "suffix": "FOREVER !",
}

mario.RegisterHelper("concat", func(val1, val2 string, options *mario.Options) string {
    return val1 + " " + val2 + " " + options.ValueStr("suffix")
})
```

Outputs:

```html
Marcel Beliveau FOREVER !
```

Helpers can get the entire current context with `options.Ctx()` that returns an `interface{}`.


#### Helper Hash Arguments

Helpers access hash arguments with `options.HashProp()` and `options.HashStr()`.

`HashProp()` returns an `interface{}` and lets the helper do the type assertions whereas `HashStr()` automatically converts the value to a `string`.

For example:

```go
source := `{{concat suffix first=a second=b}}`

ctx := map[string]interface{}{
    "a":      "Marcel",
    "b":      "Beliveau",
    "suffix": "FOREVER !",
}

mario.RegisterHelper("concat", func(suffix string, options *mario.Options) string {
    return options.HashStr("first") + " " + options.HashStr("second") + " " + suffix
})
```

Outputs:

```html
Marcel Beliveau FOREVER !
```

Helpers can get the full hash with `options.Hash()` that returns a `map[string]interface{}`.


#### Private Data

Helpers access private data variables with `options.Data()` and `options.DataStr()`.

`Data()` returns an `interface{}` and lets the helper do the type assertions whereas `DataStr()` automatically converts the value to a `string`.

Helpers can get the entire current data frame with `options.DataFrame()` that returns a `*DataFrame`.

For helpers that need to inject their own private data frame, use `options.NewDataFrame()` to create the frame and `options.FnData()` to evaluate the block with that frame.

For example:

```go
source := `{{#voodoo kind=a}}Voodoo is {{@magix}}{{/voodoo}}`

ctx := map[string]interface{}{
    "a": "awesome",
}

mario.RegisterHelper("voodoo", func(options *mario.Options) string {
    // create data frame with @magix data
    frame := options.NewDataFrame()
    frame.Set("magix", options.HashProp("kind"))

    // evaluates block with new data frame
    return options.FnData(frame)
})
```

Helpers that need to evaluate the block with a private data frame and a new context can call `options.FnCtxData()`.



## Partials

### Template Partials

You can register template partials before execution:

```go
tpl := mario.MustParse("{{> foo}} baz")
tpl.RegisterPartial("foo", "<span>bar</span>")

result := tpl.MustExec(nil)
fmt.Print(result)
```

Output:

```html
<span>bar</span> baz
```

You can register several partials at once:

```go
tpl := mario.MustParse("{{> foo}} and {{> baz}}")
tpl.RegisterPartials(map[string]string{
    "foo": "<span>bar</span>",
    "baz": "<span>bat</span>",
})

result := tpl.MustExec(nil)
fmt.Print(result)
```

Output:

```html
<span>bar</span> and <span>bat</span>
```


### Global Partials

You can registers global partials that will be accessible by all templates:

```go
mario.RegisterPartial("foo", "<span>bar</span>")

tpl := mario.MustParse("{{> foo}} baz")
result := tpl.MustExec(nil)
fmt.Print(result)
```

Or:

```go
mario.RegisterPartials(map[string]string{
    "foo": "<span>bar</span>",
    "baz": "<span>bat</span>",
})

tpl := mario.MustParse("{{> foo}} and {{> baz}}")
result := tpl.MustExec(nil)
fmt.Print(result)
```


### Dynamic Partials

It's possible to dynamically select the partial to be executed by using sub expression syntax.

For example, that template randomly evaluates the `foo` or `baz` partial:

```go
tpl := mario.MustParse("{{> (whichPartial) }}")
tpl.RegisterPartials(map[string]string{
    "foo": "<span>bar</span>",
    "baz": "<span>bat</span>",
})

ctx := map[string]interface{}{
    "whichPartial": func() string {
        rand.Seed(time.Now().UTC().UnixNano())

        names := []string{"foo", "baz"}
        return names[rand.Intn(len(names))]
    },
}

result := tpl.MustExec(ctx)
fmt.Print(result)
```


### Partial Contexts

It's possible to execute partials on a custom context by passing in the context to the partial call.

For example:

```go
tpl := mario.MustParse("User: {{> userDetails user }}")
tpl.RegisterPartial("userDetails", "{{firstname}} {{lastname}}")

ctx := map[string]interface{}{
    "user": map[string]string{
        "firstname": "Jean",
        "lastname":  "Valjean",
    },
}

result := tpl.MustExec(ctx)
fmt.Print(result)
```

Displays:

```html
User: Jean Valjean
```


### Partial Parameters

Custom data can be passed to partials through hash parameters.

For example:

```go
tpl := mario.MustParse("{{> myPartial name=hero }}")
tpl.RegisterPartial("myPartial", "My hero is {{name}}")

ctx := map[string]interface{}{
    "hero": "Goldorak",
}

result := tpl.MustExec(ctx)
fmt.Print(result)
```

Displays:

```html
My hero is Goldorak
```



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

- [raymond](https://github.com/aymerick/raymond) - golang (original library)
- [handlebars.js](http://handlebarsjs.com) - javascript
- [handlebars.java](https://github.com/jknack/handlebars.java) - java
- [handlebars.rb](https://github.com/cowboyd/handlebars.rb) - ruby
- [handlebars.php](https://github.com/XaminProject/handlebars.php) - php
- [handlebars-objc](https://github.com/Bertrand/handlebars-objc) - Objective C
- [rumblebars](https://github.com/nicolas-cherel/rumblebars) - rust
