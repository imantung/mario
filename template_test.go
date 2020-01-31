package mario_test

import (
	"testing"

	"github.com/imantung/mario"
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

func TestClone(t *testing.T) {
	t.Parallel()

	sourcePartial := `I am a {{wat}} partial`
	sourcePartial2 := `Partial for the {{wat}}`

	tpl, err := mario.Parse(sourceBasic)
	require.NoError(t, err)
	tpl.RegisterPartial("p", sourcePartial)

	if (len(tpl.Partials) != 1) || (tpl.Partials["p"] == nil) {
		t.Errorf("What?")
	}

	cloned := tpl.Clone()

	if (len(cloned.Partials) != 1) || (cloned.Partials["p"] == nil) {
		t.Errorf("Template partials must be cloned")
	}

	cloned.RegisterPartial("p2", sourcePartial2)

	if (len(cloned.Partials) != 2) || (cloned.Partials["p"] == nil) || (cloned.Partials["p2"] == nil) {
		t.Errorf("Failed to register a partial on cloned template")
	}

	if (len(tpl.Partials) != 1) || (tpl.Partials["p"] == nil) {
		t.Errorf("Modification of a cloned template MUST NOT affect original template")
	}
}

func TestTemplate_Exec(t *testing.T) {
	tpl, err := mario.Parse("<h1>{{title}}</h1><p>{{body.content}}</p>")
	require.NoError(t, err)

	output, err := tpl.Exec(map[string]interface{}{
		"title": "foo",
		"body":  map[string]string{"content": "bar"},
	})
	require.NoError(t, err)
	require.Equal(t, `<h1>foo</h1><p>bar</p>`, output)
}

func TestTemplate_MustExec(t *testing.T) {
	tpl, err := mario.Parse("<h1>{{title}}</h1><p>{{body.content}}</p>")
	require.NoError(t, err)

	output, err := tpl.Exec(map[string]interface{}{
		"title": "foo",
		"body":  map[string]string{"content": "bar"},
	})
	require.NoError(t, err)
	require.Equal(t, `<h1>foo</h1><p>bar</p>`, output)
}

func TestTemplate_ExecWith(t *testing.T) {
	// parse template
	tpl, err := mario.Parse("<h1>{{title}}</h1><p>{{#body}}{{content}} and {{@baz.bat}}{{/body}}</p>")
	require.NoError(t, err)

	// computes private data frame
	frame := mario.NewDataFrame()
	frame.Set("baz", map[string]string{"bat": "unicorns"})

	// evaluate template
	output, err := tpl.ExecWith(map[string]interface{}{
		"title": "foo",
		"body":  map[string]string{"content": "bar"},
	}, frame)
	require.NoError(t, err)
	require.Equal(t, `<h1>foo</h1><p>bar and unicorns</p>`, output)
}

func TestTemplate_PrintAST(t *testing.T) {
	tpl, err := mario.Parse("<h1>{{title}}</h1><p>{{#body}}{{content}} and {{@baz.bat}}{{/body}}</p>")
	require.NoError(t, err)
	require.Equal(t,
		"CONTENT[ '<h1>' ]\n{{ PATH:title [] }}\nCONTENT[ '</h1><p>' ]\nBLOCK:\n  PATH:body []\n  PROGRAM:\n    {{     PATH:content []\n }}\n    CONTENT[ ' and ' ]\n    {{     @PATH:baz/bat []\n }}\n  CONTENT[ '</p>' ]\n",
		tpl.PrintAST(),
	)
}
