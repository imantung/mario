package mario_test

import (
	"strings"
	"testing"

	"github.com/imantung/mario"
	"github.com/stretchr/testify/require"
)

func TestRegisterBuildInHelper(t *testing.T) {
	mario.RegisterHelper("hello", func() string {
		return "hello world"
	})
	defer mario.ResetHelpers()

	tpl := mario.Must(mario.New().Parse(`{{hello}}`))
	var b strings.Builder
	require.NoError(t, tpl.Execute(&b, Author{"Alan", "Johnson"}))
	require.Equal(t, "hello world", b.String())
}
