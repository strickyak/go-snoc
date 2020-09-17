// b.go: builtins

package snoc

import (
	"math"

	. "github.com/strickyak/yak"
)

var BuiltinSpecials = map[string]func([]Any, Env) Any{
	"quote": func(args []Any, env Env) Any {
		MustLen(args, 1)
		return args[0]
	},
	"and": func(args []Any, env Env) Any {
		z := Any(TRUE)
		for _, a := range args {
			x := Eval(a, env)
			if NullP(x) {
				return NIL
			}
			z = x
		}
		return z
	},
	"or": func(args []Any, env Env) Any {
		for _, a := range args {
			x := Eval(a, env)
			if Bool(x) {
				return x
			}
		}
		return NIL
	},
	"all": func(args []Any, env Env) Any {
		for _, a := range args {
			if NullP(Eval(a, env)) {
				return NIL
			}
		}
		return TRUE
	},
	"any": func(args []Any, env Env) Any {
		for _, a := range args {
			if Bool(Eval(a, env)) {
				return TRUE
			}
		}
		return NIL
	},
	"if": func(args []Any, env Env) Any {
		for len(args) >= 2 {
			pred := Eval(args[0], env)
			if Bool(pred) {
				return Eval(args[1], env)
			}
			args = args[2:]
		}
		MustEq(len(args), 1)
		return Eval(args[0], env)
	},
	"let": func(args []Any, env Env) Any {
		for len(args) >= 2 {
			sym, ok := args[0].(*Sym)
			Must(ok)
			env = env.SnocSnoc(Eval(args[1], env), sym)
			args = args[2:]
		}
		MustEq(len(args), 1)
		return Eval(args[0], env)
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

func TrueNil(b bool) Any {
	if b {
		return TRUE
	} else {
		return NIL
	}
}

var BuiltinPrims = map[string]func([]Any, Env) Any{
	"list": func(args []Any, env Env) Any {
		z := NIL
		for i := len(args) - 1; i >= 0; i-- {
			z = Snoc(z, args[i])
		}
		return z
	},
	"null?": func(args []Any, env Env) Any {
		MustLen(args, 1)
		return TrueNil(NullP(args[0]))
	},
	"atom?": func(args []Any, env Env) Any {
		MustLen(args, 1)
		return TrueNil(AtomP(args[0]))
	},
	"eq": func(args []Any, env Env) Any {
		MustLen(args, 2)
		return TrueNil(Eq(args[0], args[1]))
	},
	"head": func(args []Any, env Env) Any {
		MustLen(args, 1)
		return Head(args[0])
	},
	"tail": func(args []Any, env Env) Any {
		MustLen(args, 1)
		return Tail(args[0])
	},
	"1st": func(args []Any, env Env) Any {
		MustLen(args, 1)
		return Head(args[0])
	},
	"2nd": func(args []Any, env Env) Any {
		MustLen(args, 1)
		return Head(Tail(args[0]))
	},
	"3rd": func(args []Any, env Env) Any {
		MustLen(args, 1)
		return Head(Head(Tail(args[0])))
	},
	"4th": func(args []Any, env Env) Any {
		MustLen(args, 1)
		return Head(Head(Tail(Tail(args[0]))))
	},
	"5th": func(args []Any, env Env) Any {
		MustLen(args, 1)
		return Head(Head(Tail(Tail(Tail(args[0])))))
	},
	"eval": func(args []Any, env Env) Any {
		MustLen(args, 1)
		return Eval(args[0], env)
	},
	"apply": func(args []Any, env Env) Any {
		MustLen(args, 2)
		return Apply(args[0], ListToVec(args[1]), env)
	},
	"snoc": func(args []Any, env Env) Any {
		MustLen(args, 2)
		p, ok := args[0].(*Pair)
		if !ok {
			Throw(args[0], "cannot Snoc")
		}
		return Snoc(p, args[1])
	},
	"cons": func(args []Any, env Env) Any {
		MustLen(args, 2)
		p, ok := args[1].(*Pair)
		if !ok {
			Throw(args[1], "cannot Cons")
		}
		return Snoc(p, args[0])
	},
	"sum": func(args []Any, env Env) Any {
		sum := 0.0
		for _, a := range args {
			sum += ToFloat(a)
		}
		return sum
	},
	"product": func(args []Any, env Env) Any {
		product := 1.0
		for _, a := range args {
			product *= ToFloat(a)
		}
		return product
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
		Globals[k_] = &Prim{Name: k_, F: func(args []Any, env Env) Any {
			MustEq(len(args), 2)
			return fn_(args[0].(float64), args[1].(float64))
		}}
	}
	for k, fn := range BuiltinFloatingRelOps {
		k_, fn_ := k, fn // Capture an inside-loop copy.
		Globals[k_] = &Prim{Name: k_, F: func(args []Any, env Env) Any {
			MustEq(len(args), 2)
			a := args[0].(float64)
			b := args[1].(float64)
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
