// b.go: builtins

package snoc

import (
	"math"

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

var BuiltinFloatingRelOps = map[string]func(float64, float64) bool{
	"<":  func(a, b float64) bool { return a < b },
	"<=": func(a, b float64) bool { return a <= b },
	"==": func(a, b float64) bool { return a == b },
	"!=": func(a, b float64) bool { return a != b },
	">":  func(a, b float64) bool { return a > b },
	">=": func(a, b float64) bool { return a >= b },
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
		k_, fn_ := k, fn // Capture an inside-loop copy.
		Globals[k_] = &Prim{Name: k_, F: func(args []X, env Env) X {
			MustEq(len(args), 2)
			return &Float{F: fn_(args[0].(*Float).F, args[1].(*Float).F)}
		}}
	}
	for k, fn := range BuiltinFloatingRelOps {
		k_, fn_ := k, fn // Capture an inside-loop copy.
		Globals[k_] = &Prim{Name: k_, F: func(args []X, env Env) X {
			MustEq(len(args), 2)
			a := args[0].(*Float).F
			b := args[1].(*Float).F
			L("name=%q a=%g b=%g fn=>%v", k_, a, b, fn_(a, b))
			if fn_(a, b) {
				return TRUE
			}
			return NIL
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
