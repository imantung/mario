package mario_test

import (
	"testing"

	"github.com/imantung/mario"
	"github.com/stretchr/testify/require"
)

func TestRegisterBuildInHelper(t *testing.T) {
	mario.RegisterHelper("hello", func() string {
		return "hello world"
	})
	defer mario.ResetHelpers()

	template := "{{hello}}"
	ctx := Author{"Alan", "Johnson"}

	require.Equal(t, "hello world", compile(template, ctx))
}
