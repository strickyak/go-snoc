package snoc

import (
	"testing"
)

func defun(env Env, definition string) Env {
	triple := ParseText(definition, "*defun*")
	if len(triple) != 3 {
		panic(triple)
	}
	name, params, body := triple[0], triple[1], triple[2]
	nameSym, ok := name.(*Sym)
	if !ok {
		panic(name)
	}
	lambda := Snoc(Snoc(Snoc(NIL, body), params), FN)
	return env.SnocSnoc(lambda, nameSym)
}

func TestParseText(t *testing.T) {
	env := defun(NewEnv(), "add (x y) (+ x y)")
	t.Logf("ENV: %v", env)

	xs := ParseText("10 ( + 10 90) (add 100 900)", "Test1")
	for i, x := range xs {
		t.Logf("X[%d]: (%T)%v", i, x, x)
		z := Eval(x, env)
		t.Logf("Z[%d]: (%T)%v", i, z, z)
	}
}
