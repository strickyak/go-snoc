package snoc

import (
	"fmt"
	"log"
	"runtime/debug"
	"strings"
)

var Globals = make(map[string]Any)
var InternTable = make(map[string]*Sym)

var (
	// Traditionally NIL would be a special *Sym,
	// but in this Lisp it will be a special *Pair.
	// It is not Interned.  The parser will have to know
	// about this special name.  AtomP(NIL) is still true,
	// although it cannot be used as an environment key.
	NIL   = &Pair{}         // Address matters; contents do not.
	FN    = Intern("fn")    // In other Lisps, this is lambda.
	TRUE  = Intern("true")  // In other Lisps, this is T or *T*.
	DEF   = Intern("def")   // Special to the REPL; it modifies the env.
	DEFUN = Intern("defun") // Special to the REPL; it modifies the env.
)

func Intern(s string) *Sym {
	if z, ok := InternTable[s]; ok {
		return z
	}
	z2 := &Sym{S: s}
	InternTable[s] = z2
	return z2
}

// Env

func (env Env) SnocSnoc(x Any, s *Sym) Env {
	return Env{Snoc(Snoc(env.Chain, x), s)}
}

func (env Env) Get(k *Sym) (Any, bool) {
	//log.Printf("GET <<< %q", k.S)
	z, ok := Get(env.Chain, k)
	//log.Printf("GET >>> %#v, ok=%v", z, ok)
	return z, ok
}

// String methods.

func (env Env) String() string {
	return fmt.Sprintf("Env{%v}", Stringify(env.Chain))
}

func (o *Pair) String() string {
	var buf strings.Builder
	firstTime := true
	buf.WriteString("(")
	for p := o; p != NIL; p = p.T {
		if !firstTime {
			buf.WriteByte(' ')
		}
		// fmt.Fprintf(&buf, "%T ", p.H)
		buf.WriteString(Stringify(p.H))
		firstTime = false
	}
	buf.WriteString(")")
	return buf.String()
}

func (o *Sym) String() string {
	return o.S
}

func (o *Prim) String() string {
	return fmt.Sprintf("Prim(%q)", o.Name)
}

func (o *Special) String() string {
	return fmt.Sprintf("Special(%q)", o.Name)
}

// Global Lispy Funcs

func NullP(o Any) bool {
	switch t := o.(type) {
	case *Pair:
		return t == NIL
	}
	return false
}
func AtomP(o Any) bool {
	switch t := o.(type) {
	case *Pair:
		return t == NIL
	}
	return true
}
func Head(o Any) Any {
	switch t := o.(type) {
	case *Pair:
		return t.H
	}
	return Throw(o, "cannot Head")
}
func Tail(o Any) Any {
	switch t := o.(type) {
	case *Pair:
		return t.T
	}
	return Throw(o, "cannot Tail")
}
func Eq(o Any, a Any) bool {
	switch t := o.(type) {
	case *Sym:
		if b, ok := a.(*Sym); ok {
			return t == b
		}
	case string:
		if b, ok := a.(string); ok {
			return t == b
		}
	case int:
		if b, ok := a.(int); ok {
			return t == b
		}
	case float64:
		if b, ok := a.(float64); ok {
			return t == b
		}
	case *Pair:
		if b, ok := a.(*Pair); ok {
			return t == b
		}
	}
	return false
}
func Snoc(o *Pair, a Any) *Pair {
	return &Pair{H: a, T: o}
}
func Get(o *Pair, key *Sym) (Any, bool) {
	for o != NIL {
		if o.H.(*Sym) == key {
			return o.T.H, true
		}
		o = o.T.T
	}
	return nil, false
}
func ToInt(o Any) int {
	switch t := o.(type) {
	case int:
		return t
	}
	Throw(o, "cannot Int")
	return 0
}
func ToFloat(o Any) float64 {
	switch t := o.(type) {
	case float64:
		return t
	}
	Throw(o, "cannot Float")
	return 0
}
func ToStr(o Any) string {
	switch t := o.(type) {
	case string:
		return t
	}
	Throw(o, "cannot Str")
	return ""
}
func Bool(o Any) bool {
	switch t := o.(type) {
	case *Pair:
		return t != NIL
	}
	return true
}

func Stringify(o Any) string {
	switch t := o.(type) {
	case *Pair:
		return t.String()
	case *Sym:
		return t.String()
	}
	return fmt.Sprintf("%v", o)
}

func Throw(o Any, format string, args ...interface{}) Any {
	debug.PrintStack()
	var v []interface{}
	v = append(v, o)
	v = append(v, o)
	v = append(v, args...)
	log.Panicf("Exception on (%T) %v: "+format, v...)
	return NIL
}

func Eval(o Any, env Env) Any {
	//log.Printf("EVAL <<< %v ; %v", o, env)
	z := o
	switch t := o.(type) {
	case nil:
		z = o
	case *Sym:
		z2, zok := env.Get(t)
		if zok {
			z = z2
		} else {
			g, gok := Globals[t.S]
			// log.Printf("Globals %q --> (%T) %v, ok=%v", t.S, g, g, gok)
			if !gok {
				Throw(o, "cannot Eval")
			}
			z = g
		}
	case *Pair:
		if o == NIL {
			z = NIL
		} else {
			z = Apply(Eval(t.H, env), ListToVec(t.T), env)
		}
	}
	//log.Printf("EVAL >>> %v", z)
	return z
}

func Apply(o Any, args []Any, env Env) Any {
	//log.Printf("APPLY <<< %v << %v ; %v", o, args, env)
	var z Any
	switch t := o.(type) {
	case nil:
		z = Throw(t, "cannot Apply list where first is nil")
	case *Sym:
		z = Throw(t, "cannot Apply list where first is *Sym")
	case *Pair:
		z = ApplyPair(t, args, env)
	case *Prim:
		z = ApplyPrim(t, args, env)
	case *Special:
		z = ApplySpecial(t, args, env)
	default:
		z = Throw(o, "cannot Apply")
	}
	//log.Printf("APPLY >>> %v", z)
	return z
}

func ApplyPair(o *Pair, args []Any, env Env) Any {
	if o == NIL {
		Throw(o, "cannot Apply if list is nil")
	}

	if !Eq(o.H, FN) {
		Throw(o, "cannot Apply if first is not FN")
	}

	vars := ListToVec(o.T.H)
	if len(vars) != len(args) {
		Throw(o, "apply: got %d args but wanted %d", len(args), len(vars))
	}

	vals := make([]Any, len(vars))
	for i, v := range args {
		vals[i] = Eval(v, env)
	}
	for i, _ := range vars {
		env = env.SnocSnoc(vals[i], vars[i].(*Sym))
	}
	return Eval(o.T.T.H, env)
}

func ApplyPrim(o *Prim, args []Any, env Env) Any { // args are unevaluted.
	evalledArgs := make([]Any, len(args))
	for i, a := range args {
		evalledArgs[i] = Eval(a, env)
	}
	return o.F(evalledArgs, env)
}

func ApplySpecial(o *Special, args []Any, env Env) Any { // args are unevaluted.
	z := o.F(args, env)
	return z
}
