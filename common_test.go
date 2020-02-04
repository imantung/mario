package mario_test

import (
	"strings"

	"github.com/imantung/mario"
)

type testcase struct {
	template      string
	data          interface{}
	helpers       map[string]interface{}
	partials      map[string]string
	expectedError string
	expected      string
}

func compile(template string, ctx interface{}) string {
	var b strings.Builder
	if err := mario.Must(mario.New().Parse(template)).Execute(&b, ctx); err != nil {
		panic(err)
	}
	return b.String()
}
