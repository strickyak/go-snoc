package snoc

type Env struct {
	Chain X
}

type X interface {
	Eval(env Env) X
	Apply(args []X, env Env) X // args are unevaluated.

	NullP() bool
	AtomP() bool
	Head() X
	Tail() X
	Eq(a X) bool
	Snoc(a X) X
	Get(key X) X

	String() string
	Panic(rcvr X, fmt string, args ...interface{}) X
}

type XBase struct{}

type Float struct {
	XBase
	F float64
}

type Str struct {
	XBase
	S string
}

type Sym struct {
	XBase
	S string
}

type Null struct {
	Sym
}

type Pair struct {
	XBase
	H, T X
}

type Prim struct {
	XBase
	Name string
	F    func(args []X, env Env) X // args are evaluated.
}

type Special struct {
	XBase
	Name string
	F    func(args []X, env Env) X // args are unevaluated.
}
