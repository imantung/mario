package mario_test

import (
	"strings"
	"testing"

	"github.com/imantung/mario"
	"github.com/imantung/mario/ast"
	"github.com/stretchr/testify/require"
)

var sourceBasic = `<div class="entry">
  <h1>{{title}}</h1>
  <div class="body">
    {{body}}
  </div>
</div>`

var basicAST = `CONTENT[ '<div class="entry">
  <h1>' ]
{{ PATH:title [] }}
CONTENT[ '</h1>
  <div class="body">
    ' ]
{{ PATH:body [] }}
CONTENT[ '
  </div>
</div>' ]
`

func TestTemplate_Exec(t *testing.T) {
	tpl, err := mario.New().Parse("<h1>{{title}}</h1><p>{{body.content}}</p>")
	require.NoError(t, err)

	var b strings.Builder
	require.NoError(t, tpl.Execute(&b, map[string]interface{}{
		"title": "foo",
		"body":  map[string]string{"content": "bar"},
	}))
	require.Equal(t, `<h1>foo</h1><p>bar</p>`, b.String())
}

func TestTemplate_MustExec(t *testing.T) {
	tpl, err := mario.New().Parse("<h1>{{title}}</h1><p>{{body.content}}</p>")
	require.NoError(t, err)

	var b strings.Builder
	require.NoError(t, tpl.Execute(&b, map[string]interface{}{
		"title": "foo",
		"body":  map[string]string{"content": "bar"},
	}))
	require.Equal(t, `<h1>foo</h1><p>bar</p>`, b.String())
}

func TestTemplate_ExecWith(t *testing.T) {
	// parse template
	tpl, err := mario.New().Parse("<h1>{{title}}</h1><p>{{#body}}{{content}} and {{@baz.bat}}{{/body}}</p>")
	require.NoError(t, err)

	// computes private data frame
	frame := mario.NewDataFrame()
	frame.Set("baz", map[string]string{"bat": "unicorns"})

	var b strings.Builder
	require.NoError(t, tpl.ExecuteWith(&b, map[string]interface{}{
		"title": "foo",
		"body":  map[string]string{"content": "bar"},
	}, frame))
	require.Equal(t, `<h1>foo</h1><p>bar and unicorns</p>`, b.String())
}

func TestTemplate_PrintAST(t *testing.T) {
	tpl, err := mario.New().Parse("<h1>{{title}}</h1><p>{{#body}}{{content}} and {{@baz.bat}}{{/body}}</p>")
	require.NoError(t, err)
	require.Equal(t,
		"CONTENT[ '<h1>' ]\n{{ PATH:title [] }}\nCONTENT[ '</h1><p>' ]\nBLOCK:\n  PATH:body []\n  PROGRAM:\n    {{     PATH:content []\n }}\n    CONTENT[ ' and ' ]\n    {{     @PATH:baz/bat []\n }}\n  CONTENT[ '</p>' ]\n",
		ast.Print(tpl.Program()),
	)
}
