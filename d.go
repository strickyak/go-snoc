package snoc

import (
	"flag"
	"fmt"
	"log"
	"runtime/debug"
	"strings"
)

//var Globals = make(map[string]Any)
var InternTable = make(map[string]*Sym)
var FlagVerbose = flag.Bool("v", false, "verbosity")

func Log(format string, args ...interface{}) {
	if *FlagVerbose {
		log.Printf(format, args...)
	}
}

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

// String methods.

func (o *ProtoFunc) String() string {
	if o == nil {
		return "nil"
	}
	if o == o.Outer {
		return "loop"
	}
	return fmt.Sprintf("ProtoFunc{let=%v,name=%q,params=%v,values=%v}", o.IsLet, o.Name, o.Params, o.Values)
}

func (o *Func) String() string {
	if o == nil {
		return "nil"
	}
	return fmt.Sprintf("Func{let=%v,name=%q,params=%v,values=%v,body=<<<%v>>>,outer=%v}", o.IsLet, o.Name, o.Params, o.Values, o.Body, o.Outer)
}

func (o *Var) String() string {
	return fmt.Sprintf("Var{%q[%d]%v%v}", o.Sym.S, o.Slot, o.Proto.Name, o.Proto.Params)
}

func (env *Env) String() string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "Env")
	for p := env; p != nil; p = p.Up {
		fmt.Fprintf(&buf, "{%q", p.Proto.Name)
		for i, e := range p.Slots {
			fmt.Fprintf(&buf, " ")
			fmt.Fprintf(&buf, "[%d]:%v:%T", i, p.Proto.Params[i], e)
		}
		fmt.Fprintf(&buf, "}")
	}
	return buf.String()
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

func Eval(o Any, env *Env) Any {
	Log("EVAL <<< %v ; %v", o, env)
	z := o
SWITCH:
	switch t := o.(type) {
	case nil:
		panic("cannot Eval golang nil")
	case *ProtoFunc:
		z = &Func{
			Proto:  t,
			Outer:  env,      // TODO WRONG?
			Params: t.Params, // omit
			Values: t.Values, // omit
			Body:   t.Body,   // omit
			Name:   t.Name,   // omit
			IsLet:  t.IsLet,  // omit
		}
	case *Var:
		{
			for p := env; p != nil; p = p.Up {
				if p.Proto == t.Proto {
					z = p.Slots[t.Slot]
					break SWITCH
				}
			}
			Throw(o, "cannot Eval Var: %v", t)
		}
	case *Sym:
		{
			g, ok := env.Terp.Globals[t]
			// log.Printf("Globals %q --> (%T) %v, ok=%v", t.S, g, g, ok)
			if !ok {
				Throw(o, "cannot Eval symbol %q with globals %v", t.S, env.Terp.Globals)
			}
			z = g
		}
	case *Pair:
		switch {
		case o == NIL:
			z = NIL // NIL is self-evaluating.
		case t.H == FN:
			if (t.T == NIL) ||
				(t.T.T == NIL) ||
				(t.T.T.T != NIL) {
				Throw(t, "FN must have 3 elements in the list")
			}
			z = EvalLambda(ListToVec(t.T.H), t.T.T.H, env)
		default:
			z = Apply(Eval(t.H, env), ListToVec(t.T), env)
		}
	}

	// Eval never returns a *ProtoFunc; convert it into a Func.
	// TODO Probably this should not be called "Eval".
	if pf, ok := z.(*ProtoFunc); ok {
		z = Eval(pf, env)
	}

	Log("EVAL >>> %v", z)
	return z
}

func EvalLambda(params []Any, body Any, env *Env) Any {
	log.Fatalf("STOP EvalLambda")
	return nil
	/*
		pp := make([]*Sym, len(params))
		for i, e := range params {
			if p, ok := e.(*Sym); ok {
				pp[i] = p
			} else {
				Throw(e, "Fn params must be *Sym")
			}
		}
		return &Func{
			Outer:  env.Owner.Outer,
			Params: pp,
		}
	*/
}

func Apply(o Any, args []Any, env *Env) Any {
	Log("APPLY <<< %v << %v ; %v", o, args, env)
	var z Any
	switch t := o.(type) {
	case nil:
		panic("cannot Apply on golang nil")
	case *Sym:
		z = Throw(t, "cannot Apply on a symbol")
	case *Func:
		z = ApplyFunc(t, args, env)
	case *Prim:
		z = ApplyPrim(t, args, env)
	case *Special:
		z = ApplySpecial(t, args, env)
	default:
		z = Throw(o, "cannot Apply")
	}
	Log("APPLY >>> %v", z)
	return z
}

func ApplyFunc(o *Func, args []Any, env *Env) Any {
	Log("ApplyFunc << %v << %v << %v", o, args, env)
	if o.IsLet {
		if args != nil {
			Throw(o, "apply: got %d args but wanted none because it has Let Values")
		}
	} else {
		if len(args) != len(o.Params) {
			Throw(o, "apply: got %d args but wanted %d", len(args), len(o.Params))
		}
	}

	slots := make([]Any, len(o.Params))
	for i, v := range args { // For the FN case.
		slots[i] = Eval(v, env)
	}

	env2 := &Env{
		Up:    env, // dynamic or scoped?
		Proto: o.Proto,
		Slots: slots,
		Terp:  env.Terp,
	}

	var z Any
	if o.IsLet {
		lenVal := len(o.Values)
		Log("Slots len=%d", lenVal)
		for i, v := range o.Values { // For the LET case.
			Log("Slots[%d/%d] << %v", i, lenVal, v)

			tmp := Eval(v, env2)
			slots[i] = Apply(tmp, nil, env2)

			Log("Slots[%d/%d] >> %v", i, lenVal, slots[i])
		}
		Log("Body == %v; %v", o.Body, env2)
		tmp := Eval(o.Body, env2)
		z = Apply(tmp, nil, env2)
	} else {
		z = Eval(o.Body, env2)
	}
	Log("ApplyFunc >> %v", z)
	return z
}

func ApplyPrim(o *Prim, args []Any, env *Env) Any { // args are unevaluted.
	evalledArgs := make([]Any, len(args))
	for i, a := range args {
		evalledArgs[i] = Eval(a, env)
	}
	return o.F(evalledArgs, env)
}

func ApplySpecial(o *Special, args []Any, env *Env) Any { // args are unevaluted.
	z := o.F(args, env)
	return z
}
