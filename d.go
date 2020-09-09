package snoc

import (
	"fmt"
	"github.com/strickyak/yak"
	"log"
	"strings"
)

var F = fmt.Sprintf

var InternTable = make(map[string]X)
var NIL = &Null{Sym{S: "nil"}}
var FN = Intern("fn")

func init() {
	InternTable["nil"] = NIL
}

func Intern(s string) X {
	if z, ok := InternTable[s]; ok {
		return z
	}
	z := &Sym{S: s}
	InternTable[s] = z
	return z
}

// XBase

func (o XBase) Null() bool  { return false }
func (o XBase) Atom() bool  { return true }
func (o XBase) Head() X     { return o.Panic(o, "cannot Head") }
func (o XBase) Tail() X     { return o.Panic(o, "cannot Tail") }
func (o XBase) Eq(a X) bool { return o == a }
func (o XBase) Snoc(a X) X  { return o.Panic(o, "cannot Snoc") }
func (o XBase) Get(key X) X { return o.Panic(o, "cannot Get key %q", key) }

func (o XBase) String() string { return "XBase" }
func (o XBase) Panic(rcvr X, format string, args ...interface{}) X {
	var v []interface{}
	v = append(v, rcvr)
	v = append(v, rcvr)
	v = append(v, args...)
	log.Panic("Panic on %T %q: "+format, v)
	return NIL
}

func (o XBase) Eval(env X) X            { return o }
func (o XBase) Apply(args []X, env X) X { return o.Panic(o, "cannot Apply") }

// *Sym

func (o *Sym) Null() bool { return o.S == "nil" }
func (o *Sym) Head() X    { return o.Panic(o, "cannot Head") }
func (o *Sym) Tail() X    { return o.Panic(o, "cannot Tail") }
func (o *Sym) Snoc(a X) X {
	if o.Null() {
		return &Pair{H: a, T: o}
	}
	return o.Panic(o, "cannot Snoc")
}
func (o *Sym) Get(key X) X { return o.Panic(o, "cannot Get key %q", key) }

func (o *Sym) String() string { return o.S }
func (o *Sym) Panic(rcvr X, format string, args ...interface{}) X {
	var v []interface{}
	v = append(v, rcvr)
	v = append(v, rcvr)
	v = append(v, args...)
	log.Panic(o, "Panic on %T %s: "+format, v)
	return NIL
}

func (o *Sym) Eval(env X) X {
	if o.S == "nil" || o.S == "fn" {
		return o
	}
	return env.Get(o)
}

func (o *Sym) Apply(args []X, env X) X { return o.Panic(o, "cannot Apply") }

// *Null

func (o *Null) Null() bool     { return true }
func (o *Null) String() string { return o.S }

// *Pair

func (o *Pair) Null() bool { return false }
func (o *Pair) Atom() bool { return false }
func (o *Pair) Head() X    { return o.H }
func (o *Pair) Tail() X    { return o.T }

//func (o *Pair) Eq(a X) bool { return o == a }
func (o *Pair) Snoc(a X) X {
	return &Pair{H: a, T: o}
}
func (o *Pair) Get(key X) X {
	for p := X(o); p != NIL; p = p.Tail().Tail() {
		if p.Head() == key {
			return p.Tail().Head()
		}
	}
	return o.Panic(o, "cannot Get key %v", key)
}

func (o *Pair) String() string {
	var bb strings.Builder
	bb.WriteByte('(')
	p := X(o)
	for {
		p2, ok := p.(*Pair)
		if !ok {
			log.Panicf("Pair::String: got %T, wanted Pair", p)
		}
		bb.WriteString(p2.H.String())
		if p2.T == NIL {
			break
		}
		bb.WriteByte(' ')
		p = p2.T
	}
	bb.WriteByte(')')
	return bb.String()
}

func (o *Pair) Eval(env X) X {
	f := o.H.Eval(env)
	return f.Apply(ListToVec(o.T), env)
}

func (o *Pair) Apply(args []X, env X) X { // args are unevaluted.
	log.Printf("ApplyPair %q <<<<<< %v", o, args)
	yak.MustEq(o.H, FN)
	vars := ListToVec(o.T.Head())
	yak.MustEq(len(vars), len(args))
	vals := make([]X, len(vars))
	for i, v := range args {
		vals[i] = v.Eval(env)
	}
	e := env
	for i, _ := range vars {
		log.Printf("ApplyPair %v := %v", vars[i], vals[i])
		e = e.Snoc(vals[i])
		e = e.Snoc(vars[i])
	}
	z := o.T.Tail().Head().Eval(e)
	log.Printf("ApplyPair %q >>> (%T)%v", o, z, z)
	return z
}

// *Prim

func (o *Prim) Apply(args []X, env X) X { // args are unevaluted.
	log.Printf("Prim %q <<<<<< %v", o.Name, args)
	evalledArgs := make([]X, len(args))
	for i, a := range args {
		evalledArgs[i] = a.Eval(env)
	}
	log.Printf("Prim %q <<< %v", o.Name, evalledArgs)
	for i, ea := range evalledArgs {
		log.Printf("Prim %q <<< [%d] (%T)%v", o.Name, i, ea, ea)
	}
	z := o.F(evalledArgs, env)
	log.Printf("Prim %q >>> (%T)%v", o.Name, z, z)
	return z
}

func (o *Prim) Eval(env X) X { return o }
func (o *Prim) String() string {
	return F("Prim(%q)", o.Name)
}

// *Special

func (o *Special) Apply(args []X, env X) X { // args are unevaluted.
	log.Printf("Special %q <<< %v", o.Name, args)
	z := o.F(args, env)
	log.Printf("Special %q >>> (%T)%v", o.Name, z, z)
	return z
}
func (o *Special) Eval(env X) X { return o }
func (o *Special) String() string {
	return F("Special(%q)", o.Name)
}

// *Str

func (o *Str) Eval(env X) X { return o }
func (o *Str) String() string {
	return F("%q", o.S)
}

// *Float

func (o *Float) Eval(env X) X { return o }
func (o *Float) String() string {
	return F("%g", o.F)
}

// New

func NewEnv() X {
	z := X(NIL)
	prim := func(name string, f func(args []X, env X) X) {
		z = z.Snoc(&Prim{Name: name, F: f}).Snoc(Intern(name))
	}

	prim("+", func(args []X, env X) X {
		sum := 0.0
		for _, a := range args {
			sum += a.(*Float).F
		}
		return &Float{F: sum}
	})

	z = z.Snoc(ParseText(`(fn (a b) (+ a b))`, "add")[0]).Snoc(Intern("add"))

	return z
}
