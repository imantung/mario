package mario_test

import (
	"testing"

	"github.com/imantung/mario"
	"github.com/stretchr/testify/require"
)

func TestIsTrue(t *testing.T) {
	testcases := []struct {
		obj   interface{}
		truth bool
	}{
		{[0]string{}, false},
		{[1]string{"foo"}, true},
		{[]string{}, false},
		{[]string{"foo"}, true},
		{map[string]string{}, false},
		{map[string]string{"foo": "bar"}, true},
		{"", false},
		{"foo", true},
		{true, true},
		{false, false},
		{0, false},
		{10, true},
		{-10, true},
		{0.0, false},
		{10.0, true},
		{struct{}{}, true},
		{nil, false},
	}

	for _, tt := range testcases {
		require.Equal(t, tt.truth, mario.IsTrue(tt.obj), tt.obj)
	}
}

func TestStr(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		input  interface{}
		output string
	}{
		{"foo", "foo"},
		{true, "true"},
		{false, "false"},
		{25, "25"},
		{25.75, "25.75"},
		{nil, ""},
		{[]string{"foo", "bar"}, "foobar"},
		{[]interface{}{"foo", "bar"}, "foobar"},
		{[]bool{true, false}, "truefalse"},
	}

	for _, tt := range testcases {
		require.Equal(t, tt.output, mario.Str(tt.input))
	}
}

func TestSafeString(t *testing.T) {
	tpl, err := mario.New().
		WithHelperFunc("em", func() mario.SafeString {
			return mario.SafeString("<em>FOO BAR</em>")
		}).
		Parse("{{em}}")
	require.NoError(t, err)

	output, err := tpl.Execute(nil)
	require.NoError(t, err)
	require.Equal(t, `<em>FOO BAR</em>`, output)
}
