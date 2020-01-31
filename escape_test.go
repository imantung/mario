package mario_test

import (
	"testing"

	"github.com/imantung/mario"
	"github.com/stretchr/testify/require"
)

func TestEscape(t *testing.T) {
	tpl, err := mario.New().Parse("{{link url text}}")
	require.NoError(t, err)

	tpl.RegisterHelper("link", func(url string, text string) mario.SafeString {
		return mario.SafeString("<a href='" + mario.Escape(url) + "'>" + mario.Escape(text) + "</a>")
	})

	result, err := tpl.Exec(map[string]string{
		"url":  "http://www.aymerick.com/",
		"text": "This is a <em>cool</em> website",
	})
	require.NoError(t, err)
	require.Equal(t, `<a href='http://www.aymerick.com/'>This is a &lt;em&gt;cool&lt;/em&gt; website</a>`, result)

}
