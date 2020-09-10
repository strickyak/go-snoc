package snoc

import (
	"testing"

	. "github.com/strickyak/yak"
)

func defun(env Env, definition string) Env {
	triple := ParseText(definition, "*defun*")
	MustEq(len(triple), 3)
	name, params, body := triple[0], triple[1], triple[2]
	Must(name.(*Sym) != nil)
	lambda := NIL.Snoc(body).Snoc(params).Snoc(FN)
	return env.Snoc(lambda).Snoc(name)
}

func TestParseText(t *testing.T) {
	env := defun(NewEnv(), "add (x y) (+ x y)")
	L("ENV: %v", env)

	xs := ParseText("10 ( + 10 90) (add 100 900)", "Test1")
	for i, x := range xs {
		L("X[%d]: (%T)%v", i, x, x)
		z := x.Eval(env)
		L("Z[%d]: (%T)%v", i, z, z)
	}
}
