package mario_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/imantung/mario"
)

type strTest struct {
	name   string
	input  interface{}
	output string
}

var strTests = []strTest{
	{"String", "foo", "foo"},
	{"Boolean true", true, "true"},
	{"Boolean false", false, "false"},
	{"Integer", 25, "25"},
	{"Float", 25.75, "25.75"},
	{"Nil", nil, ""},
	{"[]string", []string{"foo", "bar"}, "foobar"},
	{"[]interface{} (strings)", []interface{}{"foo", "bar"}, "foobar"},
	{"[]Boolean", []bool{true, false}, "truefalse"},
}

func TestStr(t *testing.T) {
	t.Parallel()

	for _, test := range strTests {
		if res := mario.Str(test.input); res != test.output {
			t.Errorf("Failed to stringify: %s\nexpected:\n\t'%s'got:\n\t%q", test.name, test.output, res)
		}
	}
}

func TestStr2(t *testing.T) {
	output := mario.Str(3) + " foos are " + mario.Str(true) + " and " + mario.Str(-1.25) + " bars are " + mario.Str(false) + "\n"
	output += "But you know '" + mario.Str(nil) + "' John Snow\n"
	output += "map: " + mario.Str(map[string]string{"foo": "bar"}) + "\n"
	output += "array: " + mario.Str([]interface{}{true, 10, "foo", 5, "bar"})

	require.Equal(t, `3 foos are true and -1.25 bars are false
But you know '' John Snow
map: map[foo:bar]
array: true10foo5bar`, output)

}

func TestSafeString(t *testing.T) {
	tpl, err := mario.New().Parse("{{em}}")
	require.NoError(t, err)

	tpl.RegisterHelper("em", func() mario.SafeString {
		return mario.SafeString("<em>FOO BAR</em>")
	})

	output, err := tpl.Execute(nil)
	require.NoError(t, err)
	require.Equal(t, `<em>FOO BAR</em>`, output)
}
