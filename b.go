// b.go: buildtins

package snoc

import (
	"bufio"
	"io"
	"math"
	"strings"

	. "github.com/strickyak/yak"
)

func (env Env) Snoc(x X) Env {
	return Env{env.Chain.Snoc(x)}
}

func (env Env) Get(k X) X {
	return env.Chain.Get(k)
}

var BuiltinSpecials = map[string]func([]X, Env) X{
	"all": func(args []X, env Env) X {
		for _, a := range args {
			if a.Eval(env).NullP() {
				return NIL
			}
		}
		return TRUE
	},
	"any": func(args []X, env Env) X {
		for _, a := range args {
			if !a.Eval(env).NullP() {
				return TRUE
			}
		}
		return NIL
	},
	"if": func(args []X, env Env) X {
		for len(args) >= 2 {
			pred := args[0].Eval(env)
			if !pred.NullP() {
				return args[1].Eval(env)
			}
			args = args[2:]
		}
		MustEq(len(args), 1)
		return args[0].Eval(env)
	},
	"let": func(args []X, env Env) X {
		for len(args) >= 2 {
			sym, ok := args[0].(*Sym)
			Must(ok)
			env = env.Snoc(args[1].Eval(env)).Snoc(sym)
			args = args[2:]
		}
		MustEq(len(args), 1)
		return args[0].Eval(env)
	},
}

var BuiltinFloatingBinaryOps = map[string]func(float64, float64) float64{
	"+":   func(a, b float64) float64 { return a + b },
	"-":   func(a, b float64) float64 { return a - b },
	"*":   func(a, b float64) float64 { return a * b },
	"div": func(a, b float64) float64 { return a / b },
	"mod": func(a, b float64) float64 { return math.Mod(a, b) },
}

var BuiltinPrims = map[string]func([]X, Env) X{
	"sum": func(args []X, env Env) X {
		sum := 0.0
		for _, a := range args {
			sum += a.(*Float).F
		}
		return &Float{F: sum}
	},
	"product": func(args []X, env Env) X {
		product := 1.0
		for _, a := range args {
			product *= a.(*Float).F
		}
		return &Float{F: product}
	},
}

func init() {
	for k, fn := range BuiltinSpecials {
		Globals[k] = &Special{Name: k, F: fn}
	}
	for k, fn := range BuiltinPrims {
		Globals[k] = &Prim{Name: k, F: fn}
	}
	for k, fn := range BuiltinFloatingBinaryOps {
		Globals[k] = &Prim{Name: k, F: func(args []X, env Env) X {
			return &Float{F: fn(args[0].(*Float).F, args[1].(*Float).F)}
		}}
	}
	Globals["nil"] = NIL
	Globals["fn"] = FN
	Globals["def"] = DEF
	Globals["defun"] = DEFUN
	Globals["true"] = TRUE
}

func NewEnv() Env {
	return Env{NIL}
}

func TryReplParse(s string) (xs []X, ok bool) {
	defer func() {
		r := recover()
		if r != nil {
			ok = false
		}
	}()
	ok = true
	xs = ParseText(s, "*repl*")
	return
}

func TryReplEval(env Env, xs []X) (result X, newenv Env, err interface{}) {
	defer func() {
		err = recover()
	}()
	result = NIL
	for _, x := range xs {
		if p, ok := x.(*Pair); ok {
			if p.H == DEF {
				vec := ListToVec(p.T)
				MustEq(len(vec), 2)
				env = env.Snoc(vec[1]).Snoc(vec[0])
			} else if p.H == DEFUN {
				vec := ListToVec(p.T)
				MustEq(len(vec), 3)
				defun := NIL.Snoc(vec[2]).Snoc(vec[1]).Snoc(FN)
				env = env.Snoc(defun).Snoc(vec[0])
			}
		}
		result = x.Eval(env)
	}
	newenv = env
	return
}

func Repl(env Env, r io.Reader) Env {
	sc := bufio.NewScanner(r)
	var b strings.Builder
	for sc.Scan() {
		b.WriteString(sc.Text())
		s := b.String()
		xs, ok := TryReplParse(s)
		if !ok {
			continue
		}

		result, newenv, err := TryReplEval(env, xs)
		if err != nil {
			L("ERROR: %v", err)
		} else {
			L("====> %v", result)
			env = newenv
		}
	}
	return env
}
